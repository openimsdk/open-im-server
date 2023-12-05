package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"
	mongoModel "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	mysqlModel "github.com/openimsdk/open-im-server/v3/tools/data-conversion/openim/mysql/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"strconv"
)

const (
	versionTable = "dataver"
	versionKey   = "data_version"
	versionValue = 35
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "path config file")
	flag.Parse()
	if err := Main(path); err != nil {
		log.Fatal(err)
		return
	}
	os.Exit(0)
}

func InitConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &config.Config)
}

func GetMysql() (*gorm.DB, error) {
	conf := config.Config.Mysql
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Address[0], conf.Database)
	return gorm.Open(gormMysql.Open(mysqlDSN), &gorm.Config{Logger: logger.Discard})
}

func GetMongo() (*mongo.Database, error) {
	mgo, err := unrelation.NewMongo()
	if err != nil {
		return nil, err
	}
	return mgo.GetDatabase(), nil
}

func Main(path string) error {
	if err := InitConfig(path); err != nil {
		return err
	}
	if config.Config.Mysql == nil {
		return nil
	}
	mongoDB, err := GetMongo()
	if err != nil {
		return err
	}
	var version struct {
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}
	switch mongoDB.Collection(versionTable).FindOne(context.Background(), bson.M{"key": versionKey}).Decode(&version) {
	case nil:
		if ver, _ := strconv.Atoi(version.Value); ver >= versionValue {
			return nil
		}
	case mongo.ErrNoDocuments:
	default:
		return err
	}
	mysqlDB, err := GetMysql()
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1049 {
			if err := SetMongoDataVersion(mongoDB, version.Value); err != nil {
				return err
			}
			return nil // database not exist
		}
		return err
	}

	var c convert
	var tasks []func() error
	tasks = append(tasks,
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewUserMongo, c.User) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewFriendMongo, c.Friend) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewFriendRequestMongo, c.FriendRequest) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewBlackMongo, c.Black) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewGroupMongo, c.Group) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewGroupMember, c.GroupMember) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewGroupRequestMgo, c.GroupRequest) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewConversationMongo, c.Conversation) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewS3Mongo, c.Object(config.Config.Object.Enable)) },
		func() error { return NewTask(mysqlDB, mongoDB, mgo.NewLogMongo, c.Log) },
	)

	for _, task := range tasks {
		if err := task(); err != nil {
			return err
		}
	}

	if err := SetMongoDataVersion(mongoDB, version.Value); err != nil {
		return err
	}
	return nil
}

func SetMongoDataVersion(db *mongo.Database, curver string) error {
	filter := bson.M{"key": versionKey, "value": curver}
	update := bson.M{"$set": bson.M{"key": versionKey, "value": strconv.Itoa(versionValue)}}
	_, err := db.Collection(versionTable).UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	return err
}

// NewTask A mysql table B mongodb model C mongodb table
func NewTask[A interface{ TableName() string }, B any, C any](gormDB *gorm.DB, mongoDB *mongo.Database, mongoDBInit func(db *mongo.Database) (B, error), convert func(v A) C) error {
	obj, err := mongoDBInit(mongoDB)
	if err != nil {
		return err
	}
	var zero A
	tableName := zero.TableName()
	coll, err := getColl(obj)
	if err != nil {
		return fmt.Errorf("get mongo collection %s failed, err: %w", tableName, err)
	}
	var count int
	defer func() {
		log.Printf("completed convert %s total %d\n", tableName, count)
	}()
	const batch = 100
	for page := 0; ; page++ {
		res := make([]A, 0, batch)
		if err := gormDB.Limit(batch).Offset(page * batch).Find(&res).Error; err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1146 {
				return nil // table not exist
			}
			return fmt.Errorf("find mysql table %s failed, err: %w", tableName, err)
		}
		if len(res) == 0 {
			return nil
		}
		temp := make([]any, len(res))
		for i := range res {
			temp[i] = convert(res[i])
		}
		if err := insertMany(coll, temp); err != nil {
			return fmt.Errorf("insert mongo table %s failed, err: %w", tableName, err)
		}
		count += len(res)
		if len(res) < batch {
			return nil
		}
		log.Printf("current convert %s completed %d\n", tableName, count)
	}
}

func insertMany(coll *mongo.Collection, objs []any) error {
	if _, err := coll.InsertMany(context.Background(), objs); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}
	for i := range objs {
		_, err := coll.InsertOne(context.Background(), objs[i])
		switch {
		case err == nil:
		case mongo.IsDuplicateKeyError(err):
		default:
			return err
		}
	}
	return nil
}

func getColl(obj any) (_ *mongo.Collection, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("not found %+v", e)
		}
	}()
	stu := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(&mongo.Collection{}).String()
	for i := 0; i < stu.NumField(); i++ {
		field := stu.Field(i)
		if field.Type().String() == typ {
			return (*mongo.Collection)(field.UnsafePointer()), nil
		}
	}
	return nil, errors.New("not found")
}

type convert struct{}

func (convert) User(v mysqlModel.UserModel) mongoModel.UserModel {
	return mongoModel.UserModel{
		UserID:           v.UserID,
		Nickname:         v.Nickname,
		FaceURL:          v.FaceURL,
		Ex:               v.Ex,
		AppMangerLevel:   v.AppMangerLevel,
		GlobalRecvMsgOpt: v.GlobalRecvMsgOpt,
		CreateTime:       v.CreateTime,
	}
}

func (convert) Friend(v mysqlModel.FriendModel) mongoModel.FriendModel {
	return mongoModel.FriendModel{
		OwnerUserID:    v.OwnerUserID,
		FriendUserID:   v.FriendUserID,
		Remark:         v.Remark,
		CreateTime:     v.CreateTime,
		AddSource:      v.AddSource,
		OperatorUserID: v.OperatorUserID,
		Ex:             v.Ex,
	}
}

func (convert) FriendRequest(v mysqlModel.FriendRequestModel) mongoModel.FriendRequestModel {
	return mongoModel.FriendRequestModel{
		FromUserID:    v.FromUserID,
		ToUserID:      v.ToUserID,
		HandleResult:  v.HandleResult,
		ReqMsg:        v.ReqMsg,
		CreateTime:    v.CreateTime,
		HandlerUserID: v.HandlerUserID,
		HandleMsg:     v.HandleMsg,
		HandleTime:    v.HandleTime,
		Ex:            v.Ex,
	}
}

func (convert) Black(v mysqlModel.BlackModel) mongoModel.BlackModel {
	return mongoModel.BlackModel{
		OwnerUserID:    v.OwnerUserID,
		BlockUserID:    v.BlockUserID,
		CreateTime:     v.CreateTime,
		AddSource:      v.AddSource,
		OperatorUserID: v.OperatorUserID,
		Ex:             v.Ex,
	}
}

func (convert) Group(v mysqlModel.GroupModel) mongoModel.GroupModel {
	return mongoModel.GroupModel{
		GroupID:                v.GroupID,
		GroupName:              v.GroupName,
		Notification:           v.Notification,
		Introduction:           v.Introduction,
		FaceURL:                v.FaceURL,
		CreateTime:             v.CreateTime,
		Ex:                     v.Ex,
		Status:                 v.Status,
		CreatorUserID:          v.CreatorUserID,
		GroupType:              v.GroupType,
		NeedVerification:       v.NeedVerification,
		LookMemberInfo:         v.LookMemberInfo,
		ApplyMemberFriend:      v.ApplyMemberFriend,
		NotificationUpdateTime: v.NotificationUpdateTime,
		NotificationUserID:     v.NotificationUserID,
	}
}

func (convert) GroupMember(v mysqlModel.GroupMemberModel) mongoModel.GroupMemberModel {
	return mongoModel.GroupMemberModel{
		GroupID:        v.GroupID,
		UserID:         v.UserID,
		Nickname:       v.Nickname,
		FaceURL:        v.FaceURL,
		RoleLevel:      v.RoleLevel,
		JoinTime:       v.JoinTime,
		JoinSource:     v.JoinSource,
		InviterUserID:  v.InviterUserID,
		OperatorUserID: v.OperatorUserID,
		MuteEndTime:    v.MuteEndTime,
		Ex:             v.Ex,
	}
}

func (convert) GroupRequest(v mysqlModel.GroupRequestModel) mongoModel.GroupRequestModel {
	return mongoModel.GroupRequestModel{
		UserID:        v.UserID,
		GroupID:       v.GroupID,
		HandleResult:  v.HandleResult,
		ReqMsg:        v.ReqMsg,
		HandledMsg:    v.HandledMsg,
		ReqTime:       v.ReqTime,
		HandleUserID:  v.HandleUserID,
		HandledTime:   v.HandledTime,
		JoinSource:    v.JoinSource,
		InviterUserID: v.InviterUserID,
		Ex:            v.Ex,
	}
}

func (convert) Conversation(v mysqlModel.ConversationModel) mongoModel.ConversationModel {
	return mongoModel.ConversationModel{
		OwnerUserID:           v.OwnerUserID,
		ConversationID:        v.ConversationID,
		ConversationType:      v.ConversationType,
		UserID:                v.UserID,
		GroupID:               v.GroupID,
		RecvMsgOpt:            v.RecvMsgOpt,
		IsPinned:              v.IsPinned,
		IsPrivateChat:         v.IsPrivateChat,
		BurnDuration:          v.BurnDuration,
		GroupAtType:           v.GroupAtType,
		AttachedInfo:          v.AttachedInfo,
		Ex:                    v.Ex,
		MaxSeq:                v.MaxSeq,
		MinSeq:                v.MinSeq,
		CreateTime:            v.CreateTime,
		IsMsgDestruct:         v.IsMsgDestruct,
		MsgDestructTime:       v.MsgDestructTime,
		LatestMsgDestructTime: v.LatestMsgDestructTime,
	}
}

func (convert) Object(engine string) func(v mysqlModel.ObjectModel) mongoModel.ObjectModel {
	return func(v mysqlModel.ObjectModel) mongoModel.ObjectModel {
		return mongoModel.ObjectModel{
			Name:        v.Name,
			UserID:      v.UserID,
			Hash:        v.Hash,
			Engine:      engine,
			Key:         v.Key,
			Size:        v.Size,
			ContentType: v.ContentType,
			Group:       v.Cause,
			CreateTime:  v.CreateTime,
		}
	}
}

func (convert) Log(v mysqlModel.Log) mongoModel.LogModel {
	return mongoModel.LogModel{
		LogID:      v.LogID,
		Platform:   v.Platform,
		UserID:     v.UserID,
		CreateTime: v.CreateTime,
		Url:        v.Url,
		FileName:   v.FileName,
		SystemType: v.SystemType,
		Version:    v.Version,
		Ex:         v.Ex,
	}
}

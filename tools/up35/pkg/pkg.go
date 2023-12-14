package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	rtcMgo "github.com/openimsdk/open-im-server/v3/tools/up35/pkg/internal/rtc/mongo/mgo"
)

const (
	versionTable = "dataver"
	versionKey   = "data_version"
	versionValue = 35
)

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

		func() error { return NewTask(mysqlDB, mongoDB, rtcMgo.NewSignal, c.SignalModel) },
		func() error { return NewTask(mysqlDB, mongoDB, rtcMgo.NewSignalInvitation, c.SignalInvitationModel) },
		func() error { return NewTask(mysqlDB, mongoDB, rtcMgo.NewMeeting, c.Meeting) },
		func() error { return NewTask(mysqlDB, mongoDB, rtcMgo.NewMeetingInvitation, c.MeetingInvitationInfo) },
		func() error { return NewTask(mysqlDB, mongoDB, rtcMgo.NewMeetingRecord, c.MeetingVideoRecord) },
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

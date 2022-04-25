package db

import (
	"Open_IM/pkg/common/config"
	"go.mongodb.org/mongo-driver/bson"

	//"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"

	//	"context"
	//	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"time"

	"context"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB DataBases

type DataBases struct {
	MysqlDB     mysqlDB
	mgoSession  *mgo.Session
	redisPool   *redis.Pool
	mongoClient *mongo.Client
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

func init() {
	//log.NewPrivateLog(constant.LogFileName)
	var mongoClient *mongo.Client
	var err1 error
	//mysql init
	initMysqlDB()
	// mongo init
	// "mongodb://sysop:moon@localhost/records"
	uri := "mongodb://sample.host:27017/?maxPoolSize=20&w=majority"
	if config.Config.Mongo.DBUri != "" {
		// example: mongodb://$user:$password@mongo1.mongo:27017,mongo2.mongo:27017,mongo3.mongo:27017/$DBDatabase/?replicaSet=rs0&readPreference=secondary&authSource=admin&maxPoolSize=$DBMaxPoolSize
		uri = config.Config.Mongo.DBUri
	} else {
		if config.Config.Mongo.DBPassword != "" && config.Config.Mongo.DBUserName != "" {
			uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d", config.Config.Mongo.DBUserName, config.Config.Mongo.DBPassword, config.Config.Mongo.DBAddress[0],
				config.Config.Mongo.DBDatabase, config.Config.Mongo.DBMaxPoolSize)
		} else {
			uri = fmt.Sprintf("mongodb://%s/%s/?maxPoolSize=%d",
				config.Config.Mongo.DBAddress[0], config.Config.Mongo.DBDatabase,
				config.Config.Mongo.DBMaxPoolSize)
		}
	}
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(" mongo.Connect  failed, try ", utils.GetSelfFuncName(), err.Error(), uri)
		time.Sleep(time.Duration(30) * time.Second)
		mongoClient, err1 = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err1 != nil {
			fmt.Println(" mongo.Connect retry failed, panic", err.Error(), uri)
			panic(err1.Error())
		}
	}
	fmt.Println("0", utils.GetSelfFuncName(), "mongo driver client init success: ", uri)
	// mongodb create index
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	dataBase := mongoClient.Database(config.Config.Mongo.DBDatabase)

	cCommentMsgModels := []mongo.IndexModel{
		{
			Keys: bson.A{bson.D{{"create_time", -1}, {"user_id", -1}}},
		},
	}
	result, err := dataBase.Collection(cCommentMsg).Indexes().CreateMany(context.Background(), cCommentMsgModels, opts)
	if err != nil {
		fmt.Println("mongodb create cCommentMsgModels failed", result, err.Error())
	}

	cSendLogModels := []mongo.IndexModel{
		{
			Keys: bson.M{"user_id": -1, "send_time": -1},
		},
	}
	result, err = dataBase.Collection(cSendLog).Indexes().CreateMany(context.Background(), cSendLogModels, opts)
	if err != nil {
		fmt.Println("mongodb create cSendLogModels failed", result, err.Error())
	}

	cChatModels := []mongo.IndexModel{
		{
			Keys: bson.M{"uid": -1},
		},
	}
	result, err = dataBase.Collection(cChat).Indexes().CreateMany(context.Background(), cChatModels, opts)
	if err != nil {
		fmt.Println("mongodb create cChatModels failed", result, err.Error())
	}

	cWorkMomentModels := []mongo.IndexModel{
		{
			Keys: bson.M{"work_moment_id": -1},
		},
		{
			Keys: bson.M{"user_id": -1, "create_time": -1},
		},
	}
	result, err = dataBase.Collection(cWorkMoment).Indexes().CreateMany(context.Background(), cWorkMomentModels, opts)
	if err != nil {
		fmt.Println("mongodb create cWorkMomentModels failed", result, err.Error())
	}

	cTagModels := []mongo.IndexModel{
		{
			Keys: bson.M{"tag_id": -1},
		},
		{
			Keys: bson.M{"user_id": -1, "tag_id": -1},
		},
	}
	result, err = dataBase.Collection(cTag).Indexes().CreateMany(context.Background(), cTagModels, opts)
	if err != nil {
		fmt.Println("mongodb create cTagModels failed", result, err.Error())
	}
	DB.mongoClient = mongoClient

	// redis pool init
	DB.redisPool = &redis.Pool{
		MaxIdle:     config.Config.Redis.DBMaxIdle,
		MaxActive:   config.Config.Redis.DBMaxActive,
		IdleTimeout: time.Duration(config.Config.Redis.DBIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				config.Config.Redis.DBAddress,
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialDatabase(0),
				redis.DialPassword(config.Config.Redis.DBPassWord),
			)
		},
	}
}

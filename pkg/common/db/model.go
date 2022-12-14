package db

import (
	"Open_IM/pkg/common/config"
	"github.com/dtm-labs/rockscache"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"

	//"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
	go_redis "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.in/mgo.v2"
	"time"

	"context"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
	//go_redis "github.com/go-redis/redis/v8"
)

var DB DataBases

type DataBases struct {
	MysqlDB    mysqlDB
	mgoSession *mgo.Session
	//redisPool   *redis.Pool
	mongoClient *mongo.Client
	RDB         go_redis.UniversalClient
	Rc          *rockscache.Client
	WeakRc      *rockscache.Client
}

type RedisClient struct {
	client  *go_redis.Client
	cluster *go_redis.ClusterClient
	go_redis.UniversalClient
	enableCluster bool
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

func init() {
	var mongoClient *mongo.Client
	var err1 error
	fmt.Println("init mysql redis mongo ")
	defer fmt.Println("init mysql redis mongo ok ")
	initMysqlDB()
	// mongo init
	// "mongodb://sysop:moon@localhost/records"
	uri := "mongodb://sample.host:27017/?maxPoolSize=20&w=majority"
	if config.Config.Mongo.DBUri != "" {
		// example: mongodb://$user:$password@mongo1.mongo:27017,mongo2.mongo:27017,mongo3.mongo:27017/$DBDatabase/?replicaSet=rs0&readPreference=secondary&authSource=admin&maxPoolSize=$DBMaxPoolSize
		uri = config.Config.Mongo.DBUri
	} else {
		//mongodb://mongodb1.example.com:27317,mongodb2.example.com:27017/?replicaSet=mySet&authSource=authDB
		mongodbHosts := ""
		for i, v := range config.Config.Mongo.DBAddress {
			if i == len(config.Config.Mongo.DBAddress)-1 {
				mongodbHosts += v
			} else {
				mongodbHosts += v + ","
			}
		}

		if config.Config.Mongo.DBPassword != "" && config.Config.Mongo.DBUserName != "" {
			// clientOpts := options.Client().ApplyURI("mongodb://localhost:27017,localhost:27018/?replicaSet=replset")
			//mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]
			//uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d&authSource=admin&replicaSet=replset",
			uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d&authSource=admin",
				config.Config.Mongo.DBUserName, config.Config.Mongo.DBPassword, mongodbHosts,
				config.Config.Mongo.DBDatabase, config.Config.Mongo.DBMaxPoolSize)
		} else {
			uri = fmt.Sprintf("mongodb://%s/%s/?maxPoolSize=%d&authSource=admin",
				mongodbHosts, config.Config.Mongo.DBDatabase,
				config.Config.Mongo.DBMaxPoolSize)
		}
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		mongoClient, err1 = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err1 != nil {
			panic(err1.Error() + " mongo.Connect failed " + uri)
		}
	}
	// mongodb create index
	if err := createMongoIndex(mongoClient, cSendLog, false, "send_id", "-send_time"); err != nil {
		panic(err.Error() + " index create failed " + cSendLog + " send_id, -send_time")
	}
	if err := createMongoIndex(mongoClient, cChat, false, "uid"); err != nil {
		fmt.Println(err.Error() + " index create failed " + cChat + " uid ")
	}
	if err := createMongoIndex(mongoClient, cWorkMoment, true, "-create_time", "work_moment_id"); err != nil {
		panic(err.Error() + "index create failed " + cWorkMoment + " -create_time, work_moment_id")
	}
	if err := createMongoIndex(mongoClient, cWorkMoment, true, "work_moment_id"); err != nil {
		panic(err.Error() + "index create failed " + cWorkMoment + " work_moment_id ")
	}
	if err := createMongoIndex(mongoClient, cWorkMoment, false, "user_id", "-create_time"); err != nil {
		panic(err.Error() + "index create failed " + cWorkMoment + "user_id, -create_time")
	}
	if err := createMongoIndex(mongoClient, cTag, false, "user_id", "-create_time"); err != nil {
		panic(err.Error() + "index create failed " + cTag + " user_id, -create_time")
	}
	if err := createMongoIndex(mongoClient, cTag, true, "tag_id"); err != nil {
		panic(err.Error() + "index create failed " + cTag + " tag_id")
	}
	DB.mongoClient = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if config.Config.Redis.EnableCluster {
		DB.RDB = go_redis.NewClusterClient(&go_redis.ClusterOptions{
			Addrs:    config.Config.Redis.DBAddress,
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			PoolSize: 50,
		})
		_, err = DB.RDB.Ping(ctx).Result()
		if err != nil {
			fmt.Println("redis cluster failed address ", config.Config.Redis.DBAddress)
			panic(err.Error() + " redis cluster " + config.Config.Redis.DBUserName + config.Config.Redis.DBPassWord)
		}
	} else {
		DB.RDB = go_redis.NewClient(&go_redis.Options{
			Addr:     config.Config.Redis.DBAddress[0],
			Username: config.Config.Redis.DBUserName,
			Password: config.Config.Redis.DBPassWord, // no password set
			DB:       0,                              // use default DB
			PoolSize: 100,                            // 连接池大小
		})
		_, err = DB.RDB.Ping(ctx).Result()
		if err != nil {
			panic(err.Error() + " redis " + config.Config.Redis.DBAddress[0] + config.Config.Redis.DBUserName + config.Config.Redis.DBPassWord)
		}
	}
	// 强一致性缓存，当一个key被标记删除，其他请求线程会被锁住轮询直到新的key生成，适合各种同步的拉取, 如果弱一致可能导致拉取还是老数据，毫无意义
	DB.Rc = rockscache.NewClient(DB.RDB, rockscache.NewDefaultOptions())
	DB.Rc.Options.StrongConsistency = true

	// 弱一致性缓存，当一个key被标记删除，其他请求线程直接返回该key的value，适合高频并且生成很缓存很慢的情况 如大群发消息缓存的缓存
	DB.WeakRc = rockscache.NewClient(DB.RDB, rockscache.NewDefaultOptions())
	DB.WeakRc.Options.StrongConsistency = false
}

func createMongoIndex(client *mongo.Client, collection string, isUnique bool, keys ...string) error {
	db := client.Database(config.Config.Mongo.DBDatabase).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := db.Indexes()
	keysDoc := bsonx.Doc{}

	// 复合索引
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	// 创建索引
	index := mongo.IndexModel{
		Keys: keysDoc,
	}
	if isUnique == true {
		index.Options = options.Index().SetUnique(true)
	}
	result, err := indexView.CreateOne(
		context.Background(),
		index,
		opts,
	)
	if err != nil {
		return utils.Wrap(err, result)
	}
	return nil
}

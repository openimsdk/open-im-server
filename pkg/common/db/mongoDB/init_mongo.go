package mongoDB

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"strings"
	"time"
)

type Client struct {
	mongo *mongo.Client
}

func NewMongoClient(mdb *mongo.Client) *Client {
	return &Client{mongo: mdb}
}

func initMongo() *mongo.Client {
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
	log.Println("start to init mongoDB:", uri)
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err.Error() + " mongo.Connect failed " + uri)
		}
	}
	// mongodb create index
	if err := createMongoIndex(mongoClient, cSendLog, false, "send_id", "-send_time"); err != nil {
		panic(err.Error() + " index create failed " + cSendLog + " send_id, -send_time")
	}
	if err := createMongoIndex(mongoClient, cChat, false, "uid"); err != nil {
		fmt.Println(err.Error() + " index create failed " + cChat + " uid, please create index by yourself in field uid")
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
	return mongoClient
}

func createMongoIndex(client *mongo.Client, collection string, isUnique bool, keys ...string) error {
	db := client.Database(config.Config.Mongo.DBDatabase).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	indexView := db.Indexes()
	keysDoc := bsonx.Doc{}

	// create composite indexes
	for _, key := range keys {
		if strings.HasPrefix(key, "-") {
			keysDoc = keysDoc.Append(strings.TrimLeft(key, "-"), bsonx.Int32(-1))
		} else {
			keysDoc = keysDoc.Append(key, bsonx.Int32(1))
		}
	}

	// create index
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

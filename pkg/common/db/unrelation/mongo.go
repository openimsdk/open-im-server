// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unrelation

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mw/specialerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/unrelation"
)

const (
	maxRetry         = 10 // number of retries
	mongoConnTimeout = 10 * time.Second
)

type Mongo struct {
	db        *mongo.Client
	mongoConf *config.Mongo
}

// NewMongo Initialize MongoDB connection.
func NewMongo(mongoConf *config.Mongo) (*Mongo, error) {
	specialerror.AddReplace(mongo.ErrNoDocuments, errs.ErrRecordNotFound)
	uri := buildMongoURI(mongoConf)

	var mongoClient *mongo.Client
	var err error

	// Retry connecting to MongoDB
	for i := 0; i <= maxRetry; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), mongoConnTimeout)
		defer cancel()
		mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err == nil {
			if err = mongoClient.Ping(ctx, nil); err != nil {
				return nil, errs.Wrap(err, uri)
			}
			return &Mongo{db: mongoClient, mongoConf: mongoConf}, nil
		}
		if shouldRetry(err) {
			time.Sleep(time.Second) // exponential backoff could be implemented here
			continue
		}
	}
	return nil, errs.Wrap(err, uri)
}

func buildMongoURI(mongoConf *config.Mongo) string {
	uri := os.Getenv("MONGO_URI")
	if uri != "" {
		return uri
	}

	if mongoConf.Uri != "" {
		return mongoConf.Uri
	}

	username := os.Getenv("MONGO_OPENIM_USERNAME")
	password := os.Getenv("MONGO_OPENIM_PASSWORD")
	address := os.Getenv("MONGO_ADDRESS")
	port := os.Getenv("MONGO_PORT")
	database := os.Getenv("MONGO_DATABASE")
	maxPoolSize := os.Getenv("MONGO_MAX_POOL_SIZE")

	if username == "" {
		username = mongoConf.Username
	}
	if password == "" {
		password = mongoConf.Password
	}
	if address == "" {
		address = strings.Join(mongoConf.Address, ",")
	} else if port != "" {
		address = fmt.Sprintf("%s:%s", address, port)
	}
	if database == "" {
		database = mongoConf.Database
	}
	if maxPoolSize == "" {
		maxPoolSize = fmt.Sprint(mongoConf.MaxPoolSize)
	}

	if username != "" && password != "" {

		return fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%s", username, password, address, database, maxPoolSize)
	}
	return fmt.Sprintf("mongodb://%s/%s?maxPoolSize=%s", address, database, maxPoolSize)
}

func shouldRetry(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok {
		return cmdErr.Code != 13 && cmdErr.Code != 18
	}
	return true
}

// GetClient returns the MongoDB client.
func (m *Mongo) GetClient() *mongo.Client {
	return m.db
}

// GetDatabase returns the specific database from MongoDB.
func (m *Mongo) GetDatabase(database string) *mongo.Database {
	return m.db.Database(database)
}

// CreateMsgIndex creates an index for messages in MongoDB.
func (m *Mongo) CreateMsgIndex() error {
	return m.createMongoIndex(unrelation.Msg, true, "doc_id")
}

// createMongoIndex creates an index in a MongoDB collection.
func (m *Mongo) createMongoIndex(collection string, isUnique bool, keys ...string) error {
	db := m.GetDatabase(m.mongoConf.Database).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	indexView := db.Indexes()

	keysDoc := buildIndexKeys(keys)

	index := mongo.IndexModel{
		Keys: keysDoc,
	}
	if isUnique {
		index.Options = options.Index().SetUnique(true)
	}

	_, err := indexView.CreateOne(context.Background(), index, opts)
	if err != nil {
		return errs.Wrap(err, "CreateIndex")
	}
	return nil
}

// buildIndexKeys builds the BSON document for index keys.
func buildIndexKeys(keys []string) bson.D {
	keysDoc := bson.D{}
	for _, key := range keys {
		direction := 1 // default direction is ascending
		if strings.HasPrefix(key, "-") {
			direction = -1 // descending order for prefixed with "-"
			key = strings.TrimLeft(key, "-")
		}
		keysDoc = append(keysDoc, bson.E{Key: key, Value: direction})
	}
	return keysDoc
}

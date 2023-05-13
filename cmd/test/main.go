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

package main

import (
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoMsg struct {
	UID string
	Msg []string
}


func main()  {
	//"mongodb://%s:%s@%s/%s/?maxPoolSize=%d"
	uri := "mongodb://user:pass@sample.host:27017/?maxPoolSize=20&w=majority"
	DBAddress := "127.0.0.1:37017"
	DBDatabase := "new-test-db"
	Collection := "new-test-collection"
	DBMaxPoolSize := 100
	uri = fmt.Sprintf("mongodb://%s/%s/?maxPoolSize=%d",
		DBAddress,DBDatabase,
		DBMaxPoolSize)

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	filter := bson.M{"uid":"my_uid"}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	for i:=0; i < 2; i++{

		if err = mongoClient.Database(DBDatabase).Collection(Collection).FindOneAndUpdate(ctx, filter,
			bson.M{"$push": bson.M{"msg": utils.Int32ToString(int32(i))}}).Err(); err != nil{
			fmt.Println("FindOneAndUpdate failed ", i, )
			var mmsg MongoMsg
			mmsg.UID = "my_uid"
			mmsg.Msg = append(mmsg.Msg, utils.Int32ToString(int32(i)))
			_, err := mongoClient.Database(DBDatabase).Collection(Collection).InsertOne(ctx, &mmsg)
			if err != nil {
				fmt.Println("insertone failed ", err.Error(), i)
			} else{
				fmt.Println("insertone ok ", i)
			}

		}else {
			fmt.Println("FindOneAndUpdate ok ", i)
		}

	}

	var mmsg MongoMsg

	if  err = mongoClient.Database(DBDatabase).Collection(Collection).FindOne(ctx, filter).Decode(&mmsg); err != nil {
		fmt.Println("findone failed ", err.Error())
	}else{
		fmt.Println("findone ok ", mmsg.UID)
		for i, v:=range mmsg.Msg{
			fmt.Println("find value: ", i, v)
		}
	}


}

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

package mgo

import (
	"context"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/db/mongoutil"
	"testing"
)

func GetTestDriver(t *testing.T, ctx context.Context) (database.Msg, error) {

	var mongoConf config.Mongo
	mongoConf.URI = ""
	mongoConf.Address = []string{"192.168.3.212:37017"}
	mongoConf.Database = "openIM_v3"
	mongoConf.Username = "root"
	mongoConf.Password = "openIM123"
	mongoConf.MaxPoolSize = 100
	mongoConf.MaxRetry = 10

	mgocli, err := mongoutil.NewMongoDB(ctx, mongoConf.Build())
	if err != nil {
		t.Fatal("NewMongoDB failed: ", err)
		return nil, err
	}
	mongoDB := mgocli.GetDB()

	err = mongoDB.Client().Ping(ctx, nil)
	if err != nil {
		t.Fatal("Ping failed: ", err)
		return nil, err
	}

	return NewMsgMongo(mongoDB)
}

func TestMsgMgo_GetMsgDocModelByIndex(t *testing.T) {
	ctx := context.Background()

	driver, _ := GetTestDriver(t, ctx)

	conversationID := "si_10043_10070"
	msg, err := driver.GetMsgDocModelByIndex(ctx, conversationID, 0, 1)
	if err != nil {
		t.Fatal("GetMsgDocModelByIndex failed: ", err)
	}

	fmt.Println("GetMsgDocModelByIndex Ascent(1): ", msg.DocID)

	msg2, err := driver.GetMsgDocModelByIndex(ctx, conversationID, 0, -1)
	if err != nil {
		t.Fatal("GetMsgDocModelByIndex failed: ", err)
	}

	fmt.Println("GetMsgDocModelByIndex Descent(-1): ", msg2.DocID)
}

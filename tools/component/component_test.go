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
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func TestCheckMysql(t *testing.T) {
	err := mockInitCfg()
	assert.NoError(t, err, "Initialization should not produce errors")

	err = checkMysql()
	if err != nil {
		// You might expect an error if MySQL isn't running locally with the mock credentials.
		t.Logf("Expected error due to mock configuration: %v", err)
	}
}

// Mock for initCfg for testing purpose
func mockInitCfg() error {
	config.Config.Mysql.Username = "root"
	config.Config.Mysql.Password = "openIM123"
	config.Config.Mysql.Address = []string{"127.0.0.1:13306"}
	return nil
}

func TestRedis(t *testing.T) {
	config.Config.Redis.Address = []string{
		"172.16.8.142:7000",
		//"172.16.8.142:7000", "172.16.8.142:7001", "172.16.8.142:7002", "172.16.8.142:7003", "172.16.8.142:7004", "172.16.8.142:7005",
	}

	var redisClient redis.UniversalClient
	defer func() {
		if redisClient != nil {
			redisClient.Close()
		}
	}()
	if len(config.Config.Redis.Address) > 1 {
		redisClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.Config.Redis.Address,
			Username: config.Config.Redis.Username,
			Password: config.Config.Redis.Password,
		})
	} else {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Config.Redis.Address[0],
			Username: config.Config.Redis.Username,
			Password: config.Config.Redis.Password,
		})
	}
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 1000000; i++ {
		val, err := redisClient.Set(context.Background(), "b_"+strconv.Itoa(i), "test", time.Second*10).Result()
		t.Log("index", i, "resp", val, "err", err)
		if err != nil {
			return
		}
	}

}

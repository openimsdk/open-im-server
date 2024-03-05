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

package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mw/specialerror"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/redis/go-redis/v9"
)

var (
	// singleton pattern.
	redisClient redis.UniversalClient
)

const (
	maxRetry = 10 // number of retries
)

// NewRedis Initialize redis connection.
func NewRedis() (redis.UniversalClient, error) {
	if redisClient != nil {
		return redisClient, nil
	}

	// Read configuration from environment variables
	overrideConfigFromEnv()

	if len(config.Config.Redis.Address) == 0 {
		return nil, errs.Wrap(errors.New("redis address is empty"), "Redis configuration error")
	}
	specialerror.AddReplace(redis.Nil, errs.ErrRecordNotFound)
	var rdb redis.UniversalClient
	if len(config.Config.Redis.Address) > 1 || config.Config.Redis.ClusterMode {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:      config.Config.Redis.Address,
			Username:   config.Config.Redis.Username,
			Password:   config.Config.Redis.Password, // no password set
			PoolSize:   50,
			MaxRetries: maxRetry,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:       config.Config.Redis.Address[0],
			Username:   config.Config.Redis.Username,
			Password:   config.Config.Redis.Password, // no password set
			DB:         0,                            // use default DB
			PoolSize:   100,                          // connection pool size
			MaxRetries: maxRetry,
		})
	}

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = rdb.Ping(ctx).Err()
	if err != nil {
		uriFormat := "address:%v, username:%s, clusterMode:%t, enablePipeline:%t"
		errMsg := fmt.Sprintf(uriFormat, config.Config.Redis.Address, config.Config.Redis.Username, config.Config.Redis.ClusterMode, config.Config.Redis.EnablePipeline)
		return nil, errs.Wrap(err, "Redis connection failed: %s", errMsg)
	}
	redisClient = rdb
	return rdb, err
}

// overrideConfigFromEnv overrides configuration fields with environment variables if present.
func overrideConfigFromEnv() {
	if envAddr := os.Getenv("REDIS_ADDRESS"); envAddr != "" {
		if envPort := os.Getenv("REDIS_PORT"); envPort != "" {
			addresses := strings.Split(envAddr, ",")
			for i, addr := range addresses {
				addresses[i] = addr + ":" + envPort
			}
			config.Config.Redis.Address = addresses
		} else {
			config.Config.Redis.Address = strings.Split(envAddr, ",")
		}
	}

	if envUser := os.Getenv("REDIS_USERNAME"); envUser != "" {
		config.Config.Redis.Username = envUser
	}

	if envPass := os.Getenv("REDIS_PASSWORD"); envPass != "" {
		config.Config.Redis.Password = envPass
	}
}

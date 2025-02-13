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

package third

import (
	"context"
	"fmt"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/mcache"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/tools/s3/disable"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/tools/s3/aws"
	"github.com/openimsdk/tools/s3/kodo"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/s3"
	"github.com/openimsdk/tools/s3/cos"
	"github.com/openimsdk/tools/s3/minio"
	"github.com/openimsdk/tools/s3/oss"
	"google.golang.org/grpc"
)

type thirdServer struct {
	third.UnimplementedThirdServer
	thirdDatabase controller.ThirdDatabase
	s3dataBase    controller.S3Database
	defaultExpire time.Duration
	config        *Config
	s3            s3.Interface
	userClient    *rpcli.UserClient
}

type Config struct {
	RpcConfig          config.Third
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	NotificationConfig config.Notification
	Share              config.Share
	MinioConfig        config.Minio
	LocalCacheConfig   config.LocalCache
	Discovery          config.Discovery
}

func Start(ctx context.Context, config *Config, client discovery.Conn, server grpc.ServiceRegistrar) error {
	dbb := dbbuild.NewBuilder(&config.MongodbConfig, &config.RedisConfig)
	mgocli, err := dbb.Mongo(ctx)
	if err != nil {
		return err
	}
	rdb, err := dbb.Redis(ctx)
	if err != nil {
		return err
	}

	logdb, err := mgo.NewLogMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	s3db, err := mgo.NewS3Mongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	var thirdCache cache.ThirdCache
	if rdb == nil {
		tc, err := mgo.NewCacheMgo(mgocli.GetDB())
		if err != nil {
			return err
		}
		thirdCache = mcache.NewThirdCache(tc)
	} else {
		thirdCache = redis.NewThirdCache(rdb)
	}
	// Select the oss method according to the profile policy
	var o s3.Interface
	switch enable := config.RpcConfig.Object.Enable; enable {
	case "minio":
		var minioCache minio.Cache
		if rdb == nil {
			mc, err := mgo.NewCacheMgo(mgocli.GetDB())
			if err != nil {
				return err
			}
			minioCache = mcache.NewMinioCache(mc)
		} else {
			minioCache = redis.NewMinioCache(rdb)
		}
		o, err = minio.NewMinio(ctx, minioCache, *config.MinioConfig.Build())
	case "cos":
		o, err = cos.NewCos(*config.RpcConfig.Object.Cos.Build())
	case "oss":
		o, err = oss.NewOSS(*config.RpcConfig.Object.Oss.Build())
	case "kodo":
		o, err = kodo.NewKodo(*config.RpcConfig.Object.Kodo.Build())
	case "aws":
		o, err = aws.NewAws(*config.RpcConfig.Object.Aws.Build())
	case "":
		o = disable.NewDisable()
	default:
		err = fmt.Errorf("invalid object enable: %s", enable)
	}
	if err != nil {
		return err
	}
	userConn, err := client.GetConn(ctx, config.Discovery.RpcService.User)
	if err != nil {
		return err
	}
	localcache.InitLocalCache(&config.LocalCacheConfig)
	third.RegisterThirdServer(server, &thirdServer{
		thirdDatabase: controller.NewThirdDatabase(thirdCache, logdb),
		s3dataBase:    controller.NewS3Database(rdb, o, s3db),
		defaultExpire: time.Hour * 24 * 7,
		config:        config,
		s3:            o,
		userClient:    rpcli.NewUserClient(userConn),
	})
	return nil
}

func (t *thirdServer) FcmUpdateToken(ctx context.Context, req *third.FcmUpdateTokenReq) (resp *third.FcmUpdateTokenResp, err error) {
	err = t.thirdDatabase.FcmUpdateToken(ctx, req.Account, int(req.PlatformID), req.FcmToken, req.ExpireTime)
	if err != nil {
		return nil, err
	}
	return &third.FcmUpdateTokenResp{}, nil
}

func (t *thirdServer) SetAppBadge(ctx context.Context, req *third.SetAppBadgeReq) (resp *third.SetAppBadgeResp, err error) {
	err = t.thirdDatabase.SetAppBadge(ctx, req.UserID, int(req.AppUnreadCount))
	if err != nil {
		return nil, err
	}
	return &third.SetAppBadgeResp{}, nil
}

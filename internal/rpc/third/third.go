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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/tools/s3/aws"
	"github.com/openimsdk/tools/s3/kodo"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
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
	userRpcClient rpcclient.UserRpcClient
	defaultExpire time.Duration
	config        *Config
	minio         *minio.Minio
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

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
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

	// Select the oss method according to the profile policy
	enable := config.RpcConfig.Object.Enable
	var (
		o        s3.Interface
		minioCli *minio.Minio
	)
	switch enable {
	case "minio":
		minioCli, err = minio.NewMinio(ctx, redis.NewMinioCache(rdb), *config.MinioConfig.Build())
		o = minioCli
	case "cos":
		o, err = cos.NewCos(*config.RpcConfig.Object.Cos.Build())
	case "oss":
		o, err = oss.NewOSS(*config.RpcConfig.Object.Oss.Build())
	case "kodo":
		o, err = kodo.NewKodo(*config.RpcConfig.Object.Kodo.Build())
	case "aws":
		o, err = aws.NewAws(*config.RpcConfig.Object.Aws.Build())
	default:
		err = fmt.Errorf("invalid object enable: %s", enable)
	}
	if err != nil {
		return err
	}
	localcache.InitLocalCache(&config.LocalCacheConfig)
	third.RegisterThirdServer(server, &thirdServer{
		thirdDatabase: controller.NewThirdDatabase(redis.NewThirdCache(rdb), logdb),
		userRpcClient: rpcclient.NewUserRpcClient(client, config.Share.RpcRegisterName.User, config.Share.IMAdminUserID),
		s3dataBase:    controller.NewS3Database(rdb, o, s3db),
		defaultExpire: time.Hour * 24 * 7,
		config:        config,
		minio:         minioCli,
	})
	return nil
}

func (t *thirdServer) getMinioImageThumbnailKey(ctx context.Context, name string) (string, error) {
	return t.minio.GetImageThumbnailKey(ctx, name)
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

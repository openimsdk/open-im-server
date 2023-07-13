package third

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3/cos"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3/minio"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3/oss"
	"net/url"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"google.golang.org/grpc"
)

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	apiURL := config.Config.Object.ApiURL
	if apiURL == "" {
		return fmt.Errorf("api url is empty")
	}
	if _, err := url.Parse(config.Config.Object.ApiURL); err != nil {
		return err
	}
	if apiURL[len(apiURL)-1] != '/' {
		apiURL += "/"
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&relationTb.ObjectModel{}); err != nil {
		return err
	}
	// 根据配置文件策略选择 oss 方式
	enable := config.Config.Object.Enable
	var o s3.Interface
	switch config.Config.Object.Enable {
	case "minio":
		o, err = minio.NewMinio()
	case "cos":
		o, err = cos.NewCos()
	case "oss":
		o, err = oss.NewOSS()
	default:
		err = fmt.Errorf("invalid object enable: %s", enable)
	}
	if err != nil {
		return err
	}
	third.RegisterThirdServer(server, &thirdServer{
		apiURL:        apiURL,
		thirdDatabase: controller.NewThirdDatabase(cache.NewMsgCacheModel(rdb)),
		userRpcClient: rpcclient.NewUserRpcClient(client),
		s3dataBase:    controller.NewS3Database(o, relation.NewObjectInfo(db)),
		defaultExpire: time.Hour * 24 * 7,
	})
	return nil
}

type thirdServer struct {
	apiURL        string
	thirdDatabase controller.ThirdDatabase
	s3dataBase    controller.S3Database
	userRpcClient rpcclient.UserRpcClient
	defaultExpire time.Duration
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

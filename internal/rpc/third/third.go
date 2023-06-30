package third

import (
	"context"
	"net/url"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/obj"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"google.golang.org/grpc"
)

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	u, err := url.Parse(config.Config.Object.ApiURL)
	if err != nil {
		return err
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	o, err := obj.NewMinioInterface()
	if err != nil {
		return err
	}
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&relationTb.ObjectHashModel{}, &relationTb.ObjectInfoModel{}, &relationTb.ObjectPutModel{}); err != nil {
		return err
	}
	third.RegisterThirdServer(server, &thirdServer{
		thirdDatabase: controller.NewThirdDatabase(cache.NewMsgCacheModel(rdb)),
		userRpcClient: rpcclient.NewUserRpcClient(client),
		s3dataBase:    controller.NewS3Database(o, relation.NewObjectHash(db), relation.NewObjectInfo(db), relation.NewObjectPut(db), u),
	})
	return nil
}

type thirdServer struct {
	thirdDatabase controller.ThirdDatabase
	s3dataBase    controller.S3Database
	userRpcClient rpcclient.UserRpcClient
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

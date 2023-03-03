package third

import (
	"OpenIM/internal/common/check"
	"OpenIM/pkg/common/db/cache"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/db/obj"
	"OpenIM/pkg/common/db/relation"
	relationTb "OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/proto/third"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"google.golang.org/grpc"
)

func Start(client *openKeeper.ZkClient, server *grpc.Server) error {
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
		thirdDatabase: controller.NewThirdDatabase(cache.NewCache(rdb)),
		userCheck:     check.NewUserCheck(client),
		s3dataBase:    controller.NewS3Database(o, relation.NewObjectHash(db), relation.NewObjectInfo(db), relation.NewObjectPut(db)),
	})
	return nil
}

type thirdServer struct {
	thirdDatabase controller.ThirdDatabase
	s3dataBase    controller.S3Database
	userCheck     *check.UserCheck
}

func (t *thirdServer) GetSignalInvitationInfo(ctx context.Context, req *third.GetSignalInvitationInfoReq) (resp *third.GetSignalInvitationInfoResp, err error) {
	signalReq, err := t.thirdDatabase.GetSignalInvitationInfoByClientMsgID(ctx, req.ClientMsgID)
	if err != nil {
		return nil, err
	}
	resp = &third.GetSignalInvitationInfoResp{}
	resp.InvitationInfo = signalReq.Invitation
	resp.OfflinePushInfo = signalReq.OfflinePushInfo
	return resp, nil
}

func (t *thirdServer) GetSignalInvitationInfoStartApp(ctx context.Context, req *third.GetSignalInvitationInfoStartAppReq) (resp *third.GetSignalInvitationInfoStartAppResp, err error) {
	signalReq, err := t.thirdDatabase.GetAvailableSignalInvitationInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	resp = &third.GetSignalInvitationInfoStartAppResp{}
	resp.InvitationInfo = signalReq.Invitation
	resp.OfflinePushInfo = signalReq.OfflinePushInfo
	return resp, nil
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

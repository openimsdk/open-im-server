package push

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/cache"
	pbPush "OpenIM/pkg/proto/push"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"google.golang.org/grpc"
)

type pushServer struct {
	pusher *Pusher
}

func Start(client *openKeeper.ZkClient, server *grpc.Server) error {

	pbPush.RegisterPushMsgServiceServer(server, &pushServer{
		pusher: NewPusher(),
	})
}

func (r *pushServer) PushMsg(ctx context.Context, pbData *pbPush.PushMsgReq) (resp *pbPush.PushMsgResp, err error) {
	switch pbData.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = r.pusher.MsgToSuperGroupUser(ctx, pbData.SourceID, pbData.MsgData)
	default:
		err = r.pusher.MsgToUser(ctx, pbData.SourceID, pbData.MsgData)
	}
	if err != nil {
		return nil, err
	}
	return &pbPush.PushMsgResp{}, nil
}

func (r *pushServer) DelUserPushToken(ctx context.Context, req *pbPush.DelUserPushTokenReq) (resp *pbPush.DelUserPushTokenResp, err error) {
	if err = r.pusher.database.DelFcmToken(ctx, req.UserID, int(req.PlatformID)); err != nil {
		return nil, err
	}
	return &pbPush.DelUserPushTokenResp{}, nil
}

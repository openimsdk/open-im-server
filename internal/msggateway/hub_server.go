package msggateway

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msggateway"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/startrpc"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

func (s *Server) InitServer(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	msgModel := cache.NewMsgCacheModel(rdb)
	s.LongConnServer.SetDiscoveryRegistry(client)
	s.LongConnServer.SetCacheHandler(msgModel)
	msggateway.RegisterMsgGatewayServer(server, s)
	return nil
}

func (s *Server) Start() error {
	return startrpc.Start(s.rpcPort, config.Config.RpcRegisterName.OpenImMessageGatewayName, s.prometheusPort, s.InitServer)
}

type Server struct {
	rpcPort        int
	prometheusPort int
	LongConnServer LongConnServer
	pushTerminal   []int
}

func (s *Server) SetLongConnServer(LongConnServer LongConnServer) {
	s.LongConnServer = LongConnServer
}

func NewServer(rpcPort int, longConnServer LongConnServer) *Server {
	return &Server{rpcPort: rpcPort, LongConnServer: longConnServer, pushTerminal: []int{constant.IOSPlatformID, constant.AndroidPlatformID}}
}

func (s *Server) OnlinePushMsg(context context.Context, req *msggateway.OnlinePushMsgReq) (*msggateway.OnlinePushMsgResp, error) {
	panic("implement me")
}

func (s *Server) GetUsersOnlineStatus(ctx context.Context, req *msggateway.GetUsersOnlineStatusReq) (*msggateway.GetUsersOnlineStatusResp, error) {
	if !tokenverify.IsAppManagerUid(ctx) {
		return nil, errs.ErrNoPermission.Wrap("only app manager")
	}
	var resp msggateway.GetUsersOnlineStatusResp
	for _, userID := range req.UserIDs {
		clients, ok := s.LongConnServer.GetUserAllCons(userID)
		if !ok {
			continue
		}
		temp := new(msggateway.GetUsersOnlineStatusResp_SuccessResult)
		temp.UserID = userID
		for _, client := range clients {
			if client != nil {
				ps := new(msggateway.GetUsersOnlineStatusResp_SuccessDetail)
				ps.Platform = constant.PlatformIDToName(client.PlatformID)
				ps.Status = constant.OnlineStatus
				ps.ConnID = client.ctx.GetConnID()
				ps.IsBackground = client.IsBackground
				temp.Status = constant.OnlineStatus
				temp.DetailPlatformStatus = append(temp.DetailPlatformStatus, ps)
			}
		}
		if temp.Status == constant.OnlineStatus {
			resp.SuccessResult = append(resp.SuccessResult, temp)
		}
	}
	return &resp, nil
}

func (s *Server) OnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	panic("implement me")
}

func (s *Server) SuperGroupOnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	var singleUserResult []*msggateway.SingleMsgToUserResults
	for _, v := range req.PushToUserIDs {
		var resp []*msggateway.SingleMsgToUserPlatform
		tempT := &msggateway.SingleMsgToUserResults{
			UserID: v,
		}
		clients, ok := s.LongConnServer.GetUserAllCons(v)
		if !ok {
			log.ZDebug(ctx, "push user not online", "userID", v)
			tempT.Resp = resp
			singleUserResult = append(singleUserResult, tempT)
			continue
		}
		log.ZDebug(ctx, "push user online", "clients", clients, "userID", v)
		for _, client := range clients {
			if client != nil {
				temp := &msggateway.SingleMsgToUserPlatform{
					RecvID:         v,
					RecvPlatFormID: int32(client.PlatformID),
				}
				if !client.IsBackground || (client.IsBackground == true && client.PlatformID != constant.IOSPlatformID) {
					err := client.PushMessage(ctx, req.MsgData)
					if err != nil {
						temp.ResultCode = -2
						resp = append(resp, temp)
					} else {
						if utils.IsContainInt(client.PlatformID, s.pushTerminal) {
							tempT.OnlinePush = true
							prome.Inc(prome.MsgOnlinePushSuccessCounter)
							resp = append(resp, temp)
						}
					}
				} else {
					temp.ResultCode = -3
					resp = append(resp, temp)
				}
			}
		}
		tempT.Resp = resp
		singleUserResult = append(singleUserResult, tempT)
	}

	return &msggateway.OnlineBatchPushOneMsgResp{
		SinglePushResult: singleUserResult,
	}, nil
}

func (s *Server) KickUserOffline(ctx context.Context, req *msggateway.KickUserOfflineReq) (*msggateway.KickUserOfflineResp, error) {
	for _, v := range req.KickUserIDList {
		if clients, _, ok := s.LongConnServer.GetUserPlatformCons(v, int(req.PlatformID)); ok {
			for _, client := range clients {
				err := client.KickOnlineMessage()
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return &msggateway.KickUserOfflineResp{}, nil
}

func (s *Server) MultiTerminalLoginCheck(ctx context.Context, req *msggateway.MultiTerminalLoginCheckReq) (*msggateway.MultiTerminalLoginCheckResp, error) {
	//TODO implement me
	panic("implement me")
}

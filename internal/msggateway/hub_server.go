package msggateway

import (
	"OpenIM/internal/common/network"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/mw"
	"OpenIM/pkg/common/prome"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/errs"
	"OpenIM/pkg/proto/msggateway"
	"OpenIM/pkg/utils"
	"context"
	"fmt"
	"github.com/OpenIMSDK/openKeeper"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"net"
)

func (s *Server) Start() error {
	zkClient, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema, 10, "", "")
	if err != nil {
		return err
	}
	defer zkClient.Close()
	registerIP, err := network.GetRpcRegisterIP(config.Config.RpcRegisterIP)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.ListenIP, s.rpcPort))
	if err != nil {
		panic("listening err:" + err.Error())
	}
	defer listener.Close()
	var options []grpc.ServerOption
	options = append(options, mw.GrpcServer()) // ctx 中间件
	if config.Config.Prometheus.Enable {
		prome.NewGrpcRequestCounter()
		prome.NewGrpcRequestFailedCounter()
		prome.NewGrpcRequestSuccessCounter()
		options = append(options, []grpc.ServerOption{
			//grpc.UnaryInterceptor(prome.UnaryServerInterceptorPrometheus),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(options...)
	defer srv.GracefulStop()
	msggateway.RegisterMsgGatewayServer(srv, s)
	err = zkClient.Register("", registerIP, s.rpcPort)
	if err != nil {
		return err
	}
	err = srv.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

type Server struct {
	rpcPort        int
	LongConnServer LongConnServer
	pushTerminal   []int
	//rpcServer      *RpcServer
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
				ps.Platform = constant.PlatformIDToName(client.platformID)
				ps.Status = constant.OnlineStatus
				ps.ConnID = client.connID
				ps.IsBackground = client.isBackground
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
			continue
		}
		for _, client := range clients {
			if client != nil {
				temp := &msggateway.SingleMsgToUserPlatform{
					RecvID:         v,
					RecvPlatFormID: int32(client.platformID),
				}
				if !client.isBackground {
					err := client.PushMessage(ctx, req.MsgData)
					if err != nil {
						temp.ResultCode = -2
						resp = append(resp, temp)
					} else {
						if utils.IsContainInt(client.platformID, s.pushTerminal) {
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
				err := client.KickOnlineMessage(ctx)
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

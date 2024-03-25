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

package msggateway

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msggateway"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"google.golang.org/grpc"
)

func (s *Server) InitServer(ctx context.Context, config *config.GlobalConfig, disCov discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis(ctx, &config.Redis)
	if err != nil {
		return err
	}

	tokenCacheModel := cache.NewTokenCacheModel(rdb)
	s.LongConnServer.SetDiscoveryRegistry(disCov, config)
	s.LongConnServer.SetCacheHandler(tokenCacheModel)
	msggateway.RegisterMsgGatewayServer(server, s)
	return nil
}

func (s *Server) Start(ctx context.Context, conf *config.GlobalConfig) error {
	return startrpc.Start(ctx,
		s.rpcPort,
		conf.RpcRegisterName.OpenImMessageGatewayName,
		s.prometheusPort,
		conf,
		s.InitServer,
	)
}

type Server struct {
	rpcPort        int
	prometheusPort int
	LongConnServer LongConnServer
	config         *config.GlobalConfig
	pushTerminal   map[int]struct{}
}

func (s *Server) SetLongConnServer(LongConnServer LongConnServer) {
	s.LongConnServer = LongConnServer
}

func NewServer(rpcPort int, proPort int, longConnServer LongConnServer, conf *config.GlobalConfig) *Server {
	s := &Server{
		rpcPort:        rpcPort,
		prometheusPort: proPort,
		LongConnServer: longConnServer,
		pushTerminal:   make(map[int]struct{}),
		config:         conf,
	}
	s.pushTerminal[constant.IOSPlatformID] = struct{}{}
	s.pushTerminal[constant.AndroidPlatformID] = struct{}{}
	return s
}

func (s *Server) OnlinePushMsg(
	context context.Context,
	req *msggateway.OnlinePushMsgReq,
) (*msggateway.OnlinePushMsgResp, error) {
	panic("implement me")
}

func (s *Server) GetUsersOnlineStatus(
	ctx context.Context,
	req *msggateway.GetUsersOnlineStatusReq,
) (*msggateway.GetUsersOnlineStatusResp, error) {
	if !authverify.IsAppManagerUid(ctx, &s.config.Manager, &s.config.IMAdmin) {
		return nil, errs.ErrNoPermission.WrapMsg("only app manager")
	}
	var resp msggateway.GetUsersOnlineStatusResp
	for _, userID := range req.UserIDs {
		clients, ok := s.LongConnServer.GetUserAllCons(userID)
		if !ok {
			continue
		}

		uresp := new(msggateway.GetUsersOnlineStatusResp_SuccessResult)
		uresp.UserID = userID
		for _, client := range clients {
			if client == nil {
				continue
			}

			ps := new(msggateway.GetUsersOnlineStatusResp_SuccessDetail)
			ps.Platform = constant.PlatformIDToName(client.PlatformID)
			ps.Status = constant.OnlineStatus
			ps.ConnID = client.ctx.GetConnID()
			ps.Token = client.token
			ps.IsBackground = client.IsBackground
			uresp.Status = constant.OnlineStatus
			uresp.DetailPlatformStatus = append(uresp.DetailPlatformStatus, ps)
		}
		if uresp.Status == constant.OnlineStatus {
			resp.SuccessResult = append(resp.SuccessResult, uresp)
		}
	}
	return &resp, nil
}

func (s *Server) OnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	//todo implement
	return nil, nil
}

func (s *Server) SuperGroupOnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq,
) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	var singleUserResults []*msggateway.SingleMsgToUserResults
	for _, v := range req.PushToUserIDs {
		var resp []*msggateway.SingleMsgToUserPlatform
		results := &msggateway.SingleMsgToUserResults{
			UserID: v,
		}
		clients, ok := s.LongConnServer.GetUserAllCons(v)
		if !ok {
			log.ZDebug(ctx, "push user not online", "userID", v)
			results.Resp = resp
			singleUserResults = append(singleUserResults, results)
			continue
		}

		log.ZDebug(ctx, "push user online", "clients", clients, "userID", v)
		for _, client := range clients {
			if client == nil {
				continue
			}

			userPlatform := &msggateway.SingleMsgToUserPlatform{
				RecvPlatFormID: int32(client.PlatformID),
			}
			if !client.IsBackground ||
				(client.IsBackground && client.PlatformID != constant.IOSPlatformID) {
				err := client.PushMessage(ctx, req.MsgData)
				if err != nil {
					userPlatform.ResultCode = int64(errs.ErrPushMsgErr.Code())
					resp = append(resp, userPlatform)
				} else {
					if _, ok := s.pushTerminal[client.PlatformID]; ok {
						results.OnlinePush = true
						resp = append(resp, userPlatform)
					}
				}
			} else {
				userPlatform.ResultCode = int64(errs.ErrIOSBackgroundPushErr.Code())
				resp = append(resp, userPlatform)
			}
		}
		results.Resp = resp
		singleUserResults = append(singleUserResults, results)
	}

	return &msggateway.OnlineBatchPushOneMsgResp{
		SinglePushResult: singleUserResults,
	}, nil
}

func (s *Server) KickUserOffline(
	ctx context.Context,
	req *msggateway.KickUserOfflineReq,
) (*msggateway.KickUserOfflineResp, error) {
	for _, v := range req.KickUserIDList {
		clients, _, ok := s.LongConnServer.GetUserPlatformCons(v, int(req.PlatformID))
		if !ok {
			log.ZInfo(ctx, "conn not exist", "userID", v, "platformID", req.PlatformID)
			continue
		}

		for _, client := range clients {
			log.ZDebug(ctx, "kick user offline", "userID", v, "platformID", req.PlatformID, "client", client)
			if err := client.longConnServer.KickUserConn(client); err != nil {
				log.ZWarn(ctx, "kick user offline failed", err, "userID", v, "platformID", req.PlatformID)
			}
		}
		continue
	}

	return &msggateway.KickUserOfflineResp{}, nil
}

func (s *Server) MultiTerminalLoginCheck(ctx context.Context, req *msggateway.MultiTerminalLoginCheckReq) (*msggateway.MultiTerminalLoginCheckResp, error) {
	if oldClients, userOK, clientOK := s.LongConnServer.GetUserPlatformCons(req.UserID, int(req.PlatformID)); userOK {
		tempUserCtx := newTempContext()
		tempUserCtx.SetToken(req.Token)
		tempUserCtx.SetOperationID(mcontext.GetOperationID(ctx))
		client := &Client{}
		client.ctx = tempUserCtx
		client.UserID = req.UserID
		client.PlatformID = int(req.PlatformID)
		i := &kickHandler{
			clientOK:   clientOK,
			oldClients: oldClients,
			newClient:  client,
		}
		s.LongConnServer.SetKickHandlerInfo(i)
	}
	return &msggateway.MultiTerminalLoginCheckResp{}, nil
}

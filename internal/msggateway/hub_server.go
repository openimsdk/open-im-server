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

	"github.com/OpenIMSDK/tools/mcontext"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"

	"github.com/OpenIMSDK/tools/errs"
	"google.golang.org/grpc"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msggateway"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
)

func (s *Server) InitServer(disCov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	msgModel := cache.NewMsgCacheModel(rdb)
	s.LongConnServer.SetDiscoveryRegistry(disCov)
	s.LongConnServer.SetCacheHandler(msgModel)
	msggateway.RegisterMsgGatewayServer(server, s)
	return nil
}

func (s *Server) Start() error {
	return startrpc.Start(
		s.rpcPort,
		config.Config.RpcRegisterName.OpenImMessageGatewayName,
		s.prometheusPort,
		s.InitServer,
	)
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

func NewServer(rpcPort int, proPort int, longConnServer LongConnServer) *Server {
	return &Server{
		rpcPort:        rpcPort,
		prometheusPort: proPort,
		LongConnServer: longConnServer,
		pushTerminal:   []int{constant.IOSPlatformID, constant.AndroidPlatformID},
	}
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
	if !authverify.IsAppManagerUid(ctx) {
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
				ps.Token = client.token
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

func (s *Server) OnlineBatchPushOneMsg(
	ctx context.Context,
	req *msggateway.OnlineBatchPushOneMsgReq,
) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	panic("implement me")
}

func (s *Server) SuperGroupOnlineBatchPushOneMsg(
	ctx context.Context,
	req *msggateway.OnlineBatchPushOneMsgReq,
) (*msggateway.OnlineBatchPushOneMsgResp, error) {
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
				if !client.IsBackground ||
					(client.IsBackground == true && client.PlatformID != constant.IOSPlatformID) {
					err := client.PushMessage(ctx, req.MsgData)
					if err != nil {
						temp.ResultCode = -2
						resp = append(resp, temp)
					} else {
						if utils.IsContainInt(client.PlatformID, s.pushTerminal) {
							tempT.OnlinePush = true
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

func (s *Server) KickUserOffline(
	ctx context.Context,
	req *msggateway.KickUserOfflineReq,
) (*msggateway.KickUserOfflineResp, error) {
	for _, v := range req.KickUserIDList {
		if clients, _, ok := s.LongConnServer.GetUserPlatformCons(v, int(req.PlatformID)); ok {
			for _, client := range clients {
				log.ZDebug(ctx, "kick user offline", "userID", v, "platformID", req.PlatformID, "client", client)
				if err := client.longConnServer.KickUserConn(client); err != nil {
					log.ZWarn(ctx, "kick user offline failed", err, "userID", v, "platformID", req.PlatformID)
				}
			}
		} else {
			log.ZInfo(ctx, "conn not exist", "userID", v, "platformID", req.PlatformID)
		}
	}
	return &msggateway.KickUserOfflineResp{}, nil
}

func (s *Server) MultiTerminalLoginCheck(
	ctx context.Context,
	req *msggateway.MultiTerminalLoginCheckReq,
) (*msggateway.MultiTerminalLoginCheckResp, error) {
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

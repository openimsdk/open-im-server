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
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"sync/atomic"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/startrpc"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msggateway"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mq/memamq"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
)

func (s *Server) InitServer(ctx context.Context, config *Config, disCov discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	userConn, err := disCov.GetConn(ctx, config.Share.RpcRegisterName.User)
	if err != nil {
		return err
	}
	s.userClient = rpcli.NewUserClient(userConn)
	if err := s.LongConnServer.SetDiscoveryRegistry(ctx, disCov, config); err != nil {
		return err
	}
	msggateway.RegisterMsgGatewayServer(server, s)
	if s.ready != nil {
		return s.ready(s)
	}
	return nil
}

func (s *Server) Start(ctx context.Context, index int, conf *Config) error {
	return startrpc.Start(ctx, &conf.Discovery, &conf.MsgGateway.Prometheus, conf.MsgGateway.ListenIP,
		conf.MsgGateway.RPC.RegisterIP,
		conf.MsgGateway.RPC.Ports, index,
		conf.Share.RpcRegisterName.MessageGateway,
		&conf.Share,
		conf,
		s.InitServer,
	)
}

type Server struct {
	msggateway.UnimplementedMsgGatewayServer
	rpcPort        int
	LongConnServer LongConnServer
	config         *Config
	pushTerminal   map[int]struct{}
	ready          func(srv *Server) error
	queue          *memamq.MemoryQueue
	userClient     *rpcli.UserClient
}

func (s *Server) SetLongConnServer(LongConnServer LongConnServer) {
	s.LongConnServer = LongConnServer
}

func NewServer(longConnServer LongConnServer, conf *Config, ready func(srv *Server) error) *Server {
	s := &Server{
		LongConnServer: longConnServer,
		pushTerminal:   make(map[int]struct{}),
		config:         conf,
		ready:          ready,
		queue:          memamq.NewMemoryQueue(512, 1024*16),
	}
	s.pushTerminal[constant.IOSPlatformID] = struct{}{}
	s.pushTerminal[constant.AndroidPlatformID] = struct{}{}
	return s
}

func (s *Server) OnlinePushMsg(context context.Context, req *msggateway.OnlinePushMsgReq) (*msggateway.OnlinePushMsgResp, error) {
	panic("implement me")
}

func (s *Server) GetUsersOnlineStatus(ctx context.Context, req *msggateway.GetUsersOnlineStatusReq) (*msggateway.GetUsersOnlineStatusResp, error) {
	if !authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
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
			ps.PlatformID = int32(client.PlatformID)
			ps.ConnID = client.ctx.GetConnID()
			ps.Token = client.token
			ps.IsBackground = client.IsBackground
			uresp.Status = constant.Online
			uresp.DetailPlatformStatus = append(uresp.DetailPlatformStatus, ps)
		}
		if uresp.Status == constant.Online {
			resp.SuccessResult = append(resp.SuccessResult, uresp)
		}
	}
	return &resp, nil
}

func (s *Server) OnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	// todo implement
	return nil, nil
}

func (s *Server) pushToUser(ctx context.Context, userID string, msgData *sdkws.MsgData) *msggateway.SingleMsgToUserResults {
	clients, ok := s.LongConnServer.GetUserAllCons(userID)
	if !ok {
		log.ZDebug(ctx, "push user not online", "userID", userID)
		return &msggateway.SingleMsgToUserResults{
			UserID: userID,
		}
	}
	log.ZDebug(ctx, "push user online", "clients", clients, "userID", userID)
	result := &msggateway.SingleMsgToUserResults{
		UserID: userID,
		Resp:   make([]*msggateway.SingleMsgToUserPlatform, 0, len(clients)),
	}
	for _, client := range clients {
		if client == nil {
			continue
		}
		userPlatform := &msggateway.SingleMsgToUserPlatform{
			RecvPlatFormID: int32(client.PlatformID),
		}
		if !client.IsBackground ||
			(client.IsBackground && client.PlatformID != constant.IOSPlatformID) {
			err := client.PushMessage(ctx, msgData)
			if err != nil {
				log.ZWarn(ctx, "online push msg failed", err, "userID", userID, "platformID", client.PlatformID)
				userPlatform.ResultCode = int64(servererrs.ErrPushMsgErr.Code())
			} else {
				if _, ok := s.pushTerminal[client.PlatformID]; ok {
					result.OnlinePush = true
				}
			}
		} else {
			userPlatform.ResultCode = int64(servererrs.ErrIOSBackgroundPushErr.Code())
		}
		result.Resp = append(result.Resp, userPlatform)
	}
	return result
}

func (s *Server) SuperGroupOnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	if len(req.PushToUserIDs) == 0 {
		return &msggateway.OnlineBatchPushOneMsgResp{}, nil
	}
	ch := make(chan *msggateway.SingleMsgToUserResults, len(req.PushToUserIDs))
	var count atomic.Int64
	count.Add(int64(len(req.PushToUserIDs)))
	for i := range req.PushToUserIDs {
		userID := req.PushToUserIDs[i]
		err := s.queue.PushCtx(ctx, func() {
			ch <- s.pushToUser(ctx, userID, req.MsgData)
			if count.Add(-1) == 0 {
				close(ch)
			}
		})
		if err != nil {
			if count.Add(-1) == 0 {
				close(ch)
			}
			log.ZError(ctx, "pushToUser MemoryQueue failed", err, "userID", userID)
			ch <- &msggateway.SingleMsgToUserResults{
				UserID: userID,
			}
		}
	}
	resp := &msggateway.OnlineBatchPushOneMsgResp{
		SinglePushResult: make([]*msggateway.SingleMsgToUserResults, 0, len(req.PushToUserIDs)),
	}
	for {
		select {
		case <-ctx.Done():
			log.ZError(ctx, "SuperGroupOnlineBatchPushOneMsg ctx done", context.Cause(ctx))
			userIDSet := datautil.SliceSet(req.PushToUserIDs)
			for _, results := range resp.SinglePushResult {
				delete(userIDSet, results.UserID)
			}
			for userID := range userIDSet {
				resp.SinglePushResult = append(resp.SinglePushResult, &msggateway.SingleMsgToUserResults{
					UserID: userID,
				})
			}
			return resp, nil
		case res, ok := <-ch:
			if !ok {
				return resp, nil
			}
			resp.SinglePushResult = append(resp.SinglePushResult, res)
		}
	}
}

func (s *Server) KickUserOffline(ctx context.Context, req *msggateway.KickUserOfflineReq) (*msggateway.KickUserOfflineResp, error) {
	for _, v := range req.KickUserIDList {
		clients, _, ok := s.LongConnServer.GetUserPlatformCons(v, int(req.PlatformID))
		if !ok {
			log.ZDebug(ctx, "conn not exist", "userID", v, "platformID", req.PlatformID)
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

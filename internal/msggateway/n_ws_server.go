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
	"errors"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenIMSDK/protocol/msggateway"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"

	"github.com/OpenIMSDK/protocol/constant"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"

	"github.com/redis/go-redis/v9"

	"github.com/OpenIMSDK/tools/discoveryregistry"

	"github.com/go-playground/validator/v10"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
)

type LongConnServer interface {
	Run() error
	wsHandler(w http.ResponseWriter, r *http.Request)
	GetUserAllCons(userID string) ([]*Client, bool)
	GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool)
	Validate(s interface{}) error
	SetCacheHandler(cache cache.MsgModel)
	SetDiscoveryRegistry(client discoveryregistry.SvcDiscoveryRegistry)
	KickUserConn(client *Client) error
	UnRegister(c *Client)
	SetKickHandlerInfo(i *kickHandler)
	Compressor
	Encoder
	MessageHandler
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

type WsServer struct {
	port              int
	wsMaxConnNum      int64
	registerChan      chan *Client
	unregisterChan    chan *Client
	kickHandlerChan   chan *kickHandler
	clients           *UserMap
	clientPool        sync.Pool
	onlineUserNum     atomic.Int64
	onlineUserConnNum atomic.Int64
	handshakeTimeout  time.Duration
	hubServer         *Server
	validate          *validator.Validate
	cache             cache.MsgModel
	userClient        *rpcclient.UserRpcClient
	disCov            discoveryregistry.SvcDiscoveryRegistry
	Compressor
	Encoder
	MessageHandler
}
type kickHandler struct {
	clientOK   bool
	oldClients []*Client
	newClient  *Client
}

func (ws *WsServer) SetDiscoveryRegistry(disCov discoveryregistry.SvcDiscoveryRegistry) {
	ws.MessageHandler = NewGrpcHandler(ws.validate, disCov)
	u := rpcclient.NewUserRpcClient(disCov)
	ws.userClient = &u
	ws.disCov = disCov
}

func (ws *WsServer) SetUserOnlineStatus(ctx context.Context, client *Client, status int32) {
	err := ws.userClient.SetUserStatus(ctx, client.UserID, status, client.PlatformID)
	if err != nil {
		log.ZWarn(ctx, "SetUserStatus err", err)
	}
	switch status {
	case constant.Online:
		err := CallbackUserOnline(ctx, client.UserID, client.PlatformID, client.IsBackground, client.ctx.GetConnID())
		if err != nil {
			log.ZWarn(ctx, "CallbackUserOnline err", err)
		}
	case constant.Offline:
		err := CallbackUserOffline(ctx, client.UserID, client.PlatformID, client.ctx.GetConnID())
		if err != nil {
			log.ZWarn(ctx, "CallbackUserOffline err", err)
		}
	}
}

func (ws *WsServer) SetCacheHandler(cache cache.MsgModel) {
	ws.cache = cache
}

func (ws *WsServer) UnRegister(c *Client) {
	ws.unregisterChan <- c
}

func (ws *WsServer) Validate(s interface{}) error {
	return nil
}

func (ws *WsServer) GetUserAllCons(userID string) ([]*Client, bool) {
	return ws.clients.GetAll(userID)
}

func (ws *WsServer) GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool) {
	return ws.clients.Get(userID, platform)
}

func NewWsServer(opts ...Option) (*WsServer, error) {
	var config configs
	for _, o := range opts {
		o(&config)
	}
	v := validator.New()
	return &WsServer{
		port:             config.port,
		wsMaxConnNum:     config.maxConnNum,
		handshakeTimeout: config.handshakeTimeout,
		clientPool: sync.Pool{
			New: func() interface{} {
				return new(Client)
			},
		},
		registerChan:    make(chan *Client, 1000),
		unregisterChan:  make(chan *Client, 1000),
		kickHandlerChan: make(chan *kickHandler, 1000),
		validate:        v,
		clients:         newUserMap(),
		Compressor:      NewGzipCompressor(),
		Encoder:         NewGobEncoder(),
	}, nil
}

func (ws *WsServer) Run() error {
	var client *Client
	go func() {
		for {
			select {
			case client = <-ws.registerChan:
				ws.registerClient(client)
			case client = <-ws.unregisterChan:
				ws.unregisterClient(client)
			case onlineInfo := <-ws.kickHandlerChan:
				ws.multiTerminalLoginChecker(onlineInfo.clientOK, onlineInfo.oldClients, onlineInfo.newClient)
			}
		}
	}()
	http.HandleFunc("/", ws.wsHandler)
	// http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {})
	return http.ListenAndServe(":"+utils.IntToString(ws.port), nil) // Start listening
}

func (ws *WsServer) sendUserOnlineInfoToOtherNode(ctx context.Context, client *Client) error {
	conns, err := ws.disCov.GetConns(ctx, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		return err
	}
	// Online push user online message to other node
	for _, v := range conns {
		if v.Target() == ws.disCov.GetSelfConnTarget() {
			log.ZDebug(ctx, "Filter out this node", "node", v.Target())
			continue
		}
		msgClient := msggateway.NewMsgGatewayClient(v)
		_, err := msgClient.MultiTerminalLoginCheck(ctx, &msggateway.MultiTerminalLoginCheckReq{
			UserID:     client.UserID,
			PlatformID: int32(client.PlatformID), Token: client.token,
		})
		if err != nil {
			log.ZWarn(ctx, "MultiTerminalLoginCheck err", err, "node", v.Target())
			continue
		}
	}
	return nil
}

func (ws *WsServer) SetKickHandlerInfo(i *kickHandler) {
	ws.kickHandlerChan <- i
}

func (ws *WsServer) registerClient(client *Client) {
	var (
		userOK     bool
		clientOK   bool
		oldClients []*Client
	)
	oldClients, userOK, clientOK = ws.clients.Get(client.UserID, client.PlatformID)
	if !userOK {
		ws.clients.Set(client.UserID, client)
		log.ZDebug(client.ctx, "user not exist", "userID", client.UserID, "platformID", client.PlatformID)
		prommetrics.OnlineUserGauge.Add(1)
		ws.onlineUserNum.Add(1)
		ws.onlineUserConnNum.Add(1)
	} else {
		ws.multiTerminalLoginChecker(clientOK, oldClients, client)
		log.ZDebug(client.ctx, "user exist", "userID", client.UserID, "platformID", client.PlatformID)
		if clientOK {
			ws.clients.Set(client.UserID, client)
			// 已经有同平台的连接存在
			log.ZInfo(client.ctx, "repeat login", "userID", client.UserID, "platformID", client.PlatformID, "old remote addr", getRemoteAdders(oldClients))
			ws.onlineUserConnNum.Add(1)
		} else {
			ws.clients.Set(client.UserID, client)
			ws.onlineUserConnNum.Add(1)
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = ws.sendUserOnlineInfoToOtherNode(client.ctx, client)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ws.SetUserOnlineStatus(client.ctx, client, constant.Online)
	}()

	wg.Wait()

	log.ZInfo(
		client.ctx,
		"user online",
		"online user Num",
		ws.onlineUserNum.Load(),
		"online user conn Num",
		ws.onlineUserConnNum.Load(),
	)
}

func getRemoteAdders(client []*Client) string {
	var ret string
	for i, c := range client {
		if i == 0 {
			ret = c.ctx.GetRemoteAddr()
		} else {
			ret += "@" + c.ctx.GetRemoteAddr()
		}
	}
	return ret
}

func (ws *WsServer) KickUserConn(client *Client) error {
	ws.clients.deleteClients(client.UserID, []*Client{client})
	return client.KickOnlineMessage()
}

func (ws *WsServer) multiTerminalLoginChecker(clientOK bool, oldClients []*Client, newClient *Client) {
	switch config.Config.MultiLoginPolicy {
	case constant.DefalutNotKick:
	case constant.PCAndOther:
		if constant.PlatformIDToClass(newClient.PlatformID) == constant.TerminalPC {
			return
		}
		fallthrough
	case constant.AllLoginButSameTermKick:
		if clientOK {
			isDeleteUser := ws.clients.deleteClients(newClient.UserID, oldClients)
			if isDeleteUser {
				ws.onlineUserNum.Add(-1)
			}
			for _, c := range oldClients {
				err := c.KickOnlineMessage()
				if err != nil {
					log.ZWarn(c.ctx, "KickOnlineMessage", err)
				}
			}
			m, err := ws.cache.GetTokensWithoutError(
				newClient.ctx,
				newClient.UserID,
				newClient.PlatformID,
			)
			if err != nil && err != redis.Nil {
				log.ZWarn(
					newClient.ctx,
					"get token from redis err",
					err,
					"userID",
					newClient.UserID,
					"platformID",
					newClient.PlatformID,
				)
				return
			}
			if m == nil {
				log.ZWarn(
					newClient.ctx,
					"m is nil",
					errors.New("m is nil"),
					"userID",
					newClient.UserID,
					"platformID",
					newClient.PlatformID,
				)
				return
			}
			log.ZDebug(
				newClient.ctx,
				"get token from redis",
				"userID",
				newClient.UserID,
				"platformID",
				newClient.PlatformID,
				"tokenMap",
				m,
			)

			for k := range m {
				if k != newClient.ctx.GetToken() {
					m[k] = constant.KickedToken
				}
			}
			log.ZDebug(newClient.ctx, "set token map is ", "token map", m, "userID",
				newClient.UserID, "token", newClient.ctx.GetToken())
			err = ws.cache.SetTokenMapByUidPid(newClient.ctx, newClient.UserID, newClient.PlatformID, m)
			if err != nil {
				log.ZWarn(newClient.ctx, "SetTokenMapByUidPid err", err, "userID", newClient.UserID, "platformID", newClient.PlatformID)
				return
			}
		}
	}
}

func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.clients.delete(client.UserID, client.ctx.GetRemoteAddr())
	if isDeleteUser {
		ws.onlineUserNum.Add(-1)
		prommetrics.OnlineUserGauge.Dec()
	}
	ws.onlineUserConnNum.Add(-1)
	ws.SetUserOnlineStatus(client.ctx, client, constant.Offline)
	log.ZInfo(client.ctx, "user offline", "close reason", client.closedErr, "online user Num", ws.onlineUserNum, "online user conn Num",
		ws.onlineUserConnNum.Load(),
	)
}

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	connContext := newContext(w, r)
	if ws.onlineUserConnNum.Load() >= ws.wsMaxConnNum {
		httpError(connContext, errs.ErrConnOverMaxNumLimit)
		return
	}
	var (
		token         string
		userID        string
		platformIDStr string
		exists        bool
		compression   bool
	)

	token, exists = connContext.Query(Token)
	if !exists {
		httpError(connContext, errs.ErrConnArgsErr)
		return
	}
	userID, exists = connContext.Query(WsUserID)
	if !exists {
		httpError(connContext, errs.ErrConnArgsErr)
		return
	}
	platformIDStr, exists = connContext.Query(PlatformID)
	if !exists {
		httpError(connContext, errs.ErrConnArgsErr)
		return
	}
	platformID, err := strconv.Atoi(platformIDStr)
	if err != nil {
		httpError(connContext, errs.ErrConnArgsErr)
		return
	}
	if err := authverify.WsVerifyToken(token, userID, platformID); err != nil {
		httpError(connContext, err)
		return
	}
	m, err := ws.cache.GetTokensWithoutError(context.Background(), userID, platformID)
	if err != nil {
		httpError(connContext, err)
		return
	}
	if v, ok := m[token]; ok {
		switch v {
		case constant.NormalToken:
		case constant.KickedToken:
			httpError(connContext, errs.ErrTokenKicked.Wrap())
			return
		default:
			httpError(connContext, errs.ErrTokenUnknown.Wrap())
			return
		}
	} else {
		httpError(connContext, errs.ErrTokenNotExist.Wrap())
		return
	}
	wsLongConn := newGWebSocket(WebSocket, ws.handshakeTimeout)
	err = wsLongConn.GenerateLongConn(w, r)
	if err != nil {
		httpError(connContext, err)
		return
	}
	compressProtoc, exists := connContext.Query(Compression)
	if exists {
		if compressProtoc == GzipCompressionProtocol {
			compression = true
		}
	}
	compressProtoc, exists = connContext.GetHeader(Compression)
	if exists {
		if compressProtoc == GzipCompressionProtocol {
			compression = true
		}
	}
	client := ws.clientPool.Get().(*Client)
	client.ResetClient(connContext, wsLongConn, connContext.GetBackground(), compression, ws, token)
	ws.registerChan <- client
	go client.readMessage()
}

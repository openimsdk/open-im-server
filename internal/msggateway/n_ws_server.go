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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/OpenIMSDK/tools/apiresp"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msggateway"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type LongConnServer interface {
	Run() error
	wsHandler(w http.ResponseWriter, r *http.Request)
	GetUserAllCons(userID string) ([]*Client, bool)
	GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool)
	Validate(s any) error
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
	New: func() any {
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
	writeBufferSize   int
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

func (ws *WsServer) Validate(s any) error {
	//?question?
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
		writeBufferSize:  config.writeBufferSize,
		handshakeTimeout: config.handshakeTimeout,
		clientPool: sync.Pool{
			New: func() any {
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
	var (
		client *Client
		wg     errgroup.Group

		sigs = make(chan os.Signal, 1)
		done = make(chan struct{}, 1)
	)

	server := http.Server{Addr: ":" + utils.IntToString(ws.port), Handler: nil}

	wg.Go(func() error {
		for {
			select {
			case <-done:
				return nil

			case client = <-ws.registerChan:
				ws.registerClient(client)
			case client = <-ws.unregisterChan:
				ws.unregisterClient(client)
			case onlineInfo := <-ws.kickHandlerChan:
				ws.multiTerminalLoginChecker(onlineInfo.clientOK, onlineInfo.oldClients, onlineInfo.newClient)
			}
		}
	})

	wg.Go(func() error {
		http.HandleFunc("/", ws.wsHandler)
		return server.ListenAndServe()
	})

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigs

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// graceful exit operation for server
		_ = server.Shutdown(ctx)
		_ = wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil

	case <-time.After(15 * time.Second):
		return utils.Wrap1(errors.New("timeout exit"))
	}

}

var concurrentRequest = 3

func (ws *WsServer) sendUserOnlineInfoToOtherNode(ctx context.Context, client *Client) error {
	conns, err := ws.disCov.GetConns(ctx, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		return err
	}

	wg := errgroup.Group{}
	wg.SetLimit(concurrentRequest)

	// Online push user online message to other node
	for _, v := range conns {
		v := v // safe closure var
		if v.Target() == ws.disCov.GetSelfConnTarget() {
			log.ZDebug(ctx, "Filter out this node", "node", v.Target())
			continue
		}

		wg.Go(func() error {
			msgClient := msggateway.NewMsgGatewayClient(v)
			_, err := msgClient.MultiTerminalLoginCheck(ctx, &msggateway.MultiTerminalLoginCheckReq{
				UserID:     client.UserID,
				PlatformID: int32(client.PlatformID), Token: client.token,
			})
			if err != nil {
				log.ZWarn(ctx, "MultiTerminalLoginCheck err", err, "node", v.Target())
			}
			return nil
		})
	}

	_ = wg.Wait()
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
	if config.Config.Envs.Discovery == "zookeeper" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ws.sendUserOnlineInfoToOtherNode(client.ctx, client)
		}()
	}
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
		if !clientOK {
			return
		}
		ws.clients.deleteClients(newClient.UserID, oldClients)
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

func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.clients.delete(client.UserID, client.ctx.GetRemoteAddr())
	if isDeleteUser {
		ws.onlineUserNum.Add(-1)
		prommetrics.OnlineUserGauge.Dec()
	}
	ws.onlineUserConnNum.Add(-1)
	ws.SetUserOnlineStatus(client.ctx, client, constant.Offline)
	log.ZInfo(client.ctx, "user offline", "close reason", client.closedErr, "online user Num", ws.onlineUserNum.Load(), "online user conn Num",
		ws.onlineUserConnNum.Load(),
	)
}

func (ws *WsServer) ParseWSArgs(r *http.Request) (args *WSArgs, err error) {
	var v WSArgs
	defer func() {
		args = &v
	}()
	query := r.URL.Query()
	v.MsgResp, _ = strconv.ParseBool(query.Get(MsgResp))
	if ws.onlineUserConnNum.Load() >= ws.wsMaxConnNum {
		return nil, errs.ErrConnOverMaxNumLimit.Wrap("over max conn num limit")
	}
	if v.Token = query.Get(Token); v.Token == "" {
		return nil, errs.ErrConnArgsErr.Wrap("token is empty")
	}
	if v.UserID = query.Get(WsUserID); v.UserID == "" {
		return nil, errs.ErrConnArgsErr.Wrap("sendID is empty")
	}
	platformIDStr := query.Get(PlatformID)
	if platformIDStr == "" {
		return nil, errs.ErrConnArgsErr.Wrap("platformID is empty")
	}
	platformID, err := strconv.Atoi(platformIDStr)
	if err != nil {
		return nil, errs.ErrConnArgsErr.Wrap("platformID is not int")
	}
	v.PlatformID = platformID
	if err = authverify.WsVerifyToken(v.Token, v.UserID, platformID); err != nil {
		return nil, err
	}
	if query.Get(Compression) == GzipCompressionProtocol {
		v.Compression = true
	}
	if r.Header.Get(Compression) == GzipCompressionProtocol {
		v.Compression = true
	}
	m, err := ws.cache.GetTokensWithoutError(context.Background(), v.UserID, platformID)
	if err != nil {
		return nil, err
	}
	if v, ok := m[v.Token]; ok {
		switch v {
		case constant.NormalToken:
		case constant.KickedToken:
			return nil, errs.ErrTokenKicked.Wrap()
		default:
			return nil, errs.ErrTokenUnknown.Wrap(fmt.Sprintf("token status is %d", v))
		}
	} else {
		return nil, errs.ErrTokenNotExist.Wrap()
	}
	return &v, nil
}

type WSArgs struct {
	Token       string
	UserID      string
	PlatformID  int
	Compression bool
	MsgResp     bool
}

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	connContext := newContext(w, r)
	args, pErr := ws.ParseWSArgs(r)
	var wsLongConn *GWebSocket
	if args.MsgResp {
		wsLongConn = newGWebSocket(WebSocket, ws.handshakeTimeout, ws.writeBufferSize)
		if err := wsLongConn.GenerateLongConn(w, r); err != nil {
			httpError(connContext, err)
			return
		}
		data, err := json.Marshal(apiresp.ParseError(pErr))
		if err != nil {
			_ = wsLongConn.Close()
			return
		}
		if err := wsLongConn.WriteMessage(MessageText, data); err != nil {
			_ = wsLongConn.Close()
			return
		}
		if pErr != nil {
			_ = wsLongConn.Close()
			return
		}
	} else {
		if pErr != nil {
			httpError(connContext, pErr)
			return
		}
		wsLongConn = newGWebSocket(WebSocket, ws.handshakeTimeout, ws.writeBufferSize)
		if err := wsLongConn.GenerateLongConn(w, r); err != nil {
			httpError(connContext, err)
			return
		}
	}
	client := ws.clientPool.Get().(*Client)
	client.ResetClient(connContext, wsLongConn, connContext.GetBackground(), args.Compression, ws, args.Token)
	ws.registerChan <- client
	go client.readMessage()
}

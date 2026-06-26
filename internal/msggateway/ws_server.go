package msggateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/tools/apiresp"

	"github.com/go-playground/validator/v10"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	pbAuth "github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msggateway"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/stringutil"
	"golang.org/x/sync/errgroup"
)

var wsSuccessResponse, _ = json.Marshal(&apiresp.ApiResponse{})

type LongConnServer interface {
	Run(done chan error) error
	wsHandler(w http.ResponseWriter, r *http.Request)
	GetUserAllCons(userID string) ([]*Client, bool)
	GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool)
	Validate(s any) error
	SetDiscoveryRegistry(ctx context.Context, client discovery.SvcDiscoveryRegistry, config *Config) error
	KickUserConn(client *Client) error
	UnRegister(c *Client)
	SetKickHandlerInfo(i *kickHandler)
	SubUserOnlineStatus(ctx context.Context, client *Client, data *Req) ([]byte, error)
	Compressor
	MessageHandler
}

type WsServer struct {
	websocket         *websocket.Upgrader
	msgGatewayConfig  *Config
	port              int
	wsMaxConnNum      int64
	registerChan      chan *Client
	unregisterChan    chan *Client
	kickHandlerChan   chan *kickHandler
	clients           UserMap
	online            *rpccache.OnlineCache
	subscription      *Subscription
	clientPool        sync.Pool
	onlineUserNum     atomic.Int64
	onlineUserConnNum atomic.Int64
	handshakeTimeout  time.Duration
	writeBufferSize   int
	validate          *validator.Validate
	disCov            discovery.SvcDiscoveryRegistry
	Compressor
	//Encoder
	MessageHandler
	webhookClient *webhook.Client
	userClient    *rpcli.UserClient
	authClient    *rpcli.AuthClient
}

type kickHandler struct {
	clientOK   bool
	oldClients []*Client
	newClient  *Client
}

func (ws *WsServer) SetDiscoveryRegistry(ctx context.Context, disCov discovery.SvcDiscoveryRegistry, config *Config) error {
	userConn, err := disCov.GetConn(ctx, config.Share.RpcRegisterName.User)
	if err != nil {
		return err
	}
	pushConn, err := disCov.GetConn(ctx, config.Share.RpcRegisterName.Push)
	if err != nil {
		return err
	}
	authConn, err := disCov.GetConn(ctx, config.Share.RpcRegisterName.Auth)
	if err != nil {
		return err
	}
	msgConn, err := disCov.GetConn(ctx, config.Share.RpcRegisterName.Msg)
	if err != nil {
		return err
	}
	ws.userClient = rpcli.NewUserClient(userConn)
	ws.authClient = rpcli.NewAuthClient(authConn)
	ws.MessageHandler = NewGrpcHandler(ws.validate, rpcli.NewMsgClient(msgConn), rpcli.NewPushMsgServiceClient(pushConn))
	ws.disCov = disCov
	return nil
}

//func (ws *WsServer) SetUserOnlineStatus(ctx context.Context, client *Client, status int32) {
//	err := ws.userClient.SetUserStatus(ctx, client.UserID, status, client.PlatformID)
//	if err != nil {
//		log.ZWarn(ctx, "SetUserStatus err", err)
//	}
//	switch status {
//	case constant.Online:
//		ws.webhookAfterUserOnline(ctx, &ws.msgGatewayConfig.WebhooksConfig.AfterUserOnline, client.UserID, client.PlatformID, client.IsBackground, client.ctx.GetConnID())
//	case constant.Offline:
//		ws.webhookAfterUserOffline(ctx, &ws.msgGatewayConfig.WebhooksConfig.AfterUserOffline, client.UserID, client.PlatformID, client.ctx.GetConnID())
//	}
//}

func (ws *WsServer) UnRegister(c *Client) {
	ws.unregisterChan <- c
}

func (ws *WsServer) Validate(_ any) error {
	return nil
}

func (ws *WsServer) GetUserAllCons(userID string) ([]*Client, bool) {
	return ws.clients.GetAll(userID)
}

func (ws *WsServer) GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool) {
	return ws.clients.Get(userID, platform)
}

func NewWsServer(msgGatewayConfig *Config, opts ...Option) *WsServer {
	var config configs
	for _, o := range opts {
		o(&config)
	}
	//userRpcClient := rpcclient.NewUserRpcClient(client, config.Share.RpcRegisterName.User, config.Share.IMAdminUserID)
	upgrader := &websocket.Upgrader{
		HandshakeTimeout: config.handshakeTimeout,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	v := validator.New()
	return &WsServer{
		websocket:        upgrader,
		msgGatewayConfig: msgGatewayConfig,
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
		subscription:    newSubscription(),
		Compressor:      NewGzipCompressor(),
		webhookClient:   webhook.NewWebhookClient(msgGatewayConfig.WebhooksConfig.URL),
	}
}

func (ws *WsServer) Run(done chan error) error {
	var (
		client       *Client
		netErr       error
		shutdownDone = make(chan struct{}, 1)
	)

	server := http.Server{Addr: ":" + stringutil.IntToString(ws.port), Handler: nil}

	go func() {
		for {
			select {
			case <-shutdownDone:
				return
			case client = <-ws.registerChan:
				ws.registerClient(client)
			case client = <-ws.unregisterChan:
				ws.unregisterClient(client)
			case onlineInfo := <-ws.kickHandlerChan:
				ws.multiTerminalLoginChecker(onlineInfo.clientOK, onlineInfo.oldClients, onlineInfo.newClient)
			}
		}
	}()
	netDone := make(chan struct{}, 1)
	go func() {
		http.HandleFunc("/", ws.wsHandler)
		err := server.ListenAndServe()
		defer close(netDone)
		if err != nil && err != http.ErrServerClosed {
			netErr = errs.WrapMsg(err, "ws start err", server.Addr)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var err error
	select {
	case err = <-done:
		sErr := server.Shutdown(ctx)
		if sErr != nil {
			return errs.WrapMsg(sErr, "shutdown err")
		}
		close(shutdownDone)
		if err != nil {
			return err
		}
	case <-netDone:
	}
	return netErr

}

var concurrentRequest = 3

func (ws *WsServer) sendUserOnlineInfoToOtherNode(ctx context.Context, client *Client) error {
	conns, err := ws.disCov.GetConns(ctx, ws.msgGatewayConfig.Share.RpcRegisterName.MessageGateway)
	if err != nil {
		return err
	}
	if len(conns) == 0 || (len(conns) == 1 && ws.disCov.IsSelfNode(conns[0])) {
		return nil
	}

	wg := errgroup.Group{}
	wg.SetLimit(concurrentRequest)

	// Online push user online message to other node
	for _, v := range conns {
		v := v
		log.ZDebug(ctx, "sendUserOnlineInfoToOtherNode conn")
		if ws.disCov.IsSelfNode(v) {
			log.ZDebug(ctx, "Filter out this node")
			continue
		}

		wg.Go(func() error {
			msgClient := msggateway.NewMsgGatewayClient(v)
			_, err := msgClient.MultiTerminalLoginCheck(ctx, &msggateway.MultiTerminalLoginCheckReq{
				UserID:     client.UserID,
				PlatformID: int32(client.PlatformID), Token: client.token,
			})
			if err != nil {
				log.ZWarn(ctx, "MultiTerminalLoginCheck err", err)
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
			// There is already a connection to the platform
			log.ZDebug(client.ctx, "repeat login", "userID", client.UserID, "platformID",
				client.PlatformID, "old remote addr", getRemoteAdders(oldClients))
			ws.onlineUserConnNum.Add(1)
		} else {
			ws.clients.Set(client.UserID, client)
			ws.onlineUserConnNum.Add(1)
		}
	}

	log.ZDebug(client.ctx, "ws.msgGatewayConfig.Discovery.Enable", "discoveryEnable", ws.msgGatewayConfig.Discovery.Enable)

	if ws.msgGatewayConfig.Discovery.Enable != config.KUBERNETES {
		_ = ws.sendUserOnlineInfoToOtherNode(client.ctx, client)
	}

	log.ZDebug(client.ctx, "user online", "online user Num", ws.onlineUserNum.Load(), "online user conn Num", ws.onlineUserConnNum.Load())
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
	ws.clients.DeleteClients(client.UserID, []*Client{client})
	return client.KickOnlineMessage()
}

func (ws *WsServer) multiTerminalLoginChecker(clientOK bool, oldClients []*Client, newClient *Client) {
	kickTokenFunc := func(kickClients []*Client) {
		var kickTokens []string
		ws.clients.DeleteClients(newClient.UserID, kickClients)
		for _, c := range kickClients {
			kickTokens = append(kickTokens, c.token)
			err := c.KickOnlineMessage()
			if err != nil {
				log.ZWarn(c.ctx, "KickOnlineMessage", err)
			}
		}
		ctx := mcontext.WithMustInfoCtx(
			[]string{newClient.ctx.GetOperationID(), newClient.ctx.GetUserID(),
				constant.PlatformIDToName(newClient.PlatformID), newClient.ctx.GetConnID()},
		)
		if err := ws.authClient.KickTokens(ctx, kickTokens); err != nil {
			log.ZWarn(newClient.ctx, "kickTokens err", err)
		}
	}

	// If reconnect: When multiple msgGateway instances are deployed, a client may disconnect from instance A and reconnect to instance B.
	// During this process, instance A might still be executing, resulting in two clients with the same token existing simultaneously.
	// This situation needs to be filtered to prevent duplicate clients.
	checkSameTokenFunc := func(oldClients []*Client) []*Client {
		var clientsNeedToKick []*Client

		for _, c := range oldClients {
			if c.token == newClient.token {
				log.ZDebug(newClient.ctx, "token is same, not kick",
					"userID", newClient.UserID,
					"platformID", newClient.PlatformID,
					"token", newClient.token)
				continue
			}

			clientsNeedToKick = append(clientsNeedToKick, c)
		}

		return clientsNeedToKick
	}

	switch ws.msgGatewayConfig.Share.MultiLogin.Policy {
	case constant.DefalutNotKick:
	case constant.PCAndOther:
		if constant.PlatformIDToClass(newClient.PlatformID) == constant.TerminalPC {
			return
		}
		clients, ok := ws.clients.GetAll(newClient.UserID)
		clientOK = ok
		oldClients = make([]*Client, 0, len(clients))
		for _, c := range clients {
			if constant.PlatformIDToClass(c.PlatformID) == constant.TerminalPC {
				continue
			}
			oldClients = append(oldClients, c)
		}

		fallthrough
	case constant.AllLoginButSameTermKick:
		if !clientOK {
			return
		}

		oldClients = checkSameTokenFunc(oldClients)

		ws.clients.DeleteClients(newClient.UserID, oldClients)
		for _, c := range oldClients {
			err := c.KickOnlineMessage()
			if err != nil {
				log.ZWarn(c.ctx, "KickOnlineMessage", err)
			}
		}
		ctx := mcontext.WithMustInfoCtx(
			[]string{newClient.ctx.GetOperationID(), newClient.ctx.GetUserID(),
				constant.PlatformIDToName(newClient.PlatformID), newClient.ctx.GetConnID()},
		)
		req := &pbAuth.InvalidateTokenReq{
			PreservedToken: newClient.token,
			UserID:         newClient.UserID,
			PlatformID:     int32(newClient.PlatformID),
		}
		if err := ws.authClient.InvalidateToken(ctx, req); err != nil {
			log.ZWarn(newClient.ctx, "InvalidateToken err", err, "userID", newClient.UserID,
				"platformID", newClient.PlatformID)
		}
	case constant.AllLoginButSameClassKick:
		clients, ok := ws.clients.GetAll(newClient.UserID)
		if !ok {
			return
		}

		var kickClients []*Client
		for _, client := range clients {
			if constant.PlatformIDToClass(client.PlatformID) == constant.PlatformIDToClass(newClient.PlatformID) {
				kickClients = append(kickClients, client)
			}
		}
		kickClients = checkSameTokenFunc(kickClients)

		kickTokenFunc(kickClients)
	}
}

func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.clients.DeleteClients(client.UserID, []*Client{client})
	if isDeleteUser {
		ws.onlineUserNum.Add(-1)
		prommetrics.OnlineUserGauge.Dec()
	}
	ws.onlineUserConnNum.Add(-1)
	ws.subscription.DelClient(client)
	//ws.SetUserOnlineStatus(client.ctx, client, constant.Offline)
	log.ZDebug(client.ctx, "user offline", "close reason", client.closedErr, "online user Num",
		ws.onlineUserNum.Load(), "online user conn Num",
		ws.onlineUserConnNum.Load(),
	)
}

// validateRespWithRequest checks if the response matches the expected userID and platformID.
func (ws *WsServer) validateRespWithRequest(ctx *UserConnContext, resp *pbAuth.ParseTokenResp) error {
	userID := ctx.GetUserID()
	platformID := ctx.GetPlatformID()
	if resp.UserID != userID {
		return servererrs.ErrTokenInvalid.WrapMsg(fmt.Sprintf("token uid %s != userID %s", resp.UserID, userID))
	}
	if int(resp.PlatformID) != platformID {
		return servererrs.ErrTokenInvalid.WrapMsg(fmt.Sprintf("token platform %d != platformID %d", resp.PlatformID, platformID))
	}
	return nil
}

func (ws *WsServer) handlerError(ctx *UserConnContext, w http.ResponseWriter, r *http.Request, err error) {
	if !ctx.ShouldSendResp() {
		httpError(ctx, err)
		return
	}
	// the browser cannot get the response of upgrade failure
	data, err := json.Marshal(apiresp.ParseError(err))
	if err != nil {
		log.ZError(ctx, "json marshal failed", err)
		return
	}
	conn, upgradeErr := ws.websocket.Upgrade(w, r, nil)
	if upgradeErr != nil {
		log.ZWarn(ctx, "websocket upgrade failed", upgradeErr, "respErr", err, "resp", string(data))
		return
	}
	defer conn.Close()
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.ZWarn(ctx, "WriteMessage failed", err, "respErr", err, "resp", string(data))
		return
	}
}

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new connection context
	connContext := newContext(w, r)

	// Check if the current number of online user connections exceeds the maximum limit
	if ws.onlineUserConnNum.Load() >= ws.wsMaxConnNum {
		// If it exceeds the maximum connection number, return an error via HTTP and stop processing
		ws.handlerError(connContext, w, r, servererrs.ErrConnOverMaxNumLimit.WrapMsg("over max conn num limit"))
		return
	}

	// Parse essential arguments (e.g., user ID, Token)
	err := connContext.ParseEssentialArgs()
	if err != nil {
		// If there's an error during parsing, return an error via HTTP and stop processing
		ws.handlerError(connContext, w, r, err)
		return
	}

	// Call the authentication client to parse the Token obtained from the context
	resp, err := ws.authClient.ParseToken(connContext, connContext.GetToken())
	if err != nil {
		ws.handlerError(connContext, w, r, err)
		return
	}

	// Validate the authentication response matches the request (e.g., user ID and platform ID)
	err = ws.validateRespWithRequest(connContext, resp)
	if err != nil {
		// If validation fails, return an error via HTTP and stop processing
		ws.handlerError(connContext, w, r, err)
		return
	}
	conn, err := ws.websocket.Upgrade(w, r, nil)
	if err != nil {
		log.ZWarn(connContext, "websocket upgrade failed", err)
		return
	}
	if connContext.ShouldSendResp() {
		if err := conn.WriteMessage(websocket.TextMessage, wsSuccessResponse); err != nil {
			log.ZWarn(connContext, "WriteMessage first response", err)
			return
		}
	}
	log.ZDebug(connContext, "new conn", "token", connContext.GetToken())

	var pingInterval time.Duration
	if connContext.GetPlatformID() == constant.WebPlatformID {
		pingInterval = pingPeriod
	}

	// Retrieve a client object from the client pool, reset its state, and associate it with the current WebSocket long connection
	client := ws.clientPool.Get().(*Client)
	client.ResetClient(connContext, NewWebSocketClientConn(conn, maxMessageSize, pongWait, pingInterval), ws)

	// Register the client with the server and start message processing
	ws.registerChan <- client
	go client.readMessage()
}

package msggateway

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type LongConnServer interface {
	Run() error
	wsHandler(w http.ResponseWriter, r *http.Request)
	GetUserAllCons(userID string) ([]*Client, bool)
	GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool)
	Validate(s interface{}) error
	SetCacheHandler(cache cache.MsgModel)
	SetDiscoveryRegistry(client discoveryregistry.SvcDiscoveryRegistry)
	UnRegister(c *Client)
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
	onlineUserNum     int64
	onlineUserConnNum int64
	handshakeTimeout  time.Duration
	hubServer         *Server
	validate          *validator.Validate
	cache             cache.MsgModel
	Compressor
	Encoder
	MessageHandler
}
type kickHandler struct {
	clientOK   bool
	oldClients []*Client
	newClient  *Client
}

func (ws *WsServer) SetDiscoveryRegistry(client discoveryregistry.SvcDiscoveryRegistry) {
	ws.MessageHandler = NewGrpcHandler(ws.validate, client)
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
	if config.port < 1024 {
		return nil, errors.New("port not allow to listen")

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
				ws.multiTerminalLoginChecker(onlineInfo)
			}
		}
	}()
	http.HandleFunc("/", ws.wsHandler)
	// http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {})
	return http.ListenAndServe(":"+utils.IntToString(ws.port), nil) //Start listening
}

func (ws *WsServer) registerClient(client *Client) {
	var (
		userOK     bool
		clientOK   bool
		oldClients []*Client
	)
	ws.clients.Set(client.UserID, client)
	oldClients, userOK, clientOK = ws.clients.Get(client.UserID, client.PlatformID)
	if !userOK {
		log.ZDebug(client.ctx, "user not exist", "userID", client.UserID, "platformID", client.PlatformID)
		atomic.AddInt64(&ws.onlineUserNum, 1)
		atomic.AddInt64(&ws.onlineUserConnNum, 1)

	} else {
		i := &kickHandler{
			clientOK:   clientOK,
			oldClients: oldClients,
			newClient:  client,
		}
		ws.kickHandlerChan <- i
		log.ZDebug(client.ctx, "user exist", "userID", client.UserID, "platformID", client.PlatformID)
		if clientOK { //已经有同平台的连接存在
			log.ZInfo(client.ctx, "repeat login", "userID", client.UserID, "platformID", client.PlatformID, "old remote addr", getRemoteAdders(oldClients))
			atomic.AddInt64(&ws.onlineUserConnNum, 1)
		} else {
			atomic.AddInt64(&ws.onlineUserConnNum, 1)
		}
	}
	log.ZInfo(client.ctx, "user online", "online user Num", ws.onlineUserNum, "online user conn Num", ws.onlineUserConnNum)
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

func (ws *WsServer) multiTerminalLoginChecker(info *kickHandler) {
	switch config.Config.MultiLoginPolicy {
	case constant.DefalutNotKick:
	case constant.PCAndOther:
		if constant.PlatformIDToClass(info.newClient.PlatformID) == constant.TerminalPC {
			return
		}
		fallthrough
	case constant.AllLoginButSameTermKick:
		if info.clientOK {
			for _, c := range info.oldClients {
				err := c.KickOnlineMessage()
				if err != nil {
					log.ZError(c.ctx, "KickOnlineMessage", err)
				}
			}
		}
	}
}

func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.clients.delete(client.UserID, client.ctx.GetRemoteAddr())
	if isDeleteUser {
		atomic.AddInt64(&ws.onlineUserNum, -1)
	}
	atomic.AddInt64(&ws.onlineUserConnNum, -1)
	log.ZInfo(client.ctx, "user offline", "close reason", client.closedErr, "online user Num", ws.onlineUserNum, "online user conn Num", ws.onlineUserConnNum)
}

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	connContext := newContext(w, r)
	if ws.onlineUserConnNum >= ws.wsMaxConnNum {
		httpError(connContext, errs.ErrConnOverMaxNumLimit)
		return
	}
	var (
		token       string
		userID      string
		platformID  string
		exists      bool
		compression bool
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
	platformID, exists = connContext.Query(PlatformID)
	if !exists || utils.StringToInt(platformID) == 0 {
		httpError(connContext, errs.ErrConnArgsErr)
		return
	}
	err := tokenverify.WsVerifyToken(token, userID, platformID)
	if err != nil {
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
	client.ResetClient(connContext, wsLongConn, connContext.GetBackground(), compression, ws)
	ws.registerChan <- client
	go client.readMessage()
}

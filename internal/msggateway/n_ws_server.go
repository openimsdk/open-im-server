package msggateway

import (
	context2 "context"
	"errors"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

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
	//SetMessageHandler(msgRpcClient *rpcclient.MsgClient)
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
	clients           *UserMap
	clientPool        sync.Pool
	onlineUserNum     int64
	onlineUserConnNum int64
	handshakeTimeout  time.Duration
	hubServer         *Server
	validate          *validator.Validate
	Compressor
	Encoder
	MessageHandler
}

func (ws *WsServer) SetDiscoveryRegistry(client discoveryregistry.SvcDiscoveryRegistry) {
	ws.MessageHandler = NewGrpcHandler(ws.validate, client)
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
		registerChan:   make(chan *Client, 1000),
		unregisterChan: make(chan *Client, 1000),
		validate:       v,
		clients:        newUserMap(),
		Compressor:     NewGzipCompressor(),
		Encoder:        NewGobEncoder(),
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
			}
		}
	}()
	http.HandleFunc("/", ws.wsHandler)
	// http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {})
	return http.ListenAndServe(":"+utils.IntToString(ws.port), nil) //Start listening
}

func (ws *WsServer) registerClient(client *Client) {
	var (
		userOK   bool
		clientOK bool
		cli      []*Client
	)
	cli, userOK, clientOK = ws.clients.Get(client.UserID, client.PlatformID)
	if !userOK {
		log.ZDebug(client.ctx, "user not exist", "userID", client.UserID, "platformID", client.PlatformID)
		ws.clients.Set(client.UserID, client)
		atomic.AddInt64(&ws.onlineUserNum, 1)
		atomic.AddInt64(&ws.onlineUserConnNum, 1)

	} else {
		log.ZDebug(client.ctx, "user exist", "userID", client.UserID, "platformID", client.PlatformID)
		if clientOK { //已经有同平台的连接存在
			ws.clients.Set(client.UserID, client)
			ws.multiTerminalLoginChecker(cli)
			log.ZInfo(client.ctx, "repeat login", "userID", client.UserID, "platformID", client.PlatformID, "old remote addr", getRemoteAdders(cli))
			atomic.AddInt64(&ws.onlineUserConnNum, 1)
		} else {
			ws.clients.Set(client.UserID, client)
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
			ret += " @ " + c.ctx.GetRemoteAddr()
		}
	}
	return ret
}

func (ws *WsServer) multiTerminalLoginChecker(client []*Client) {

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
	context := newContext(w, r)
	if ws.onlineUserConnNum >= ws.wsMaxConnNum {
		httpError(context, errs.ErrConnOverMaxNumLimit)
		return
	}
	var (
		token       string
		userID      string
		platformID  string
		exists      bool
		compression bool
	)

	token, exists = context.Query(Token)
	if !exists {
		httpError(context, errs.ErrConnArgsErr)
		return
	}
	userID, exists = context.Query(WsUserID)
	if !exists {
		httpError(context, errs.ErrConnArgsErr)
		return
	}
	platformID, exists = context.Query(PlatformID)
	if !exists {
		httpError(context, errs.ErrConnArgsErr)
		return
	}
	log.ZDebug(context2.Background(), "conn", "platformID", platformID)
	err := tokenverify.WsVerifyToken(token, userID, platformID)
	if err != nil {
		httpError(context, err)
		return
	}
	wsLongConn := newGWebSocket(WebSocket, ws.handshakeTimeout)
	err = wsLongConn.GenerateLongConn(w, r)
	if err != nil {
		httpError(context, err)
		return
	}
	compressProtoc, exists := context.Query(Compression)
	if exists {
		if compressProtoc == GzipCompressionProtocol {
			compression = true
		}
	}
	compressProtoc, exists = context.GetHeader(Compression)
	if exists {
		if compressProtoc == GzipCompressionProtocol {
			compression = true
		}
	}
	client := ws.clientPool.Get().(*Client)
	client.ResetClient(context, wsLongConn, compression, ws)
	ws.registerChan <- client
	go client.readMessage()
}

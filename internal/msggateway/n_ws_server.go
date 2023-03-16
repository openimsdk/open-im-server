package msggateway

import (
	"errors"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/go-playground/validator/v10"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type LongConnServer interface {
	Run() error
	wsHandler(w http.ResponseWriter, r *http.Request)
	GetUserAllCons(userID string) ([]*Client, bool)
	GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool)
	Validate(s interface{}) error
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
	port                            int
	wsMaxConnNum                    int64
	registerChan                    chan *Client
	unregisterChan                  chan *Client
	clients                         *UserMap
	clientPool                      sync.Pool
	onlineUserNum                   int64
	onlineUserConnNum               int64
	handshakeTimeout                time.Duration
	readBufferSize, WriteBufferSize int
	validate                        *validator.Validate
	Compressor
	Encoder
	MessageHandler
}

func (ws *WsServer) UnRegister(c *Client) {
	ws.unregisterChan <- c
}

func (ws *WsServer) Validate(s interface{}) error {
	//TODO implement me
	panic("implement me")
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
		readBufferSize:   config.messageMaxMsgLength,
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
		MessageHandler: NewGrpcHandler(v, nil),
		//handler:  NewGrpcHandler(validate),
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
	http.HandleFunc("/", ws.wsHandler)                              //Get request from client to handle by wsHandler
	return http.ListenAndServe(":"+utils.IntToString(ws.port), nil) //Start listening
}

func (ws *WsServer) registerClient(client *Client) {
	var (
		userOK   bool
		clientOK bool
		cli      []*Client
	)
	cli, userOK, clientOK = ws.clients.Get(client.userID, client.platformID)
	if !userOK {
		ws.clients.Set(client.userID, client)
		atomic.AddInt64(&ws.onlineUserNum, 1)
		atomic.AddInt64(&ws.onlineUserConnNum, 1)
		fmt.Println("R在线用户数量:", ws.onlineUserNum)
		fmt.Println("R在线用户连接数量:", ws.onlineUserConnNum)
	} else {
		if clientOK { //已经有同平台的连接存在
			ws.clients.Set(client.userID, client)
			ws.multiTerminalLoginChecker(cli)
		} else {
			ws.clients.Set(client.userID, client)
			atomic.AddInt64(&ws.onlineUserConnNum, 1)
			fmt.Println("R在线用户数量:", ws.onlineUserNum)
			fmt.Println("R在线用户连接数量:", ws.onlineUserConnNum)
		}
	}

}

func (ws *WsServer) multiTerminalLoginChecker(client []*Client) {

}
func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.clients.delete(client.userID, client.platformID)
	if isDeleteUser {
		atomic.AddInt64(&ws.onlineUserNum, -1)
	}
	atomic.AddInt64(&ws.onlineUserConnNum, -1)
	fmt.Println("R在线用户数量:", ws.onlineUserNum)
	fmt.Println("R在线用户连接数量:", ws.onlineUserConnNum)
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
	err := tokenverify.WsVerifyToken(token, userID, platformID)
	if err != nil {
		httpError(context, err)
		return
	}
	wsLongConn := newGWebSocket(WebSocket, ws.handshakeTimeout, ws.readBufferSize)
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

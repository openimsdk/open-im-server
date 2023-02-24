package new

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/utils"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

type LongConnServer interface {
	Run() error
}

type Server struct {
	rpcPort        int
	wsMaxConnNum   int
	longConnServer *LongConnServer
	//rpcServer      *RpcServer
}
type WsServer struct {
	port                            int
	wsMaxConnNum                    int64
	wsUpGrader                      *websocket.Upgrader
	registerChan                    chan *Client
	unregisterChan                  chan *Client
	clients                         *UserMap
	clientPool                      sync.Pool
	onlineUserNum                   int64
	onlineUserConnNum               int64
	gzipCompressor                  Compressor
	encoder                         Encoder
	handler                         MessageHandler
	handshakeTimeout                time.Duration
	readBufferSize, WriteBufferSize int
	validate                        *validator.Validate
}

func newWsServer(opts ...Option) (*WsServer, error) {
	var config configs
	for _, o := range opts {
		o(&config)
	}
	if config.port < 1024 {
		return nil, errors.New("port not allow to listen")

	}
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
		validate: validator.New(),
		clients:  newUserMap(),
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
		cli      *Client
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

func (ws *WsServer) multiTerminalLoginChecker(client *Client) {

}
func (ws *WsServer) unregisterClient(client *Client) {
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
		httpError(context, constant.ErrConnOverMaxNumLimit)
		return
	}
	var (
		token       string
		userID      string
		platformID  string
		exists      bool
		compression bool
		compressor  Compressor
	)

	token, exists = context.Query(TOKEN)
	if !exists {
		httpError(context, constant.ErrConnArgsErr)
		return
	}
	userID, exists = context.Query(WS_USERID)
	if !exists {
		httpError(context, constant.ErrConnArgsErr)
		return
	}
	platformID, exists = context.Query(PLATFORM_ID)
	if !exists {
		httpError(context, constant.ErrConnArgsErr)
		return
	}
	err := tokenverify.WsVerifyToken(token, userID, platformID)
	if err != nil {
		httpError(context, err)
		return
	}
	wsLongConn := newGWebSocket(constant.WebSocket, ws.handshakeTimeout, ws.readBufferSize)
	err = wsLongConn.GenerateLongConn(w, r)
	if err != nil {
		httpError(context, err)
		return
	}
	compressProtoc, exists := context.Query(COMPRESSION)
	if exists {
		if compressProtoc == GZIP_COMPRESSION_PROTOCAL {
			compression = true
			compressor = ws.gzipCompressor
		}
	}
	compressProtoc, exists = context.GetHeader(COMPRESSION)
	if exists {
		if compressProtoc == GZIP_COMPRESSION_PROTOCAL {
			compression = true
			compressor = ws.gzipCompressor
		}
	}
	client := ws.clientPool.Get().(*Client)
	client.ResetClient(context, wsLongConn, compression, compressor, ws.encoder, ws.handler, ws.unregisterChan, ws.validate)
	ws.registerChan <- client
	go client.readMessage()
}

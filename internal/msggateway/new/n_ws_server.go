package new

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"open_im_sdk/pkg/utils"
	"sync"
	"time"
)

type LongConnServer interface {
	Run() error
}

type Server struct {
	rpcPort        int
	wsMaxConnNum   int
	longConnServer *LongConnServer
	rpcServer      *RpcServer
}
type WsServer struct {
	port              int
	wsMaxConnNum      int
	wsUpGrader        *websocket.Upgrader
	registerChan      chan *Client
	unregisterChan    chan *Client
	clients           *UserMap
	clientPool        sync.Pool
	onlineUserNum     int64
	onlineUserConnNum int64
	compressor        Compressor
	handler           MessageHandler
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
		port:         config.port,
		wsMaxConnNum: config.maxConnNum,
		wsUpGrader: &websocket.Upgrader{
			HandshakeTimeout: config.handshakeTimeout,
			ReadBufferSize:   config.messageMaxMsgLength,
			CheckOrigin:      func(r *http.Request) bool { return true },
		},
		clientPool: sync.Pool{
			New: func() interface{} {
				return new(Client)
			},
		},
	}, nil
}
func (ws *WsServer) Run() error {
	http.HandleFunc("/", ws.wsHandler)                              //Get request from client to handle by wsHandler
	return http.ListenAndServe(":"+utils.IntToString(ws.port), nil) //Start listening

}
func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	context := newContext(w, r)
	if isPass, compression := ws.headerCheck(w, r, operationID); isPass {
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			log.Error(operationID, "upgrade http conn err", err.Error(), query)
			return
		} else {
			newConn := &UserConn{conn, new(sync.Mutex), utils.StringToInt32(query["platformID"][0]), 0, compression, query["sendID"][0], false, query["token"][0], conn.RemoteAddr().String() + "_" + strconv.Itoa(int(utils.GetCurrentTimestampByMill()))}
			userCount++
			ws.addUserConn(query["sendID"][0], utils.StringToInt(query["platformID"][0]), newConn, query["token"][0], newConn.connID, operationID)
			go ws.readMsg(newConn)
		}
	} else {
		log.Error(operationID, "headerCheck failed ")
	}
}

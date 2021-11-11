package gate

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/log"
	"Open_IM/src/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WServer struct {
	wsAddr       string
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsConnToUser sync.Map
	wsUserToConn sync.Map
}

func (ws *WServer) onInit(wsPort int) {
	ip := utils.ServerIP
	ws.wsAddr = ip + ":" + utils.IntToString(wsPort)
	ws.wsMaxConnNum = config.Config.LongConnSvr.WebsocketMaxConnNum
	ws.wsConnToUser = sync.Map{}
	ws.wsUserToConn = sync.Map{}
	ws.wsUpGrader = &websocket.Upgrader{
		HandshakeTimeout: time.Duration(config.Config.LongConnSvr.WebsocketTimeOut) * time.Second,
		ReadBufferSize:   config.Config.LongConnSvr.WebsocketMaxMsgLen,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
}

func (ws *WServer) run() {
	http.HandleFunc("/", ws.wsHandler)         //Get request from client to handle by wsHandler
	err := http.ListenAndServe(ws.wsAddr, nil) //Start listening
	if err != nil {
		log.ErrorByKv("Ws listening err", "", "err", err.Error())
	}
}

func (ws *WServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	if ws.headerCheck(w, r) {
		query := r.URL.Query()
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			log.ErrorByKv("upgrade http conn err", "", "err", err)
			return
		} else {
			//Connection mapping relationship,
			//userID+" "+platformID->conn
			SendID := query["sendID"][0] + " " + utils.PlatformIDToName(int32(utils.StringToInt64(query["platformID"][0])))
			ws.addUserConn(SendID, conn)
			go ws.readMsg(conn)
		}
	}
}

func (ws *WServer) readMsg(conn *websocket.Conn) {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.ErrorByKv("WS ReadMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn), "error", err)
			ws.delUserConn(conn)
			return
		} else {
			log.ErrorByKv("test", "", "msgType", msgType, "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn))
		}
		ws.msgParse(conn, msg)
		//ws.writeMsg(conn, 1, chat)
	}

}

func (ws *WServer) writeMsg(conn *websocket.Conn, a int, msg []byte) error {
	rwLock.Lock()
	defer rwLock.Unlock()
	return conn.WriteMessage(a, msg)

}
func (ws *WServer) addUserConn(uid string, conn *websocket.Conn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	if v, ok := ws.wsUserToConn.Load(uid); ok {
		oldConn := v.(*websocket.Conn)
		err := oldConn.Close()
		ws.wsConnToUser.Delete(oldConn)
		if err != nil {
			log.ErrorByKv("close err", "", "uid", uid, "conn", conn)
		}
	} else {
		log.InfoByKv("this user is first login", "", "uid", uid)
	}

	ws.wsConnToUser.Store(conn, uid)
	ws.wsUserToConn.Store(uid, conn)
	log.WarnByKv("WS Add operation", "", "wsUser added", ws.wsUserToConn, "uid", uid, "online_num", ws.onlineNum())

}

func (ws *WServer) delUserConn(conn *websocket.Conn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	var uidPlatform string
	if v, ok := ws.wsConnToUser.Load(conn); ok {
		uid := v.(string)
		uidPlatform = uid
		if _, ok := ws.wsUserToConn.Load(uid); ok {
			ws.wsUserToConn.Delete(uid)
			log.WarnByKv("WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "uid", uid, "online_num", ws.onlineNum())
		} else {
			log.WarnByKv("uid not exist", "", "wsUser deleted", ws.wsUserToConn, "uid", uid, "online_num", ws.onlineNum())
		}
		ws.wsConnToUser.Delete(conn)
	}
	err := conn.Close()
	if err != nil {
		log.ErrorByKv("close err", "", "uid", uidPlatform, "conn", conn)
	}

}

func (ws *WServer) getUserConn(uid string) *websocket.Conn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if v, ok := ws.wsUserToConn.Load(uid); ok {
		conn := v.(*websocket.Conn)
		return conn
	}
	return nil
}
func (ws *WServer) getUserUid(conn *websocket.Conn) string {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if v, ok := ws.wsConnToUser.Load(conn); ok {
		uid := v.(string)
		return uid
	}
	return ""
}
func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request) bool {
	status := http.StatusUnauthorized
	query := r.URL.Query()
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		if !utils.VerifyToken(query["token"][0], query["sendID"][0]) {
			log.ErrorByKv("Token verify failed", "", "query", query)
			w.Header().Set("Sec-Websocket-Version", "13")
			http.Error(w, http.StatusText(status), status)
			return false
		} else {
			log.InfoByKv("Connection Authentication Success", "", "token", query["token"][0], "userID", query["sendID"][0])
			return true
		}
	} else {
		log.ErrorByKv("Args err", "", "query", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		http.Error(w, http.StatusText(status), status)
		return false
	}

}

func (ws *WServer) onlineNum() int {
	var count int
	ws.wsUserToConn.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

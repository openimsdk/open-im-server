package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"bytes"
	"encoding/gob"
	"github.com/garyburd/redigo/redis"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type UserConn struct {
	*websocket.Conn
	w            *sync.Mutex
	PushedMaxSeq uint32
}
type WServer struct {
	wsAddr       string
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsConnToUser map[*UserConn]map[int]string
	wsUserToConn map[string]map[int]*UserConn
}

func (ws *WServer) onInit(wsPort int) {
	ws.wsAddr = ":" + utils.IntToString(wsPort)
	ws.wsMaxConnNum = config.Config.LongConnSvr.WebsocketMaxConnNum
	ws.wsConnToUser = make(map[*UserConn]map[int]string)
	ws.wsUserToConn = make(map[string]map[int]*UserConn)
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
		panic("Ws listening err:" + err.Error())
	}
}

func (ws *WServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	if ws.headerCheck(w, r) {
		query := r.URL.Query()
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			log.Error("", "upgrade http conn err", err, query)
			return
		} else {
			//Connection mapping relationship,
			//userID+" "+platformID->conn
			//Initialize a lock for each user
			newConn := &UserConn{conn, new(sync.Mutex)}
			userCount++
			ws.addUserConn(query["sendID"][0], utils.StringToInt(query["platformID"][0]), newConn, query["token"][0])
			go ws.readMsg(newConn)
		}
	}
}

func (ws *WServer) readMsg(conn *UserConn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if messageType == websocket.PingMessage {
			log.NewInfo("", "this is a  pingMessage")
		}
		if err != nil {
			uid, platform := ws.getUserUid(conn)
			log.Error("", "WS ReadMsg error", "userIP", conn.RemoteAddr().String(), "userUid", uid, "platform", platform, "error", err.Error())
			userCount--
			ws.delUserConn(conn)
			return
		} else {
			//log.ErrorByKv("test", "", "msgType", msgType, "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn))
		}
		ws.msgParse(conn, msg)
		//ws.writeMsg(conn, 1, chat)
	}

}
func (ws *WServer) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	return conn.WriteMessage(a, msg)

}
func (ws *WServer) MultiTerminalLoginChecker(uid string, platformID int, newConn *UserConn, token string, operationID string) {
	switch config.Config.MultiLoginPolicy {
	case constant.AllLoginButSameTermKick:
		if oldConnMap, ok := ws.wsUserToConn[uid]; ok { // user->map[platform->conn]
			if oldConn, ok := oldConnMap[platformID]; ok {
				log.NewDebug(operationID, uid, platformID, "kick old conn")
				ws.sendKickMsg(oldConn, newConn)
				m, err := db.DB.GetTokenMapByUidPid(uid, constant.PlatformIDToName(platformID))
				if err != nil && err != redis.ErrNil {
					log.NewError(operationID, "get token from redis err", err.Error())
					return
				}
				if m == nil {
					log.NewError(operationID, "get token from redis err", "m is nil")
					return
				}
				for k, _ := range m {
					if k != token {
						m[k] = constant.KickedToken
					}
				}
				log.NewDebug(operationID, "get map is ", m)
				err = db.DB.SetTokenMapByUidPid(uid, platformID, m)
				if err != nil {
					log.NewError(operationID, "SetTokenMapByUidPid err", err.Error())
					return
				}
				err = oldConn.Close()
				delete(oldConnMap, platformID)
				ws.wsUserToConn[uid] = oldConnMap
				if len(oldConnMap) == 0 {
					delete(ws.wsUserToConn, uid)
				}
				delete(ws.wsConnToUser, oldConn)
				if err != nil {
					log.NewError(operationID, "conn close err", err.Error(), uid, platformID)
				}

			} else {
				log.NewWarn(operationID, "abnormal uid-conn  ", uid, platformID, oldConnMap[platformID])
			}

		} else {
			log.NewDebug(operationID, "no other conn", ws.wsUserToConn, uid, platformID)
		}

	case constant.SingleTerminalLogin:
	case constant.WebAndOther:
	}
}
func (ws *WServer) sendKickMsg(oldConn, newConn *UserConn) {
	mReply := Resp{
		ReqIdentifier: constant.WSKickOnlineMsg,
		ErrCode:       constant.ErrTokenInvalid.ErrCode,
		ErrMsg:        constant.ErrTokenInvalid.ErrMsg,
	}
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(mReply)
	if err != nil {
		log.NewError(mReply.OperationID, mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg, "Encode Msg error", oldConn.RemoteAddr().String(), newConn.RemoteAddr().String(), err.Error())
		return
	}
	err = ws.writeMsg(oldConn, websocket.BinaryMessage, b.Bytes())
	if err != nil {
		log.NewError(mReply.OperationID, mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg, "WS WriteMsg error", oldConn.RemoteAddr().String(), newConn.RemoteAddr().String(), err.Error())
	}
}
func (ws *WServer) addUserConn(uid string, platformID int, conn *UserConn, token string) {
	rwLock.Lock()
	defer rwLock.Unlock()
	operationID := utils.OperationIDGenerator()
	callbackResp := callbackUserOnline(operationID, uid, platformID, token)
	if callbackResp.ErrCode != 0 {
		log.NewError(operationID, utils.GetSelfFuncName(), "callbackUserOnline resp:", callbackResp)
	}
	ws.MultiTerminalLoginChecker(uid, platformID, conn, token, operationID)
	if oldConnMap, ok := ws.wsUserToConn[uid]; ok {
		oldConnMap[platformID] = conn
		ws.wsUserToConn[uid] = oldConnMap
		log.Debug(operationID, "user not first come in, add conn ", uid, platformID, conn, oldConnMap)
	} else {
		i := make(map[int]*UserConn)
		i[platformID] = conn
		ws.wsUserToConn[uid] = i
		log.Debug(operationID, "user first come in, new user, conn", uid, platformID, conn, ws.wsUserToConn[uid])
	}
	if oldStringMap, ok := ws.wsConnToUser[conn]; ok {
		oldStringMap[platformID] = uid
		ws.wsConnToUser[conn] = oldStringMap
	} else {
		i := make(map[int]string)
		i[platformID] = uid
		ws.wsConnToUser[conn] = i
	}
	count := 0
	for _, v := range ws.wsUserToConn {
		count = count + len(v)
	}
	log.Debug(operationID, "WS Add operation", "", "wsUser added", ws.wsUserToConn, "connection_uid", uid, "connection_platform", constant.PlatformIDToName(platformID), "online_user_num", len(ws.wsUserToConn), "online_conn_num", count)
}

func (ws *WServer) delUserConn(conn *UserConn) {
	rwLock.Lock()
	defer rwLock.Unlock()
	operationID := utils.OperationIDGenerator()
	var uid string
	var platform int
	if oldStringMap, ok := ws.wsConnToUser[conn]; ok {
		for k, v := range oldStringMap {
			platform = k
			uid = v
		}
		if oldConnMap, ok := ws.wsUserToConn[uid]; ok {
			delete(oldConnMap, platform)
			ws.wsUserToConn[uid] = oldConnMap
			if len(oldConnMap) == 0 {
				delete(ws.wsUserToConn, uid)
			}
			count := 0
			for _, v := range ws.wsUserToConn {
				count = count + len(v)
			}
			log.Debug(operationID, "WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "disconnection_uid", uid, "disconnection_platform", platform, "online_user_num", len(ws.wsUserToConn), "online_conn_num", count)
		} else {
			log.Debug(operationID, "WS delete operation", "", "wsUser deleted", ws.wsUserToConn, "disconnection_uid", uid, "disconnection_platform", platform, "online_user_num", len(ws.wsUserToConn))
		}
		delete(ws.wsConnToUser, conn)

	}
	err := conn.Close()
	if err != nil {
		log.Error(operationID, " close err", "", "uid", uid, "platform", platform)
	}
	callbackResp := callbackUserOffline(operationID, uid, platform)
	if callbackResp.ErrCode != 0 {
		log.NewError(operationID, utils.GetSelfFuncName(), "callbackUserOffline failed", callbackResp)
	}
}

func (ws *WServer) getUserConn(uid string, platform int) *UserConn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if connMap, ok := ws.wsUserToConn[uid]; ok {
		if conn, flag := connMap[platform]; flag {
			return conn
		}
	}
	return nil
}
func (ws *WServer) getUserAllCons(uid string) map[int]*UserConn {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if connMap, ok := ws.wsUserToConn[uid]; ok {
		return connMap
	}
	return nil
}

func (ws *WServer) getUserUid(conn *UserConn) (uid string, platform int) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	if stringMap, ok := ws.wsConnToUser[conn]; ok {
		for k, v := range stringMap {
			platform = k
			uid = v
		}
		return uid, platform
	}
	return "", 0
}
func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request) bool {
	status := http.StatusUnauthorized
	query := r.URL.Query()
	operationID := ""
	if len(query["operationID"]) != 0 {
		operationID = query["operationID"][0]
	}
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		if ok, err, msg := token_verify.WsVerifyToken(query["token"][0], query["sendID"][0], query["platformID"][0], operationID); !ok {
			//	e := err.(*constant.ErrInfo)
			log.Error(operationID, "Token verify failed ", "query ", query, msg, err.Error())
			w.Header().Set("Sec-Websocket-Version", "13")
			w.Header().Set("ws_err_msg", err.Error())
			http.Error(w, err.Error(), status)
			return false
		} else {
			log.Info(operationID, "Connection Authentication Success", "", "token", query["token"][0], "userID", query["sendID"][0])
			return true
		}
	} else {
		log.Error(operationID, "Args err", "query", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		w.Header().Set("ws_err_msg", "args err, need token, sendID, platformID")
		http.Error(w, http.StatusText(status), status)
		return false
	}
}

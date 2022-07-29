package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"strings"

	//"gopkg.in/errgo.v2/errors"
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
	query := r.URL.Query()
	operationID := ""
	if len(query["operationID"]) != 0 {
		operationID = query["operationID"][0]
	} else {
		operationID = utils.OperationIDGenerator()
	}
	log.Debug(operationID, utils.GetSelfFuncName(), " args: ", query)
	if ws.headerCheck(w, r, operationID) {
		conn, err := ws.wsUpGrader.Upgrade(w, r, nil) //Conn is obtained through the upgraded escalator
		if err != nil {
			log.Error(operationID, "upgrade http conn err", err.Error(), query)
			return
		} else {
			newConn := &UserConn{conn, new(sync.Mutex), 0}
			userCount++
			ws.addUserConn(query["sendID"][0], utils.StringToInt(query["platformID"][0]), newConn, query["token"][0], operationID)
			go ws.readMsg(newConn)
		}
	} else {
		log.Error(operationID, "headerCheck failed ")
	}
}

func (ws *WServer) readMsg(conn *UserConn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if messageType == websocket.PingMessage {
			log.NewInfo("", "this is a  pingMessage")
		}
		if err != nil {
			log.Error("", "WS ReadMsg error ", " userIP", conn.RemoteAddr().String(), "userUid", "platform", "error", err.Error())
			userCount--
			ws.delUserConn(conn)
			return
		}
		ws.msgParse(conn, msg)
	}
}

func (ws *WServer) SetWriteTimeout(conn *UserConn, timeout int) {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (ws *WServer) writeMsg(conn *UserConn, a int, msg []byte) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(60) * time.Second))
	return conn.WriteMessage(a, msg)
}

func (ws *WServer) SetWriteTimeoutWriteMsg(conn *UserConn, a int, msg []byte, timeout int) error {
	conn.w.Lock()
	defer conn.w.Unlock()
	conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
	return conn.WriteMessage(a, msg)
}

func (ws *WServer) MultiTerminalLoginRemoteChecker(userID string, platformID int32, token string, operationID string) {
	grpcCons := getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImRelayName)
	for _, v := range grpcCons {
		if v.Target() == rpcSvr.target {
			log.Debug(operationID, "Filter out this node ", rpcSvr.target)
			continue
		}
		client := pbRelay.NewRelayClient(v)
		req := &pbRelay.MultiTerminalLoginCheckReq{OperationID: operationID, PlatformID: platformID, UserID: userID, Token: token}
		log.NewInfo(operationID, "MultiTerminalLoginCheckReq ", client, req.String())
		resp, err := client.MultiTerminalLoginCheck(context.Background(), req)
		if err != nil {
			log.Error(operationID, "MultiTerminalLoginCheck failed ", err.Error())
		}
		if resp.ErrCode != 0 {
			log.Error(operationID, "MultiTerminalLoginCheck errCode, errMsg: ", resp.ErrCode, resp.ErrMsg)
		}
	}
}

func (ws *WServer) MultiTerminalLoginCheckerWithLock(uid string, platformID int, token string, operationID string) {
	rwLock.Lock()
	defer rwLock.Unlock()
	switch config.Config.MultiLoginPolicy {
	case constant.AllLoginButSameTermKick:
		if oldConnMap, ok := ws.wsUserToConn[uid]; ok { // user->map[platform->conn]
			if oldConn, ok := oldConnMap[platformID]; ok {
				log.NewDebug(operationID, uid, platformID, "kick old conn")
				m, err := db.DB.GetTokenMapByUidPid(uid, constant.PlatformIDToName(platformID))
				if err != nil && err != go_redis.Nil {
					log.NewError(operationID, "get token from redis err", err.Error(), uid, constant.PlatformIDToName(platformID))
					return
				}
				if m == nil {
					log.NewError(operationID, "get token from redis err", "m is nil", uid, constant.PlatformIDToName(platformID))
					return
				}
				log.NewDebug(operationID, "get token map is ", m, uid, constant.PlatformIDToName(platformID))

				for k, _ := range m {
					if k != token {
						m[k] = constant.KickedToken
					}
				}
				log.NewDebug(operationID, "set token map is ", m, uid, constant.PlatformIDToName(platformID))
				err = db.DB.SetTokenMapByUidPid(uid, platformID, m)
				if err != nil {
					log.NewError(operationID, "SetTokenMapByUidPid err", err.Error(), uid, platformID, m)
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

func (ws *WServer) MultiTerminalLoginChecker(uid string, platformID int, newConn *UserConn, token string, operationID string) {
	switch config.Config.MultiLoginPolicy {
	case constant.AllLoginButSameTermKick:
		if oldConnMap, ok := ws.wsUserToConn[uid]; ok { // user->map[platform->conn]
			if oldConn, ok := oldConnMap[platformID]; ok {
				log.NewDebug(operationID, uid, platformID, "kick old conn")
				//	ws.sendKickMsg(oldConn, newConn)
				m, err := db.DB.GetTokenMapByUidPid(uid, constant.PlatformIDToName(platformID))
				if err != nil && err != go_redis.Nil {
					log.NewError(operationID, "get token from redis err", err.Error(), uid, constant.PlatformIDToName(platformID))
					return
				}
				if m == nil {
					log.NewError(operationID, "get token from redis err", "m is nil", uid, constant.PlatformIDToName(platformID))
					return
				}
				log.NewDebug(operationID, "get token map is ", m, uid, constant.PlatformIDToName(platformID))

				for k, _ := range m {
					if k != token {
						m[k] = constant.KickedToken
					}
				}
				log.NewDebug(operationID, "set token map is ", m, uid, constant.PlatformIDToName(platformID))
				err = db.DB.SetTokenMapByUidPid(uid, platformID, m)
				if err != nil {
					log.NewError(operationID, "SetTokenMapByUidPid err", err.Error(), uid, platformID, m)
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
		log.NewError(mReply.OperationID, mReply.ReqIdentifier, mReply.ErrCode, mReply.ErrMsg, "sendKickMsg WS WriteMsg error", oldConn.RemoteAddr().String(), newConn.RemoteAddr().String(), err.Error())
	}
}

func (ws *WServer) addUserConn(uid string, platformID int, conn *UserConn, token string, operationID string) {
	rwLock.Lock()
	defer rwLock.Unlock()
	log.Info(operationID, utils.GetSelfFuncName(), " args: ", uid, platformID, conn, token)
	callbackResp := callbackUserOnline(operationID, uid, platformID, token)
	if callbackResp.ErrCode != 0 {
		log.NewError(operationID, utils.GetSelfFuncName(), "callbackUserOnline resp:", callbackResp)
	}
	//go ws.MultiTerminalLoginRemoteChecker(uid, int32(platformID), token, operationID)
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

//func (ws *WServer) getUserUid(conn *UserConn) (uid string, platform int) {
//	rwLock.RLock()
//	defer rwLock.RUnlock()
//
//	if stringMap, ok := ws.wsConnToUser[conn]; ok {
//		for k, v := range stringMap {
//			platform = k
//			uid = v
//		}
//		return uid, platform
//	}
//	return "", 0
//}
func (ws *WServer) headerCheck(w http.ResponseWriter, r *http.Request, operationID string) bool {
	status := http.StatusUnauthorized
	query := r.URL.Query()
	if len(query["token"]) != 0 && len(query["sendID"]) != 0 && len(query["platformID"]) != 0 {
		if ok, err, msg := token_verify.WsVerifyToken(query["token"][0], query["sendID"][0], query["platformID"][0], operationID); !ok {
			if errors.Is(err, constant.ErrTokenExpired) {
				status = int(constant.ErrTokenExpired.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenInvalid) {
				status = int(constant.ErrTokenInvalid.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenMalformed) {
				status = int(constant.ErrTokenMalformed.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenNotValidYet) {
				status = int(constant.ErrTokenNotValidYet.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenUnknown) {
				status = int(constant.ErrTokenUnknown.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenKicked) {
				status = int(constant.ErrTokenKicked.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenDifferentPlatformID) {
				status = int(constant.ErrTokenDifferentPlatformID.ErrCode)
			}
			if errors.Is(err, constant.ErrTokenDifferentUserID) {
				status = int(constant.ErrTokenDifferentUserID.ErrCode)
			}
			//switch errors.Cause(err) {
			//case constant.ErrTokenExpired:
			//	status = int(constant.ErrTokenExpired.ErrCode)
			//case constant.ErrTokenInvalid:
			//	status = int(constant.ErrTokenInvalid.ErrCode)
			//case constant.ErrTokenMalformed:
			//	status = int(constant.ErrTokenMalformed.ErrCode)
			//case constant.ErrTokenNotValidYet:
			//	status = int(constant.ErrTokenNotValidYet.ErrCode)
			//case constant.ErrTokenUnknown:
			//	status = int(constant.ErrTokenUnknown.ErrCode)
			//case constant.ErrTokenKicked:
			//	status = int(constant.ErrTokenKicked.ErrCode)
			//case constant.ErrTokenDifferentPlatformID:
			//	status = int(constant.ErrTokenDifferentPlatformID.ErrCode)
			//case constant.ErrTokenDifferentUserID:
			//	status = int(constant.ErrTokenDifferentUserID.ErrCode)
			//}

			log.Error(operationID, "Token verify failed ", "query ", query, msg, err.Error(), "status: ", status)
			w.Header().Set("Sec-Websocket-Version", "13")
			w.Header().Set("ws_err_msg", err.Error())
			http.Error(w, err.Error(), status)
			return false
		} else {
			log.Info(operationID, "Connection Authentication Success", "", "token ", query["token"][0], "userID ", query["sendID"][0], "platformID ", query["platformID"][0])
			return true
		}
	} else {
		status = int(constant.ErrArgs.ErrCode)
		log.Error(operationID, "Args err ", "query ", query)
		w.Header().Set("Sec-Websocket-Version", "13")
		errMsg := "args err, need token, sendID, platformID"
		w.Header().Set("ws_err_msg", errMsg)
		http.Error(w, errMsg, status)
		return false
	}
}

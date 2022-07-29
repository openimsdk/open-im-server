package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbRelay "Open_IM/pkg/proto/relay"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"net"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type RPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	platformList    []int
	pushTerminal    []int
	target          string
}

func (r *RPCServer) onInit(rpcPort int) {
	r.rpcPort = rpcPort
	r.rpcRegisterName = config.Config.RpcRegisterName.OpenImRelayName
	r.etcdSchema = config.Config.Etcd.EtcdSchema
	r.etcdAddr = config.Config.Etcd.EtcdAddr
	r.platformList = genPlatformArray()
	r.pushTerminal = []int{constant.IOSPlatformID, constant.AndroidPlatformID}
}
func (r *RPCServer) run() {
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(r.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + r.rpcRegisterName)
	}
	defer listener.Close()
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	pbRelay.RegisterRelayServer(srv, r)

	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	err = getcdv3.RegisterEtcd4Unique(r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName, 10)
	if err != nil {
		log.Error("", "register push message rpc to etcd err", "", "err", err.Error(), r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName)
	}
	r.target = getcdv3.GetTarget(r.etcdSchema, rpcRegisterIP, r.rpcPort, r.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.Error("", "push message rpc listening err", "", "err", err.Error())
		return
	}
}
func (r *RPCServer) OnlinePushMsg(_ context.Context, in *pbRelay.OnlinePushMsgReq) (*pbRelay.OnlinePushMsgResp, error) {
	log.NewInfo(in.OperationID, "PushMsgToUser is arriving", in.String())
	var resp []*pbRelay.SingleMsgToUserPlatform
	msgBytes, _ := proto.Marshal(in.MsgData)
	mReply := Resp{
		ReqIdentifier: constant.WSPushMsg,
		OperationID:   in.OperationID,
		Data:          msgBytes,
	}
	var replyBytes bytes.Buffer
	enc := gob.NewEncoder(&replyBytes)
	err := enc.Encode(mReply)
	if err != nil {
		log.NewError(in.OperationID, "data encode err", err.Error())
	}
	var tag bool
	recvID := in.PushToUserID
	for _, v := range r.platformList {
		if conn := ws.getUserConn(recvID, v); conn != nil {
			tag = true
			resultCode := sendMsgToUser(conn, replyBytes.Bytes(), in, v, recvID)
			temp := &pbRelay.SingleMsgToUserPlatform{
				ResultCode:     resultCode,
				RecvID:         recvID,
				RecvPlatFormID: int32(v),
			}
			resp = append(resp, temp)
		} else {
			temp := &pbRelay.SingleMsgToUserPlatform{
				ResultCode:     -1,
				RecvID:         recvID,
				RecvPlatFormID: int32(v),
			}
			resp = append(resp, temp)
		}
	}
	if !tag {
		log.NewDebug(in.OperationID, "push err ,no matched ws conn not in map", in.String())
	}
	return &pbRelay.OnlinePushMsgResp{
		Resp: resp,
	}, nil
}
func (r *RPCServer) GetUsersOnlineStatus(_ context.Context, req *pbRelay.GetUsersOnlineStatusReq) (*pbRelay.GetUsersOnlineStatusResp, error) {
	log.NewInfo(req.OperationID, "rpc GetUsersOnlineStatus arrived server", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "no permission GetUsersOnlineStatus ", req.OpUserID)
		return &pbRelay.GetUsersOnlineStatusResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	var resp pbRelay.GetUsersOnlineStatusResp
	for _, userID := range req.UserIDList {
		temp := new(pbRelay.GetUsersOnlineStatusResp_SuccessResult)
		temp.UserID = userID
		userConnMap := ws.getUserAllCons(userID)
		for platform, userConn := range userConnMap {
			if userConn != nil {
				ps := new(pbRelay.GetUsersOnlineStatusResp_SuccessDetail)
				ps.Platform = constant.PlatformIDToName(platform)
				ps.Status = constant.OnlineStatus
				temp.Status = constant.OnlineStatus
				temp.DetailPlatformStatus = append(temp.DetailPlatformStatus, ps)
			}
		}

		if temp.Status == constant.OnlineStatus {
			resp.SuccessResult = append(resp.SuccessResult, temp)
		}
	}
	log.NewInfo(req.OperationID, "GetUsersOnlineStatus rpc return ", resp.String())
	return &resp, nil
}

func (r *RPCServer) SuperGroupOnlineBatchPushOneMsg(_ context.Context, req *pbRelay.OnlineBatchPushOneMsgReq) (*pbRelay.OnlineBatchPushOneMsgResp, error) {
	log.NewInfo(req.OperationID, "BatchPushMsgToUser is arriving", req.String())
	var singleUserResult []*pbRelay.SingelMsgToUserResultList
	//r.GetBatchMsgForPush(req.OperationID,req.MsgData,req.PushToUserIDList,)
	msgBytes, _ := proto.Marshal(req.MsgData)
	mReply := Resp{
		ReqIdentifier: constant.WSPushMsg,
		OperationID:   req.OperationID,
		Data:          msgBytes,
	}
	var replyBytes bytes.Buffer
	enc := gob.NewEncoder(&replyBytes)
	err := enc.Encode(mReply)
	if err != nil {
		log.NewError(req.OperationID, "data encode err", err.Error())
	}
	for _, v := range req.PushToUserIDList {
		var resp []*pbRelay.SingleMsgToUserPlatform
		tempT := &pbRelay.SingelMsgToUserResultList{
			UserID: v,
		}
		userConnMap := ws.getUserAllCons(v)
		for platform, userConn := range userConnMap {
			if userConn != nil {
				resultCode := sendMsgBatchToUser(userConn, replyBytes.Bytes(), req, platform, v)
				if resultCode == 0 && utils.IsContainInt(platform, r.pushTerminal) {
					tempT.OnlinePush = true
					log.Info(req.OperationID, "PushSuperMsgToUser is success By Ws", "args", req.String(), "recvPlatForm", constant.PlatformIDToName(platform), "recvID", v)
					temp := &pbRelay.SingleMsgToUserPlatform{
						ResultCode:     resultCode,
						RecvID:         v,
						RecvPlatFormID: int32(platform),
					}
					resp = append(resp, temp)
				}

			}
		}
		tempT.Resp = resp
		singleUserResult = append(singleUserResult, tempT)

	}

	return &pbRelay.OnlineBatchPushOneMsgResp{
		SinglePushResult: singleUserResult,
	}, nil
}
func (r *RPCServer) OnlineBatchPushOneMsg(_ context.Context, req *pbRelay.OnlineBatchPushOneMsgReq) (*pbRelay.OnlineBatchPushOneMsgResp, error) {
	log.NewInfo(req.OperationID, "BatchPushMsgToUser is arriving", req.String())
	var singleUserResult []*pbRelay.SingelMsgToUserResultList

	for _, v := range req.PushToUserIDList {
		var resp []*pbRelay.SingleMsgToUserPlatform
		tempT := &pbRelay.SingelMsgToUserResultList{
			UserID: v,
		}
		userConnMap := ws.getUserAllCons(v)
		var platformList []int
		for k, _ := range userConnMap {
			platformList = append(platformList, k)
		}
		log.Debug(req.OperationID, "GetSingleUserMsgForPushPlatforms begin", req.MsgData.Seq, v, platformList, req.MsgData.String())
		needPushMapList := r.GetSingleUserMsgForPushPlatforms(req.OperationID, req.MsgData, v, platformList)
		log.Debug(req.OperationID, "GetSingleUserMsgForPushPlatforms end", req.MsgData.Seq, v, platformList, len(needPushMapList))
		for platform, list := range needPushMapList {
			if list != nil {
				log.Debug(req.OperationID, "needPushMapList ", "userID: ", v, "platform: ", platform, "push msg num:")
				//for _, v := range list {
				//	log.Debug(req.OperationID, "req.MsgData.MsgDataList begin", "len: ", len(req.MsgData.MsgDataList), v.String())
				//	req.MsgData.MsgDataList = append(req.MsgData.MsgDataList, v)
				//	log.Debug(req.OperationID, "req.MsgData.MsgDataList end", "len: ", len(req.MsgData.MsgDataList))
				//}
				msgBytes, err := proto.Marshal(list)
				if err != nil {
					log.Error(req.OperationID, "proto marshal err", err.Error())
					continue
				}
				req.MsgData.MsgDataList = msgBytes
				//req.MsgData.MsgDataList = append(req.MsgData.MsgDataList, v)
				log.Debug(req.OperationID, "r.encodeWsData  no string")
				//log.Debug(req.OperationID, "r.encodeWsData  data0 list ", req.MsgData.MsgDataList[0].String())

				log.Debug(req.OperationID, "r.encodeWsData  ", req.MsgData.String())
				replyBytes, err := r.encodeWsData(req.MsgData, req.OperationID)
				if err != nil {
					log.Error(req.OperationID, "encodeWsData failed ", req.MsgData.String())
					continue
				}
				log.Debug(req.OperationID, "encodeWsData", "len: ", replyBytes.Len())
				resultCode := sendMsgBatchToUser(userConnMap[platform], replyBytes.Bytes(), req, platform, v)
				if resultCode == 0 && utils.IsContainInt(platform, r.pushTerminal) {
					tempT.OnlinePush = true
					log.Info(req.OperationID, "PushSuperMsgToUser is success By Ws", "args", req.String(), "recv PlatForm", constant.PlatformIDToName(platform), "recvID", v)
					temp := &pbRelay.SingleMsgToUserPlatform{
						ResultCode:     resultCode,
						RecvID:         v,
						RecvPlatFormID: int32(platform),
					}
					resp = append(resp, temp)
				}
			} else {
				if utils.IsContainInt(platform, r.pushTerminal) {
					tempT.OnlinePush = true
					temp := &pbRelay.SingleMsgToUserPlatform{
						ResultCode:     0,
						RecvID:         v,
						RecvPlatFormID: int32(platform),
					}
					resp = append(resp, temp)
				}
			}
		}
		tempT.Resp = resp
		singleUserResult = append(singleUserResult, tempT)
	}
	return &pbRelay.OnlineBatchPushOneMsgResp{
		SinglePushResult: singleUserResult,
	}, nil
}
func (r *RPCServer) encodeWsData(wsData *sdk_ws.MsgData, operationID string) (bytes.Buffer, error) {
	log.Debug(operationID, "encodeWsData begin", wsData.String())
	msgBytes, err := proto.Marshal(wsData)
	if err != nil {
		log.NewError(operationID, "Marshal", err.Error())
		return bytes.Buffer{}, utils.Wrap(err, "")
	}
	log.Debug(operationID, "encodeWsData begin", wsData.String())
	mReply := Resp{
		ReqIdentifier: constant.WSPushMsg,
		OperationID:   operationID,
		Data:          msgBytes,
	}
	var replyBytes bytes.Buffer
	enc := gob.NewEncoder(&replyBytes)
	err = enc.Encode(mReply)
	if err != nil {
		log.NewError(operationID, "data encode err", err.Error())
		return bytes.Buffer{}, utils.Wrap(err, "")
	}
	return replyBytes, nil
}

func (r *RPCServer) KickUserOffline(_ context.Context, req *pbRelay.KickUserOfflineReq) (*pbRelay.KickUserOfflineResp, error) {
	log.NewInfo(req.OperationID, "KickUserOffline is arriving", req.String())
	for _, v := range req.KickUserIDList {
		oldConnMap := ws.getUserAllCons(v)
		if conn, ok := oldConnMap[int(req.PlatformID)]; ok { // user->map[platform->conn]
			log.NewWarn(req.OperationID, "send kick msg, close connection ", req.PlatformID, v)
			ws.sendKickMsg(conn, &UserConn{})
			conn.Close()
		}
		log.NewWarn(req.OperationID, "SetTokenKicked ", v, req.PlatformID, req.OperationID)
		SetTokenKicked(v, int(req.PlatformID), req.OperationID)
	}
	return &pbRelay.KickUserOfflineResp{}, nil
}

func (r *RPCServer) MultiTerminalLoginCheck(ctx context.Context, req *pbRelay.MultiTerminalLoginCheckReq) (*pbRelay.MultiTerminalLoginCheckResp, error) {
	ws.MultiTerminalLoginCheckerWithLock(req.UserID, int(req.PlatformID), req.Token, req.OperationID)
	return &pbRelay.MultiTerminalLoginCheckResp{}, nil
}

func sendMsgToUser(conn *UserConn, bMsg []byte, in *pbRelay.OnlinePushMsgReq, RecvPlatForm int, RecvID string) (ResultCode int64) {
	err := ws.writeMsg(conn, websocket.BinaryMessage, bMsg)
	if err != nil {
		log.NewError(in.OperationID, "PushMsgToUser is failed By Ws", "Addr", conn.RemoteAddr().String(),
			"error", err, "senderPlatform", constant.PlatformIDToName(int(in.MsgData.SenderPlatformID)), "recvPlatform", RecvPlatForm, "args", in.String(), "recvID", RecvID)
		ResultCode = -2
		return ResultCode
	} else {
		log.NewDebug(in.OperationID, "PushMsgToUser is success By Ws", "args", in.String(), "recvPlatForm", RecvPlatForm, "recvID", RecvID)
		ResultCode = 0
		return ResultCode
	}

}
func sendMsgBatchToUser(conn *UserConn, bMsg []byte, in *pbRelay.OnlineBatchPushOneMsgReq, RecvPlatForm int, RecvID string) (ResultCode int64) {
	err := ws.writeMsg(conn, websocket.BinaryMessage, bMsg)
	if err != nil {
		log.NewError(in.OperationID, "PushMsgToUser is failed By Ws", "Addr", conn.RemoteAddr().String(),
			"error", err, "senderPlatform", constant.PlatformIDToName(int(in.MsgData.SenderPlatformID)), "recv Platform", RecvPlatForm, "args", in.String(), "recvID", RecvID)
		ResultCode = -2
		return ResultCode
	} else {
		log.NewDebug(in.OperationID, "PushMsgToUser is success By Ws", "args", in.String(), "recv PlatForm", RecvPlatForm, "recvID", RecvID)
		ResultCode = 0
		return ResultCode
	}

}
func genPlatformArray() (array []int) {
	for i := 1; i <= constant.LinuxPlatformID; i++ {
		array = append(array, i)
	}
	return array
}

package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/msg"
	pbRtc "Open_IM/pkg/proto/rtc"
	sdk_ws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"runtime"
	"strings"
)

func (ws *WServer) msgParse(conn *UserConn, binaryMsg []byte) {
	b := bytes.NewBuffer(binaryMsg)
	m := Req{}
	dec := gob.NewDecoder(b)
	err := dec.Decode(&m)
	if err != nil {
		log.NewError("", "ws Decode  err", err.Error())
		err = conn.Close()
		if err != nil {
			log.NewError("", "ws close err", err.Error())
		}
		return
	}
	if err := validate.Struct(m); err != nil {
		log.NewError("", "ws args validate  err", err.Error())
		ws.sendErrMsg(conn, 201, err.Error(), m.ReqIdentifier, m.MsgIncr, m.OperationID)
		return
	}
	log.NewInfo(m.OperationID, "Basic Info Authentication Success", m.SendID, m.MsgIncr, m.ReqIdentifier)

	switch m.ReqIdentifier {
	case constant.WSGetNewestSeq:
		log.NewInfo(m.OperationID, "getSeqReq ", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.getSeqReq(conn, &m)
	case constant.WSSendMsg:
		log.NewInfo(m.OperationID, "sendMsgReq ", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.sendMsgReq(conn, &m)
	case constant.WSSendSignalMsg:
		log.NewInfo(m.OperationID, "sendSignalMsgReq ", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.sendSignalMsgReq(conn, &m)
	case constant.WSPullMsgBySeqList:
		log.NewInfo(m.OperationID, "pullMsgBySeqListReq ", m.SendID, m.MsgIncr, m.ReqIdentifier)
		ws.pullMsgBySeqListReq(conn, &m)
	case constant.WsLogoutMsg:
		log.NewInfo(m.OperationID, "conn.Close()", m.SendID, m.MsgIncr, m.ReqIdentifier)
	//	conn.Close()
	default:
		log.Error(m.OperationID, "ReqIdentifier failed ", m.SendID, m.MsgIncr, m.ReqIdentifier)
	}
	log.NewInfo(m.OperationID, "goroutine num is ", runtime.NumGoroutine())
}

func (ws *WServer) getSeqReq(conn *UserConn, m *Req) {
	log.NewInfo(m.OperationID, "Ws call success to getNewSeq", m.MsgIncr, m.SendID, m.ReqIdentifier)
	nReply := new(sdk_ws.GetMaxAndMinSeqResp)
	isPass, errCode, errMsg, data := ws.argsValidate(m, constant.WSGetNewestSeq, m.OperationID)
	log.Info(m.OperationID, "argsValidate ", isPass, errCode, errMsg)
	if isPass {
		rpcReq := sdk_ws.GetMaxAndMinSeqReq{}
		rpcReq.GroupIDList = data.(sdk_ws.GetMaxAndMinSeqReq).GroupIDList
		rpcReq.UserID = m.SendID
		rpcReq.OperationID = m.OperationID
		log.Debug(m.OperationID, "Ws call success to getMaxAndMinSeq", m.SendID, m.ReqIdentifier, m.MsgIncr, data.(sdk_ws.GetMaxAndMinSeqReq).GroupIDList)
		grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, rpcReq.OperationID)
		if grpcConn == nil {
			errMsg := rpcReq.OperationID + "getcdv3.GetConn == nil"
			nReply.ErrCode = 500
			nReply.ErrMsg = errMsg
			log.NewError(rpcReq.OperationID, errMsg)
			ws.getSeqResp(conn, m, nReply)
			return
		}
		msgClient := pbChat.NewMsgClient(grpcConn)
		rpcReply, err := msgClient.GetMaxAndMinSeq(context.Background(), &rpcReq)
		if err != nil {
			nReply.ErrCode = 500
			nReply.ErrMsg = err.Error()
			log.Error(rpcReq.OperationID, "rpc call failed to GetMaxAndMinSeq ", nReply.String())
			ws.getSeqResp(conn, m, nReply)
		} else {
			log.NewInfo(rpcReq.OperationID, "rpc call success to getSeqReq", rpcReply.String())
			ws.getSeqResp(conn, m, rpcReply)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		log.Error(m.OperationID, "argsValidate failed send resp: ", nReply.String())
		ws.getSeqResp(conn, m, nReply)
	}
}

func (ws *WServer) getSeqResp(conn *UserConn, m *Req, pb *sdk_ws.GetMaxAndMinSeqResp) {

	b, _ := proto.Marshal(pb)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		OperationID:   m.OperationID,
		Data:          b,
	}
	log.Debug(m.OperationID, "getSeqResp come  here req: ", pb.String(), "send resp: ",
		mReply.ReqIdentifier, mReply.MsgIncr, mReply.ErrCode, mReply.ErrMsg)
	ws.sendMsg(conn, mReply)
}

func (ws *WServer) pullMsgBySeqListReq(conn *UserConn, m *Req) {
	log.NewInfo(m.OperationID, "Ws call success to pullMsgBySeqListReq start", m.SendID, m.ReqIdentifier, m.MsgIncr, string(m.Data))
	nReply := new(sdk_ws.PullMessageBySeqListResp)
	isPass, errCode, errMsg, data := ws.argsValidate(m, constant.WSPullMsgBySeqList, m.OperationID)
	if isPass {
		rpcReq := sdk_ws.PullMessageBySeqListReq{}
		rpcReq.SeqList = data.(sdk_ws.PullMessageBySeqListReq).SeqList
		rpcReq.UserID = m.SendID
		rpcReq.OperationID = m.OperationID
		rpcReq.GroupSeqList = data.(sdk_ws.PullMessageBySeqListReq).GroupSeqList
		log.NewInfo(m.OperationID, "Ws call success to pullMsgBySeqListReq middle", m.SendID, m.ReqIdentifier, m.MsgIncr, data.(sdk_ws.PullMessageBySeqListReq).SeqList)
		grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, m.OperationID)
		if grpcConn == nil {
			errMsg := rpcReq.OperationID + "getcdv3.GetConn == nil"
			nReply.ErrCode = 500
			nReply.ErrMsg = errMsg
			log.NewError(rpcReq.OperationID, errMsg)
			ws.pullMsgBySeqListResp(conn, m, nReply)
			return
		}
		msgClient := pbChat.NewMsgClient(grpcConn)
		reply, err := msgClient.PullMessageBySeqList(context.Background(), &rpcReq)
		if err != nil {
			log.NewError(rpcReq.OperationID, "pullMsgBySeqListReq err", err.Error())
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.pullMsgBySeqListResp(conn, m, nReply)
		} else {
			log.NewInfo(rpcReq.OperationID, "rpc call success to pullMsgBySeqListReq", reply.String(), len(reply.List))
			ws.pullMsgBySeqListResp(conn, m, reply)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		ws.pullMsgBySeqListResp(conn, m, nReply)
	}
}
func (ws *WServer) pullMsgBySeqListResp(conn *UserConn, m *Req, pb *sdk_ws.PullMessageBySeqListResp) {
	log.NewInfo(m.OperationID, "pullMsgBySeqListResp come  here ", pb.String())
	c, _ := proto.Marshal(pb)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		OperationID:   m.OperationID,
		Data:          c,
	}
	log.NewInfo(m.OperationID, "pullMsgBySeqListResp all data  is ", mReply.ReqIdentifier, mReply.MsgIncr, mReply.ErrCode, mReply.ErrMsg,
		len(mReply.Data))

	ws.sendMsg(conn, mReply)

}
func (ws *WServer) sendMsgReq(conn *UserConn, m *Req) {
	sendMsgAllCountLock.Lock()
	sendMsgAllCount++
	sendMsgAllCountLock.Unlock()
	log.NewInfo(m.OperationID, "Ws call success to sendMsgReq start", m.MsgIncr, m.ReqIdentifier, m.SendID)

	nReply := new(pbChat.SendMsgResp)
	isPass, errCode, errMsg, pData := ws.argsValidate(m, constant.WSSendMsg, m.OperationID)
	if isPass {
		data := pData.(sdk_ws.MsgData)
		pbData := pbChat.SendMsgReq{
			Token:       m.Token,
			OperationID: m.OperationID,
			MsgData:     &data,
		}
		log.NewInfo(m.OperationID, "Ws call success to sendMsgReq middle", m.ReqIdentifier, m.SendID, m.MsgIncr, data.String())
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, m.OperationID)
		if etcdConn == nil {
			errMsg := m.OperationID + "getcdv3.GetConn == nil"
			nReply.ErrCode = 500
			nReply.ErrMsg = errMsg
			log.NewError(m.OperationID, errMsg)
			ws.sendMsgResp(conn, m, nReply)
			return
		}
		client := pbChat.NewMsgClient(etcdConn)
		reply, err := client.SendMsg(context.Background(), &pbData)
		if err != nil {
			log.NewError(pbData.OperationID, "UserSendMsg err", err.Error())
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.sendMsgResp(conn, m, nReply)
		} else {
			log.NewInfo(pbData.OperationID, "rpc call success to sendMsgReq", reply.String())
			ws.sendMsgResp(conn, m, reply)
		}

	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		ws.sendMsgResp(conn, m, nReply)
	}

}
func (ws *WServer) sendMsgResp(conn *UserConn, m *Req, pb *pbChat.SendMsgResp) {
	var mReplyData sdk_ws.UserSendMsgResp
	mReplyData.ClientMsgID = pb.GetClientMsgID()
	mReplyData.ServerMsgID = pb.GetServerMsgID()
	mReplyData.SendTime = pb.GetSendTime()
	b, _ := proto.Marshal(&mReplyData)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		OperationID:   m.OperationID,
		Data:          b,
	}
	ws.sendMsg(conn, mReply)
}

func (ws *WServer) sendSignalMsgReq(conn *UserConn, m *Req) {
	log.NewInfo(m.OperationID, "Ws call success to sendSignalMsgReq start", m.MsgIncr, m.ReqIdentifier, m.SendID, string(m.Data))
	nReply := new(pbChat.SendMsgResp)
	isPass, errCode, errMsg, pData := ws.argsValidate(m, constant.WSSendSignalMsg, m.OperationID)
	if isPass {
		signalResp := pbRtc.SignalResp{}
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImRealTimeCommName, m.OperationID)
		if etcdConn == nil {
			errMsg := m.OperationID + "getcdv3.GetConn == nil"
			log.NewError(m.OperationID, errMsg)
			ws.sendSignalMsgResp(conn, 204, errMsg, m, &signalResp)
			return
		}
		rtcClient := pbRtc.NewRtcServiceClient(etcdConn)
		req := &pbRtc.SignalMessageAssembleReq{
			SignalReq:   pData.(*pbRtc.SignalReq),
			OperationID: m.OperationID,
		}
		respPb, err := rtcClient.SignalMessageAssemble(context.Background(), req)
		if err != nil {
			log.NewError(m.OperationID, utils.GetSelfFuncName(), "SignalMessageAssemble", err.Error(), config.Config.RpcRegisterName.OpenImRealTimeCommName)
			ws.sendSignalMsgResp(conn, 204, "grpc SignalMessageAssemble failed: "+err.Error(), m, &signalResp)
			return
		}
		signalResp.Payload = respPb.SignalResp.Payload
		msgData := sdk_ws.MsgData{}
		utils.CopyStructFields(&msgData, respPb.MsgData)
		log.NewInfo(m.OperationID, utils.GetSelfFuncName(), respPb.String())
		if respPb.IsPass {
			pbData := pbChat.SendMsgReq{
				Token:       m.Token,
				OperationID: m.OperationID,
				MsgData:     &msgData,
			}
			log.NewInfo(m.OperationID, utils.GetSelfFuncName(), "pbData: ", pbData)
			log.NewInfo(m.OperationID, "Ws call success to sendSignalMsgReq middle", m.ReqIdentifier, m.SendID, m.MsgIncr, msgData)
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, m.OperationID)
			if etcdConn == nil {
				errMsg := m.OperationID + "getcdv3.GetConn == nil"
				log.NewError(m.OperationID, errMsg)
				ws.sendSignalMsgResp(conn, 200, errMsg, m, &signalResp)
				return
			}
			client := pbChat.NewMsgClient(etcdConn)
			reply, err := client.SendMsg(context.Background(), &pbData)
			if err != nil {
				log.NewError(pbData.OperationID, utils.GetSelfFuncName(), "rpc sendMsg err", err.Error())
				nReply.ErrCode = 200
				nReply.ErrMsg = err.Error()
				ws.sendSignalMsgResp(conn, 200, err.Error(), m, &signalResp)
			} else {
				log.NewInfo(pbData.OperationID, "rpc call success to sendMsgReq", reply.String(), signalResp.String(), m)
				ws.sendSignalMsgResp(conn, 0, "", m, &signalResp)
			}
		} else {
			log.NewError(m.OperationID, utils.GetSelfFuncName(), respPb.IsPass, respPb.CommonResp.ErrCode, respPb.CommonResp.ErrMsg)
			ws.sendSignalMsgResp(conn, respPb.CommonResp.ErrCode, respPb.CommonResp.ErrMsg, m, &signalResp)
		}
	} else {
		ws.sendSignalMsgResp(conn, errCode, errMsg, m, nil)
	}

}
func (ws *WServer) sendSignalMsgResp(conn *UserConn, errCode int32, errMsg string, m *Req, pb *pbRtc.SignalResp) {
	// := make(map[string]interface{})
	log.Debug(m.OperationID, "sendSignalMsgResp is", pb.String())
	b, _ := proto.Marshal(pb)
	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       errCode,
		ErrMsg:        errMsg,
		OperationID:   m.OperationID,
		Data:          b,
	}
	ws.sendMsg(conn, mReply)
}
func (ws *WServer) sendMsg(conn *UserConn, mReply interface{}) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(mReply)
	if err != nil {
		//	uid, platform := ws.getUserUid(conn)
		log.NewError(mReply.(Resp).OperationID, mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, "Encode Msg error", conn.RemoteAddr().String(), err.Error())
		return
	}
	err = ws.writeMsg(conn, websocket.BinaryMessage, b.Bytes())
	if err != nil {
		//	uid, platform := ws.getUserUid(conn)
		log.NewError(mReply.(Resp).OperationID, mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, "ws writeMsg error", conn.RemoteAddr().String(), err.Error())
	} else {
		log.Debug(mReply.(Resp).OperationID, mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, "ws write response success")
	}
}
func (ws *WServer) sendErrMsg(conn *UserConn, errCode int32, errMsg string, reqIdentifier int32, msgIncr string, operationID string) {
	mReply := Resp{
		ReqIdentifier: reqIdentifier,
		MsgIncr:       msgIncr,
		ErrCode:       errCode,
		ErrMsg:        errMsg,
		OperationID:   operationID,
	}
	ws.sendMsg(conn, mReply)
}

func SetTokenKicked(userID string, platformID int, operationID string) {
	m, err := db.DB.GetTokenMapByUidPid(userID, constant.PlatformIDToName(platformID))
	if err != nil {
		log.Error(operationID, "GetTokenMapByUidPid failed ", err.Error(), userID, constant.PlatformIDToName(platformID))
		return
	}
	for k, _ := range m {
		m[k] = constant.KickedToken
	}
	err = db.DB.SetTokenMapByUidPid(userID, platformID, m)
	if err != nil {
		log.Error(operationID, "SetTokenMapByUidPid failed ", err.Error(), userID, constant.PlatformIDToName(platformID))
		return
	}
}

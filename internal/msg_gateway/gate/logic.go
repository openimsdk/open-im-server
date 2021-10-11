package gate

import (
	pbChat "Open_IM/pkg/proto/chat"
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	"Open_IM/src/common/log"
	"Open_IM/src/grpc-etcdv3/getcdv3"
	"Open_IM/src/utils"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"strings"
)

func (ws *WServer) msgParse(conn *websocket.Conn, jsonMsg []byte) {
	//ws online debug data
	//{"ReqIdentifier":1001,"Token":"123","SendID":"c4ca4238a0b923820dcc509a6f75849b","Time":"123","OperationID":"123","MsgIncr":0}
	//{"ReqIdentifier":1002,"Token":"123","SendID":"c4ca4238a0b923820dcc509a6f75849b","Time":"123","OperationID":"123","MsgIncr":0,"SeqBegin":1,"SeqEnd":6}
	//{"ReqIdentifier":1003,"Token":"123","SendID":"c4ca4238a0b923820dcc509a6f75849b",
	//"RecvID":"a87ff679a2f3e71d9181a67b7542122c","ClientMsgID":"2343","Time":"147878787","OperationID":
	//"123","MsgIncr":0,"SubMsgType":101,"MsgType":100,"MsgFrom":1,"Content":"sdfsdf"}
	m := Req{}
	if err := json.Unmarshal(jsonMsg, &m); err != nil {
		log.ErrorByKv("ws json Unmarshal err", "", "err", err.Error())
		ws.sendErrMsg(conn, 200, err.Error())
		return
	}
	if err := validate.Struct(m); err != nil {
		log.ErrorByKv("ws args validate  err", "", "err", err.Error())
		ws.sendErrMsg(conn, 201, err.Error())
		return
	}

	if !utils.VerifyToken(m.Token, m.SendID) {
		ws.sendErrMsg(conn, 202, "token validate err")
		return
	}
	log.InfoByKv("Basic Info Authentication Success", m.OperationID, "reqIdentifier", m.ReqIdentifier, "sendID", m.SendID)

	switch m.ReqIdentifier {
	case constant.WSGetNewestSeq:
		ws.newestSeqReq(conn, &m)
	case constant.WSPullMsg:
		ws.pullMsgReq(conn, &m)
	case constant.WSSendMsg:
		ws.sendMsgReq(conn, &m)
	default:
	}
}
func (ws *WServer) newestSeqResp(conn *websocket.Conn, m *Req, pb *pbChat.GetNewSeqResp) {
	mReply := make(map[string]interface{})
	mData := make(map[string]interface{})
	mReply["reqIdentifier"] = m.ReqIdentifier
	mReply["msgIncr"] = m.MsgIncr
	mReply["errCode"] = pb.GetErrCode()
	mReply["errMsg"] = pb.GetErrMsg()
	mData["seq"] = pb.GetSeq()
	mReply["data"] = mData
	ws.sendMsg(conn, mReply)
}
func (ws *WServer) newestSeqReq(conn *websocket.Conn, m *Req) {
	log.InfoByKv("Ws call success to getNewSeq", m.OperationID, "Parameters", m)
	pbData := pbChat.GetNewSeqReq{}
	pbData.UserID = m.SendID
	pbData.OperationID = m.OperationID
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	if grpcConn == nil {
		log.ErrorByKv("get grpcConn err", pbData.OperationID, "args", m)
	}
	msgClient := pbChat.NewChatClient(grpcConn)
	reply, err := msgClient.GetNewSeq(context.Background(), &pbData)
	if err != nil {
		log.ErrorByKv("rpc call failed to getNewSeq", pbData.OperationID, "err", err, "pbData", pbData.String())
		return
	}
	log.InfoByKv("rpc call success to getNewSeq", pbData.OperationID, "replyData", reply.String())
	ws.newestSeqResp(conn, m, reply)

}

func (ws *WServer) pullMsgResp(conn *websocket.Conn, m *Req, pb *pbChat.PullMessageResp) {
	mReply := make(map[string]interface{})
	msg := make(map[string]interface{})
	mReply["reqIdentifier"] = m.ReqIdentifier
	mReply["msgIncr"] = m.MsgIncr
	mReply["errCode"] = pb.GetErrCode()
	mReply["errMsg"] = pb.GetErrMsg()
	//空切片
	if v := pb.GetSingleUserMsg(); v != nil {
		msg["single"] = v
	} else {
		msg["single"] = []pbChat.GatherFormat{}
	}
	if v := pb.GetGroupUserMsg(); v != nil {
		msg["group"] = v
	} else {
		msg["group"] = []pbChat.GatherFormat{}
	}
	msg["maxSeq"] = pb.GetMaxSeq()
	msg["minSeq"] = pb.GetMinSeq()
	mReply["data"] = msg
	ws.sendMsg(conn, mReply)

}

func (ws *WServer) pullMsgReq(conn *websocket.Conn, m *Req) {
	log.InfoByKv("Ws call success to pullMsgReq", m.OperationID, "Parameters", m)
	reply := new(pbChat.PullMessageResp)
	isPass, errCode, errMsg, data := ws.argsValidate(m, constant.WSPullMsg)
	if isPass {
		pbData := pbChat.PullMessageReq{}
		pbData.UserID = m.SendID
		pbData.OperationID = m.OperationID
		pbData.SeqBegin = data.(SeqData).SeqBegin
		pbData.SeqEnd = data.(SeqData).SeqEnd
		grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
		msgClient := pbChat.NewChatClient(grpcConn)
		reply, err := msgClient.PullMessage(context.Background(), &pbData)
		if err != nil {
			log.ErrorByKv("PullMessage error", pbData.OperationID, "err", err.Error())
			return
		}
		log.InfoByKv("rpc call success to pullMsgRep", pbData.OperationID, "ReplyArgs", reply.String(), "maxSeq", reply.GetMaxSeq(),
			"MinSeq", reply.GetMinSeq(), "singLen", len(reply.GetSingleUserMsg()), "groupLen", len(reply.GetGroupUserMsg()))
		ws.pullMsgResp(conn, m, reply)
	} else {
		reply.ErrCode = errCode
		reply.ErrMsg = errMsg
		ws.pullMsgResp(conn, m, reply)
	}
}

func (ws *WServer) sendMsgResp(conn *websocket.Conn, m *Req, pb *pbChat.UserSendMsgResp) {
	mReply := make(map[string]interface{})
	mReplyData := make(map[string]interface{})
	mReply["reqIdentifier"] = m.ReqIdentifier
	mReply["msgIncr"] = m.MsgIncr
	mReply["errCode"] = pb.GetErrCode()
	mReply["errMsg"] = pb.GetErrMsg()
	mReplyData["clientMsgID"] = pb.GetClientMsgID()
	mReplyData["serverMsgID"] = pb.GetServerMsgID()
	mReply["data"] = mReplyData
	ws.sendMsg(conn, mReply)
}

func (ws *WServer) sendMsgReq(conn *websocket.Conn, m *Req) {
	log.InfoByKv("Ws call success to sendMsgReq", m.OperationID, "Parameters", m)
	reply := new(pbChat.UserSendMsgResp)
	isPass, errCode, errMsg, pData := ws.argsValidate(m, constant.WSSendMsg)
	if isPass {
		data := pData.(MsgData)
		pbData := pbChat.UserSendMsgReq{
			ReqIdentifier: m.ReqIdentifier,
			Token:         m.Token,
			SendID:        m.SendID,
			OperationID:   m.OperationID,
			PlatformID:    data.PlatformID,
			SessionType:   data.SessionType,
			MsgFrom:       data.MsgFrom,
			ContentType:   data.ContentType,
			RecvID:        data.RecvID,
			ForceList:     data.ForceList,
			Content:       data.Content,
			Options:       utils.MapToJsonString(data.Options),
			ClientMsgID:   data.ClientMsgID,
			OffLineInfo:   utils.MapToJsonString(data.OfflineInfo),
		}
		log.InfoByKv("Ws call success to sendMsgReq", m.OperationID, "Parameters", m)
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
		client := pbChat.NewChatClient(etcdConn)
		log.Info("", "", "api UserSendMsg call, api call rpc...")
		reply, _ := client.UserSendMsg(context.Background(), &pbData)
		log.Info("", "", "api UserSendMsg call end..., [data: %s] [reply: %s]", pbData.String(), reply.String())
		ws.sendMsgResp(conn, m, reply)
	} else {
		reply.ErrCode = errCode
		reply.ErrMsg = errMsg
		ws.sendMsgResp(conn, m, reply)
	}

}

func (ws *WServer) sendMsg(conn *websocket.Conn, mReply map[string]interface{}) {
	bMsg, _ := json.Marshal(mReply)
	err := ws.writeMsg(conn, websocket.TextMessage, bMsg)
	if err != nil {
		log.ErrorByKv("WS WriteMsg error", "", "userIP", conn.RemoteAddr().String(), "userUid", ws.getUserUid(conn), "error", err, "mReply", mReply)
	}
}
func (ws *WServer) sendErrMsg(conn *websocket.Conn, errCode int32, errMsg string) {
	mReply := make(map[string]interface{})
	mReply["errCode"] = errCode
	mReply["errMsg"] = errMsg
	ws.sendMsg(conn, mReply)
}

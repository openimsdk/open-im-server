package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"runtime"
	"strings"
)

func (ws *WServer) msgParse(conn *UserConn, binaryMsg []byte) {
	//ws online debug data
	//{"ReqIdentifier":1001,"Token":"123","SendID":"c4ca4238a0b923820dcc509a6f75849b","Time":"123","OperationID":"123","MsgIncr":0}
	//{"ReqIdentifier":1002,"Token":"123","SendID":"c4ca4238a0b923820dcc509a6f75849b","Time":"123","OperationID":"123","MsgIncr":0,"SeqBegin":1,"SeqEnd":6}
	//{"ReqIdentifier":1003,"Token":"123","SendID":"c4ca4238a0b923820dcc509a6f75849b",
	//"RecvID":"a87ff679a2f3e71d9181a67b7542122c","ClientMsgID":"2343","Time":"147878787","OperationID":
	//"123","MsgIncr":0,"SubMsgType":101,"MsgType":100,"MsgFrom":1,"Content":"sdfsdf"}
	b := bytes.NewBuffer(binaryMsg)
	m := Req{}
	dec := gob.NewDecoder(b)
	err := dec.Decode(&m)
	if err != nil {
		log.ErrorByKv("ws json Unmarshal err", "", "err", err.Error())
		ws.sendErrMsg(conn, 200, err.Error(), constant.WSDataError, "", "")
		err = conn.Close()
		if err != nil {
			log.NewError("", "ws close err", err.Error())
		}
		return
	}
	if err := validate.Struct(m); err != nil {
		log.ErrorByKv("ws args validate  err", "", "err", err.Error())
		ws.sendErrMsg(conn, 201, err.Error(), m.ReqIdentifier, m.MsgIncr, m.OperationID)
		return
	}
	//if !utils.VerifyToken(m.Token, m.SendID) {
	//	ws.sendErrMsg(conn, 202, "token validate err", m.ReqIdentifier, m.MsgIncr,m.OperationID)
	//	return
	//}
	log.InfoByKv("Basic Info Authentication Success", m.OperationID, "reqIdentifier", m.ReqIdentifier, "sendID", m.SendID, "msgIncr", m.MsgIncr)

	switch m.ReqIdentifier {
	case constant.WSGetNewestSeq:
		go ws.getSeqReq(conn, &m)
	case constant.WSPullMsg:
		go ws.pullMsgReq(conn, &m)
	case constant.WSSendMsg:
		sendTime := utils.GetCurrentTimestampByNano()
		go ws.sendMsgReq(conn, &m, sendTime)
	case constant.WSPullMsgBySeqList:
		go ws.pullMsgBySeqListReq(conn, &m)
	default:
	}
	log.NewInfo("", "goroutine num is ", runtime.NumGoroutine())
}
func (ws *WServer) getSeqReq(conn *UserConn, m *Req) {
	log.NewInfo(m.OperationID, "Ws call success to getNewSeq", m.MsgIncr, m.SendID, m.ReqIdentifier)
	pbData := pbChat.GetMaxAndMinSeqReq{}
	nReply := new(pbChat.GetMaxAndMinSeqResp)
	pbData.UserID = m.SendID
	pbData.OperationID = m.OperationID
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	if grpcConn == nil {
		log.ErrorByKv("get grpcConn err", pbData.OperationID, "args", m)
	}
	msgClient := pbChat.NewChatClient(grpcConn)
	reply, err := msgClient.GetMaxAndMinSeq(context.Background(), &pbData)
	if err != nil {
		log.ErrorByKv("rpc call failed to getSeqReq", pbData.OperationID, "err", err, "pbData", pbData.String())
		nReply.ErrCode = 200
		nReply.ErrMsg = err.Error()
		ws.getSeqResp(conn, m, nReply)
	} else {
		log.InfoByKv("rpc call success to getSeqReq", pbData.OperationID, "replyData", reply.String())
		ws.getSeqResp(conn, m, reply)
	}
}
func (ws *WServer) getSeqResp(conn *UserConn, m *Req, pb *pbChat.GetMaxAndMinSeqResp) {
	var mReplyData open_im_sdk.GetMaxAndMinSeqResp
	mReplyData.MaxSeq = pb.GetMaxSeq()
	mReplyData.MinSeq = pb.GetMinSeq()
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
func (ws *WServer) pullMsgReq(conn *UserConn, m *Req) {
	log.NewInfo(m.OperationID, "Ws call success to pullMsgReq", m.ReqIdentifier, m.MsgIncr, m.SendID)
	nReply := new(pbChat.PullMessageResp)
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
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.pullMsgResp(conn, m, nReply)
		} else {
			log.InfoByKv("rpc call success to pullMsgRep", pbData.OperationID, "ReplyArgs", reply.String(), "maxSeq", reply.GetMaxSeq(),
				"MinSeq", reply.GetMinSeq(), "singLen", len(reply.GetSingleUserMsg()), "groupLen", len(reply.GetGroupUserMsg()))
			ws.pullMsgResp(conn, m, reply)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		ws.pullMsgResp(conn, m, nReply)
	}
}
func (ws *WServer) pullMsgResp(conn *UserConn, m *Req, pb *pbChat.PullMessageResp) {
	log.NewInfo(m.OperationID, "pullMsgResp come  here ", pb.String())
	var mReplyData open_im_sdk.PullMessageBySeqListResp
	a, err := json.Marshal(pb.SingleUserMsg)
	if err != nil {
		log.NewError(m.OperationID, "GetSingleUserMsg,json marshal,err", err.Error())
	}
	log.NewInfo(m.OperationID, "pullMsgResp json is ", len(pb.SingleUserMsg))
	err = json.Unmarshal(a, &mReplyData.SingleUserMsg)
	if err != nil {
		log.NewError(m.OperationID, "SingleUserMsg,json Unmarshal,err", err.Error())
	}
	b, err := json.Marshal(pb.GroupUserMsg)
	if err != nil {
		log.NewError(m.OperationID, "mReplyData,json marshal,err", err.Error())
	}
	err = json.Unmarshal(b, &mReplyData.GroupUserMsg)
	if err != nil {
		log.NewError(m.OperationID, "test SingleUserMsg,json Unmarshal,err", err.Error())
	}
	c, err := proto.Marshal(&mReplyData)
	log.NewInfo(m.OperationID, "test info is ", len(mReplyData.SingleUserMsg), mReplyData.SingleUserMsg)

	mReply := Resp{
		ReqIdentifier: m.ReqIdentifier,
		MsgIncr:       m.MsgIncr,
		ErrCode:       pb.GetErrCode(),
		ErrMsg:        pb.GetErrMsg(),
		OperationID:   m.OperationID,
		Data:          c,
	}
	log.NewInfo(m.OperationID, "pullMsgResp all data  is ", mReply.ReqIdentifier, mReply.MsgIncr, mReply.ErrCode, mReply.ErrMsg,
		len(mReply.Data))

	ws.sendMsg(conn, mReply)

}
func (ws *WServer) pullMsgBySeqListReq(conn *UserConn, m *Req) {
	log.NewInfo(m.OperationID, "Ws call success to pullMsgBySeqListReq start", m.SendID, m.ReqIdentifier, m.MsgIncr)
	nReply := new(pbChat.PullMessageResp)
	isPass, errCode, errMsg, data := ws.argsValidate(m, constant.WSPullMsgBySeqList)
	log.NewInfo(m.OperationID, "Ws call success to pullMsgBySeqListReq middle", m.SendID, m.ReqIdentifier, m.MsgIncr, data.(open_im_sdk.PullMessageBySeqListReq).SeqList)
	if isPass {
		pbData := pbChat.PullMessageBySeqListReq{}
		pbData.SeqList = data.(open_im_sdk.PullMessageBySeqListReq).SeqList
		pbData.UserID = m.SendID
		pbData.OperationID = m.OperationID
		grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
		msgClient := pbChat.NewChatClient(grpcConn)
		reply, err := msgClient.PullMessageBySeqList(context.Background(), &pbData)
		if err != nil {
			log.NewError(pbData.OperationID, "pullMsgBySeqListReq err", err.Error())
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.pullMsgResp(conn, m, nReply)
		} else {
			log.NewInfo(pbData.OperationID, "rpc call success to pullMsgBySeqListReq", reply.String(), reply.GetMaxSeq(), reply.GetMinSeq(), len(reply.GetSingleUserMsg()), len(reply.GetGroupUserMsg()))
			ws.pullMsgResp(conn, m, reply)
		}
	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		ws.pullMsgResp(conn, m, nReply)
	}
}
func (ws *WServer) sendMsgReq(conn *UserConn, m *Req, sendTime int64) {
	log.NewInfo(m.OperationID, "Ws call success to sendMsgReq start", m.MsgIncr, m.ReqIdentifier, m.SendID, sendTime)
	nReply := new(pbChat.UserSendMsgResp)
	isPass, errCode, errMsg, pData := ws.argsValidate(m, constant.WSSendMsg)
	if isPass {
		data := pData.(open_im_sdk.UserSendMsgReq)
		pbData := pbChat.UserSendMsgReq{
			ReqIdentifier:  m.ReqIdentifier,
			Token:          m.Token,
			SendID:         m.SendID,
			OperationID:    m.OperationID,
			PlatformID:     data.PlatformID,
			SessionType:    data.SessionType,
			MsgFrom:        data.MsgFrom,
			ContentType:    data.ContentType,
			RecvID:         data.RecvID,
			ForceList:      data.ForceList,
			SenderNickName: data.SenderNickName,
			SenderFaceURL:  data.SenderFaceURL,
			Content:        data.Content,
			Options:        utils.MapIntToJsonString(data.Options),
			ClientMsgID:    data.ClientMsgID,
			SendTime:       sendTime,
		}
		log.NewInfo(m.OperationID, "Ws call success to sendMsgReq middle", m.ReqIdentifier, m.SendID, m.MsgIncr, data)
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
		client := pbChat.NewChatClient(etcdConn)
		reply, err := client.UserSendMsg(context.Background(), &pbData)
		if err != nil {
			log.NewError(pbData.OperationID, "UserSendMsg err", err.Error())
			nReply.ErrCode = 200
			nReply.ErrMsg = err.Error()
			ws.sendMsgResp(conn, m, nReply, sendTime)
		} else {
			log.NewInfo(pbData.OperationID, "rpc call success to sendMsgReq", reply.String())
			ws.sendMsgResp(conn, m, reply, sendTime)
		}

	} else {
		nReply.ErrCode = errCode
		nReply.ErrMsg = errMsg
		ws.sendMsgResp(conn, m, nReply, sendTime)
	}

}
func (ws *WServer) sendMsgResp(conn *UserConn, m *Req, pb *pbChat.UserSendMsgResp, sendTime int64) {
	// := make(map[string]interface{})

	var mReplyData open_im_sdk.UserSendMsgResp
	mReplyData.ClientMsgID = pb.GetClientMsgID()
	mReplyData.ServerMsgID = pb.GetServerMsgID()
	mReplyData.SendTime = sendTime
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

func (ws *WServer) sendMsg(conn *UserConn, mReply interface{}) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(mReply)
	if err != nil {
		log.NewError(mReply.(Resp).OperationID, mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, "Encode Msg error", conn.RemoteAddr().String(), ws.getUserUid(conn), err.Error())
		return
	}
	err = ws.writeMsg(conn, websocket.BinaryMessage, b.Bytes())
	if err != nil {
		log.NewError(mReply.(Resp).OperationID, mReply.(Resp).ReqIdentifier, mReply.(Resp).ErrCode, mReply.(Resp).ErrMsg, "WS WriteMsg error", conn.RemoteAddr().String(), ws.getUserUid(conn), err.Error())
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

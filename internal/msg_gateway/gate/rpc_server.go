package gate

import (
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	"Open_IM/src/common/log"
	"Open_IM/src/grpc-etcdv3/getcdv3"
	"Open_IM/src/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"net"
	"strings"
)

type RPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func (r *RPCServer) onInit(rpcPort int) {
	r.rpcPort = rpcPort
	r.rpcRegisterName = config.Config.RpcRegisterName.OpenImOnlineMessageRelayName
	r.etcdSchema = config.Config.Etcd.EtcdSchema
	r.etcdAddr = config.Config.Etcd.EtcdAddr
}
func (r *RPCServer) run() {
	ip := utils.ServerIP
	registerAddress := ip + ":" + utils.IntToString(r.rpcPort)
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.ErrorByArgs(fmt.Sprintf("fail to listening consumer, err:%v\n", err))
		return
	}
	defer listener.Close()
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	pbRelay.RegisterOnlineMessageRelayServiceServer(srv, r)
	err = getcdv3.RegisterEtcd4Unique(r.etcdSchema, strings.Join(r.etcdAddr, ","), ip, r.rpcPort, r.rpcRegisterName, 10)
	if err != nil {
		log.ErrorByKv("register push message rpc to etcd err", "", "err", err.Error())
	}
	err = srv.Serve(listener)
	if err != nil {
		log.ErrorByKv("push message rpc listening err", "", "err", err.Error())
		return
	}
}
func (r *RPCServer) MsgToUser(_ context.Context, in *pbRelay.MsgToUserReq) (*pbRelay.MsgToUserResp, error) {
	log.InfoByKv("PushMsgToUser is arriving", in.OperationID, "args", in.String())
	var resp []*pbRelay.SingleMsgToUser
	var RecvID string
	msg := make(map[string]interface{})
	mReply := make(map[string]interface{})
	mReply["reqIdentifier"] = constant.WSPushMsg
	mReply["errCode"] = 0
	mReply["errMsg"] = ""
	msg["sendID"] = in.SendID
	msg["recvID"] = in.RecvID
	msg["msgFrom"] = in.MsgFrom
	msg["contentType"] = in.ContentType
	msg["sessionType"] = in.SessionType
	msg["senderNickName"] = in.SenderNickName
	msg["senderFaceUrl"] = in.SenderFaceURL
	msg["clientMsgID"] = in.ClientMsgID
	msg["serverMsgID"] = in.ServerMsgID
	msg["content"] = in.Content
	msg["seq"] = in.RecvSeq
	msg["sendTime"] = in.SendTime
	msg["senderPlatformID"] = in.PlatformID
	mReply["data"] = msg
	bMsg, _ := json.Marshal(mReply)
	switch in.GetSessionType() {
	case constant.SingleChatType:
		RecvID = in.GetRecvID()
	case constant.GroupChatType:
		RecvID = strings.Split(in.GetRecvID(), " ")[0]
	}
	log.InfoByKv("test", in.OperationID, "wsUserToConn", ws.wsUserToConn)
	for key, conn := range ws.wsUserToConn {
		UIDAndPID := strings.Split(key, " ")
		if UIDAndPID[0] == RecvID {
			resultCode := sendMsgToUser(conn, bMsg, in, UIDAndPID[1], UIDAndPID[0])
			temp := &pbRelay.SingleMsgToUser{
				ResultCode:     resultCode,
				RecvID:         UIDAndPID[0],
				RecvPlatFormID: utils.PlatformNameToID(UIDAndPID[1]),
			}
			resp = append(resp, temp)
		}
	}
	//switch in.GetContentType() {
	//case constant.SyncSenderMsg:
	//	log.InfoByKv("come sync", in.OperationID, "args", in.String())
	//	RecvID = in.GetSendID()
	//	if in.MsgFrom != constant.SysMsgType {
	//		for key, conn := range ws.wsUserToConn {
	//			UIDAndPID := strings.Split(key, " ")
	//			if UIDAndPID[0] == RecvID && utils.PlatformIDToName(in.GetPlatformID()) != UIDAndPID[1] {
	//				resultCode := sendMsgToUser(conn, bMsg, in, UIDAndPID[1], UIDAndPID[0])
	//				temp := &pbRelay.SingleMsgToUser{
	//					ResultCode:     resultCode,
	//					RecvID:         UIDAndPID[0],
	//					RecvPlatFormID: utils.PlatformNameToID(UIDAndPID[1]),
	//				}
	//				resp = append(resp, temp)
	//			}
	//
	//		}
	//	}
	//default:
	//	log.InfoByKv("not come sync", in.OperationID, "args", in.String())
	//	switch in.SessionType {
	//	case constant.SingleChatType:
	//		log.InfoByKv("come single", in.OperationID, "args", in.String())
	//		RecvID = in.GetRecvID()
	//	case constant.GroupChatType:
	//		RecvID = strings.Split(in.GetRecvID(), " ")[0]
	//	default:
	//	}
	//	log.InfoByKv("come for range", in.OperationID, "args", in.String())
	//
	//	for key, conn := range ws.wsUserToConn {
	//		UIDAndPID := strings.Split(key, " ")
	//		if UIDAndPID[0] == RecvID {
	//			resultCode := sendMsgToUser(conn, bMsg, in, UIDAndPID[1], UIDAndPID[0])
	//			temp := &pbRelay.SingleMsgToUser{
	//				ResultCode:     resultCode,
	//				RecvID:         UIDAndPID[0],
	//				RecvPlatFormID: utils.PlatformNameToID(UIDAndPID[1]),
	//			}
	//			resp = append(resp, temp)
	//		}
	//	}
	//}
	return &pbRelay.MsgToUserResp{
		Resp: resp,
	}, nil
}

func sendMsgToUser(conn *websocket.Conn, bMsg []byte, in *pbRelay.MsgToUserReq, RecvPlatForm, RecvID string) (ResultCode int64) {
	err := ws.writeMsg(conn, websocket.TextMessage, bMsg)
	if err != nil {
		log.ErrorByKv("PushMsgToUser is failed By Ws", "", "Addr", conn.RemoteAddr().String(),
			"error", err, "senderPlatform", utils.PlatformIDToName(in.PlatformID), "recvPlatform", RecvPlatForm, "args", in.String(), "recvID", RecvID)
		ResultCode = -2
		return ResultCode
	} else {
		log.InfoByKv("PushMsgToUser is success By Ws", in.OperationID, "args", in.String())
		ResultCode = 0
		return ResultCode
	}

}

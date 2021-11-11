package gate

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbRelay "Open_IM/pkg/proto/relay"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"strings"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
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
	msg := open_im_sdk.MsgData{
		SendID:           in.SendID,
		RecvID:           in.RecvID,
		MsgFrom:          in.MsgFrom,
		ContentType:      in.ContentType,
		SessionType:      in.SessionType,
		SenderNickName:   in.SenderNickName,
		SenderFaceURL:    in.SenderFaceURL,
		ClientMsgID:      in.ClientMsgID,
		ServerMsgID:      in.ServerMsgID,
		Content:          in.Content,
		Seq:              in.RecvSeq,
		SendTime:         in.SendTime,
		SenderPlatformID: in.PlatformID,
	}
	msgBytes, _ := proto.Marshal(&msg)
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
	switch in.GetSessionType() {
	case constant.SingleChatType:
		RecvID = in.GetRecvID()
	case constant.GroupChatType:
		RecvID = strings.Split(in.GetRecvID(), " ")[0]
	}
	var tag bool
	var UIDAndPID []string
	userIDList := genUidPlatformArray(RecvID)
	for _, v := range userIDList {
		UIDAndPID = strings.Split(v, " ")
		if conn := ws.getUserConn(v); conn != nil {
			tag = true
			resultCode := sendMsgToUser(conn, replyBytes.Bytes(), in, UIDAndPID[1], UIDAndPID[0])
			temp := &pbRelay.SingleMsgToUser{
				ResultCode:     resultCode,
				RecvID:         UIDAndPID[0],
				RecvPlatFormID: utils.PlatformNameToID(UIDAndPID[1]),
			}
			resp = append(resp, temp)
		} else {
			temp := &pbRelay.SingleMsgToUser{
				ResultCode:     -1,
				RecvID:         UIDAndPID[0],
				RecvPlatFormID: utils.PlatformNameToID(UIDAndPID[1]),
			}
			resp = append(resp, temp)
		}
	}
	if !tag {
		log.NewError(in.OperationID, "push err ,no matched ws conn not in map", in.String())
	}
	return &pbRelay.MsgToUserResp{
		Resp: resp,
	}, nil
}

func sendMsgToUser(conn *UserConn, bMsg []byte, in *pbRelay.MsgToUserReq, RecvPlatForm, RecvID string) (ResultCode int64) {
	err := ws.writeMsg(conn, websocket.BinaryMessage, bMsg)
	if err != nil {
		log.ErrorByKv("PushMsgToUser is failed By Ws", "", "Addr", conn.RemoteAddr().String(),
			"error", err, "senderPlatform", utils.PlatformIDToName(in.PlatformID), "recvPlatform", RecvPlatForm, "args", in.String(), "recvID", RecvID)
		ResultCode = -2
		return ResultCode
	} else {
		log.InfoByKv("PushMsgToUser is success By Ws", in.OperationID, "args", in.String(), "recvPlatForm", RecvPlatForm, "recvID", RecvID)
		ResultCode = 0
		return ResultCode
	}

}
func genUidPlatformArray(uid string) (array []string) {
	for i := 1; i <= utils.LinuxPlatformID; i++ {
		array = append(array, uid+" "+utils.PlatformIDToName(int32(i)))
	}
	return array
}

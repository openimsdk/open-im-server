package msggateway

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/prome"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/proto/msggateway"
	"OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"bytes"
	"context"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
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

func initPrometheus() {
	prome.NewMsgRecvTotalCounter()
	prome.NewGetNewestSeqTotalCounter()
	prome.NewPullMsgBySeqListTotalCounter()
	prome.NewMsgOnlinePushSuccessCounter()
	prome.NewOnlineUserGauges()
	//prome.NewSingleChatMsgRecvSuccessCounter()
	//prome.NewGroupChatMsgRecvSuccessCounter()
	//prome.NewWorkSuperGroupChatMsgRecvSuccessCounter()
}

func (r *RPCServer) onInit(rpcPort int) {
	r.rpcPort = rpcPort
	r.rpcRegisterName = config.Config.RpcRegisterName.OpenImMessageGatewayName
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
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		prome.NewGrpcRequestCounter()
		prome.NewGrpcRequestFailedCounter()
		prome.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(prome.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	msggateway.RegisterMsgGatewayServer(srv, r)

	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	err = rpc.RegisterEtcd4Unique(r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName, 10)
	if err != nil {
		log.Error("", "register push message rpc to etcd err", "", "err", err.Error(), r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName)
		panic(utils.Wrap(err, "register msg_gataway module  rpc to etcd err"))
	}
	r.target = rpc.GetTarget(r.etcdSchema, rpcRegisterIP, r.rpcPort, r.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.Error("", "push message rpc listening err", "", "err", err.Error())
		return
	}
}
func (r *RPCServer) OnlinePushMsg(ctx context.Context, in *msggateway.OnlinePushMsgReq) (*msggateway.OnlinePushMsgResp, error) {
	log.NewInfo(in.OperationID, "PushMsgToUser is arriving", in.String())
	var resp []*msggateway.SingleMsgToUserPlatform
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
			temp := &msggateway.SingleMsgToUserPlatform{
				ResultCode:     resultCode,
				RecvID:         recvID,
				RecvPlatFormID: int32(v),
			}
			resp = append(resp, temp)
		} else {
			temp := &msggateway.SingleMsgToUserPlatform{
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
	return &msggateway.OnlinePushMsgResp{
		Resp: resp,
	}, nil
}
func (r *RPCServer) GetUsersOnlineStatus(_ context.Context, req *msggateway.GetUsersOnlineStatusReq) (*msggateway.GetUsersOnlineStatusResp, error) {
	log.NewInfo(req.OperationID, "rpc GetUsersOnlineStatus arrived server", req.String())
	if !tokenverify.IsManagerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "no permission GetUsersOnlineStatus ", req.OpUserID)
		return &msggateway.GetUsersOnlineStatusResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	var resp msggateway.GetUsersOnlineStatusResp
	for _, userID := range req.UserIDList {
		temp := new(msggateway.GetUsersOnlineStatusResp_SuccessResult)
		temp.UserID = userID
		userConnMap := ws.getUserAllCons(userID)
		for platform, userConn := range userConnMap {
			if userConn != nil {
				ps := new(msggateway.GetUsersOnlineStatusResp_SuccessDetail)
				ps.Platform = constant.PlatformIDToName(platform)
				ps.Status = constant.OnlineStatus
				ps.ConnID = userConn.connID
				ps.IsBackground = userConn.IsBackground
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

func (r *RPCServer) SuperGroupOnlineBatchPushOneMsg(_ context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	log.NewInfo(req.OperationID, "BatchPushMsgToUser is arriving", req.String())
	var singleUserResult []*msggateway.SingleMsgToUserResultList
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
		var resp []*msggateway.SingleMsgToUserPlatform
		tempT := &msggateway.SingleMsgToUserResultList{
			UserID: v,
		}
		userConnMap := ws.getUserAllCons(v)
		for platform, userConn := range userConnMap {
			if userConn != nil {
				temp := &msggateway.SingleMsgToUserPlatform{
					RecvID:         v,
					RecvPlatFormID: int32(platform),
				}
				if !userConn.IsBackground {
					resultCode := sendMsgBatchToUser(userConn, replyBytes.Bytes(), req, platform, v)
					if resultCode == 0 && utils.IsContainInt(platform, r.pushTerminal) {
						tempT.OnlinePush = true
						prome.Inc(prome.MsgOnlinePushSuccessCounter)
						log.Info(req.OperationID, "PushSuperMsgToUser is success By Ws", "args", req.String(), "recvPlatForm", constant.PlatformIDToName(platform), "recvID", v)
						temp.ResultCode = resultCode
						resp = append(resp, temp)
					}
				} else {
					temp.ResultCode = -2
					resp = append(resp, temp)
				}
			}
		}
		tempT.Resp = resp
		singleUserResult = append(singleUserResult, tempT)
	}

	return &msggateway.OnlineBatchPushOneMsgResp{
		SinglePushResult: singleUserResult,
	}, nil
}
func (r *RPCServer) OnlineBatchPushOneMsg(_ context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	log.NewInfo(req.OperationID, "BatchPushMsgToUser is arriving", req.String())
	var singleUserResult []*msggateway.SingleMsgToUserResultList

	for _, v := range req.PushToUserIDList {
		var resp []*msggateway.SingleMsgToUserPlatform
		tempT := &msggateway.SingleMsgToUserResultList{
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
					temp := &msggateway.SingleMsgToUserPlatform{
						ResultCode:     resultCode,
						RecvID:         v,
						RecvPlatFormID: int32(platform),
					}
					resp = append(resp, temp)
				}
			} else {
				if utils.IsContainInt(platform, r.pushTerminal) {
					tempT.OnlinePush = true
					temp := &msggateway.SingleMsgToUserPlatform{
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
	return &msggateway.OnlineBatchPushOneMsgResp{
		SinglePushResult: singleUserResult,
	}, nil
}
func (r *RPCServer) encodeWsData(wsData *sdkws.MsgData, operationID string) (bytes.Buffer, error) {
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

func (r *RPCServer) KickUserOffline(_ context.Context, req *msggateway.KickUserOfflineReq) (*msggateway.KickUserOfflineResp, error) {
	log.NewInfo(req.OperationID, "KickUserOffline is arriving", req.String())
	for _, v := range req.KickUserIDList {
		log.NewWarn(req.OperationID, "SetTokenKicked ", v, req.PlatformID, req.OperationID)
		SetTokenKicked(v, int(req.PlatformID), req.OperationID)
		oldConnMap := ws.getUserAllCons(v)
		if conn, ok := oldConnMap[int(req.PlatformID)]; ok { // user->map[platform->conn]
			log.NewWarn(req.OperationID, "send kick msg, close connection ", req.PlatformID, v)
			ws.sendKickMsg(conn)
			conn.Close()
		}
	}
	return &msggateway.KickUserOfflineResp{}, nil
}

func (r *RPCServer) MultiTerminalLoginCheck(ctx context.Context, req *msggateway.MultiTerminalLoginCheckReq) (*msggateway.MultiTerminalLoginCheckResp, error) {

	ws.MultiTerminalLoginCheckerWithLock(req.UserID, int(req.PlatformID), req.Token, req.OperationID)
	return &msggateway.MultiTerminalLoginCheckResp{}, nil
}

func sendMsgToUser(conn *UserConn, bMsg []byte, in *msggateway.OnlinePushMsgReq, RecvPlatForm int, RecvID string) (ResultCode int64) {
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
func sendMsgBatchToUser(conn *UserConn, bMsg []byte, in *msggateway.OnlineBatchPushOneMsgReq, RecvPlatForm int, RecvID string) (ResultCode int64) {
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

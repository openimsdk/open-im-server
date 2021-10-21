package group

import (
	"Open_IM/internal/push/content_struct"
	"Open_IM/internal/push/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"
)

type groupServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewGroupServer(port int) *groupServer {
	log.NewPrivateLog("group")
	return &groupServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImGroupName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}
func (s *groupServer) Run() {
	log.Info("", "", "rpc group init....")

	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.InfoByArgs("listen network failed,err=%s", err.Error())
		return
	}
	log.Info("", "", "listen network success, address = %s", registerAddress)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbGroup.RegisterGroupServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.ErrorByArgs("get etcd failed,err=%s", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.ErrorByArgs("listen rpc_group error,err=%s", err.Error())
		return
	}
	log.Info("", "", "rpc create group init success")
}

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (*pbGroup.CreateGroupResp, error) {
	log.InfoByArgs("rpc create group is server,args=%s", req.String())
	var (
		groupId string
	)
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.CreateGroupResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	//Time stamp + MD5 to generate group chat id
	groupId = utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	err = im_mysql_model.InsertIntoGroup(groupId, req.GroupName, req.Introduction, req.Notification, req.FaceUrl, req.Ex)
	if err != nil {
		log.ErrorByKv("create group chat failed", req.OperationID, "err=%s", err.Error())
		return &pbGroup.CreateGroupResp{ErrorCode: config.ErrCreateGroup.ErrCode, ErrorMsg: config.ErrCreateGroup.ErrMsg}, nil
	}

	isMagagerFlag := 0
	tokenUid := claims.UID

	if utils.IsContain(tokenUid, config.Config.Manager.AppManagerUid) {
		isMagagerFlag = 1
	}

	if isMagagerFlag == 0 {
		//Add the group owner to the group first, otherwise the group creation will fail
		us, err := im_mysql_model.FindUserByUID(claims.UID)
		if err != nil {
			log.Error("", req.OperationID, "find userInfo failed", err.Error())
			return &pbGroup.CreateGroupResp{ErrorCode: config.ErrCreateGroup.ErrCode, ErrorMsg: config.ErrCreateGroup.ErrMsg}, nil
		}
		err = im_mysql_model.InsertIntoGroupMember(groupId, claims.UID, us.Name, us.Icon, constant.GroupOwner)
		if err != nil {
			log.Error("", req.OperationID, "create group chat failed,err=%s", err.Error())
			return &pbGroup.CreateGroupResp{ErrorCode: config.ErrCreateGroup.ErrCode, ErrorMsg: config.ErrCreateGroup.ErrMsg}, nil
		}

		err = db.DB.AddGroupMember(groupId, claims.UID)
		if err != nil {
			log.Error("", "", "create mongo group member failed, db.DB.AddGroupMember fail [err: %s]", err.Error())
			return &pbGroup.CreateGroupResp{ErrorCode: config.ErrCreateGroup.ErrCode, ErrorMsg: config.ErrCreateGroup.ErrMsg}, nil
		}
	}

	//Binding group id and member id
	for _, user := range req.MemberList {
		us, err := im_mysql_model.FindUserByUID(user.Uid)
		if err != nil {
			log.Error("", req.OperationID, "find userInfo failed,uid=%s", user.Uid, err.Error())
			continue
		}
		err = im_mysql_model.InsertIntoGroupMember(groupId, user.Uid, us.Name, us.Icon, user.SetRole)
		if err != nil {
			log.ErrorByArgs("pull %s to group %s failed,err=%s", user.Uid, groupId, err.Error())
		}
		err = db.DB.AddGroupMember(groupId, user.Uid)
		if err != nil {
			log.Error("", "", "add mongo group member failed, db.DB.AddGroupMember fail [err: %s]", err.Error())
		}
	}

	if isMagagerFlag == 1 {

		//type NotificationContent struct {
		//	IsDisplay   int32  `json:"isDisplay"`
		//	DefaultTips string `json:"defaultTips"`
		//	Detail      string `json:"detail"`
		//}	n := NotificationContent{
		//		IsDisplay:   1,
		//		DefaultTips: "You have joined the group chat:" + createGroupResp.Data.GroupName,
		//		Detail:      createGroupResp.Data.GroupId,
		//	}

		////Push message when create group chat
		n := content_struct.NotificationContent{1, req.GroupName, groupId}
		logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
			SendID:      claims.UID,
			RecvID:      groupId,
			Content:     n.ContentToString(),
			SendTime:    utils.GetCurrentTimestampByNano(),
			MsgFrom:     constant.SysMsgType,     //Notification message identification
			ContentType: constant.CreateGroupTip, //Add friend flag
			SessionType: constant.GroupChatType,
			OperationID: req.OperationID,
		})
	}

	log.Info(req.Token, req.OperationID, "rpc create group success return")
	return &pbGroup.CreateGroupResp{GroupID: groupId}, nil
}

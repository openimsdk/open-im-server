package group

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	"Open_IM/src/common/db"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	"Open_IM/src/grpc-etcdv3/getcdv3"
	pbGroup "Open_IM/src/proto/group"
	"Open_IM/src/utils"
	"context"
	"fmt"
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
	err = im_mysql_model.InsertIntoGroup(groupId, req.GroupName, req.Introduction, req.Notification, req.FaceUrl)
	if err != nil {
		log.ErrorByKv("create group chat failed", req.OperationID, "err=%s", err.Error())
		return &pbGroup.CreateGroupResp{ErrorCode: config.ErrCreateGroup.ErrCode, ErrorMsg: config.ErrCreateGroup.ErrMsg}, nil
	}

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
	////Push message when create group chat
	//logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
	//	SendID:      claims.UID,
	//	RecvID:      groupId,
	//	Content:     content_struct.NewContentStructString(0, "", req.String()),
	//	SendTime:    utils.GetCurrentTimestampBySecond(),
	//	MsgFrom:     constant.SysMsgType,     //Notification message identification
	//	ContentType: constant.CreateGroupTip, //Add friend flag
	//	SessionType: constant.GroupChatType,
	//	OperationID: req.OperationID,
	//})
	log.Info(req.Token, req.OperationID, "rpc create group success return")

	t := db.DB.GetGroupMember(groupId)
	fmt.Println("YYYYYYYYYYYYYYYYYYYYYYYYYYYYYY")
	fmt.Println(t)

	return &pbGroup.CreateGroupResp{GroupID: groupId}, nil
}

package group

import (
	"Open_IM/internal/rpc/chat"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
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
	log.NewInfo(req.OperationID, "CreateGroup, args=%s", req.String())
	var (
		groupId string
	)
	//Parse token, to find current user information
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.NewError(req.OperationID, "ParseToken failed, ", err.Error(), req.String())
		return &pbGroup.CreateGroupResp{ErrorCode: constant.ErrParseToken.ErrCode, ErrorMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	//Time stamp + MD5 to generate group chat id
	groupId = utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	err = im_mysql_model.InsertIntoGroup(groupId, req.GroupName, req.Introduction, req.Notification, req.FaceUrl, req.Ex)
	if err != nil {
		log.NewError(req.OperationID, "InsertIntoGroup failed, ", err.Error(), req.String())
		return &pbGroup.CreateGroupResp{ErrorCode: constant.ErrCreateGroup.ErrCode, ErrorMsg: constant.ErrCreateGroup.ErrMsg}, nil
	}

	isManagerFlag := 0
	tokenUid := claims.UID

	if utils.IsContain(tokenUid, config.Config.Manager.AppManagerUid) {
		isManagerFlag = 1
	}

	us, err := im_mysql_model.FindUserByUID(claims.UID)
	if err != nil {
		log.Error("", req.OperationID, "find userInfo failed", err.Error())
		return &pbGroup.CreateGroupResp{ErrorCode: constant.ErrCreateGroup.ErrCode, ErrorMsg: constant.ErrCreateGroup.ErrMsg}, nil
	}

	if isManagerFlag == 0 {
		//Add the group owner to the group first, otherwise the group creation will fail
		err = im_mysql_model.InsertIntoGroupMember(groupId, claims.UID, us.Nickname, us.FaceUrl, constant.GroupOwner)
		if err != nil {
			log.Error("", req.OperationID, "create group chat failed,err=%s", err.Error())
			return &pbGroup.CreateGroupResp{ErrorCode: constant.ErrCreateGroup.ErrCode, ErrorMsg: constant.ErrCreateGroup.ErrMsg}, nil
		}

		err = db.DB.AddGroupMember(groupId, claims.UID)
		if err != nil {
			log.NewError(req.OperationID, "AddGroupMember failed ", err.Error(), groupId, claims.UID)
			return &pbGroup.CreateGroupResp{ErrorCode: constant.ErrCreateGroup.ErrCode, ErrorMsg: constant.ErrCreateGroup.ErrMsg}, nil
		}
	}

	//Binding group id and member id
	for _, user := range req.MemberList {
		us, err := im_mysql_model.FindUserByUID(user.Uid)
		if err != nil {
			log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), user.Uid)
			continue
		}
		err = im_mysql_model.InsertIntoGroupMember(groupId, user.Uid, us.Nickname, us.FaceUrl, user.SetRole)
		if err != nil {
			log.ErrorByArgs("InsertIntoGroupMember failed", user.Uid, groupId, err.Error())
		}
		err = db.DB.AddGroupMember(groupId, user.Uid)
		if err != nil {
			log.Error("", "", "add mongo group member failed, db.DB.AddGroupMember fail [err: %s]", err.Error())
		}
	}

	if isManagerFlag == 1 {

	}
	group, err := im_mysql_model.FindGroupInfoByGroupId(groupId)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed ", err.Error(), groupId)
		return &pbGroup.CreateGroupResp{GroupID: groupId}, nil
	}
	memberList, err := im_mysql_model.FindGroupMemberListByGroupId(groupId)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupMemberListByGroupId failed ", err.Error(), groupId)
		return &pbGroup.CreateGroupResp{GroupID: groupId}, nil
	}
	chat.GroupCreatedNotification(req.OperationID, us, group, memberList)
	log.NewInfo(req.OperationID, "GroupCreatedNotification, rpc CreateGroup success return ", groupId)

	return &pbGroup.CreateGroupResp{GroupID: groupId}, nil
}

package group

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	cp "Open_IM/pkg/common/utils"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"math/big"
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
	log.NewPrivateLog(constant.LogFileName)
	return &groupServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImGroupName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *groupServer) Run() {
	log.NewInfo("", "group rpc start ")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("", "listen network success, ", address, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbGroup.RegisterGroupServer(srv, s)

	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("", "RegisterEtcd failed ", err.Error())
		return
	}
	log.Info("", "RegisterEtcd ", s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("", "group rpc success")
}

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (*pbGroup.CreateGroupResp, error) {
	log.NewInfo(req.OperationID, "CreateGroup, args ", req.String())
	if !token_verify.CheckAccess(req.OpUserID, req.OwnerUserID) {
		log.NewError(req.OperationID, "CheckAccess false ", req.OpUserID, req.OwnerUserID)
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	groupId := req.GroupInfo.GroupID
	if groupId == "" {
		groupId = utils.Md5(req.OperationID + strconv.FormatInt(time.Now().UnixNano(), 10))
		bi := big.NewInt(0)
		bi.SetString(groupId[0:8], 16)
		groupId = bi.String()
	}
	//to group
	groupInfo := db.Group{}
	utils.CopyStructFields(&groupInfo, req.GroupInfo)
	groupInfo.CreatorUserID = req.OpUserID
	groupInfo.GroupID = groupId

	if groupInfo.NotificationUpdateTime.Unix() < 0 {
		groupInfo.NotificationUpdateTime = utils.UnixSecondToTime(0)
	}
	err := imdb.InsertIntoGroup(groupInfo)
	if err != nil {
		log.NewError(req.OperationID, "InsertIntoGroup failed, ", err.Error(), groupInfo)
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, http.WrapError(constant.ErrDB)
	}
	var okUserIDList []string
	resp := &pbGroup.CreateGroupResp{GroupInfo: &open_im_sdk.GroupInfo{}}
	groupMember := db.GroupMember{}
	us := &db.User{}
	if req.OwnerUserID != "" {
		us, err = imdb.GetUserByUserID(req.OwnerUserID)
		if err != nil {
			log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), req.OwnerUserID)
			return &pbGroup.CreateGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, http.WrapError(constant.ErrDB)
		}
		//to group member
		groupMember = db.GroupMember{GroupID: groupId, RoleLevel: constant.GroupOwner, OperatorUserID: req.OpUserID, JoinSource: constant.JoinByInvitation, InviterUserID: req.OpUserID}
		utils.CopyStructFields(&groupMember, us)
		err = imdb.InsertIntoGroupMember(groupMember)
		if err != nil {
			log.NewError(req.OperationID, "InsertIntoGroupMember failed ", err.Error(), groupMember)
			return &pbGroup.CreateGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, http.WrapError(constant.ErrDB)
		}
	}
	if req.GroupInfo.GroupType != constant.SuperGroup {
		//to group member
		for _, user := range req.InitMemberList {
			us, err := imdb.GetUserByUserID(user.UserID)
			if err != nil {
				log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), user.UserID)
				continue
			}
			if user.RoleLevel == constant.GroupOwner {
				log.NewError(req.OperationID, "only one owner, failed ", user)
				continue
			}
			groupMember.RoleLevel = user.RoleLevel
			groupMember.JoinSource = constant.JoinByInvitation
			groupMember.InviterUserID = req.OpUserID
			utils.CopyStructFields(&groupMember, us)
			err = imdb.InsertIntoGroupMember(groupMember)
			if err != nil {
				log.NewError(req.OperationID, "InsertIntoGroupMember failed ", err.Error(), groupMember)
				continue
			}
			okUserIDList = append(okUserIDList, user.UserID)
		}
		group, err := imdb.GetGroupInfoByGroupID(groupId)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), groupId)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = err.Error()
			return resp, nil
		}
		utils.CopyStructFields(resp.GroupInfo, group)
		memberCount, err := imdb.GetGroupMemberNumByGroupID(groupId)
		resp.GroupInfo.MemberCount = uint32(memberCount)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupMemberNumByGroupID failed ", err.Error(), groupId)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = err.Error()
			return resp, nil
		}
		if req.OwnerUserID != "" {
			resp.GroupInfo.OwnerUserID = req.OwnerUserID
			okUserIDList = append(okUserIDList, req.OwnerUserID)
		}
		// superGroup stored in mongodb
	} else {
		for _, v := range req.InitMemberList {
			okUserIDList = append(okUserIDList, v.UserID)
		}
		if err := db.DB.CreateSuperGroup(groupId, okUserIDList, len(okUserIDList)); err != nil {
			log.NewError(req.OperationID, "GetGroupMemberNumByGroupID failed ", err.Error(), groupId)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = err.Error() + ": CreateSuperGroup failed"
			return resp, nil
		}
	}

	if len(okUserIDList) != 0 {
		log.NewInfo(req.OperationID, "rpc CreateGroup return ", resp.String())
		if req.GroupInfo.GroupType != constant.SuperGroup {
			for _, userID := range okUserIDList {
				if err := rocksCache.DelJoinedGroupIDListFromCache(userID); err != nil {
					log.NewWarn(req.OperationID, utils.GetSelfFuncName(), userID, err.Error())
				}
			}
			chat.GroupCreatedNotification(req.OperationID, req.OpUserID, groupId, okUserIDList)
		} else {
			for _, userID := range okUserIDList {
				if err := rocksCache.DelJoinedSuperGroupIDListFromCache(userID); err != nil {
					log.NewWarn(req.OperationID, utils.GetSelfFuncName(), userID, err.Error())
				}
			}
			go func() {
				for _, v := range okUserIDList {
					chat.SuperGroupNotification(req.OperationID, v, v)
				}
			}()
		}
		return resp, nil
	} else {
		log.NewInfo(req.OperationID, "rpc CreateGroup return ", resp.String())
		return resp, nil
	}
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	log.NewInfo(req.OperationID, "GetJoinedGroupList, args ", req.String())
	if !token_verify.CheckAccess(req.OpUserID, req.FromUserID) {
		log.NewError(req.OperationID, "CheckAccess false ", req.OpUserID, req.FromUserID)
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	joinedGroupList, err := rocksCache.GetJoinedGroupIDListFromCache(req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetJoinedGroupIDListFromCache failed", err.Error(), req.FromUserID)
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "joinedGroupList: ", joinedGroupList)
	var resp pbGroup.GetJoinedGroupListResp
	for _, v := range joinedGroupList {
		var groupNode open_im_sdk.GroupInfo
		num, err := imdb.GetGroupMemberNumByGroupID(v)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), v)
			continue
		}
		owner, err2 := imdb.GetGroupOwnerInfoByGroupID(v)
		if err2 != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err2.Error(), v)
			continue
		}
		group, err := rocksCache.GetGroupInfoFromCache(v)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), v)
			continue
		}
		if group.Status == constant.GroupStatusDismissed {
			log.NewError(req.OperationID, "constant.GroupStatusDismissed ", group)
			continue
		}
		utils.CopyStructFields(&groupNode, group)
		groupNode.CreateTime = uint32(group.CreateTime.Unix())
		groupNode.NotificationUpdateTime = uint32(group.NotificationUpdateTime.Unix())
		if group.NotificationUpdateTime.Unix() < 0 {
			groupNode.NotificationUpdateTime = 0
		}

		groupNode.MemberCount = uint32(num)
		groupNode.OwnerUserID = owner.UserID
		resp.GroupList = append(resp.GroupList, &groupNode)
		log.NewDebug(req.OperationID, "joinedGroup ", groupNode)
	}
	log.NewInfo(req.OperationID, "GetJoinedGroupList rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	log.NewInfo(req.OperationID, "InviteUserToGroup args ", req.String())
	if !imdb.IsExistGroupMember(req.GroupID, req.OpUserID) && !token_verify.IsManagerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "no permission InviteUserToGroup ", req.GroupID, req.OpUserID)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}, nil
	}
	var resp pbGroup.InviteUserToGroupResp
	if groupInfo.NeedVerification == constant.AllNeedVerification &&
		!imdb.IsGroupOwnerAdmin(req.GroupID, req.OpUserID) && !token_verify.IsManagerUserID(req.OpUserID) {
		var resp pbGroup.InviteUserToGroupResp
		joinReq := pbGroup.JoinGroupReq{}
		for _, v := range req.InvitedUserIDList {
			var groupRequest db.GroupRequest
			groupRequest.UserID = v
			groupRequest.GroupID = req.GroupID
			groupRequest.JoinSource = constant.JoinByInvitation
			groupRequest.InviterUserID = req.OpUserID
			err = imdb.InsertIntoGroupRequest(groupRequest)
			if err != nil {
				var resultNode pbGroup.Id2Result
				resultNode.Result = -1
				resultNode.UserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)

				continue
				log.NewError(req.OperationID, "InsertIntoGroupRequest failed ", err.Error(), groupRequest)
				//	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			} else {
				var resultNode pbGroup.Id2Result
				resultNode.Result = 0
				resultNode.UserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				joinReq.GroupID = req.GroupID
				joinReq.OperationID = req.OperationID
				joinReq.OpUserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				chat.JoinGroupApplicationNotification(&joinReq)
			}
		}
		log.NewInfo(req.OperationID, "InviteUserToGroup rpc return ", resp)
		return &resp, nil
	}
	//
	//from User:  invite: applicant
	//to user:  invite: invited
	var okUserIDList []string
	if groupInfo.GroupType != constant.SuperGroup {
		for _, v := range req.InvitedUserIDList {
			var resultNode pbGroup.Id2Result
			resultNode.UserID = v
			resultNode.Result = 0
			toUserInfo, err := imdb.GetUserByUserID(v)
			if err != nil {
				log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), v)
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}

			if imdb.IsExistGroupMember(req.GroupID, v) {
				log.NewError(req.OperationID, "IsExistGroupMember ", req.GroupID, v)
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}
			var toInsertInfo db.GroupMember
			utils.CopyStructFields(&toInsertInfo, toUserInfo)
			toInsertInfo.GroupID = req.GroupID
			toInsertInfo.RoleLevel = constant.GroupOrdinaryUsers
			toInsertInfo.OperatorUserID = req.OpUserID
			toInsertInfo.InviterUserID = req.OpUserID
			toInsertInfo.JoinSource = constant.JoinByInvitation
			err = imdb.InsertIntoGroupMember(toInsertInfo)
			if err != nil {
				log.NewError(req.OperationID, "InsertIntoGroupMember failed ", req.GroupID, toUserInfo.UserID, toUserInfo.Nickname, toUserInfo.FaceURL)
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}
			okUserIDList = append(okUserIDList, v)
			err = db.DB.AddGroupMember(req.GroupID, toUserInfo.UserID)
			if err != nil {
				log.NewError(req.OperationID, "AddGroupMember failed ", err.Error(), req.GroupID, toUserInfo.UserID)
			}
			resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
		}
		var haveConUserID []string
		conversations, err := imdb.GetConversationsByConversationIDMultipleOwner(okUserIDList, utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType))
		if err != nil {
			log.NewError(req.OperationID, "GetConversationsByConversationIDMultipleOwner failed ", err.Error(), req.GroupID, constant.GroupChatType)
		}
		for _, v := range conversations {
			haveConUserID = append(haveConUserID, v.OwnerUserID)
		}
		var reqPb pbUser.SetConversationReq
		var c pbUser.Conversation
		for _, v := range conversations {
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = v.OwnerUserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.RecvMsgOpt = v.RecvMsgOpt
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsPinned = v.IsPinned
			c.AttachedInfo = v.AttachedInfo
			c.IsPrivateChat = v.IsPrivateChat
			c.GroupAtType = v.GroupAtType
			c.IsNotInGroup = false
			c.Ex = v.Ex
			reqPb.Conversation = &c
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetConn == nil"
				log.NewError(req.OperationID, errMsg)
				return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}, nil
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v.OwnerUserID)
			} else {
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v.OwnerUserID)
			}
		}
		for _, v := range utils.DifferenceString(haveConUserID, okUserIDList) {
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = v
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = false
			reqPb.Conversation = &c
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetConn == nil"
				log.NewError(req.OperationID, errMsg)
				return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}, nil
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v)
			} else {
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v)
			}
		}
	} else {
		okUserIDList = req.InvitedUserIDList
		if err := db.DB.AddUserToSuperGroup(req.GroupID, req.InvitedUserIDList); err != nil {
			log.NewError(req.OperationID, "AddUserToSuperGroup failed ", req.GroupID, err)
			return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
		}
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}, nil
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
		GroupID:     req.GroupID,
		OperationID: req.OperationID,
	})
	if err != nil {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}

	if groupInfo.GroupType != constant.SuperGroup {
		for _, userID := range okUserIDList {
			err = rocksCache.DelJoinedGroupIDListFromCache(userID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), userID)
			}
		}
		if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		}
		chat.MemberInvitedNotification(req.OperationID, req.GroupID, req.OpUserID, req.Reason, okUserIDList)
	} else {
		for _, v := range req.InvitedUserIDList {
			if err := rocksCache.DelJoinedSuperGroupIDListFromCache(v); err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			}
		}
		go func() {
			for _, v := range req.InvitedUserIDList {
				chat.SuperGroupNotification(req.OperationID, v, v)
			}
		}()
	}

	log.NewInfo(req.OperationID, "InviteUserToGroup rpc return ", resp)
	return &resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	log.NewInfo(req.OperationID, "GetGroupAllMember, args ", req.String())
	var resp pbGroup.GetGroupAllMemberResp
	//groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	groupInfo, err := rocksCache.GetGroupInfoFromCache(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = constant.ErrDB.ErrMsg
		return &resp, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := rocksCache.GetAllGroupMembersInfoFromCache(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = constant.ErrDB.ErrMsg
			return &resp, nil
		}
		for _, v := range memberList {
			var node open_im_sdk.GroupMemberFullInfo
			cp.GroupMemberDBCopyOpenIM(&node, v)
			resp.MemberList = append(resp.MemberList, &node)
		}
	}
	log.NewInfo(req.OperationID, "GetGroupAllMember rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	log.NewInfo(req.OperationID, "GetGroupMemberList args ", req.String())
	var resp pbGroup.GetGroupMemberListResp
	memberList, err := imdb.GetGroupMemberByGroupID(req.GroupID, req.Filter, req.NextSeq, 30)
	if err != nil {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = constant.ErrDB.ErrMsg
		log.NewError(req.OperationID, "GetGroupMemberByGroupId failed,", req.GroupID, req.Filter, req.NextSeq, 30)
		return &resp, nil
	}

	for _, v := range memberList {
		var node open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&node, &v)
		resp.MemberList = append(resp.MemberList, &node)
	}
	//db operate  get db sorted by join time
	if int32(len(memberList)) < 30 {
		resp.NextSeq = 0
	} else {
		resp.NextSeq = req.NextSeq + int32(len(memberList))
	}
	resp.ErrCode = 0
	log.NewInfo(req.OperationID, "GetGroupMemberList rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) getGroupUserLevel(groupID, userID string) (int, error) {
	opFlag := 0
	if !token_verify.IsManagerUserID(userID) {
		opInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
		if err != nil {
			return opFlag, utils.Wrap(err, "")
		}
		if opInfo.RoleLevel == constant.GroupOrdinaryUsers {
			opFlag = 0
		} else if opInfo.RoleLevel == constant.GroupOwner {
			opFlag = 2 //owner
		} else {
			opFlag = 3 //admin
		}
	} else {
		opFlag = 1 //app manager
	}
	return opFlag, nil
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupInfoByGroupID", req.GroupID, err.Error())
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	var okUserIDList []string
	var resp pbGroup.KickGroupMemberResp
	if groupInfo.GroupType != constant.SuperGroup {
		opFlag := 0
		if !token_verify.IsManagerUserID(req.OpUserID) {
			opInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.OpUserID)
			if err != nil {
				errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID  failed " + err.Error() + req.GroupID + req.OpUserID
				log.Error(req.OperationID, errMsg)
				return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
			}
			if opInfo.RoleLevel == constant.GroupOrdinaryUsers {
				errMsg := req.OperationID + " opInfo.RoleLevel == constant.GroupOrdinaryUsers " + opInfo.UserID + opInfo.GroupID
				log.Error(req.OperationID, errMsg)
				return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
			} else if opInfo.RoleLevel == constant.GroupOwner {
				opFlag = 2 //owner
			} else {
				opFlag = 3 //admin
			}
		} else {
			opFlag = 1 //app manager
		}

		//op is group owner?
		if len(req.KickedUserIDList) == 0 {
			log.NewError(req.OperationID, "failed, kick list 0")
			return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}, nil
		}

		//remove
		for _, v := range req.KickedUserIDList {
			kickedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, v)
			if err != nil {
				log.NewError(req.OperationID, " GetGroupMemberInfoByGroupIDAndUserID failed ", req.GroupID, v, err.Error())
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}

			if kickedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
				log.Error(req.OperationID, "is constant.GroupAdmin, can't kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}
			if kickedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
				log.NewDebug(req.OperationID, "is constant.GroupOwner, can't kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}

			err = imdb.RemoveGroupMember(req.GroupID, v)
			if err != nil {
				log.NewError(req.OperationID, "RemoveGroupMember failed ", err.Error(), req.GroupID, v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
			} else {
				log.NewDebug(req.OperationID, "kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: 0})
				okUserIDList = append(okUserIDList, v)
			}
		}
		var reqPb pbUser.SetConversationReq
		var c pbUser.Conversation
		for _, v := range okUserIDList {
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = v
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = true
			reqPb.Conversation = &c
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetConn == nil"
				log.NewError(req.OperationID, errMsg)
				resp.ErrCode = constant.ErrInternal.ErrCode
				resp.ErrMsg = errMsg
				return &resp, nil
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v)
			} else {
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v)
			}
		}
	} else {
		okUserIDList = req.KickedUserIDList
		if err := db.DB.RemoverUserFromSuperGroup(req.GroupID, okUserIDList); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, req.KickedUserIDList, err.Error())
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = constant.ErrDB.ErrMsg
			return &resp, nil
		}
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}, nil
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
		GroupID:     req.GroupID,
		OperationID: req.OperationID,
	})
	if err != nil {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}

	if groupInfo.GroupType != constant.SuperGroup {
		for _, userID := range okUserIDList {
			err = rocksCache.DelJoinedGroupIDListFromCache(userID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), userID)
			}
		}
		if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		}
		chat.MemberKickedNotification(req, okUserIDList)
	} else {
		for _, userID := range okUserIDList {
			err = rocksCache.DelJoinedSuperGroupIDListFromCache(userID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), userID)
			}
		}
		go func() {
			for _, v := range req.KickedUserIDList {
				chat.SuperGroupNotification(req.OperationID, v, v)
			}
		}()

	}

	log.NewInfo(req.OperationID, "GetGroupMemberList rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	log.NewInfo(req.OperationID, "GetGroupMembersInfo args ", req.String())
	var resp pbGroup.GetGroupMembersInfoResp
	resp.MemberList = []*open_im_sdk.GroupMemberFullInfo{}
	groupMembers, err := rocksCache.GetAllGroupMembersInfoFromCache(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, err.Error())
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = constant.ErrDB.ErrMsg
		return &resp, nil
	}
	for _, member := range groupMembers {
		if utils.IsContain(member.UserID, req.MemberList) {
			var memberNode open_im_sdk.GroupMemberFullInfo
			utils.CopyStructFields(&memberNode, member)
			memberNode.JoinTime = int32(member.JoinTime.Unix())
			resp.MemberList = append(resp.MemberList, &memberNode)
		}
	}

	resp.ErrCode = 0
	log.NewInfo(req.OperationID, "GetGroupMembersInfo rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) GetGroupApplicationList(_ context.Context, req *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	log.NewInfo(req.OperationID, "GetGroupApplicationList args ", req.String())
	reply, err := imdb.GetGroupApplicationList(req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupApplicationList failed ", err.Error(), req.FromUserID)
		return &pbGroup.GetGroupApplicationListResp{ErrCode: 701, ErrMsg: "GetGroupApplicationList failed"}, nil
	}

	log.NewDebug(req.OperationID, "GetGroupApplicationList reply ", reply)
	resp := pbGroup.GetGroupApplicationListResp{}
	for _, v := range reply {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		group, err := imdb.GetGroupInfoByGroupID(v.GroupID)
		if err != nil {
			log.Error(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), v.GroupID)
			continue
		}
		if group.Status == constant.GroupStatusDismissed {
			log.Debug(req.OperationID, "group constant.GroupStatusDismissed  ", group.GroupID)
			continue
		}
		user, err := imdb.GetUserByUserID(v.UserID)
		if err != nil {
			log.Error(req.OperationID, "GetUserByUserID failed ", err.Error(), v.UserID)
			continue
		}

		cp.GroupRequestDBCopyOpenIM(&node, &v)
		cp.UserDBCopyOpenIMPublicUser(node.UserInfo, user)
		cp.GroupDBCopyOpenIM(node.GroupInfo, group)
		log.NewDebug(req.OperationID, "node ", node, "v ", v)
		resp.GroupRequestList = append(resp.GroupRequestList, &node)
	}
	log.NewInfo(req.OperationID, "GetGroupMembersInfo rpc return ", resp)
	return &resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	log.NewInfo(req.OperationID, "GetGroupsInfo args ", req.String())
	groupsInfoList := make([]*open_im_sdk.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromRedis, err := rocksCache.GetGroupInfoFromCache(groupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), groupID)
			continue
		}
		var groupInfo open_im_sdk.GroupInfo
		cp.GroupDBCopyOpenIM(&groupInfo, groupInfoFromRedis)
		//groupInfo.NeedVerification

		groupInfo.NeedVerification = groupInfoFromRedis.NeedVerification
		groupsInfoList = append(groupsInfoList, &groupInfo)
	}

	resp := pbGroup.GetGroupsInfoResp{GroupInfoList: groupsInfoList}
	log.NewInfo(req.OperationID, "GetGroupsInfo rpc return  ", resp.String())
	return &resp, nil
}

func (s *groupServer) GroupApplicationResponse(_ context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	log.NewInfo(req.OperationID, "GroupApplicationResponse args ", req.String())

	groupRequest := db.GroupRequest{}
	utils.CopyStructFields(&groupRequest, req)
	groupRequest.UserID = req.FromUserID
	groupRequest.HandleUserID = req.OpUserID
	groupRequest.HandledTime = time.Now()
	if !token_verify.IsManagerUserID(req.OpUserID) && !imdb.IsGroupOwnerAdmin(req.GroupID, req.OpUserID) {
		log.NewError(req.OperationID, "IsManagerUserID IsGroupOwnerAdmin false ", req.GroupID, req.OpUserID)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	err := imdb.UpdateGroupRequest(groupRequest)
	if err != nil {
		log.NewError(req.OperationID, "GroupApplicationResponse failed ", err.Error(), groupRequest)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	request, err := imdb.GetGroupRequestByGroupIDAndUserID(req.GroupID, req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, "GroupApplicationResponse failed ", err.Error(), req.GroupID, req.FromUserID)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	if req.HandleResult == constant.GroupResponseAgree {
		user, err := imdb.GetUserByUserID(req.FromUserID)
		if err != nil {
			log.NewError(req.OperationID, "GroupApplicationResponse failed ", err.Error(), req.FromUserID)
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		member := db.GroupMember{}
		member.GroupID = req.GroupID
		member.UserID = req.FromUserID
		member.RoleLevel = constant.GroupOrdinaryUsers
		member.OperatorUserID = req.OpUserID
		member.FaceURL = user.FaceURL
		member.Nickname = user.Nickname
		member.JoinSource = request.JoinSource
		member.InviterUserID = request.InviterUserID
		err = imdb.InsertIntoGroupMember(member)
		if err != nil {
			log.NewError(req.OperationID, "GroupApplicationResponse failed ", err.Error(), member)
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		var reqPb pbUser.SetConversationReq
		reqPb.OperationID = req.OperationID
		var c pbUser.Conversation
		conversation, err := imdb.GetConversation(req.FromUserID, utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType))
		if err != nil {
			c.OwnerUserID = req.FromUserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = false
		} else {
			c.OwnerUserID = conversation.OwnerUserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.RecvMsgOpt = conversation.RecvMsgOpt
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsPinned = conversation.IsPinned
			c.AttachedInfo = conversation.AttachedInfo
			c.IsPrivateChat = conversation.IsPrivateChat
			c.GroupAtType = conversation.GroupAtType
			c.IsNotInGroup = false
			c.Ex = conversation.Ex
		}
		reqPb.Conversation = &c
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
		}
		client := pbUser.NewUserClient(etcdConn)
		respPb, err := client.SetConversation(context.Background(), &reqPb)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error())
		} else {
			log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String())
		}

		etcdCacheConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
		if etcdCacheConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
		}
		cacheClient := pbCache.NewCacheClient(etcdCacheConn)
		cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{OperationID: req.OperationID, GroupID: req.GroupID})
		if err != nil {
			log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
			return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}

		group, err := rocksCache.GetGroupInfoFromCache(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, err.Error())
		}
		if group != nil {
			if group.GroupType != constant.SuperGroup {
				if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
					log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
				}
			}
		}

		chat.GroupApplicationAcceptedNotification(req)
		chat.MemberEnterNotification(req)
	} else if req.HandleResult == constant.GroupResponseRefuse {
		chat.GroupApplicationRejectedNotification(req)
	} else {
		log.Error(req.OperationID, "HandleResult failed ", req.HandleResult)
		return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
	}

	log.NewInfo(req.OperationID, "rpc GroupApplicationResponse return ", pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{}})
	return &pbGroup.GroupApplicationResponseResp{CommonResp: &pbGroup.CommonResp{}}, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.JoinGroupResp, error) {
	log.NewInfo(req.OperationID, "JoinGroup args ", req.String())
	_, err := imdb.GetUserByUserID(req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), req.OpUserID)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	if groupInfo.NeedVerification == constant.Directly {
		if groupInfo.GroupType != constant.SuperGroup {
			us, err := imdb.GetUserByUserID(req.OpUserID)
			if err != nil {
				log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), req.OpUserID)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			//to group member
			groupMember := db.GroupMember{GroupID: req.GroupID, RoleLevel: constant.GroupOrdinaryUsers, OperatorUserID: req.OpUserID}
			utils.CopyStructFields(&groupMember, us)
			err = imdb.InsertIntoGroupMember(groupMember)
			if err != nil {
				log.NewError(req.OperationID, "InsertIntoGroupMember failed ", err.Error(), groupMember)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetConn == nil"
				log.NewError(req.OperationID, errMsg)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: constant.ErrInternal.ErrMsg}}, nil
			}
			cacheClient := pbCache.NewCacheClient(etcdConn)
			cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
				GroupID:     req.GroupID,
				OperationID: req.OperationID,
			})
			if err != nil {
				log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			if cacheResp.CommonResp.ErrCode != 0 {
				log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			//for _, userID := range okUserIDList {
			//	err = rocksCache.DelJoinedGroupIDListFromCache(userID)
			//	if err != nil {
			//		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), userID)
			//	}
			//}
			err = rocksCache.DelJoinedGroupIDListFromCache(req.OpUserID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			}

			err = rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			}

			chat.MemberEnterDirectlyNotification(req.GroupID, req.OpUserID, req.OperationID)
			log.NewInfo(req.OperationID, "JoinGroup rpc return ")
			return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
		} else {
			log.Error(req.OperationID, "JoinGroup rpc failed, group type:  ", groupInfo.GroupType, "not support directly")
			return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
		}
	}

	var groupRequest db.GroupRequest
	groupRequest.UserID = req.OpUserID
	groupRequest.ReqMsg = req.ReqMessage
	groupRequest.GroupID = req.GroupID
	groupRequest.JoinSource = req.JoinSource
	err = imdb.InsertIntoGroupRequest(groupRequest)
	if err != nil {
		log.NewError(req.OperationID, "InsertIntoGroupRequest failed ", err.Error(), groupRequest)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	//_, err = imdb.GetGroupMemberListByGroupIDAndRoleLevel(req.GroupID, constant.GroupOwner)
	//if err != nil {
	//	log.NewError(req.OperationID, "GetGroupMemberListByGroupIDAndRoleLevel failed ", err.Error(), req.GroupID, constant.GroupOwner)
	//	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
	//}

	chat.JoinGroupApplicationNotification(req)
	log.NewInfo(req.OperationID, "JoinGroup rpc return ")
	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	log.NewInfo(req.OperationID, "QuitGroup args ", req.String())
	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "ReduceGroupMemberFromCache rpc call failed ", err.Error())
		return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		_, err = imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.OpUserID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupMemberInfoByGroupIDAndUserID failed ", err.Error(), req.GroupID, req.OpUserID)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}

		err = imdb.DeleteGroupMemberByGroupIDAndUserID(req.GroupID, req.OpUserID)
		if err != nil {
			log.NewError(req.OperationID, "DeleteGroupMemberByGroupIdAndUserId failed ", err.Error(), req.GroupID, req.OpUserID)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}

		err = db.DB.DelGroupMember(req.GroupID, req.OpUserID)
		if err != nil {
			log.NewError(req.OperationID, "DelGroupMember failed ", req.GroupID, req.OpUserID)
			//	return &pbGroup.CommonResp{ErrorCode: constant.ErrQuitGroup.ErrCode, ErrorMsg: constant.ErrQuitGroup.ErrMsg}, nil
		}
		//modify quitter conversation info
		var reqPb pbUser.SetConversationReq
		var c pbUser.Conversation
		reqPb.OperationID = req.OperationID
		c.OwnerUserID = req.OpUserID
		c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
		c.ConversationType = constant.GroupChatType
		c.GroupID = req.GroupID
		c.IsNotInGroup = true
		reqPb.Conversation = &c
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
		}
		client := pbUser.NewUserClient(etcdConn)
		respPb, err := client.SetConversation(context.Background(), &reqPb)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error())
		} else {
			log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String())
		}
	} else {
		okUserIDList := []string{req.OpUserID}
		if err := db.DB.RemoverUserFromSuperGroup(req.GroupID, okUserIDList); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, okUserIDList, err.Error())
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: constant.ErrInternal.ErrMsg}}, nil
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
		GroupID:     req.GroupID,
		OperationID: req.OperationID,
	})
	if err != nil {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
		return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
		return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	if groupInfo.GroupType != constant.SuperGroup {
		if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		}
		if err := rocksCache.DelJoinedGroupIDListFromCache(req.OpUserID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.OpUserID)
		}
		chat.MemberQuitNotification(req)
	} else {
		if err := rocksCache.DelJoinedSuperGroupIDListFromCache(req.OpUserID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.OpUserID)
		}
		chat.SuperGroupNotification(req.OperationID, req.OpUserID, req.OpUserID)
	}
	log.NewInfo(req.OperationID, "rpc QuitGroup return ", pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}})
	return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func hasAccess(req *pbGroup.SetGroupInfoReq) bool {
	if utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
		return true
	}
	groupUserInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupInfoForSet.GroupID, req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMemberInfoByGroupIDAndUserID failed, ", err.Error(), req.GroupInfoForSet.GroupID, req.OpUserID)
		return false

	}
	if groupUserInfo.RoleLevel == constant.GroupOwner || groupUserInfo.RoleLevel == constant.GroupAdmin {
		return true
	}
	return false
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.SetGroupInfoResp, error) {
	log.NewInfo(req.OperationID, "SetGroupInfo args ", req.String())
	if !hasAccess(req) {
		log.NewError(req.OperationID, "no access ", req)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	group, err := imdb.GetGroupInfoByGroupID(req.GroupInfoForSet.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), req.GroupInfoForSet.GroupID)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, http.WrapError(constant.ErrDB)
	}

	if group.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	////bitwise operators: 0001:groupName; 0010:Notification  0100:Introduction; 1000:FaceUrl; 10000:owner
	var changedType int32
	groupName := ""
	notification := ""
	introduction := ""
	faceURL := ""
	if group.GroupName != req.GroupInfoForSet.GroupName && req.GroupInfoForSet.GroupName != "" {
		changedType = 1
		groupName = req.GroupInfoForSet.GroupName
	}
	if group.Notification != req.GroupInfoForSet.Notification && req.GroupInfoForSet.Notification != "" {
		changedType = changedType | (1 << 1)
		notification = req.GroupInfoForSet.Notification
	}
	if group.Introduction != req.GroupInfoForSet.Introduction && req.GroupInfoForSet.Introduction != "" {
		changedType = changedType | (1 << 2)
		introduction = req.GroupInfoForSet.Introduction
	}
	if group.FaceURL != req.GroupInfoForSet.FaceURL && req.GroupInfoForSet.FaceURL != "" {
		changedType = changedType | (1 << 3)
		faceURL = req.GroupInfoForSet.FaceURL
	}

	if req.GroupInfoForSet.NeedVerification != nil {
		changedType = changedType | (1 << 4)
		m := make(map[string]interface{})
		m["need_verification"] = req.GroupInfoForSet.NeedVerification.Value
		if err := imdb.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			log.NewError(req.OperationID, "UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, http.WrapError(constant.ErrDB)
		}
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		changedType = changedType | (1 << 5)
		m := make(map[string]interface{})
		m["look_member_info"] = req.GroupInfoForSet.LookMemberInfo.Value
		if err := imdb.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			log.NewError(req.OperationID, "UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, http.WrapError(constant.ErrDB)
		}
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		changedType = changedType | (1 << 6)
		m := make(map[string]interface{})
		m["apply_member_friend"] = req.GroupInfoForSet.ApplyMemberFriend.Value
		if err := imdb.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			log.NewError(req.OperationID, "UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, http.WrapError(constant.ErrDB)
		}
	}
	//
	//if req.RoleLevel != nil {
	//
	//}
	//only administrators can set group information
	var groupInfo db.Group
	utils.CopyStructFields(&groupInfo, req.GroupInfoForSet)
	if req.GroupInfoForSet.Notification != "" {
		groupInfo.NotificationUserID = req.OpUserID
		groupInfo.NotificationUpdateTime = time.Now()
	}
	err = imdb.SetGroupInfo(groupInfo)
	if err != nil {
		log.NewError(req.OperationID, "SetGroupInfo failed ", err.Error(), groupInfo)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, http.WrapError(constant.ErrDB)
	}
	if err := rocksCache.DelGroupInfoFromCache(req.GroupInfoForSet.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DelGroupInfoFromCache failed ", err.Error(), req.GroupInfoForSet.GroupID)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, http.WrapError(constant.ErrDB)
	}

	log.NewInfo(req.OperationID, "SetGroupInfo rpc return ", pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}})
	if changedType != 0 {
		chat.GroupInfoSetNotification(req.OperationID, req.OpUserID, req.GroupInfoForSet.GroupID, groupName, notification, introduction, faceURL, req.GroupInfoForSet.NeedVerification)
	}
	if req.GroupInfoForSet.Notification != "" {
		//get group member user id
		getGroupMemberIDListFromCacheReq := &pbCache.GetGroupMemberIDListFromCacheReq{OperationID: req.OperationID, GroupID: req.GroupInfoForSet.GroupID}
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, http.WrapError(constant.ErrInternal)
		}
		client := pbCache.NewCacheClient(etcdConn)
		cacheResp, err := client.GetGroupMemberIDListFromCache(context.Background(), getGroupMemberIDListFromCacheReq)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupMemberIDListFromCache rpc call failed ", err.Error())
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			log.NewError(req.OperationID, "GetGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
		}
		var conversationReq pbConversation.ModifyConversationFieldReq

		conversation := pbConversation.Conversation{
			OwnerUserID:      req.OpUserID,
			ConversationID:   utils.GetConversationIDBySessionType(req.GroupInfoForSet.GroupID, constant.GroupChatType),
			ConversationType: constant.GroupChatType,
			GroupID:          req.GroupInfoForSet.GroupID,
		}
		conversationReq.Conversation = &conversation
		conversationReq.OperationID = req.OperationID
		conversationReq.FieldType = constant.FieldGroupAtType
		conversation.GroupAtType = constant.GroupNotification
		conversationReq.UserIDList = cacheResp.UserIDList
		nEtcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImConversationName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, http.WrapError(constant.ErrInternal)
		}
		nClient := pbConversation.NewConversationClient(nEtcdConn)
		conversationReply, err := nClient.ModifyConversationField(context.Background(), &conversationReq)
		if err != nil {
			log.NewError(conversationReq.OperationID, "ModifyConversationField rpc failed, ", conversationReq.String(), err.Error())
		} else if conversationReply.CommonResp.ErrCode != 0 {
			log.NewError(conversationReq.OperationID, "ModifyConversationField rpc failed, ", conversationReq.String(), conversationReply.String())
		}
	}
	return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
}

func (s *groupServer) TransferGroupOwner(_ context.Context, req *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	log.NewInfo(req.OperationID, "TransferGroupOwner ", req.String())

	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	if req.OldOwnerUserID == req.NewOwnerUserID {
		log.NewError(req.OperationID, "same owner ", req.OldOwnerUserID, req.NewOwnerUserID)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
	}
	groupMemberInfo := db.GroupMember{GroupID: req.GroupID, UserID: req.OldOwnerUserID, RoleLevel: constant.GroupOrdinaryUsers}
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		log.NewError(req.OperationID, "UpdateGroupMemberInfo failed ", groupMemberInfo)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	groupMemberInfo = db.GroupMember{GroupID: req.GroupID, UserID: req.NewOwnerUserID, RoleLevel: constant.GroupOwner}
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		log.NewError(req.OperationID, "UpdateGroupMemberInfo failed ", groupMemberInfo)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, err.Error())
	}
	chat.GroupOwnerTransferredNotification(req)
	return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil

}

func (s *groupServer) GetGroupById(_ context.Context, req *pbGroup.GetGroupByIdReq) (*pbGroup.GetGroupByIdResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetGroupByIdResp{CMSGroup: &pbGroup.CMSGroup{
		GroupInfo: &open_im_sdk.GroupInfo{},
	}}
	group, err := imdb.GetGroupById(req.GroupId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupById error", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	resp.CMSGroup.GroupInfo = &open_im_sdk.GroupInfo{
		GroupID:       group.GroupID,
		GroupName:     group.GroupName,
		FaceURL:       group.FaceURL,
		OwnerUserID:   group.CreatorUserID,
		MemberCount:   0,
		Status:        group.Status,
		CreatorUserID: group.CreatorUserID,
		GroupType:     group.GroupType,
		CreateTime:    uint32(group.CreateTime.Unix()),
	}
	groupMember, err := imdb.GetGroupMaster(group.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMaster", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	resp.CMSGroup.GroupMasterName = groupMember.Nickname
	resp.CMSGroup.GroupMasterId = groupMember.UserID
	resp.CMSGroup.GroupInfo.CreatorUserID = group.CreatorUserID
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *groupServer) GetGroup(_ context.Context, req *pbGroup.GetGroupReq) (*pbGroup.GetGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetGroupResp{
		CMSGroups: []*pbGroup.CMSGroup{},
	}
	groups, err := imdb.GetGroupsByName(req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupsByName error", req.String())
		return resp, http.WrapError(constant.ErrDB)
	}
	nums, err := imdb.GetGroupsCountNum(db.Group{GroupName: req.GroupName})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupsCountNum error", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	resp.GroupNums = nums
	resp.Pagination = &open_im_sdk.RequestPagination{
		PageNumber: req.Pagination.PageNumber,
		ShowNumber: req.Pagination.ShowNumber,
	}
	for _, v := range groups {
		groupMember, err := imdb.GetGroupMaster(v.GroupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMaster error", err.Error())
			continue
		}
		resp.CMSGroups = append(resp.CMSGroups, &pbGroup.CMSGroup{
			GroupInfo: &open_im_sdk.GroupInfo{
				GroupID:       v.GroupID,
				GroupName:     v.GroupName,
				FaceURL:       v.FaceURL,
				OwnerUserID:   v.CreatorUserID,
				Status:        v.Status,
				CreatorUserID: v.CreatorUserID,
				CreateTime:    uint32(v.CreateTime.Unix()),
			},
			GroupMasterName: groupMember.Nickname,
			GroupMasterId:   groupMember.UserID,
		})
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *groupServer) GetGroups(_ context.Context, req *pbGroup.GetGroupsReq) (*pbGroup.GetGroupsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetGroups ", req.String())
	resp := &pbGroup.GetGroupsResp{
		CMSGroups:  []*pbGroup.CMSGroup{},
		Pagination: &open_im_sdk.RequestPagination{},
	}
	groups, err := imdb.GetGroups(int(req.Pagination.PageNumber), int(req.Pagination.ShowNumber))
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroups error", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	groupsCountNum, err := imdb.GetGroupsCountNum(db.Group{})
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "groupsCountNum ", groupsCountNum)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupsCountNum", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	resp.GroupNum = int32(groupsCountNum)
	resp.Pagination.PageNumber = req.Pagination.PageNumber
	resp.Pagination.ShowNumber = req.Pagination.ShowNumber
	for _, v := range groups {
		groupMember, err := imdb.GetGroupMaster(v.GroupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMaster failed", err.Error(), v)
			continue
		}
		resp.CMSGroups = append(resp.CMSGroups, &pbGroup.CMSGroup{
			GroupInfo: &open_im_sdk.GroupInfo{
				GroupID:       v.GroupID,
				GroupName:     v.GroupName,
				FaceURL:       v.FaceURL,
				OwnerUserID:   v.CreatorUserID,
				Status:        v.Status,
				CreatorUserID: v.CreatorUserID,
				CreateTime:    uint32(v.CreateTime.Unix()),
			},
			GroupMasterId:   groupMember.UserID,
			GroupMasterName: groupMember.Nickname,
		})
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetGroups ", resp.String())
	return resp, nil
}

func (s *groupServer) OperateGroupStatus(_ context.Context, req *pbGroup.OperateGroupStatusReq) (*pbGroup.OperateGroupStatusResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.String())
	resp := &pbGroup.OperateGroupStatusResp{}
	if err := imdb.OperateGroupStatus(req.GroupId, req.Status); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "OperateGroupStatus", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	return resp, nil
}

func (s *groupServer) DeleteGroup(_ context.Context, req *pbGroup.DeleteGroupReq) (*pbGroup.DeleteGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.String())
	resp := &pbGroup.DeleteGroupResp{}
	if err := imdb.DeleteGroup(req.GroupId); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DeleteGroup error", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	return resp, nil
}

func (s *groupServer) OperateUserRole(_ context.Context, req *pbGroup.OperateUserRoleReq) (*pbGroup.OperateUserRoleResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args:", req.String())
	resp := &pbGroup.OperateUserRoleResp{}
	oldOwnerUserID, err := imdb.GetGroupMaster(req.GroupId)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMaster failed", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		return resp, http.WrapError(constant.ErrInternal)
	}
	client := pbGroup.NewGroupClient(etcdConn)
	var reqPb pbGroup.TransferGroupOwnerReq
	reqPb.OperationID = req.OperationID
	reqPb.NewOwnerUserID = req.UserId
	reqPb.GroupID = req.GroupId
	reqPb.OpUserID = "cms admin"
	reqPb.OldOwnerUserID = oldOwnerUserID.UserID
	reply, err := client.TransferGroupOwner(context.Background(), &reqPb)
	if reply.CommonResp.ErrCode != 0 || err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "TransferGroupOwner rpc failed")
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		}
	}
	return resp, nil
}

func (s *groupServer) GetGroupMembersCMS(_ context.Context, req *pbGroup.GetGroupMembersCMSReq) (*pbGroup.GetGroupMembersCMSResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args:", req.String())
	resp := &pbGroup.GetGroupMembersCMSResp{}
	groupMembers, err := imdb.GetGroupMembersByGroupIdCMS(req.GroupId, req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMembersByGroupIdCMS Error", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	groupMembersCount, err := imdb.GetGroupMembersCount(req.GroupId, req.UserName)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMembersCMS Error", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	log.NewInfo(req.OperationID, groupMembersCount)
	resp.MemberNums = int32(groupMembersCount)
	for _, groupMember := range groupMembers {
		resp.Members = append(resp.Members, &open_im_sdk.GroupMemberFullInfo{
			GroupID:    req.GroupId,
			UserID:     groupMember.UserID,
			RoleLevel:  groupMember.RoleLevel,
			JoinTime:   int32(groupMember.JoinTime.Unix()),
			Nickname:   groupMember.Nickname,
			FaceURL:    groupMember.FaceURL,
			JoinSource: groupMember.JoinSource,
		})
	}
	resp.Pagination = &open_im_sdk.ResponsePagination{
		CurrentPage: req.Pagination.PageNumber,
		ShowNumber:  req.Pagination.ShowNumber,
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp:", resp.String())
	return resp, nil
}

func (s *groupServer) RemoveGroupMembersCMS(_ context.Context, req *pbGroup.RemoveGroupMembersCMSReq) (*pbGroup.RemoveGroupMembersCMSResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args:", req.String())
	resp := &pbGroup.RemoveGroupMembersCMSResp{}
	for _, userId := range req.UserIds {
		err := imdb.RemoveGroupMember(req.GroupId, userId)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			resp.Failed = append(resp.Failed, userId)
		} else {
			resp.Success = append(resp.Success, userId)
		}
	}
	reqKick := &pbGroup.KickGroupMemberReq{
		GroupID:          req.GroupId,
		KickedUserIDList: resp.Success,
		Reason:           "admin kick",
		OperationID:      req.OperationID,
		OpUserID:         req.OpUserId,
	}
	var reqPb pbUser.SetConversationReq
	var c pbUser.Conversation
	for _, v := range resp.Success {
		reqPb.OperationID = req.OperationID
		c.OwnerUserID = v
		c.ConversationID = utils.GetConversationIDBySessionType(req.GroupId, constant.GroupChatType)
		c.ConversationType = constant.GroupChatType
		c.GroupID = req.GroupId
		c.IsNotInGroup = true
		reqPb.Conversation = &c
		etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetConn == nil"
			log.NewError(req.OperationID, errMsg)
			return resp, http.WrapError(constant.ErrInternal)
		}
		client := pbUser.NewUserClient(etcdConn)
		respPb, err := client.SetConversation(context.Background(), &reqPb)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v)
		} else {
			log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v)
		}
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		return resp, http.WrapError(constant.ErrDB)
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
		GroupID:     req.GroupId,
		OperationID: req.OperationID,
	})
	if err != nil {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
		return resp, http.WrapError(constant.ErrDB)
	}
	if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupId); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupId)
	}

	chat.MemberKickedNotification(reqKick, resp.Success)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	return resp, nil
}

func (s *groupServer) AddGroupMembersCMS(_ context.Context, req *pbGroup.AddGroupMembersCMSReq) (*pbGroup.AddGroupMembersCMSResp, error) {
	log.NewInfo(req.OperationId, utils.GetSelfFuncName(), "args:", req.String())
	resp := &pbGroup.AddGroupMembersCMSResp{}
	for _, userId := range req.UserIds {
		if isExist := imdb.IsExistGroupMember(req.GroupId, userId); isExist {
			log.NewError(req.OperationId, utils.GetSelfFuncName(), "user is exist in group", userId, req.GroupId)
			resp.Failed = append(resp.Failed, userId)
			continue
		}
		user, err := imdb.GetUserByUserID(userId)
		if err != nil {
			log.NewError(req.OperationId, utils.GetSelfFuncName(), "GetUserByUserID", err.Error())
			resp.Failed = append(resp.Failed, userId)
			continue
		}
		groupMember := db.GroupMember{
			GroupID:        req.GroupId,
			UserID:         userId,
			Nickname:       user.Nickname,
			FaceURL:        "",
			RoleLevel:      1,
			JoinTime:       time.Time{},
			JoinSource:     constant.JoinByAdmin,
			OperatorUserID: "CmsAdmin",
			Ex:             "",
		}
		if err := imdb.InsertIntoGroupMember(groupMember); err != nil {
			log.NewError(req.OperationId, utils.GetSelfFuncName(), "InsertIntoGroupMember failed", req.String())
			resp.Failed = append(resp.Failed, userId)
		} else {
			resp.Success = append(resp.Success, userId)
		}
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationId)
	if etcdConn == nil {
		errMsg := req.OperationId + "getcdv3.GetConn == nil"
		log.NewError(req.OperationId, errMsg)
		return resp, http.WrapError(constant.ErrDB)
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
		GroupID:     req.GroupId,
		OperationID: req.OperationId,
	})
	if err != nil {
		log.NewError(req.OperationId, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
		return resp, http.WrapError(constant.ErrDB)
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationId, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
		return resp, http.WrapError(constant.ErrDB)
	}
	if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupId); err != nil {
		log.NewError(req.OperationId, utils.GetSelfFuncName(), err.Error(), req.GroupId)
	}

	chat.MemberInvitedNotification(req.OperationId, req.GroupId, req.OpUserId, "admin add you to group", resp.Success)
	return resp, nil
}

func (s *groupServer) GetUserReqApplicationList(_ context.Context, req *pbGroup.GetUserReqApplicationListReq) (*pbGroup.GetUserReqApplicationListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetUserReqApplicationListResp{}
	groupRequests, err := imdb.GetUserReqGroupByUserID(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserReqGroupByUserID failed ", err.Error())
		resp.CommonResp = &pbGroup.CommonResp{
			ErrCode: constant.ErrDB.ErrCode,
			ErrMsg:  constant.ErrDB.ErrMsg,
		}
		return resp, nil
	}
	for _, groupReq := range groupRequests {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		group, err := imdb.GetGroupInfoByGroupID(groupReq.GroupID)
		if err != nil {
			log.Error(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), groupReq.GroupID)
			continue
		}
		user, err := imdb.GetUserByUserID(groupReq.UserID)
		if err != nil {
			log.Error(req.OperationID, "GetUserByUserID failed ", err.Error(), groupReq.UserID)
			continue
		}
		cp.GroupRequestDBCopyOpenIM(&node, &groupReq)
		cp.UserDBCopyOpenIMPublicUser(node.UserInfo, user)
		cp.GroupDBCopyOpenIM(node.GroupInfo, group)
		resp.GroupRequestList = append(resp.GroupRequestList, &node)
	}
	resp.CommonResp = &pbGroup.CommonResp{
		ErrCode: 0,
		ErrMsg:  "",
	}
	return resp, nil
}

func (s *groupServer) DismissGroup(ctx context.Context, req *pbGroup.DismissGroupReq) (*pbGroup.DismissGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) && !imdb.IsGroupOwnerAdmin(req.GroupID, req.OpUserID) {
		log.NewError(req.OperationID, "verify failed ", req.OpUserID, req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	err := imdb.OperateGroupStatus(req.GroupID, constant.GroupStatusDismissed)
	if err != nil {
		log.NewError(req.OperationID, "OperateGroupStatus failed ", req.GroupID, constant.GroupStatusDismissed)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := imdb.GetGroupMemberListByGroupID(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupMemberListByGroupID failed,", err.Error(), req.GroupID)
		}
		//modify quitter conversation info
		var reqPb pbUser.SetConversationReq
		var c pbUser.Conversation
		for _, v := range memberList {
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = v.UserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = true
			reqPb.Conversation = &c
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetConn == nil"
				log.NewError(req.OperationID, errMsg)
				return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v.UserID)
			} else {
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v.UserID)
			}
		}
		err = imdb.DeleteGroupMemberByGroupID(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "DeleteGroupMemberByGroupID failed ", req.GroupID)
			return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		chat.GroupDismissedNotification(req)
	} else {
		err = db.DB.DeleteSuperGroup(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "DeleteGroupMemberByGroupID failed ", req.GroupID)
			return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil

		}
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 500, ErrMsg: errMsg}}, nil
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
		GroupID:     req.GroupID,
		OperationID: req.OperationID,
	})
	if err != nil {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, nil
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationID, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: cacheResp.CommonResp.ErrCode, ErrMsg: cacheResp.CommonResp.ErrMsg}}, nil
	}
	if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

//  rpc MuteGroupMember(MuteGroupMemberReq) returns(MuteGroupMemberResp);
//  rpc CancelMuteGroupMember(CancelMuteGroupMemberReq) returns(CancelMuteGroupMemberResp);
//  rpc MuteGroup(MuteGroupReq) returns(MuteGroupResp);
//  rpc CancelMuteGroup(CancelMuteGroupReq) returns(CancelMuteGroupResp);

func (s *groupServer) MuteGroupMember(ctx context.Context, req *pbGroup.MuteGroupMemberReq) (*pbGroup.MuteGroupMemberResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())
	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
	if err != nil {
		errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.UserID + err.Error()
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.UserID
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.UserID
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	groupMemberInfo := db.GroupMember{GroupID: req.GroupID, UserID: req.UserID}

	groupMemberInfo.MuteEndTime = time.Unix(int64(time.Now().Second())+int64(req.MutedSeconds), time.Now().UnixNano())
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		log.Error(req.OperationID, "UpdateGroupMemberInfo failed ", err.Error(), groupMemberInfo)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
	}
	chat.GroupMemberMutedNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, req.MutedSeconds)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbGroup.CancelMuteGroupMemberReq) (*pbGroup.CancelMuteGroupMemberResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
	if err != nil {
		errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.UserID + err.Error()
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.UserID
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.UserID
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	groupMemberInfo := db.GroupMember{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(0, 0)
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		log.Error(req.OperationID, "UpdateGroupMemberInfo failed ", err.Error(), groupMemberInfo)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
	}
	chat.GroupMemberCancelMutedNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) MuteGroup(ctx context.Context, req *pbGroup.MuteGroupReq) (*pbGroup.MuteGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	//mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
	//if err != nil {
	//	errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}

	err = imdb.OperateGroupStatus(req.GroupID, constant.GroupStatusMuted)
	if err != nil {
		log.Error(req.OperationID, "OperateGroupStatus failed ", err.Error(), req.GroupID, constant.GroupStatusMuted)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if err := rocksCache.DelGroupInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
	}
	chat.GroupMutedNotification(req.OperationID, req.OpUserID, req.GroupID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) CancelMuteGroup(ctx context.Context, req *pbGroup.CancelMuteGroupReq) (*pbGroup.CancelMuteGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	//mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.)
	//if err != nil {
	//	errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	log.Debug(req.OperationID, "UpdateGroupInfoDefaultZero ", req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	err = imdb.UpdateGroupInfoDefaultZero(req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	if err != nil {
		log.Error(req.OperationID, "UpdateGroupInfoDefaultZero failed ", err.Error(), req.GroupID)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if err := rocksCache.DelGroupInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
	}
	chat.GroupCancelMutedNotification(req.OperationID, req.OpUserID, req.GroupID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) SetGroupMemberNickname(ctx context.Context, req *pbGroup.SetGroupMemberNicknameReq) (*pbGroup.SetGroupMemberNicknameResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())
	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " verify failed " + req.OpUserID + req.GroupID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	groupMemberInfo := db.GroupMember{}
	groupMemberInfo.UserID = req.UserID
	groupMemberInfo.GroupID = req.GroupID
	if req.Nickname == "" {
		userNickname, err := imdb.GetUserNameByUserID(groupMemberInfo.UserID)
		if err != nil {
			errMsg := req.OperationID + " GetUserNameByUserID failed " + err.Error()
			log.Error(req.OperationID, errMsg)
			return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		groupMemberInfo.Nickname = userNickname
	} else {
		groupMemberInfo.Nickname = req.Nickname
	}
	err := imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		errMsg := req.OperationID + " UpdateGroupMemberInfo failed " + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
	}
	chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) SetGroupMemberInfo(ctx context.Context, req *pbGroup.SetGroupMemberInfoReq) (resp *pbGroup.SetGroupMemberInfoResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbGroup.SetGroupMemberInfoResp{CommonResp: &pbGroup.CommonResp{}}
	groupMember := db.GroupMember{
		GroupID: req.GroupID,
		UserID:  req.UserID,
	}
	m := make(map[string]interface{})
	if req.RoleLevel != nil {
		m["role_level"] = req.RoleLevel.Value
	}
	if req.FaceURL != nil {
		m["user_group_face_url"] = req.FaceURL.Value
	}
	if req.Nickname != nil {
		m["nickname"] = req.Nickname.Value
	}
	if req.Ex != nil {
		m["ex"] = req.Ex.Value
	}
	err = imdb.UpdateGroupMemberInfoByMap(groupMember, m)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetGroupMemberInfo failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + ":" + err.Error()
		return resp, nil
	}
	if err := rocksCache.DelAllGroupMembersInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
	}
	if req.RoleLevel != nil {
		switch req.RoleLevel.Value {
		case constant.GroupOrdinaryUsers:
			//msg.GroupMemberRoleLevelChangeNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, constant.GroupMemberSetToOrdinaryUserNotification)
			chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
		case constant.GroupAdmin, constant.GroupOwner:
			//msg.GroupMemberRoleLevelChangeNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, constant.GroupMemberSetToAdminNotification)
			chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
		}
	} else {
		chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

package group

import (
	"Open_IM/internal/push/content_struct"
	"Open_IM/internal/push/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	pbChat "Open_IM/pkg/proto/chat"
	"encoding/json"

	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"

	"fmt"
)

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.GetJoinedGroupListResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	log.Info(claims.UID, req.OperationID, "recv req: ", req.String())

	joinedGroupList, err := imdb.GetJoinedGroupIdListByMemberId(claims.UID)
	if err != nil {
		log.Error(claims.UID, req.OperationID, "GetJoinedGroupIdListByMemberId failed, err: ", err.Error())
		return &pbGroup.GetJoinedGroupListResp{ErrorCode: config.ErrParam.ErrCode, ErrorMsg: config.ErrParam.ErrMsg}, nil
	}

	var resp pbGroup.GetJoinedGroupListResp

	for _, v := range joinedGroupList {
		var groupNode pbGroup.GroupInfo
		num := imdb.GetGroupMemberNumByGroupId(v.GroupId)
		owner := imdb.GetGroupOwnerByGroupId(v.GroupId)
		group, err := imdb.FindGroupInfoByGroupId(v.GroupId)
		if num > 0 && owner != "" && err == nil {
			groupNode.GroupId = v.GroupId
			groupNode.FaceUrl = group.FaceUrl
			groupNode.CreateTime = uint64(group.CreateTime.Unix())
			groupNode.GroupName = group.Name
			groupNode.Introduction = group.Introduction
			groupNode.Notification = group.Notification
			groupNode.OwnerId = owner
			groupNode.MemberCount = uint32(int32(num))
			resp.GroupList = append(resp.GroupList, &groupNode)
		}
		log.Info(claims.UID, req.OperationID, "member num: ", num, "owner: ", owner)
	}
	resp.ErrorCode = 0
	return &resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.InviteUserToGroupResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	log.Info(claims.UID, req.OperationID, "recv req: ", req.String())
	//	if !imdb.IsExistGroupMember(req.GroupID, claims.UID) &&  claims.UID != config.Config.AppManagerUid

	if !imdb.IsExistGroupMember(req.GroupID, claims.UID) && !utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		log.Error(req.Token, req.OperationID, "err= invite user not in group")
		return &pbGroup.InviteUserToGroupResp{ErrorCode: config.ErrAccess.ErrCode, ErrorMsg: config.ErrAccess.ErrMsg}, nil
	}

	groupInfoFromMysql, err := imdb.FindGroupInfoByGroupId(req.GroupID)
	if err != nil || groupInfoFromMysql == nil {
		log.NewError(req.OperationID, "get group info error", req.GroupID, req.UidList)
		return &pbGroup.InviteUserToGroupResp{ErrorCode: config.ErrAccess.ErrCode, ErrorMsg: config.ErrAccess.ErrMsg}, nil
	}

	//
	//from User:  invite: applicant
	//to user:  invite: invited
	//to application
	var resp pbGroup.InviteUserToGroupResp
	/*
		fromUserInfo, err := imdb.FindUserByUID(claims.UID)
		if err != nil {
			log.Error(claims.UID, req.OperationID, "FindUserByUID failed, err: ", err.Error())
			return &pbGroup.InviteUserToGroupResp{ErrorCode: config.ErrParam.ErrCode, ErrorMsg: config.ErrParam.ErrMsg}, nil
		}*/
	var nicknameList string
	for _, v := range req.UidList {
		var resultNode pbGroup.Id2Result
		resultNode.UId = v
		resultNode.Result = 0
		toUserInfo, err := imdb.FindUserByUID(v)
		if err != nil {
			log.Error(v, req.OperationID, "FindUserByUID failed, err: ", err.Error())
			resultNode.Result = -1
			resp.Id2Result = append(resp.Id2Result, &resultNode)
			continue
		}

		if imdb.IsExistGroupMember(req.GroupID, v) {
			log.Error(v, req.OperationID, "user has already in group")
			resultNode.Result = -1
			resp.Id2Result = append(resp.Id2Result, &resultNode)
			continue
		}

		err = imdb.InsertGroupMember(req.GroupID, toUserInfo.UID, toUserInfo.Name, toUserInfo.Icon, 0)
		if err != nil {
			log.Error(v, req.OperationID, "InsertGroupMember failed, ", err.Error(), "params: ",
				req.GroupID, toUserInfo.UID, toUserInfo.Name, toUserInfo.Icon)
			resultNode.Result = -1
			resp.Id2Result = append(resp.Id2Result, &resultNode)
			continue
		}
		err = db.DB.AddGroupMember(req.GroupID, toUserInfo.UID)
		if err != nil {
			log.Error("", "", "add mongo group member failed, db.DB.AddGroupMember fail [err: %s]", err.Error())
		}
		nicknameList = nicknameList + toUserInfo.Name + " "
		resp.Id2Result = append(resp.Id2Result, &resultNode)
	}
	resp.ErrorCode = 0
	resp.ErrorMsg = "ok"

	//if claims.UID == config.Config.AppManagerUid
	if utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		m, _ := imdb.FindUserByUID(claims.UID)
		var iu inviteUserToGroupReq
		iu.GroupID = req.GroupID
		iu.OperationID = req.OperationID
		iu.Reason = req.Reason
		iu.UidList = req.UidList
		n := content_struct.NotificationContent{1, nicknameList + "  invited into the group chat by " + m.Name, iu.ContentToString()}
		logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
			SendID:      claims.UID,
			RecvID:      req.GroupID,
			Content:     n.ContentToString(),
			SendTime:    utils.GetCurrentTimestampByNano(),
			MsgFrom:     constant.UserMsgType,
			ContentType: constant.InviteUserToGroupTip,
			SessionType: constant.GroupChatType,
			OperationID: req.OperationID,
		})
	}

	return &resp, nil
}

type inviteUserToGroupReq struct {
	GroupID     string   `json:"groupID"`
	UidList     []string `json:"uidList"`
	Reason      string   `json:"reason"`
	OperationID string   `json:"operationID"`
}

func (c *inviteUserToGroupReq) ContentToString() string {
	data, _ := json.Marshal(c)
	dataString := string(data)
	return dataString
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	//claims, err := utils.ParseToken(req.Token)
	//if err != nil {
	//	log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
	//	if req.Token != config.Config.Secret {
	//		return &pbGroup.GetGroupAllMemberResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	//	}
	//}

	var resp pbGroup.GetGroupAllMemberResp
	resp.ErrorCode = 0
	memberList, err := imdb.FindGroupMemberListByGroupId(req.GroupID)
	if err != nil {
		resp.ErrorCode = config.ErrDb.ErrCode
		resp.ErrorMsg = err.Error()
		log.NewError(req.OperationID, "FindGroupMemberListByGroupId failed,", err.Error(), req.GroupID)
		return &resp, nil
	}

	for _, v := range memberList {
		var node pbGroup.GroupMemberFullInfo
		node.Role = v.AdministratorLevel
		node.NickName = v.NickName
		node.UserId = v.Uid
		node.FaceUrl = v.UserGroupFaceUrl
		node.JoinTime = uint64(v.JoinTime.Unix())
		resp.MemberList = append(resp.MemberList, &node)
	}

	resp.ErrorCode = 0
	return &resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.GetGroupMemberListResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	//	log.Info(claims.UID, req.OperationID, "recv req: ", req.String())
	fmt.Println("req: ", req.GroupID)
	var resp pbGroup.GetGroupMemberListResp
	resp.ErrorCode = 0
	memberList, err := imdb.GetGroupMemberByGroupId(req.GroupID, req.Filter, req.NextSeq, 30)
	if err != nil {
		resp.ErrorCode = config.ErrDb.ErrCode
		resp.ErrorMsg = err.Error()
		log.Error(claims.UID, req.OperationID, "GetGroupMemberByGroupId failed, ", err.Error(), "params: ", req.GroupID, req.Filter, req.NextSeq)
		return &resp, nil
	}

	for _, v := range memberList {
		var node pbGroup.GroupMemberFullInfo
		node.Role = v.AdministratorLevel
		node.NickName = v.NickName
		node.UserId = v.Uid
		//	node.FaceUrl =
		node.JoinTime = uint64(v.JoinTime.Unix())
		resp.MemberList = append(resp.MemberList, &node)
	}
	//db operate  get db sorted by join time
	if int32(len(memberList)) < 30 {
		resp.NextSeq = 0
	} else {
		resp.NextSeq = req.NextSeq + int32(len(memberList))
	}

	resp.ErrorCode = 0
	return &resp, nil
}

type groupMemberFullInfo struct {
	GroupId  string `json:"groupID"`
	UserId   string `json:"userId"`
	Role     int    `json:"role"`
	JoinTime uint64 `json:"joinTime"`
	NickName string `json:"nickName"`
	FaceUrl  string `json:"faceUrl"`
}

type kickGroupMemberApiReq struct {
	GroupID     string                `json:"groupID"`
	UidListInfo []groupMemberFullInfo `json:"uidListInfo"`
	Reason      string                `json:"reason"`
	OperationID string                `json:"operationID"`
}

func (c *kickGroupMemberApiReq) ContentToString() string {
	data, _ := json.Marshal(c)
	dataString := string(data)
	return dataString
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.KickGroupMemberResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	log.Info(claims.UID, req.OperationID, "recv req: ", req.String())

	ownerList, err := imdb.GetOwnerManagerByGroupId(req.GroupID)
	if err != nil {
		log.Error(claims.UID, req.OperationID, req.GroupID, "GetOwnerManagerByGroupId, ", err.Error())
		return &pbGroup.KickGroupMemberResp{ErrorCode: config.ErrParam.ErrCode, ErrorMsg: config.ErrParam.ErrMsg}, nil
	}
	//op is group owner?
	var flag = 0
	for _, v := range ownerList {
		if v.Uid == claims.UID {
			flag = 1
			break
		}
	}
	if flag != 1 {
		//	if claims.UID == config.Config.AppManagerUid {
		if utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
			flag = 1
		}
	}

	if flag != 1 {
		log.Error(claims.UID, req.OperationID, "no access kick")
		return &pbGroup.KickGroupMemberResp{ErrorCode: config.ErrAccess.ErrCode, ErrorMsg: config.ErrAccess.ErrMsg}, nil
	}

	if len(req.UidListInfo) == 0 {
		log.Error(claims.UID, req.OperationID, "kick list 0")
		return &pbGroup.KickGroupMemberResp{ErrorCode: config.ErrParam.ErrCode, ErrorMsg: config.ErrParam.ErrMsg}, nil
	}
	//remove
	var resp pbGroup.KickGroupMemberResp
	for _, v := range req.UidListInfo {
		//owner cant kicked
		if v.UserId == claims.UID {
			log.Error(claims.UID, req.OperationID, v.UserId, "cant kick owner")
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v.UserId, Result: -1})
			continue
		}
		err := imdb.RemoveGroupMember(req.GroupID, v.UserId)
		if err != nil {
			log.Error(claims.UID, req.OperationID, v.UserId, req.GroupID, "RemoveGroupMember failed ", err.Error())
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v.UserId, Result: -1})
		} else {
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v.UserId, Result: 0})
		}

		err = db.DB.DelGroupMember(req.GroupID, v.UserId)
		if err != nil {
			log.Error("", "", "delete mongo group member failed, db.DB.DelGroupMember fail [err: %s]", err.Error())
		}

	}
	var kq kickGroupMemberApiReq

	kq.GroupID = req.GroupID
	kq.OperationID = req.OperationID
	kq.Reason = req.Reason

	var gf groupMemberFullInfo
	for _, v := range req.UidListInfo {
		gf.UserId = v.UserId
		gf.GroupId = req.GroupID
		kq.UidListInfo = append(kq.UidListInfo, gf)
	}

	n := content_struct.NotificationContent{1, req.GroupID, kq.ContentToString()}

	if utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		log.Info("", req.OperationID, claims.UID, req.GroupID)
		logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
			SendID:      claims.UID,
			RecvID:      req.GroupID,
			Content:     n.ContentToString(),
			SendTime:    utils.GetCurrentTimestampByNano(),
			MsgFrom:     constant.UserMsgType,
			ContentType: constant.KickGroupMemberTip,
			SessionType: constant.GroupChatType,
			OperationID: req.OperationID,
		})

		for _, v := range req.UidListInfo {
			log.Info("", req.OperationID, claims.UID, v.UserId)
			logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
				SendID:      claims.UID,
				RecvID:      v.UserId,
				Content:     n.ContentToString(),
				SendTime:    utils.GetCurrentTimestampBySecond(),
				MsgFrom:     constant.UserMsgType,
				ContentType: constant.KickGroupMemberTip,
				SessionType: constant.SingleChatType,
				OperationID: req.OperationID,
			})
		}
	}
	resp.ErrorCode = 0
	return &resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.GetGroupMembersInfoResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	log.InfoByKv(claims.UID, req.OperationID, "param: ", req.MemberList)
	var resp pbGroup.GetGroupMembersInfoResp

	for _, v := range req.MemberList {
		var memberNode pbGroup.GroupMemberFullInfo
		memberInfo, err := imdb.GetMemberInfoById(req.GroupID, v)
		memberNode.UserId = v
		fmt.Println("id : ", memberNode.UserId)
		if err != nil {
			log.Error(claims.UID, req.OperationID, req.GroupID, v, "GetMemberInfoById failed, ", err.Error())
			//error occurs, only id is valid
			resp.MemberList = append(resp.MemberList, &memberNode)
			continue
		}
		user, err := imdb.FindUserByUID(v)
		if err == nil && user != nil {
			memberNode.FaceUrl = user.Icon
			memberNode.JoinTime = uint64(memberInfo.JoinTime.Unix())
			memberNode.UserId = user.UID
			memberNode.NickName = memberInfo.NickName
			memberNode.Role = memberInfo.AdministratorLevel
		}
		resp.MemberList = append(resp.MemberList, &memberNode)
	}
	resp.ErrorCode = 0
	return &resp, nil
}

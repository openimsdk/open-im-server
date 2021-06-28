package group

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	imdb "Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbChat "Open_IM/src/proto/chat"
	pbGroup "Open_IM/src/proto/group"
	"Open_IM/src/push/logic"
	"Open_IM/src/utils"
	"context"
	"encoding/json"
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
	//
	//from User:  invite: applicant
	//to user:  invite: invited
	//to application
	var resp pbGroup.InviteUserToGroupResp
	fromUserInfo, err := imdb.FindUserByUID(claims.UID)
	if err != nil {
		log.Error(claims.UID, req.OperationID, "FindUserByUID failed, err: ", err.Error())
		return &pbGroup.InviteUserToGroupResp{ErrorCode: config.ErrParam.ErrCode, ErrorMsg: config.ErrParam.ErrMsg}, nil
	}

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
		err = imdb.InsertGroupRequest(req.GroupID, fromUserInfo.UID, fromUserInfo.Name, fromUserInfo.Icon, toUserInfo.UID, req.Reason, "invited", 1)
		if err != nil {
			log.Error(v, req.OperationID, "InsertGroupRequest failed, err: ", err.Error(), "params: ",
				req.GroupID, fromUserInfo.UID, fromUserInfo.Name, fromUserInfo.Icon, toUserInfo.UID, req.Reason)
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
		resp.Id2Result = append(resp.Id2Result, &resultNode)
	}
	resp.ErrorCode = 0
	resp.ErrorMsg = "ok"

	var chatMsg pbChat.WSToMsgSvrChatMsg
	chatMsg.SendID = claims.UID
	chatMsg.RecvID = req.GroupID
	content, _ := json.Marshal(req)
	chatMsg.Content = string(content)
	chatMsg.SendTime = utils.GetCurrentTimestampBySecond()
	chatMsg.MsgFrom = constant.UserMsgType
	chatMsg.ContentType = constant.InviteUserToGroupTip
	chatMsg.SessionType = constant.GroupChatType
	logic.SendMsgByWS(&chatMsg)

	return &resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
	}
	if req.Token != config.Config.Secret {
		return &pbGroup.GetGroupAllMemberResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}

	var resp pbGroup.GetGroupAllMemberResp
	resp.ErrorCode = 0
	memberList, err := imdb.FindGroupMemberListByGroupId(req.GroupID)
	if err != nil {
		resp.ErrorCode = config.ErrDb.ErrCode
		resp.ErrorMsg = err.Error()
		log.Error(claims.UID, req.OperationID, "FindGroupMemberListByGroupId failed, ", err.Error(), "params: ", req.GroupID)
		return &resp, nil
	}

	for _, v := range memberList {
		var node pbGroup.GroupMemberFullInfo
		node.Role = v.AdministratorLevel
		node.NickName = v.NickName
		node.UserId = v.Uid
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
		log.Error(claims.UID, req.OperationID, "no access kick")
		return &pbGroup.KickGroupMemberResp{ErrorCode: config.ErrAccess.ErrCode, ErrorMsg: config.ErrAccess.ErrMsg}, nil
	}

	if len(req.UidList) == 0 {
		log.Error(claims.UID, req.OperationID, "kick list 0")
		return &pbGroup.KickGroupMemberResp{ErrorCode: config.ErrParam.ErrCode, ErrorMsg: config.ErrParam.ErrMsg}, nil
	}
	//remove
	var resp pbGroup.KickGroupMemberResp
	for _, v := range req.UidList {
		//owner cant kicked
		if v == claims.UID {
			log.Error(claims.UID, req.OperationID, v, "cant kick owner")
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v, Result: -1})
			continue
		}
		err := imdb.RemoveGroupMember(req.GroupID, v)
		if err != nil {
			log.Error(claims.UID, req.OperationID, v, req.GroupID, "RemoveGroupMember failed ", err.Error())
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v, Result: -1})
		} else {
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v, Result: 0})
		}
	}

	var chatMsg pbChat.WSToMsgSvrChatMsg
	chatMsg.SendID = claims.UID
	chatMsg.RecvID = req.GroupID
	content, _ := json.Marshal(req)
	chatMsg.Content = string(content)
	chatMsg.SendTime = utils.GetCurrentTimestampBySecond()
	chatMsg.MsgFrom = constant.UserMsgType
	chatMsg.ContentType = constant.KickGroupMemberTip
	chatMsg.SessionType = constant.GroupChatType
	logic.SendMsgByWS(&chatMsg)

	for _, v := range req.UidList {
		kickChatMsg := chatMsg
		kickChatMsg.RecvID = v
		kickChatMsg.SendTime = utils.GetCurrentTimestampBySecond()
		kickChatMsg.SessionType = constant.SingleChatType
		logic.SendMsgByWS(&kickChatMsg)
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

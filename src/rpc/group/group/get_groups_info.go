package group

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbGroup "Open_IM/src/proto/group"
	"Open_IM/src/utils"
	"context"
)

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	log.Info(req.Token, req.OperationID, "rpc get group info is server,args=%s", req.String())
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.GetGroupsInfoResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	groupsInfoList := make([]*pbGroup.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		//	Get group information to check whether the current user is in the current group
		_, err = im_mysql_model.FindGroupMemberInfoByGroupIdAndUserId(groupID, claims.UID)
		if err != nil {
			log.Error(req.Token, req.OperationID, "find group member failed,err=%s", err.Error())
			continue
		}
		groupInfoFromMysql, err := im_mysql_model.FindGroupInfoByGroupId(groupID)
		if err != nil {
			log.Error(req.Token, req.OperationID, "find group info failed,err=%s", err.Error())
			continue
		}
		var groupInfo pbGroup.GroupInfo
		groupInfo.GroupId = groupID
		groupInfo.GroupName = groupInfoFromMysql.Name
		groupInfo.Introduction = groupInfoFromMysql.Introduction
		groupInfo.Notification = groupInfoFromMysql.Notification
		groupInfo.FaceUrl = groupInfoFromMysql.FaceUrl
		groupInfo.MemberCount = uint32(im_mysql_model.GetGroupMemberNumByGroupId(groupID))
		groupInfo.CreateTime = uint64(groupInfoFromMysql.CreateTime.Unix())

		groupsInfoList = append(groupsInfoList, &groupInfo)
	}
	log.Info(req.Token, req.OperationID, "rpc get groupsInfo success return")
	return &pbGroup.GetGroupsInfoResp{Data: groupsInfoList}, nil
}

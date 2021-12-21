package group

import (
	"Open_IM/internal/rpc/chat"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	pbGroup "Open_IM/pkg/proto/group"
	"context"
)

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc set group info is server,args=%s", req.String())

	//Parse token, to find current user information
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.CommonResp{ErrorCode: constant.ErrParseToken.ErrCode, ErrorMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	groupUserInfo, err := im_mysql_model.FindGroupMemberInfoByGroupIdAndUserId(req.GroupID, claims.UID)
	if err != nil {
		log.Error("", req.OperationID, "your are not in the group,can not change this group info,err=%s", err.Error())
		return &pbGroup.CommonResp{ErrorCode: constant.ErrSetGroupInfo.ErrCode, ErrorMsg: constant.ErrSetGroupInfo.ErrMsg}, nil
	}
	if groupUserInfo.AdministratorLevel == constant.OrdinaryMember {
		return &pbGroup.CommonResp{ErrorCode: constant.ErrSetGroupInfo.ErrCode, ErrorMsg: constant.ErrAccess.ErrMsg}, nil
	}
	group, err := im_mysql_model.FindGroupInfoByGroupId(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed, ", err.Error(), req.GroupID)
		return &pbGroup.CommonResp{ErrorCode: constant.ErrSetGroupInfo.ErrCode, ErrorMsg: constant.ErrAccess.ErrMsg}, nil
	}
	////bitwise operators: 1:groupName; 10:Notification  100:Introduction; 1000:FaceUrl
	var changedType int32
	if group.GroupName != req.GroupName && req.GroupName != "" {
		changedType = 1
	}
	if group.Notification != req.Notification && req.Notification != "" {
		changedType = changedType | (1 << 1)
	}
	if group.Introduction != req.Introduction && req.Introduction != "" {
		changedType = changedType | (1 << 2)
	}
	if group.FaceUrl != req.FaceUrl && req.FaceUrl != "" {
		changedType = changedType | (1 << 3)
	}
	//only administrators can set group information
	if err = im_mysql_model.SetGroupInfo(req.GroupID, req.GroupName, req.Introduction, req.Notification, req.FaceUrl, ""); err != nil {
		return &pbGroup.CommonResp{ErrorCode: constant.ErrSetGroupInfo.ErrCode, ErrorMsg: constant.ErrSetGroupInfo.ErrMsg}, nil
	}

	if changedType != 0 {
		chat.GroupInfoChangedNotification(req.OperationID, claims.UID, changedType, group, groupUserInfo)
	}

	return &pbGroup.CommonResp{}, nil
}

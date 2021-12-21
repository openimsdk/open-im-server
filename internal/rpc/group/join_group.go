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

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.CommonResp, error) {
	log.NewInfo(req.Token, req.OperationID, "JoinGroup args ", req.String())
	//Parse token, to find current user information
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.NewError(req.OperationID, "ParseToken failed", err.Error(), req.String())
		return &pbGroup.CommonResp{ErrorCode: constant.ErrParseToken.ErrCode, ErrorMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	applicationUserInfo, err := im_mysql_model.FindUserByUID(claims.UID)
	if err != nil {
		log.NewError(req.OperationID, "FindUserByUID failed", err.Error(), claims.UID)
		return &pbGroup.CommonResp{ErrorCode: constant.ErrSearchUserInfo.ErrCode, ErrorMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
	}

	_, err = im_mysql_model.FindGroupRequestUserInfoByGroupIDAndUid(req.GroupID, claims.UID)
	if err == nil {
		err = im_mysql_model.DelGroupRequest(req.GroupID, claims.UID, "0")
	}

	if err = im_mysql_model.InsertIntoGroupRequest(req.GroupID, claims.UID, "0", req.Message, applicationUserInfo.Nickname, applicationUserInfo.FaceUrl); err != nil {
		log.Error(req.Token, req.OperationID, "Insert into group request failed,er=%s", err.Error())
		return &pbGroup.CommonResp{ErrorCode: constant.ErrJoinGroupApplication.ErrCode, ErrorMsg: constant.ErrJoinGroupApplication.ErrMsg}, nil
	}

	memberList, err := im_mysql_model.FindGroupMemberListByGroupIdAndFilterInfo(req.GroupID, constant.GroupOwner)
	if len(memberList) == 0 {
		log.NewError(req.OperationID, "FindGroupMemberListByGroupIdAndFilterInfo failed ", req.GroupID, constant.GroupOwner, err)
		return &pbGroup.CommonResp{ErrorCode: 0, ErrorMsg: ""}, nil
	}
	group, err := im_mysql_model.FindGroupInfoByGroupId(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed ", req.GroupID)
		return &pbGroup.CommonResp{ErrorCode: 0, ErrorMsg: ""}, nil
	}
	chat.ReceiveJoinApplicationNotification(req.OperationID, memberList[0].UserID, applicationUserInfo, group)

	log.NewInfo(req.OperationID, "ReceiveJoinApplicationNotification rpc JoinGroup success return")
	return &pbGroup.CommonResp{ErrorCode: 0, ErrorMsg: ""}, nil
}

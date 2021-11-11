package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
)

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc join group is server,args=%s", req.String())
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	applicationUserInfo, err := im_mysql_model.FindUserByUID(claims.UID)
	if err != nil {
		log.Error(req.Token, req.OperationID, "No this user,err=%s", err.Error())
		return &pbGroup.CommonResp{ErrorCode: config.ErrSearchUserInfo.ErrCode, ErrorMsg: config.ErrSearchUserInfo.ErrMsg}, nil
	}

	_, err = im_mysql_model.FindGroupRequestUserInfoByGroupIDAndUid(req.GroupID, claims.UID)
	if err == nil {
		err = im_mysql_model.DelGroupRequest(req.GroupID, claims.UID, "0")
	}

	log.Info(req.Token, req.OperationID, "args: ", req.GroupID, claims.UID, "0", req.Message, applicationUserInfo.Name, applicationUserInfo.Icon)

	if err = im_mysql_model.InsertIntoGroupRequest(req.GroupID, claims.UID, "0", req.Message, applicationUserInfo.Name, applicationUserInfo.Icon); err != nil {
		log.Error(req.Token, req.OperationID, "Insert into group request failed,er=%s", err.Error())
		return &pbGroup.CommonResp{ErrorCode: config.ErrJoinGroupApplication.ErrCode, ErrorMsg: config.ErrJoinGroupApplication.ErrMsg}, nil
	}
	////Find the the group owner
	//groupCreatorInfo, err := im_mysql_model.FindGroupMemberListByGroupIdAndFilterInfo(req.GroupID, constant.GroupCreator)
	//if err != nil {
	//	log.Error(req.Token, req.OperationID, "find group creator failed", err.Error())
	//} else {
	//	//Push message when join group chat
	//	logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
	//		SendID:      claims.UID,
	//		RecvID:      groupCreatorInfo[0].Uid,
	//		Content:     content_struct.NewContentStructString(0, "", req.String()),
	//		SendTime:    utils.GetCurrentTimestampBySecond(),
	//		MsgFrom:     constant.SysMsgType,
	//		ContentType: constant.JoinGroupTip,
	//		SessionType: constant.SingleChatType,
	//		OperationID: req.OperationID,
	//	})
	//}

	log.Info(req.Token, req.OperationID, "rpc join group success return")
	return &pbGroup.CommonResp{}, nil
}

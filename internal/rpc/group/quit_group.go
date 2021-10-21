package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
)

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.CommonResp, error) {
	log.InfoByArgs("rpc quit group is server,args:", req.String())
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	log.InfoByKv("args:", req.OperationID, req.GetGroupID(), claims.UID)
	//Check to see  whether there is a user in the group.
	_, err = im_mysql_model.FindGroupMemberInfoByGroupIdAndUserId(req.GroupID, claims.UID)
	if err != nil {
		log.Error(req.Token, req.OperationID, "no such group or you are not in the group,err=%s", err.Error(), req.OperationID, req.GroupID, claims.UID)
		return &pbGroup.CommonResp{ErrorCode: config.ErrQuitGroup.ErrCode, ErrorMsg: config.ErrQuitGroup.ErrMsg}, nil
	}
	//After the user's verification is successful, user will quit the group chat.
	err = im_mysql_model.DeleteGroupMemberByGroupIdAndUserId(req.GroupID, claims.UID)
	if err != nil {
		log.ErrorByArgs("this user exit the group failed,err=%s", err.Error(), req.OperationID, req.GroupID, claims.UID)
		return &pbGroup.CommonResp{ErrorCode: config.ErrQuitGroup.ErrCode, ErrorMsg: config.ErrQuitGroup.ErrMsg}, nil
	}

	err = db.DB.DelGroupMember(req.GroupID, claims.UID)
	if err != nil {
		log.Error("", "", "delete mongo group member failed, db.DB.DelGroupMember fail [err: %s]", err.Error())
		return &pbGroup.CommonResp{ErrorCode: config.ErrQuitGroup.ErrCode, ErrorMsg: config.ErrQuitGroup.ErrMsg}, nil
	}
	////Push message when quit group chat
	//jsonInfo, _ := json.Marshal(req)
	//logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
	//	SendID:      claims.UID,
	//	RecvID:      req.GroupID,
	//	Content:     string(jsonInfo),
	//	SendTime:    utils.GetCurrentTimestampBySecond(),
	//	MsgFrom:     constant.SysMsgType,
	//	ContentType: constant.QuitGroupTip,
	//	SessionType: constant.GroupChatType,
	//	OperationID: req.OperationID,
	//})
	log.Info(req.Token, req.OperationID, "rpc quit group is success return")

	return &pbGroup.CommonResp{}, nil
}

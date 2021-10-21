package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
)

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc set group info is server,args=%s", req.String())

	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	groupUserInfo, err := im_mysql_model.FindGroupMemberInfoByGroupIdAndUserId(req.GroupID, claims.UID)
	if err != nil {
		log.Error("", req.OperationID, "your are not in the group,can not change this group info,err=%s", err.Error())
		return &pbGroup.CommonResp{ErrorCode: config.ErrSetGroupInfo.ErrCode, ErrorMsg: config.ErrSetGroupInfo.ErrMsg}, nil
	}
	if groupUserInfo.AdministratorLevel == constant.OrdinaryMember {
		return &pbGroup.CommonResp{ErrorCode: config.ErrSetGroupInfo.ErrCode, ErrorMsg: config.ErrAccess.ErrMsg}, nil
	}
	//only administrators can set group information
	if err = im_mysql_model.SetGroupInfo(req.GroupID, req.GroupName, req.Introduction, req.Notification, req.FaceUrl, ""); err != nil {
		return &pbGroup.CommonResp{ErrorCode: config.ErrSetGroupInfo.ErrCode, ErrorMsg: config.ErrSetGroupInfo.ErrMsg}, nil
	}
	////Push message when set group info
	//jsonInfo, _ := json.Marshal(req)
	//logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
	//	SendID:      claims.UID,
	//	RecvID:      req.GroupID,
	//	Content:     string(jsonInfo),
	//	SendTime:    utils.GetCurrentTimestampBySecond(),
	//	MsgFrom:     constant.SysMsgType,
	//	ContentType: constant.SetGroupInfoTip,
	//	SessionType: constant.GroupChatType,
	//	OperationID: req.OperationID,
	//})
	return &pbGroup.CommonResp{}, nil
}

package group

import (
	"Open_IM/src/common/constant"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbChat "Open_IM/src/proto/chat"
	"Open_IM/src/proto/group"
	"Open_IM/src/push/content_struct"
	"Open_IM/src/push/logic"
	"Open_IM/src/utils"
	"context"
	"encoding/json"
)

func (s *groupServer) GroupApplicationResponse(_ context.Context, pb *group.GroupApplicationResponseReq) (*group.GroupApplicationResponseResp, error) {
	log.Info("", "", "rpc GroupApplicationResponse call start..., [pb: %s]", pb.String())

	reply, err := im_mysql_model.GroupApplicationResponse(pb)
	if err != nil {
		log.Error("", "", "rpc GroupApplicationResponse call..., im_mysql_model.GroupApplicationResponse fail [pb: %s] [err: %s]", pb.String(), err.Error())
		return nil, err
	}
	log.Info("", "", "rpc GroupApplicationResponse call..., im_mysql_model.GroupApplicationResponse")

	if pb.HandleResult == 1 {
		logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
			SendID:      pb.OwnerID,
			RecvID:      pb.GroupID,
			Content:     pb.String(),
			SendTime:    utils.GetCurrentTimestampBySecond(),
			MsgFrom:     constant.SysMsgType,
			ContentType: constant.AcceptGroupApplicationTip,
			SessionType: constant.GroupChatType,
			OperationID: pb.OperationID,
		})
	}

	var recvID string
	if pb.ToUserID == "0" {
		recvID = pb.FromUserID
	} else {
		recvID = pb.ToUserID
	}

	ownerUser, err := im_mysql_model.FindGroupMemberInfoByGroupIdAndUserId(pb.GroupID, pb.OwnerID)
	if err != nil {
		return nil, err
	}

	agreeOrReject := content_struct.AgreeOrRejectGroupMember{
		GroupId:  ownerUser.GroupId,
		UserId:   ownerUser.Uid,
		Role:     int(ownerUser.AdministratorLevel),
		JoinTime: uint64(ownerUser.JoinTime.Unix()),
		NickName: ownerUser.NickName,
		FaceUrl:  ownerUser.UserGroupFaceUrl,
		Reason:   pb.HandledMsg,
	}
	bAgreeOrReject, err := json.Marshal(agreeOrReject)
	if err != nil {
		return nil, err
	}

	if pb.HandleResult == 1 {
		logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
			SendID:      pb.OwnerID,
			RecvID:      recvID,
			Content:     string(bAgreeOrReject),
			SendTime:    utils.GetCurrentTimestampBySecond(),
			MsgFrom:     constant.SysMsgType,
			ContentType: constant.AcceptGroupApplicationResultTip,
			SessionType: constant.SingleChatType,
			OperationID: pb.OperationID,
		})
	} else {
		logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
			SendID:      pb.OwnerID,
			RecvID:      recvID,
			Content:     string(bAgreeOrReject),
			SendTime:    utils.GetCurrentTimestampBySecond(),
			MsgFrom:     constant.SysMsgType,
			ContentType: constant.RefuseGroupApplicationResultTip,
			SessionType: constant.SingleChatType,
			OperationID: pb.OperationID,
		})
	}

	return reply, nil
}

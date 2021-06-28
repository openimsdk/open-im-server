package group

import (
	"Open_IM/src/common/constant"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbChat "Open_IM/src/proto/chat"
	"Open_IM/src/proto/group"
	"Open_IM/src/push/logic"
	"Open_IM/src/utils"
	"context"
)

func (s *groupServer) TransferGroupOwner(_ context.Context, pb *group.TransferGroupOwnerReq) (*group.TransferGroupOwnerResp, error) {
	log.Info("", "", "rpc TransferGroupOwner call start..., [pb: %s]", pb.String())

	reply, err := im_mysql_model.TransferGroupOwner(pb)
	if err != nil {
		log.Error("", "", "rpc TransferGroupOwner call..., im_mysql_model.TransferGroupOwner fail [pb: %s] [err: %s]", pb.String(), err.Error())
		return nil, err
	}
	log.Info("", "", "rpc TransferGroupOwner call..., im_mysql_model.TransferGroupOwner")

	logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
		SendID:      pb.OldOwner,
		RecvID:      pb.GroupID,
		Content:     pb.String(),
		SendTime:    utils.GetCurrentTimestampBySecond(),
		MsgFrom:     constant.UserMsgType,
		ContentType: constant.TransferGroupOwnerTip,
		SessionType: constant.GroupChatType,
		OperationID: pb.OperationID,
	})

	return reply, nil
}

package logic

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
)

func saveUserChat(uid string, pbMsg *pbMsg.MsgSvrToPushSvrChatMsg) error {
	time := utils.GetCurrentTimestampByMill()
	seq, err := db.DB.IncrUserSeq(uid)
	if err != nil {
		log.NewError(pbMsg.OperationID, "data insert to redis err", err.Error(), pbMsg.String())
		return err
	}
	pbMsg.RecvSeq = seq
	log.NewInfo(pbMsg.OperationID, "IncrUserSeq cost time", utils.GetCurrentTimestampByMill()-time)
	return db.DB.SaveUserChat(uid, pbMsg.SendTime, pbMsg)
}

func getGroupList(groupID string) ([]string, error) {
	return im_mysql_model.SelectGroupList(groupID)
}

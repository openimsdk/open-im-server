package logic

import (
	"Open_IM/src/common/db"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbMsg "Open_IM/src/proto/chat"
	"Open_IM/src/utils"
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

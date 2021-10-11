package logic

import (
	"Open_IM/src/common/db"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	pbMsg "Open_IM/src/proto/chat"
)

func saveUserChat(uid string, pbMsg *pbMsg.MsgSvrToPushSvrChatMsg) error {
	seq, err := db.DB.IncrUserSeq(uid)
	if err != nil {
		return err
	}
	pbMsg.RecvSeq = seq
	return db.DB.SaveUserChat(uid, pbMsg.SendTime, pbMsg)
}

func getGroupList(groupID string) ([]string, error) {
	return im_mysql_model.SelectGroupList(groupID)
}

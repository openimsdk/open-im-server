package logic

import (
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
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

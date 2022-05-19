package logic

import (
	"Open_IM/pkg/common/db"
	pbMsg "Open_IM/pkg/proto/chat"
)

func saveUserChat(uid string, msg *pbMsg.MsgDataToMQ) error {
	//time := utils.GetCurrentTimestampByMill()
	//seq, err := db.DB.IncrUserSeq(uid)
	//if err != nil {
	//	log.NewError(msg.OperationID, "data insert to redis err", err.Error(), msg.String())
	//	return err
	//}
	//msg.MsgData.Seq = uint32(seq)
	pbSaveData := pbMsg.MsgDataToDB{}
	pbSaveData.MsgData = msg.MsgData
	//log.NewInfo(msg.OperationID, "IncrUserSeq cost time", utils.GetCurrentTimestampByMill()-time)
	return db.DB.SaveUserChatMongo2(uid, pbSaveData.MsgData.SendTime, &pbSaveData)
	//	return db.DB.SaveUserChatMongo2(uid, pbSaveData.MsgData.SendTime, &pbSaveData)
}

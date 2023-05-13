// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logic

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
)

func saveUserChat(uid string, msg *pbMsg.MsgDataToMQ) error {
	time := utils.GetCurrentTimestampByMill()
	seq, err := db.DB.IncrUserSeq(uid)
	if err != nil {
		log.NewError(msg.OperationID, "data insert to redis err", err.Error(), msg.String())
		return err
	}
	msg.MsgData.Seq = uint32(seq)
	pbSaveData := pbMsg.MsgDataToDB{}
	pbSaveData.MsgData = msg.MsgData
	log.NewInfo(msg.OperationID, "IncrUserSeq cost time", utils.GetCurrentTimestampByMill()-time)
	return db.DB.SaveUserChatMongo2(uid, pbSaveData.MsgData.SendTime, &pbSaveData)
	//	return db.DB.SaveUserChatMongo2(uid, pbSaveData.MsgData.SendTime, &pbSaveData)
}

func saveUserChatList(userID string, msgList []*pbMsg.MsgDataToMQ, operationID string) (error, uint64) {
	log.Info(operationID, utils.GetSelfFuncName(), "args ", userID, len(msgList))
	//return db.DB.BatchInsertChat(userID, msgList, operationID)
	return db.DB.BatchInsertChat2Cache(userID, msgList, operationID)
}

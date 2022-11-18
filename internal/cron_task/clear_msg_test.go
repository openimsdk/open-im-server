package cronTask

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"testing"
	"time"
)

func getMsgListFake(num int) []*pbMsg.MsgDataToMQ {
	var msgList []*pbMsg.MsgDataToMQ
	for i := 1; i < num; i++ {
		msgList = append(msgList, &pbMsg.MsgDataToMQ{
			Token:       "tk",
			OperationID: "operationID",
			MsgData: &server_api_params.MsgData{
				SendID:           "sendID1",
				RecvID:           "recvID1",
				GroupID:          "",
				ClientMsgID:      "xxx",
				ServerMsgID:      "xxx",
				SenderPlatformID: 1,
				SenderNickname:   "testNickName",
				SenderFaceURL:    "testFaceURL",
				SessionType:      1,
				MsgFrom:          100,
				ContentType:      101,
				Content:          []byte("testFaceURL"),
				Seq:              uint32(i),
				SendTime:         time.Now().Unix(),
				CreateTime:       time.Now().Unix(),
				Status:           1,
			},
		})
	}
}

func TestDeleteMongoMsgAndResetRedisSeq(t *testing.T) {
	operationID := getCronTaskOperationID()
	testUID1 := "test_del_id1"
	//testUID2 := "test_del_id2"
	//testUID3 := "test_del_id3"
	//testUID4 := "test_del_id4"
	//testUID5 := "test_del_id5"
	//testUID6 := "test_del_id6"
	testUserIDList := []string{testUID1}

	db.DB.SetUserMaxSeq(testUID1, 500)
	db.DB.BatchInsertChat2DB(testUID1, getMsgListFake(500), testUID1+"-"+operationID, 500)

	//db.DB.SetUserMaxSeq(testUID1, 6000)
	//db.DB.BatchInsertChat2DB()
	//
	//db.DB.SetUserMaxSeq(testUID1, 4999)
	//db.DB.BatchInsertChat2DB()
	//
	//db.DB.SetUserMaxSeq(testUID1, 30000)
	//db.DB.BatchInsertChat2DB()
	//
	//db.DB.SetUserMaxSeq(testUID1, 9999)
	//db.DB.BatchInsertChat2DB()

	for _, userID := range testUserIDList {
		operationID = userID + "-" + operationID
		if err := DeleteMongoMsgAndResetRedisSeq(operationID, userID); err != nil {
			t.Error("checkMaxSeqWithMongo failed", userID)
		}
		if err := checkMaxSeqWithMongo(operationID, userID, constant.WriteDiffusion); err != nil {
			t.Error("checkMaxSeqWithMongo failed", userID)
		}
	}

	testWorkingGroupIDList := []string{"test_del_id1", "test_del_id2", "test_del_id3", "test_del_id4", "test_del_id5"}
	for _, groupID := range testWorkingGroupIDList {
		operationID = groupID + "-" + operationID
		log.NewDebug(operationID, utils.GetSelfFuncName(), "groupID:", groupID, "userIDList:", testUserIDList)
		if err := ResetUserGroupMinSeq(operationID, groupID, testUserIDList); err != nil {
			t.Error("checkMaxSeqWithMongo failed", groupID)
		}
		if err := checkMaxSeqWithMongo(operationID, groupID, constant.ReadDiffusion); err != nil {
			t.Error("checkMaxSeqWithMongo failed", groupID)
		}
	}
}

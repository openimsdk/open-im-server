package cronTask

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"testing"
)

func TestDeleteMongoMsgAndResetRedisSeq(t *testing.T) {
	operationID := getCronTaskOperationID()
	testUserIDList := []string{"test_del_id1", "test_del_id2", "test_del_id3", "test_del_id4", "test_del_id5"}
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

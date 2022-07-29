package db

import (
	"Open_IM/pkg/common/constant"
	pbChat "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"context"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SetTokenMapByUidPid(t *testing.T) {
	m := make(map[string]int, 0)
	m["test1"] = 1
	m["test2"] = 2
	m["2332"] = 4
	err := DB.SetTokenMapByUidPid("1234", 2, m)
	assert.Nil(t, err)

}
func Test_GetTokenMapByUidPid(t *testing.T) {
	m, err := DB.GetTokenMapByUidPid("1234", "Android")
	assert.Nil(t, err)
	fmt.Println(m)
}

//func TestDataBases_GetMultiConversationMsgOpt(t *testing.T) {
//	m, err := DB.GetMultiConversationMsgOpt("fg", []string{"user", "age", "color"})
//	assert.Nil(t, err)
//	fmt.Println(m)
//}
func Test_GetKeyTTL(t *testing.T) {
	ctx := context.Background()
	key := flag.String("key", "key", "key value")
	flag.Parse()
	ttl, err := DB.RDB.TTL(ctx, *key).Result()
	assert.Nil(t, err)
	fmt.Println(ttl)
}
func Test_HGetAll(t *testing.T) {
	ctx := context.Background()
	key := flag.String("key", "key", "key value")
	flag.Parse()
	ttl, err := DB.RDB.TTL(ctx, *key).Result()
	assert.Nil(t, err)
	fmt.Println(ttl)
}

func Test_NewSetMessageToCache(t *testing.T) {
	var msg pbChat.MsgDataToMQ
	m := make(map[string]bool)
	var offlinePush server_api_params.OfflinePushInfo
	offlinePush.Title = "3"
	offlinePush.Ex = "34"
	offlinePush.IOSPushSound = "+1"
	offlinePush.IOSBadgeCount = true
	m[constant.IsPersistent] = true
	m[constant.IsHistory] = true
	var data server_api_params.MsgData
	uid := "test_uid"
	data.Seq = 11
	data.ClientMsgID = "23jwhjsdf"
	data.SendID = "111"
	data.RecvID = "222"
	data.Content = []byte{1, 2, 3, 4, 5, 6, 7}
	data.Seq = 1212
	data.Options = m
	data.OfflinePushInfo = &offlinePush
	data.AtUserIDList = []string{"1212", "23232"}
	msg.MsgData = &data
	messageList := []*pbChat.MsgDataToMQ{&msg}
	err := DB.SetMessageToCache(messageList, uid, "cacheTest")
	assert.Nil(t, err)

}
func Test_NewGetMessageListBySeq(t *testing.T) {
	var msg pbChat.MsgDataToMQ
	var data server_api_params.MsgData
	uid := "test_uid"
	data.Seq = 11
	data.ClientMsgID = "23jwhjsdf"
	msg.MsgData = &data

	seqMsg, failedSeqList, err := DB.GetMessageListBySeq(uid, []uint32{1212}, "cacheTest")
	assert.Nil(t, err)
	fmt.Println(seqMsg, failedSeqList)

}
func Test_SetUserGlobalMsgRecvOpt(t *testing.T) {
	var opt int32
	uid := "test_uid"
	opt = 1
	err := DB.SetUserGlobalMsgRecvOpt(uid, opt)
	assert.Nil(t, err)
}
func Test_GetUserGlobalMsgRecvOpt(t *testing.T) {
	uid := "test_uid"
	opt, err := DB.GetUserGlobalMsgRecvOpt(uid)
	assert.Nil(t, err)
	fmt.Println("get opt", opt)
}
func Test_JudgeAccountEXISTS(t *testing.T) {
	uid := "test_uid"
	b, err := DB.JudgeAccountEXISTS(uid)
	assert.Nil(t, err)
	fmt.Println(b)
}
func Test_SetAccountCode(t *testing.T) {
	uid := "test_uid"
	code := 666666
	err := DB.SetAccountCode(uid, code, 100)
	assert.Nil(t, err)
}
func Test_GetAccountCode(t *testing.T) {
	uid := "test_uid"
	code, err := DB.GetAccountCode(uid)
	assert.Nil(t, err)
	fmt.Println(code)
}
func Test_SetFcmToken(t *testing.T) {
	uid := "test_uid"
	token := "dfnWBtOjSj-XIZnUvDlegv:APA91bG09XTtiXfpE6U7gUVMOhnKcUkNCv4WHn0UZr2clUi-tS1jEH-HiCEW8GIAhjLIGcfUJ6NIKteC023ZxDH7J0PJ5sTxoup3fHDUPLU7KgQoZS4tPyFqCbZ6bRB7esDPEnD1n_s0"
	platformID := 2
	err := DB.SetFcmToken(uid, platformID, token, 0)
	assert.Nil(t, err)
}
func Test_GetFcmToken(t *testing.T) {
	uid := "test_uid"
	platformID := 2
	token, err := DB.GetFcmToken(uid, platformID)
	assert.Nil(t, err)
	fmt.Println("token is :", token)
}

//func Test_GetGroupMemberList(t *testing.T) {
//	groupID := "3791742301"
//	list, err := DB.GetGroupMemberIDListFromCache(groupID)
//	assert.Nil(t, err)
//	fmt.Println(list)
//}

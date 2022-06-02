package db

import (
	"Open_IM/pkg/common/constant"
	pbChat "Open_IM/pkg/proto/chat"
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
	_ = DB.SetTokenMapByUidPid("1234", 2, m)

}
func Test_GetTokenMapByUidPid(t *testing.T) {
	m, err := DB.GetTokenMapByUidPid("1234", "Android")
	assert.Nil(t, err)
	fmt.Println(m)
}

func TestDataBases_GetMultiConversationMsgOpt(t *testing.T) {
	m, err := DB.GetMultiConversationMsgOpt("fg", []string{"user", "age", "color"})
	assert.Nil(t, err)
	fmt.Println(m)
}
func Test_GetKeyTTL(t *testing.T) {
	ctx := context.Background()
	key := flag.String("key", "key", "key value")
	flag.Parse()
	ttl, err := DB.rdb.TTL(ctx, *key).Result()
	assert.Nil(t, err)
	fmt.Println(ttl)
}
func Test_HGetAll(t *testing.T) {
	ctx := context.Background()
	key := flag.String("key", "key", "key value")
	flag.Parse()
	ttl, err := DB.rdb.TTL(ctx, *key).Result()
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
	err := DB.NewSetMessageToCache(messageList, uid, "cacheTest")
	assert.Nil(t, err)

}
func Test_NewGetMessageListBySeq(t *testing.T) {
	var msg pbChat.MsgDataToMQ
	var data server_api_params.MsgData
	uid := "test_uid"
	data.Seq = 11
	data.ClientMsgID = "23jwhjsdf"
	msg.MsgData = &data

	seqMsg, failedSeqList, err := DB.NewGetMessageListBySeq(uid, []uint32{11}, "cacheTest")
	assert.Nil(t, err)
	fmt.Println(seqMsg, failedSeqList)

}

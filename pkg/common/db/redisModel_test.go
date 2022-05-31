package db

import (
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
	var data server_api_params.MsgData
	uid := "test_uid"
	data.Seq = 11
	data.ClientMsgID = "23jwhjsdf"
	msg.MsgData = &data
	messageList := []*pbChat.MsgDataToMQ{&msg}
	err := DB.NewSetMessageToCache(messageList, uid, "cacheTest")
	//err := DB.rdb.HMSet(context.Background(), "12", map[string]interface{}{"1": 2}).Err()
	assert.Nil(t, err)

}

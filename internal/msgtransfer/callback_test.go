package msgtransfer

import (
	"context"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
)

func Test_FilterAfterMsg_when_configuredForTypingAndNotification(t *testing.T) {
	allowed := &config.AfterConfig{AllowedTypes: []string{"0-2147483647"}}
	denied := &config.AfterConfig{
		AllowedTypes: []string{"0-2147483647"},
		DeniedTypes:  []string{"0-2147483647"},
	}
	typingMsg := &sdkws.MsgData{RecvID: "recipient", ContentType: constant.Typing}
	notificationMsg := &sdkws.MsgData{RecvID: "recipient", ContentType: constant.GroupCreatedNotification}

	assert.True(t, webhook.FilterAfterMsg(typingMsg, allowed, "recipient"))
	assert.True(t, webhook.FilterAfterMsg(notificationMsg, allowed, "recipient"))
	assert.False(t, webhook.FilterAfterMsg(typingMsg, denied, "recipient"))
	assert.False(t, webhook.FilterAfterMsg(notificationMsg, denied, "recipient"))
}

func Test_FilterAfterMsg_when_customAllowedAndDeniedIntervals(t *testing.T) {
	// Given: after-config with allowed and denied intervals
	afterConf := &config.AfterConfig{
		AllowedTypes: []string{"101-105", "201"},
		DeniedTypes:  []string{"103"},
	}

	// When & Then: test matching allowed
	assert.True(t, webhook.FilterAfterMsg(&sdkws.MsgData{RecvID: "recipient", ContentType: 101}, afterConf, "recipient"))
	assert.True(t, webhook.FilterAfterMsg(&sdkws.MsgData{RecvID: "recipient", ContentType: 201}, afterConf, "recipient"))

	// When & Then: test matching denied (103 is in 101-105, but also in denied 103)
	assert.False(t, webhook.FilterAfterMsg(&sdkws.MsgData{RecvID: "recipient", ContentType: 103}, afterConf, "recipient"))

	// When & Then: test not in allowed (300)
	assert.False(t, webhook.FilterAfterMsg(&sdkws.MsgData{RecvID: "recipient", ContentType: 300}, afterConf, "recipient"))
}

func Test_FilterAfterMsg_when_attentionIDsConfigured(t *testing.T) {
	groupMsg := &sdkws.MsgData{
		SendID:      "sender",
		RecvID:      "recipient",
		GroupID:     "group",
		SessionType: constant.ReadGroupChatType,
		ContentType: constant.Picture,
	}

	// When: attention matches the sender
	assert.True(t, webhook.FilterAfterMsg(groupMsg, &config.AfterConfig{AttentionIds: []string{"sender"}}, "group"))

	// When: attention matches the group target
	assert.True(t, webhook.FilterAfterMsg(groupMsg, &config.AfterConfig{AttentionIds: []string{"group"}}, "group"))

	// Then: attention does not match the recipient for a group message
	assert.False(t, webhook.FilterAfterMsg(groupMsg, &config.AfterConfig{AttentionIds: []string{"recipient"}}, "group"))
}

func Test_ToCommonCallback_and_GetContent(t *testing.T) {
	// Given: a text message and a tips message
	ctx := context.Background()
	textMsg := &sdkws.MsgData{
		SendID:      "sender_1",
		RecvID:      "recv_1",
		ServerMsgID: "server_msg_1",
		ClientMsgID: "client_msg_1",
		ContentType: constant.Picture,
		Content:     []byte("https://example.com/pic.jpg"),
	}

	tips := &sdkws.TipsComm{JsonDetail: `{"detail":"group created"}`}
	tipsBytes, err := proto.Marshal(tips)
	assert.NoError(t, err)
	notiMsg := &sdkws.MsgData{
		SendID:      "sender_1",
		GroupID:     "group_1",
		ServerMsgID: "server_msg_2",
		ClientMsgID: "client_msg_2",
		ContentType: constant.GroupCreatedNotification,
		Content:     tipsBytes,
	}

	// When: converting to common callback
	textCb := toCommonCallback(ctx, textMsg, cbapi.CallbackAfterMsgSaveDBCommand)
	notiContent := GetContent(notiMsg)

	// Then: contents and fields are correctly mapped
	assert.Equal(t, "sender_1", textCb.SendID)
	assert.Equal(t, "https://example.com/pic.jpg", textCb.Content)
	assert.Equal(t, cbapi.CallbackAfterMsgSaveDBCommand, textCb.CallbackCommand)
	assert.Equal(t, `{"detail":"group created"}`, notiContent)
}

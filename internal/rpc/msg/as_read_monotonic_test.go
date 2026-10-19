package msg

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/open-im-server/v3/pkg/notification"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/conversation"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mq/memamq"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type receiptDatabase struct {
	controller.CommonMsgDatabase
	calls   int
	initial int64
}

func (*receiptDatabase) GetMaxSeq(context.Context, string) (int64, error) { return 346, nil }
func (db *receiptDatabase) GetHasReadSeq(context.Context, string, string) (int64, error) {
	return db.initial, nil
}
func (db *receiptDatabase) SetHasReadSeq(context.Context, string, string, int64) (int64, error) {
	db.calls++
	return 346, nil // Another device advanced after the RPC's initial read.
}
func (*receiptDatabase) MarkSingleChatMsgsAsRead(context.Context, string, string, []int64) error {
	return nil
}

type receiptConversationClient struct {
	conversation.ConversationClient
	sessionType int32
}

func (c *receiptConversationClient) GetConversation(_ context.Context, req *conversation.GetConversationReq, _ ...grpc.CallOption) (*conversation.GetConversationResp, error) {
	return &conversation.GetConversationResp{Conversation: &conversation.Conversation{
		ConversationID: req.ConversationID, OwnerUserID: req.OwnerUserID,
		ConversationType: c.sessionType, UserID: "peer", GroupID: "group",
	}}, nil
}

func TestReadSeqReceiptsUseEffectiveProgress(t *testing.T) {
	for _, method := range []string{"set", "messages", "single conversation", "group conversation"} {
		for _, requested := range []int64{200, 346} {
			for _, initial := range []int64{100, 346} {
				t.Run(fmt.Sprintf("%s/initial=%d/requested=%d", method, initial, requested), func(t *testing.T) {
					ctx := mcontext.WithOpUserIDContext(context.Background(), "user")
					db := &receiptDatabase{initial: initial}
					queue := memamq.NewMemoryQueue(1, 8)
					defer queue.Stop()
					sent := make(chan *sdkws.MsgData, 1)
					notifier := notification.NewNotificationSender(&config.Notification{}, notification.WithQueue(queue),
						notification.WithLocalSendMsg(func(_ context.Context, req *pbmsg.SendMsgReq) (*pbmsg.SendMsgResp, error) {
							sent <- req.MsgData
							return &pbmsg.SendMsgResp{}, nil
						}))
					sessionType := int32(constant.SingleChatType)
					if method == "group conversation" {
						sessionType = constant.ReadGroupChatType
					}
					client := &rpcli.ConversationClient{ConversationClient: &receiptConversationClient{sessionType: sessionType}}
					server := &msgServer{MsgDatabase: db, notificationSender: notifier, config: &Config{},
						ConversationLocalCache: rpccache.NewConversationLocalCache(client, &config.LocalCache{}, nil),
						webhookClient:          webhook.NewWebhookClient("", queue),
					}
					var err error
					switch method {
					case "set":
						_, err = server.SetConversationHasReadSeq(ctx, &pbmsg.SetConversationHasReadSeqReq{UserID: "user", ConversationID: "group", HasReadSeq: requested})
					case "messages":
						_, err = server.MarkMsgsAsRead(ctx, &pbmsg.MarkMsgsAsReadReq{UserID: "user", ConversationID: "group", Seqs: []int64{requested}})
					default:
						_, err = server.MarkConversationAsRead(ctx, &pbmsg.MarkConversationAsReadReq{UserID: "user", ConversationID: "group", HasReadSeq: requested})
					}
					require.NoError(t, err)
					require.Equal(t, 1, db.calls, "equal reports must still allow persistence/notification repair")
					select {
					case message := <-sent:
						require.EqualValues(t, constant.HasReadReceipt, message.ContentType)
						var elem sdkws.NotificationElem
						require.NoError(t, json.Unmarshal(message.Content, &elem))
						var tips sdkws.MarkAsReadTips
						require.NoError(t, json.Unmarshal([]byte(elem.Detail), &tips))
						require.EqualValues(t, 346, tips.HasReadSeq)
						require.Equal(t, "group", tips.ConversationID)
					case <-time.After(3 * time.Second):
						t.Fatal("read receipt was not emitted")
					}
				})
			}
		}
	}
}

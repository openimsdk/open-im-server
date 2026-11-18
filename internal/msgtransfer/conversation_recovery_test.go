package msgtransfer

import (
	"context"
	"errors"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/protocol/constant"
	pbconv "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type conversationRecoveryConn struct {
	grpc.ClientConnInterface
	attempts    int
	hasDeadline bool
}

func (c *conversationRecoveryConn) Invoke(ctx context.Context, method string, _, _ any, _ ...grpc.CallOption) error {
	if method != pbconv.Conversation_CreateSingleChatConversations_FullMethodName {
		return errors.New("unexpected RPC")
	}
	c.attempts++
	_, c.hasDeadline = ctx.Deadline()
	if c.attempts == 1 {
		return errors.New("conversation RPC unavailable")
	}
	return nil
}

type conversationRecoveryDiscovery struct {
	discovery.Conn
	conn grpc.ClientConnInterface
}

func (d conversationRecoveryDiscovery) GetConn(context.Context, string, ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	return d.conn, nil
}

type conversationRecoveryStorage struct {
	controller.MsgTransferDatabase
	seq       int64
	persisted int
	pushed    int
}

func (s *conversationRecoveryStorage) BatchInsertChat2Cache(_ context.Context, _ string, msgs []*sdkws.MsgData) (int64, bool, map[string]int64, error) {
	previous := s.seq
	s.seq++
	msgs[0].Seq = s.seq
	return previous, previous == 0, map[string]int64{msgs[0].SendID: s.seq}, nil
}

func (*conversationRecoveryStorage) SetHasReadSeqs(context.Context, string, map[string]int64) error {
	return nil
}

func (s *conversationRecoveryStorage) MsgToMongoMQ(context.Context, string, string, []*sdkws.MsgData, int64) error {
	s.persisted++
	return nil
}

func (s *conversationRecoveryStorage) MsgToPushMQ(context.Context, string, string, *sdkws.MsgData) error {
	s.pushed++
	return nil
}

func TestSingleChatConversationCreationRetriesAfterSequenceAllocation(t *testing.T) {
	for _, name := range []string{"same worker", "replacement worker"} {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			conn := &conversationRecoveryConn{}
			storage := &conversationRecoveryStorage{}
			newHandler := func() *OnlineHistoryRedisConsumerHandler {
				handler, err := NewOnlineHistoryRedisConsumerHandler(ctx, conversationRecoveryDiscovery{conn: conn}, &Config{}, storage)
				require.NoError(t, err)
				return handler
			}
			send := func(handler *OnlineHistoryRedisConsumerHandler) {
				msg := &sdkws.MsgData{SessionType: constant.SingleChatType, SendID: "alice", RecvID: "bob"}
				handler.handleMsg(ctx, "alice_bob", "si_alice_bob", []*ContextMsg{{ctx: ctx, message: msg}}, nil)
			}
			handler := newHandler()
			send(handler)
			require.Equal(t, 1, conn.attempts)
			require.True(t, conn.hasDeadline)
			require.Equal(t, 1, storage.persisted)
			require.Equal(t, 1, storage.pushed)

			if name == "replacement worker" {
				handler = newHandler()
			}
			send(handler)
			require.Equal(t, 2, conn.attempts)
			send(handler)
			require.Equal(t, 2, conn.attempts, "healthy conversations should not perform the RPC for every batch")
			require.Equal(t, 3, storage.persisted)
			require.Equal(t, 3, storage.pushed)
		})
	}
}

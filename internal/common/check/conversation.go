package check

import (
	"OpenIM/pkg/common/config"
	discoveryRegistry "OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/proto/conversation"
	pbConversation "OpenIM/pkg/proto/conversation"
	"context"
	"google.golang.org/grpc"
)

type ConversationChecker struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func NewConversationChecker(zk discoveryRegistry.SvcDiscoveryRegistry) *ConversationChecker {
	return &ConversationChecker{zk: zk}
}

func (c *ConversationChecker) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) error {
	cc, err := c.getConn()
	if err != nil {
		return err
	}
	_, err = conversation.NewConversationClient(cc).ModifyConversationField(ctx, req)
	return err
}

func (c *ConversationChecker) getConn() (*grpc.ClientConn, error) {
	return c.zk.GetConn(config.Config.RpcRegisterName.OpenImConversationName)
}

func (c *ConversationChecker) GetSingleConversationRecvMsgOpt(ctx context.Context, userID, conversationID string) (int32, error) {
	panic("implement me")
}

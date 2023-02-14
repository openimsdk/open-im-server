package check

import (
	"Open_IM/pkg/common/config"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	"Open_IM/pkg/proto/conversation"
	pbConversation "Open_IM/pkg/proto/conversation"
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

package check

import (
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	pbConversation "Open_IM/pkg/proto/conversation"
	"context"
)

type ConversationChecker struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func (c *ConversationChecker) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) (resp *pbConversation.ModifyConversationFieldResp, err error) {
	return
}

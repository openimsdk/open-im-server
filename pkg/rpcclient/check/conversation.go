package check

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
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
	cc, err := c.getConn()
	if err != nil {
		return 0, err
	}
	var req conversation.GetConversationReq
	req.OwnerUserID = userID
	req.ConversationID = conversationID
	sConversation, err := conversation.NewConversationClient(cc).GetConversation(ctx, &req)
	if err != nil {
		return 0, err
	}
	return sConversation.GetConversation().RecvMsgOpt, err
}

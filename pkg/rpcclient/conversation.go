package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type ConversationClient struct {
	*MetaClient
}

func NewConversationClient(zk discoveryRegistry.SvcDiscoveryRegistry) *ConversationClient {
	return &ConversationClient{NewMetaClient(zk, config.Config.RpcRegisterName.OpenImConversationName)}
}

func (c *ConversationClient) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) error {
	cc, err := c.getConn()
	if err != nil {
		return err
	}
	_, err = conversation.NewConversationClient(cc).ModifyConversationField(ctx, req)
	return err
}

func (c *ConversationClient) GetSingleConversationRecvMsgOpt(ctx context.Context, userID, conversationID string) (int32, error) {
	cc, err := c.getConn()
	if err != nil {
		return 0, err
	}
	var req conversation.GetConversationReq
	req.OwnerUserID = userID
	req.ConversationID = conversationID
	conversation, err := conversation.NewConversationClient(cc).GetConversation(ctx, &req)
	if err != nil {
		return 0, err
	}
	return conversation.GetConversation().RecvMsgOpt, err
}

func (c *ConversationClient) SingleChatFirstCreateConversation(ctx context.Context, recvID, sendID string) error {
	conversation := new(pbConversation.Conversation)
	conversationID := utils.GetConversationIDBySessionType(constant.SingleChatType, recvID, sendID)
	conversation.ConversationType = constant.SingleChatType
	conversation2 := proto.Clone(conversation).(*pbConversation.Conversation)
	conversation.OwnerUserID = sendID
	conversation.UserID = recvID
	conversation.ConversationID = conversationID
	conversation2.OwnerUserID = recvID
	conversation2.UserID = sendID
	conversation2.ConversationID = conversationID
	log.ZDebug(ctx, "create single conversation", "conversation", conversation, "conversation2", conversation2)
	return c.CreateConversationsWithoutNotification(ctx, []*pbConversation.Conversation{conversation, conversation2})
}

func (c *ConversationClient) GroupChatFirstCreateConversation(ctx context.Context, groupID string, userIDs []string) error {
	var conversations []*pbConversation.Conversation
	for _, v := range userIDs {
		conversation := pbConversation.Conversation{ConversationType: constant.SuperGroupChatType, GroupID: groupID, OwnerUserID: v, ConversationID: utils.GetConversationIDBySessionType(constant.SuperGroupChatType, groupID)}
		conversations = append(conversations, &conversation)
	}
	log.ZDebug(ctx, "create group conversation", "conversations", conversations)
	return c.CreateConversationsWithoutNotification(ctx, conversations)
}

func (c *ConversationClient) CreateConversationsWithoutNotification(ctx context.Context, conversations []*pbConversation.Conversation) error {
	cc, err := c.getConn()
	if err != nil {
		return err
	}
	_, err = conversation.NewConversationClient(cc).CreateConversationsWithoutNotification(ctx, &pbConversation.CreateConversationsWithoutNotificationReq{Conversations: conversations})
	return err
}

func (c *ConversationClient) DelConversations(ctx context.Context) {

}

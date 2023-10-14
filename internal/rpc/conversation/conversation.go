package conversation

import (
	"github.com/apache/dubbo-go"
	pbconv "github.com/OpenIMSDK/protocol/conversation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/internal/cache"
	"github.com/openimsdk/open-im-server/v3/internal/controller"
)

type conversationServer struct {
	userRpcClient  *rpcclient.UserRpcClient
	RegisterCenter discoveryregistry.SvcDiscoveryRegistry
	conversationDatabase controller.ConversationDatabase
}

func Start(client discoveryregistry.SvcDiscoveryRegistry) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	userRpcClient := rpcclient.NewUserRpcClient(client)
	dubbo.RegisterProviderService(&conversationServer{
		userRpcClient:  &userRpcClient,
		RegisterCenter: client,
		conversationDatabase: controller.NewConversationDatabase(
			rdb,
			&userRpcClient,
		),
	})
	return nil
}

func (s *conversationServer) SendMessage(ctx context.Context, req *pbconv.SendMessageReq) (*pbconv.SendMessageResp, error) {
	// Use the userRpcClient and conversationDatabase to send a message
	message, err := s.conversationDatabase.SendMessage(req.ConversationId, req.Message)
	if err != nil {
		return nil, err
	}
	return &pbconv.SendMessageResp{Message: message}, nil
}

func (s *conversationServer) GetConversation(ctx context.Context, req *pbconv.GetConversationReq) (*pbconv.GetConversationResp, error) {
	// Use the conversationDatabase to get a conversation
	conversation, err := s.conversationDatabase.GetConversation(req.ConversationId)
	if err != nil {
		return nil, err
	}
	return &pbconv.GetConversationResp{Conversation: conversation}, nil
}

func (s *conversationServer) DeleteConversation(ctx context.Context, req *pbconv.DeleteConversationReq) (*pbconv.DeleteConversationResp, error) {
	// Use the conversationDatabase to delete a conversation
	err := s.conversationDatabase.DeleteConversation(req.ConversationId)
	if err != nil {
		return nil, err
	}
	return &pbconv.DeleteConversationResp{}, nil
}

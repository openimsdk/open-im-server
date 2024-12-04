package conversation

import (
	"context"

	"github.com/openimsdk/protocol/conversation"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/util/hashutil"
)

func (c *conversationServer) GetFullOwnerConversationIDs(ctx context.Context, req *conversation.GetFullOwnerConversationIDsReq) (*conversation.GetFullOwnerConversationIDsResp, error) {
	vl, err := c.conversationDatabase.FindMaxConversationUserVersionCache(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	conversationIDs, err := c.conversationDatabase.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	idHash := hashutil.IdHash(conversationIDs)
	if req.IdHash == idHash {
		conversationIDs = nil
	}
	return &conversation.GetFullOwnerConversationIDsResp{
		Version:         idHash,
		VersionID:       vl.ID.Hex(),
		Equal:           req.IdHash == idHash,
		ConversationIDs: conversationIDs,
	}, nil
}

func (c *conversationServer) GetIncrementalConversation(ctx context.Context, req *conversation.GetIncrementalConversationReq) (*conversation.GetIncrementalConversationResp, error) {
	opt := incrversion.Option[*conversation.Conversation, conversation.GetIncrementalConversationResp]{
		Ctx:             ctx,
		VersionKey:      req.UserID,
		VersionID:       req.VersionID,
		VersionNumber:   req.Version,
		Version:         c.conversationDatabase.FindConversationUserVersion,
		CacheMaxVersion: c.conversationDatabase.FindMaxConversationUserVersionCache,
		Find: func(ctx context.Context, conversationIDs []string) ([]*conversation.Conversation, error) {
			return c.getConversations(ctx, req.UserID, conversationIDs)
		},
		Resp: func(version *model.VersionLog, delIDs []string, insertList, updateList []*conversation.Conversation, full bool) *conversation.GetIncrementalConversationResp {
			return &conversation.GetIncrementalConversationResp{
				VersionID: version.ID.Hex(),
				Version:   uint64(version.Version),
				Full:      full,
				Delete:    delIDs,
				Insert:    insertList,
				Update:    updateList,
			}
		},
	}
	return opt.Build()
}

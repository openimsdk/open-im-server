package conversation

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/incrversion"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/util/hashutil"
	"github.com/openimsdk/protocol/conversation"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/utils/datautil"
)

var (
	readInactiveConversationFilterEnabled  = true
	readInactiveConversationCountThreshold = 5000
	readInactiveConversationDuration       = int64((30 * 24 * time.Hour) / time.Millisecond)
)

func (c *conversationServer) GetFullOwnerConversationIDs(ctx context.Context, req *conversation.GetFullOwnerConversationIDsReq) (*conversation.GetFullOwnerConversationIDsResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	vl, err := c.conversationDatabase.FindMaxConversationUserVersionCache(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	conversationIDs, err := c.conversationDatabase.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if shouldExcludeReadInactiveConversations(len(conversationIDs)) {
		conversationIDs, err = c.excludeReadInactiveConversations(ctx, req.UserID, conversationIDs, readInactiveConversationDuration)
		if err != nil {
			return nil, err
		}
	}
	total := int64(len(conversationIDs))
	idHash := hashutil.IdHash(conversationIDs)
	if req.IdHash == idHash {
		conversationIDs = nil
	} else if validPagination(req.GetPagination()) {
		conversationIDs = datautil.Paginate(
			conversationIDs,
			int(req.GetPagination().GetPageNumber()),
			int(req.GetPagination().GetShowNumber()),
		)
	}
	return &conversation.GetFullOwnerConversationIDsResp{
		Version:         uint64(vl.Version),
		VersionID:       vl.ID.Hex(),
		Equal:           req.IdHash == idHash,
		ConversationIDs: conversationIDs,
		Total:           total,
	}, nil
}

func validPagination(pagination *sdkws.RequestPagination) bool {
	return pagination != nil && pagination.GetPageNumber() > 0 && pagination.GetShowNumber() > 0
}

func shouldExcludeReadInactiveConversations(conversationCount int) bool {
	return readInactiveConversationFilterEnabled &&
		readInactiveConversationDuration > 0 &&
		conversationCount > readInactiveConversationCountThreshold
}

func (c *conversationServer) excludeReadInactiveConversations(ctx context.Context, userID string, conversationIDs []string, inactiveDuration int64) ([]string, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	pinnedConversationIDs, err := c.conversationDatabase.GetPinnedConversationIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	pinned := datautil.SliceSet(pinnedConversationIDs)
	seqs, err := c.msgClient.GetConversationsFullSyncSeqs(ctx, &pbmsg.GetConversationsFullSyncSeqsReq{
		UserID:          userID,
		ConversationIDs: conversationIDs,
	})
	if err != nil {
		return nil, err
	}
	expireBefore := time.Now().UnixMilli() - inactiveDuration
	filteredConversationIDs := make([]string, 0, len(conversationIDs))
	for _, conversationID := range conversationIDs {
		if _, ok := pinned[conversationID]; ok {
			filteredConversationIDs = append(filteredConversationIDs, conversationID)
			continue
		}
		seq := seqs.GetSeqs()[conversationID]
		if seq == nil || !isReadInactiveConversation(seq, expireBefore) {
			filteredConversationIDs = append(filteredConversationIDs, conversationID)
		}
	}
	return filteredConversationIDs, nil
}

func isReadInactiveConversation(seq *pbmsg.FullSyncSeqs, expireBefore int64) bool {
	if seq.GetMaxSeq() == 0 || seq.GetUserMinSeq() > seq.GetMaxSeq() {
		return true
	}
	return seq.GetHasReadSeq() >= seq.GetMaxSeq() &&
		seq.GetMaxSeqTime() > 0 &&
		seq.GetMaxSeqTime() < expireBefore
}

func (c *conversationServer) GetIncrementalConversation(ctx context.Context, req *conversation.GetIncrementalConversationReq) (*conversation.GetIncrementalConversationResp, error) {
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
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

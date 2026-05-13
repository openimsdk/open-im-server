package group

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/constant"
	pbconversation "github.com/openimsdk/protocol/conversation"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/mcontext"
)

// PinGroup pins the group conversation for the current operator
// (syncs conversation.isPinned = true via the conversation service).
func (s *groupServer) PinGroup(ctx context.Context, req *pbgroup.PinGroupReq) (*pbgroup.PinGroupResp, error) {
	if err := s.setGroupConversationPinned(ctx, req.GroupID, true); err != nil {
		return nil, err
	}
	return &pbgroup.PinGroupResp{}, nil
}

// UnpinGroup unpins the group conversation for the current operator
// (syncs conversation.isPinned = false via the conversation service).
func (s *groupServer) UnpinGroup(ctx context.Context, req *pbgroup.UnpinGroupReq) (*pbgroup.UnpinGroupResp, error) {
	if err := s.setGroupConversationPinned(ctx, req.GroupID, false); err != nil {
		return nil, err
	}
	return &pbgroup.UnpinGroupResp{}, nil
}

func (s *groupServer) setGroupConversationPinned(ctx context.Context, groupID string, pinned bool) error {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		return servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	// Must be a member of the group to pin / unpin its conversation.
	if _, err := s.db.TakeGroupMember(ctx, groupID, opUserID); err != nil {
		return err
	}
	conv := &pbconversation.ConversationReq{
		ConversationID:   msgprocessor.GetConversationIDBySessionType(constant.ReadGroupChatType, groupID),
		ConversationType: constant.ReadGroupChatType,
		GroupID:          groupID,
		IsPinned:         &wrapperspb.BoolValue{Value: pinned},
	}
	return s.conversationClient.SetConversations(ctx, []string{opUserID}, conv)
}

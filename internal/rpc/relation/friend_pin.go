package relation

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/util/conversationutil"
	"github.com/openimsdk/protocol/constant"
	conversationpb "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/log"
)

// PinFriend pins a friend conversation for ownerUserID
// (syncs friend.is_pinned = true and conversation.isPinned = true).
func (s *friendServer) PinFriend(ctx context.Context, req *relation.PinFriendReq) (*relation.PinFriendResp, error) {
	if err := s.setFriendConversationPinned(ctx, req.OwnerUserID, req.FriendUserID, true); err != nil {
		return nil, err
	}
	return &relation.PinFriendResp{}, nil
}

// UnpinFriend unpins a friend conversation for ownerUserID
// (syncs friend.is_pinned = false and conversation.isPinned = false).
func (s *friendServer) UnpinFriend(ctx context.Context, req *relation.UnpinFriendReq) (*relation.UnpinFriendResp, error) {
	if err := s.setFriendConversationPinned(ctx, req.OwnerUserID, req.FriendUserID, false); err != nil {
		return nil, err
	}
	return &relation.UnpinFriendResp{}, nil
}

func (s *friendServer) setFriendConversationPinned(ctx context.Context, ownerUserID, friendUserID string, pinned bool) error {
	if err := authverify.CheckAccessV3(ctx, ownerUserID, s.config.Share.IMAdminUserID); err != nil {
		return err
	}
	if _, err := s.db.FindFriendsWithError(ctx, ownerUserID, []string{friendUserID}); err != nil {
		return err
	}
	if err := s.db.UpdateFriends(ctx, ownerUserID, []string{friendUserID}, map[string]any{"is_pinned": pinned}); err != nil {
		return err
	}
	convID := conversationutil.GenConversationIDForSingle(ownerUserID, friendUserID)
	if err := s.conversationClient.SetConversations(ctx, []string{ownerUserID},
		&conversationpb.ConversationReq{
			ConversationID:   convID,
			ConversationType: constant.SingleChatType,
			UserID:           friendUserID,
			IsPinned:         &wrapperspb.BoolValue{Value: pinned},
		}); err != nil {
		log.ZWarn(ctx, "sync conversation isPinned failed", err,
			"ownerUserID", ownerUserID, "friendUserID", friendUserID, "isPinned", pinned)
	}
	s.notificationSender.FriendsInfoUpdateNotification(ctx, ownerUserID, []string{friendUserID})
	return nil
}

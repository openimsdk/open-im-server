package jssdk

import (
	"context"
	"sort"

	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/log"

	"github.com/gin-gonic/gin"

	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/jssdk"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

const (
	maxGetActiveConversation     = 500
	defaultGetActiveConversation = 100
)

func NewJSSdkApi(userClient *rpcli.UserClient, relationClient *rpcli.RelationClient, groupClient *rpcli.GroupClient,
	conversationClient *rpcli.ConversationClient, msgClient *rpcli.MsgClient) *JSSdk {
	return &JSSdk{
		userClient:         userClient,
		relationClient:     relationClient,
		groupClient:        groupClient,
		conversationClient: conversationClient,
		msgClient:          msgClient,
	}
}

type JSSdk struct {
	userClient         *rpcli.UserClient
	relationClient     *rpcli.RelationClient
	groupClient        *rpcli.GroupClient
	conversationClient *rpcli.ConversationClient
	msgClient          *rpcli.MsgClient
}

func (x *JSSdk) GetActiveConversations(c *gin.Context) {
	call(c, x.getActiveConversations)
}

func (x *JSSdk) GetConversations(c *gin.Context) {
	call(c, x.getConversations)
}

func (x *JSSdk) fillConversations(ctx context.Context, conversations []*jssdk.ConversationMsg) error {
	if len(conversations) == 0 {
		return nil
	}
	var (
		userIDs  []string
		groupIDs []string
	)
	for _, c := range conversations {
		if c.Conversation.GroupID == "" {
			userIDs = append(userIDs, c.Conversation.UserID)
		} else {
			groupIDs = append(groupIDs, c.Conversation.GroupID)
		}
	}
	var (
		userMap   map[string]*sdkws.UserInfo
		friendMap map[string]*relation.FriendInfoOnly
		groupMap  map[string]*sdkws.GroupInfo
	)
	if len(userIDs) > 0 {
		users, err := x.userClient.GetUsersInfo(ctx, userIDs)
		if err != nil {
			return err
		}
		friends, err := x.relationClient.GetFriendsInfo(ctx, conversations[0].Conversation.OwnerUserID, userIDs)
		if err != nil {
			return err
		}
		userMap = datautil.SliceToMap(users, (*sdkws.UserInfo).GetUserID)
		friendMap = datautil.SliceToMap(friends, (*relation.FriendInfoOnly).GetFriendUserID)
	}
	if len(groupIDs) > 0 {
		groups, err := x.groupClient.GetGroupsInfo(ctx, groupIDs)
		if err != nil {
			return err
		}
		groupMap = datautil.SliceToMap(groups, (*sdkws.GroupInfo).GetGroupID)
	}
	for _, c := range conversations {
		if c.Conversation.GroupID == "" {
			c.User = userMap[c.Conversation.UserID]
			c.Friend = friendMap[c.Conversation.UserID]
		} else {
			c.Group = groupMap[c.Conversation.GroupID]
		}
	}
	return nil
}

func (x *JSSdk) getActiveConversations(ctx context.Context, req *jssdk.GetActiveConversationsReq) (*jssdk.GetActiveConversationsResp, error) {
	if req.Count <= 0 || req.Count > maxGetActiveConversation {
		req.Count = defaultGetActiveConversation
	}
	req.OwnerUserID = mcontext.GetOpUserID(ctx)
	conversationIDs, err := x.conversationClient.GetConversationIDs(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	if len(conversationIDs) == 0 {
		return &jssdk.GetActiveConversationsResp{}, nil
	}

	activeConversation, err := x.msgClient.GetActiveConversation(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}
	if len(activeConversation) == 0 {
		return &jssdk.GetActiveConversationsResp{}, nil
	}
	readSeq, err := x.msgClient.GetHasReadSeqs(ctx, conversationIDs, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	sortConversations := sortActiveConversations{
		Conversation: activeConversation,
	}
	if len(activeConversation) > 1 {
		pinnedConversationIDs, err := x.conversationClient.GetPinnedConversationIDs(ctx, req.OwnerUserID)
		if err != nil {
			return nil, err
		}
		sortConversations.PinnedConversationIDs = datautil.SliceSet(pinnedConversationIDs)
	}
	sort.Sort(&sortConversations)
	sortList := sortConversations.Top(int(req.Count))
	conversations, err := x.conversationClient.GetConversations(ctx, datautil.Slice(sortList, func(c *msg.ActiveConversation) string {
		return c.ConversationID
	}), req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	msgs, err := x.msgClient.GetSeqMessage(ctx, req.OwnerUserID, datautil.Slice(sortList, func(c *msg.ActiveConversation) *msg.ConversationSeqs {
		return &msg.ConversationSeqs{
			ConversationID: c.ConversationID,
			Seqs:           []int64{c.MaxSeq},
		}
	}))
	if err != nil {
		return nil, err
	}
	x.checkMessagesAndGetLastMessage(ctx, req.OwnerUserID, msgs)
	conversationMap := datautil.SliceToMap(conversations, func(c *conversation.Conversation) string {
		return c.ConversationID
	})
	resp := make([]*jssdk.ConversationMsg, 0, len(sortList))
	for _, c := range sortList {
		conv, ok := conversationMap[c.ConversationID]
		if !ok {
			continue
		}
		if msgList, ok := msgs[c.ConversationID]; ok && len(msgList.Msgs) > 0 {
			resp = append(resp, &jssdk.ConversationMsg{
				Conversation: conv,
				LastMsg:      msgList.Msgs[0],
				MaxSeq:       c.MaxSeq,
				ReadSeq:      readSeq[c.ConversationID],
			})
		}

	}
	if err := x.fillConversations(ctx, resp); err != nil {
		return nil, err
	}
	var unreadCount int64
	for _, c := range activeConversation {
		count := c.MaxSeq - readSeq[c.ConversationID]
		if count > 0 {
			unreadCount += count
		}
	}
	return &jssdk.GetActiveConversationsResp{
		Conversations: resp,
		UnreadCount:   unreadCount,
	}, nil
}

func (x *JSSdk) getConversations(ctx context.Context, req *jssdk.GetConversationsReq) (*jssdk.GetConversationsResp, error) {
	req.OwnerUserID = mcontext.GetOpUserID(ctx)
	conversations, err := x.conversationClient.GetConversations(ctx, req.ConversationIDs, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return &jssdk.GetConversationsResp{}, nil
	}
	req.ConversationIDs = datautil.Slice(conversations, func(c *conversation.Conversation) string {
		return c.ConversationID
	})
	maxSeqs, err := x.msgClient.GetMaxSeqs(ctx, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	readSeqs, err := x.msgClient.GetHasReadSeqs(ctx, req.ConversationIDs, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	conversationSeqs := make([]*msg.ConversationSeqs, 0, len(conversations))
	for _, c := range conversations {
		if seq := maxSeqs[c.ConversationID]; seq > 0 {
			conversationSeqs = append(conversationSeqs, &msg.ConversationSeqs{
				ConversationID: c.ConversationID,
				Seqs:           []int64{seq},
			})
		}
	}
	var msgs map[string]*sdkws.PullMsgs
	if len(conversationSeqs) > 0 {
		msgs, err = x.msgClient.GetSeqMessage(ctx, req.OwnerUserID, conversationSeqs)
		if err != nil {
			return nil, err
		}
	}
	x.checkMessagesAndGetLastMessage(ctx, req.OwnerUserID, msgs)
	resp := make([]*jssdk.ConversationMsg, 0, len(conversations))
	for _, c := range conversations {
		if msgList, ok := msgs[c.ConversationID]; ok && len(msgList.Msgs) > 0 {
			resp = append(resp, &jssdk.ConversationMsg{
				Conversation: c,
				LastMsg:      msgList.Msgs[0],
				MaxSeq:       maxSeqs[c.ConversationID],
				ReadSeq:      readSeqs[c.ConversationID],
			})
		}

	}
	if err := x.fillConversations(ctx, resp); err != nil {
		return nil, err
	}
	var unreadCount int64
	for conversationID, maxSeq := range maxSeqs {
		count := maxSeq - readSeqs[conversationID]
		if count > 0 {
			unreadCount += count
		}
	}
	return &jssdk.GetConversationsResp{
		Conversations: resp,
		UnreadCount:   unreadCount,
	}, nil
}

// This function checks whether the latest MaxSeq message is valid.
// If not, it needs to fetch a valid message again.
func (x *JSSdk) checkMessagesAndGetLastMessage(ctx context.Context, userID string, messages map[string]*sdkws.PullMsgs) {
	var conversationIDs []string

	for conversationID, message := range messages {
		allInValid := true
		for _, data := range message.Msgs {
			if data.Status < constant.MsgStatusHasDeleted {
				allInValid = false
				break
			}
		}

		// when the conversation has been deleted by the user, the length of message.Msgs is empty
		if allInValid && len(message.Msgs) > 0 {
			conversationIDs = append(conversationIDs, conversationID)
		}
	}
	if len(conversationIDs) > 0 {
		resp, err := x.msgClient.GetLastMessage(ctx, &msg.GetLastMessageReq{
			UserID:          userID,
			ConversationIDs: conversationIDs,
		})
		if err != nil {
			log.ZError(ctx, "fetchLatestValidMessages", err, "conversationIDs", conversationIDs)
			return
		}
		for conversationID, message := range resp.Msgs {
			messages[conversationID] = &sdkws.PullMsgs{Msgs: []*sdkws.MsgData{message}}
		}
	}

}

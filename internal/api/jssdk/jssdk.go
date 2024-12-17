package jssdk

import (
	"context"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/jssdk"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

const (
	maxGetActiveConversation     = 500
	defaultGetActiveConversation = 100
)

func NewJSSdkApi() *JSSdk {
	return &JSSdk{}
}

type JSSdk struct {
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
		users, err := field(ctx, user.GetDesignateUsersCaller.Invoke, &user.GetDesignateUsersReq{UserIDs: userIDs}, (*user.GetDesignateUsersResp).GetUsersInfo)
		if err != nil {
			return err
		}
		friends, err := field(ctx, relation.GetFriendInfoCaller.Invoke, &relation.GetFriendInfoReq{OwnerUserID: conversations[0].Conversation.OwnerUserID, FriendUserIDs: userIDs}, (*relation.GetFriendInfoResp).GetFriendInfos)
		if err != nil {
			return err
		}
		userMap = datautil.SliceToMap(users, (*sdkws.UserInfo).GetUserID)
		friendMap = datautil.SliceToMap(friends, (*relation.FriendInfoOnly).GetFriendUserID)
	}
	if len(groupIDs) > 0 {
		resp, err := group.GetGroupsInfoCaller.Invoke(ctx, &group.GetGroupsInfoReq{GroupIDs: groupIDs})
		if err != nil {
			return err
		}
		groupMap = datautil.SliceToMap(resp.GroupInfos, (*sdkws.GroupInfo).GetGroupID)
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
	conversationIDs, err := field(ctx, conversation.GetConversationIDsCaller.Invoke,
		&conversation.GetConversationIDsReq{UserID: req.OwnerUserID}, (*conversation.GetConversationIDsResp).GetConversationIDs)
	if err != nil {
		return nil, err
	}
	if len(conversationIDs) == 0 {
		return &jssdk.GetActiveConversationsResp{}, nil
	}
	readSeq, err := field(ctx, msg.GetHasReadSeqsCaller.Invoke,
		&msg.GetHasReadSeqsReq{UserID: req.OwnerUserID, ConversationIDs: conversationIDs}, (*msg.SeqsInfoResp).GetMaxSeqs)
	if err != nil {
		return nil, err
	}
	activeConversation, err := field(ctx, msg.GetActiveConversationCaller.Invoke,
		&msg.GetActiveConversationReq{ConversationIDs: conversationIDs}, (*msg.GetActiveConversationResp).GetConversations)
	if err != nil {
		return nil, err
	}
	if len(activeConversation) == 0 {
		return &jssdk.GetActiveConversationsResp{}, nil
	}
	sortConversations := sortActiveConversations{
		Conversation: activeConversation,
	}
	if len(activeConversation) > 1 {
		pinnedConversationIDs, err := field(ctx, conversation.GetPinnedConversationIDsCaller.Invoke,
			&conversation.GetPinnedConversationIDsReq{UserID: req.OwnerUserID}, (*conversation.GetPinnedConversationIDsResp).GetConversationIDs)
		if err != nil {
			return nil, err
		}
		sortConversations.PinnedConversationIDs = datautil.SliceSet(pinnedConversationIDs)
	}
	sort.Sort(&sortConversations)
	sortList := sortConversations.Top(int(req.Count))
	conversations, err := field(ctx, conversation.GetConversationsCaller.Invoke,
		&conversation.GetConversationsReq{
			OwnerUserID: req.OwnerUserID,
			ConversationIDs: datautil.Slice(sortList, func(c *msg.ActiveConversation) string {
				return c.ConversationID
			})}, (*conversation.GetConversationsResp).GetConversations)
	if err != nil {
		return nil, err
	}
	msgs, err := field(ctx, msg.GetSeqMessageCaller.Invoke,
		&msg.GetSeqMessageReq{
			UserID: req.OwnerUserID,
			Conversations: datautil.Slice(sortList, func(c *msg.ActiveConversation) *msg.ConversationSeqs {
				return &msg.ConversationSeqs{
					ConversationID: c.ConversationID,
					Seqs:           []int64{c.MaxSeq},
				}
			}),
		}, (*msg.GetSeqMessageResp).GetMsgs)
	if err != nil {
		return nil, err
	}
	conversationMap := datautil.SliceToMap(conversations, func(c *conversation.Conversation) string {
		return c.ConversationID
	})
	resp := make([]*jssdk.ConversationMsg, 0, len(sortList))
	for _, c := range sortList {
		conv, ok := conversationMap[c.ConversationID]
		if !ok {
			continue
		}
		var lastMsg *sdkws.MsgData
		if msgList, ok := msgs[c.ConversationID]; ok && len(msgList.Msgs) > 0 {
			lastMsg = msgList.Msgs[0]
		}
		resp = append(resp, &jssdk.ConversationMsg{
			Conversation: conv,
			LastMsg:      lastMsg,
			MaxSeq:       c.MaxSeq,
			ReadSeq:      readSeq[c.ConversationID],
		})
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
	conversations, err := field(ctx, conversation.GetConversationsCaller.Invoke, &conversation.GetConversationsReq{OwnerUserID: req.OwnerUserID, ConversationIDs: req.ConversationIDs}, (*conversation.GetConversationsResp).GetConversations)
	if err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return &jssdk.GetConversationsResp{}, nil
	}
	req.ConversationIDs = datautil.Slice(conversations, func(c *conversation.Conversation) string {
		return c.ConversationID
	})
	maxSeqs, err := field(ctx, msg.GetMaxSeqsCaller.Invoke,
		&msg.GetMaxSeqsReq{ConversationIDs: req.ConversationIDs}, (*msg.SeqsInfoResp).GetMaxSeqs)
	if err != nil {
		return nil, err
	}
	readSeqs, err := field(ctx, msg.GetHasReadSeqsCaller.Invoke,
		&msg.GetHasReadSeqsReq{UserID: req.OwnerUserID, ConversationIDs: req.ConversationIDs}, (*msg.SeqsInfoResp).GetMaxSeqs)
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
		msgs, err = field(ctx, msg.GetSeqMessageCaller.Invoke,
			&msg.GetSeqMessageReq{UserID: req.OwnerUserID, Conversations: conversationSeqs}, (*msg.GetSeqMessageResp).GetMsgs)
		if err != nil {
			return nil, err
		}
	}
	resp := make([]*jssdk.ConversationMsg, 0, len(conversations))
	for _, c := range conversations {
		var lastMsg *sdkws.MsgData
		if msgList, ok := msgs[c.ConversationID]; ok && len(msgList.Msgs) > 0 {
			lastMsg = msgList.Msgs[0]
		}
		resp = append(resp, &jssdk.ConversationMsg{
			Conversation: c,
			LastMsg:      lastMsg,
			MaxSeq:       maxSeqs[c.ConversationID],
			ReadSeq:      readSeqs[c.ConversationID],
		})
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

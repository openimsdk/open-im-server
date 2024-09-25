package jssdk

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"sort"
)

const (
	maxGetActiveConversation     = 500
	defaultGetActiveConversation = 100
)

func NewJSSdkApi(msg msg.MsgClient, conv conversation.ConversationClient) *JSSdk {
	return &JSSdk{
		msg:  msg,
		conv: conv,
	}
}

type JSSdk struct {
	msg  msg.MsgClient
	conv conversation.ConversationClient
}

func (x *JSSdk) GetActiveConversations(c *gin.Context) {
	call(c, x.getActiveConversations)
}

func (x *JSSdk) GetConversations(c *gin.Context) {
	call(c, x.getConversations)
}

func (x *JSSdk) getActiveConversations(ctx *gin.Context) (*ConversationsResp, error) {
	req, err := a2r.ParseRequest[ActiveConversationsReq](ctx)
	if err != nil {
		return nil, err
	}
	if req.Count <= 0 || req.Count > maxGetActiveConversation {
		req.Count = defaultGetActiveConversation
	}
	opUserID := mcontext.GetOpUserID(ctx)
	conversationIDs, err := field(ctx, x.conv.GetConversationIDs,
		&conversation.GetConversationIDsReq{UserID: opUserID}, (*conversation.GetConversationIDsResp).GetConversationIDs)
	if err != nil {
		return nil, err
	}
	if len(conversationIDs) == 0 {
		return &ConversationsResp{}, nil
	}
	readSeq, err := field(ctx, x.msg.GetHasReadSeqs,
		&msg.GetHasReadSeqsReq{UserID: opUserID, ConversationIDs: conversationIDs}, (*msg.SeqsInfoResp).GetMaxSeqs)
	if err != nil {
		return nil, err
	}
	activeConversation, err := field(ctx, x.msg.GetActiveConversation,
		&msg.GetActiveConversationReq{ConversationIDs: conversationIDs}, (*msg.GetActiveConversationResp).GetConversations)
	if err != nil {
		return nil, err
	}
	if len(activeConversation) == 0 {
		return &ConversationsResp{}, nil
	}
	sortConversations := sortActiveConversations{
		Conversation: activeConversation,
	}
	if len(activeConversation) > 1 {
		pinnedConversationIDs, err := field(ctx, x.conv.GetPinnedConversationIDs,
			&conversation.GetPinnedConversationIDsReq{UserID: opUserID}, (*conversation.GetPinnedConversationIDsResp).GetConversationIDs)
		if err != nil {
			return nil, err
		}
		sortConversations.PinnedConversationIDs = datautil.SliceSet(pinnedConversationIDs)
	}
	sort.Sort(&sortConversations)
	sortList := sortConversations.Top(req.Count)
	conversations, err := field(ctx, x.conv.GetConversations,
		&conversation.GetConversationsReq{
			OwnerUserID: opUserID,
			ConversationIDs: datautil.Slice(sortList, func(c *msg.ActiveConversation) string {
				return c.ConversationID
			})}, (*conversation.GetConversationsResp).GetConversations)
	if err != nil {
		return nil, err
	}
	msgs, err := field(ctx, x.msg.GetSeqMessage,
		&msg.GetSeqMessageReq{
			UserID: opUserID,
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
	resp := make([]ConversationMsg, 0, len(sortList))
	for _, c := range sortList {
		conv, ok := conversationMap[c.ConversationID]
		if !ok {
			continue
		}
		var lastMsg *sdkws.MsgData
		if msgList, ok := msgs[c.ConversationID]; ok && len(msgList.Msgs) > 0 {
			lastMsg = msgList.Msgs[0]
		}
		resp = append(resp, ConversationMsg{
			Conversation: conv,
			LastMsg:      lastMsg,
			MaxSeq:       c.MaxSeq,
			ReadSeq:      readSeq[c.ConversationID],
		})
	}
	var unreadCount int64
	for _, c := range activeConversation {
		count := c.MaxSeq - readSeq[c.ConversationID]
		if count > 0 {
			unreadCount += count
		}
	}
	return &ConversationsResp{
		Conversations: resp,
		UnreadCount:   unreadCount,
	}, nil
}

func (x *JSSdk) getConversations(ctx *gin.Context) (*ConversationsResp, error) {
	req, err := a2r.ParseRequest[conversation.GetConversationsReq](ctx)
	if err != nil {
		return nil, err
	}
	req.OwnerUserID = mcontext.GetOpUserID(ctx)
	conversations, err := field(ctx, x.conv.GetConversations, req, (*conversation.GetConversationsResp).GetConversations)
	if err != nil {
		return nil, err
	}
	if len(conversations) == 0 {
		return &ConversationsResp{}, nil
	}
	req.ConversationIDs = datautil.Slice(conversations, func(c *conversation.Conversation) string {
		return c.ConversationID
	})
	maxSeqs, err := field(ctx, x.msg.GetMaxSeqs,
		&msg.GetMaxSeqsReq{ConversationIDs: req.ConversationIDs}, (*msg.SeqsInfoResp).GetMaxSeqs)
	if err != nil {
		return nil, err
	}
	readSeqs, err := field(ctx, x.msg.GetHasReadSeqs,
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
		msgs, err = field(ctx, x.msg.GetSeqMessage,
			&msg.GetSeqMessageReq{UserID: req.OwnerUserID, Conversations: conversationSeqs}, (*msg.GetSeqMessageResp).GetMsgs)
		if err != nil {
			return nil, err
		}
	}
	resp := make([]ConversationMsg, 0, len(conversations))
	for _, c := range conversations {
		var lastMsg *sdkws.MsgData
		if msgList, ok := msgs[c.ConversationID]; ok && len(msgList.Msgs) > 0 {
			lastMsg = msgList.Msgs[0]
		}
		resp = append(resp, ConversationMsg{
			Conversation: c,
			LastMsg:      lastMsg,
			MaxSeq:       maxSeqs[c.ConversationID],
			ReadSeq:      readSeqs[c.ConversationID],
		})
	}
	var unreadCount int64
	for conversationID, maxSeq := range maxSeqs {
		count := maxSeq - readSeqs[conversationID]
		if count > 0 {
			unreadCount += count
		}
	}
	return &ConversationsResp{
		Conversations: resp,
		UnreadCount:   unreadCount,
	}, nil
}

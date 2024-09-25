package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
	"sort"
)

const limitGetActiveConversation = 100

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

func (x *JSSdk) GetActiveConversation(c *gin.Context) {
	call(c, x.getActiveConversation)
}

func (x *JSSdk) GetConversations(c *gin.Context) {
	call(c, x.getConversations)
}

func (x *JSSdk) getActiveConversation(ctx *gin.Context) ([]ConversationMsg, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	conversationIDs, err := field(ctx, x.conv.GetConversationIDs,
		&conversation.GetConversationIDsReq{UserID: opUserID}, (*conversation.GetConversationIDsResp).GetConversationIDs)
	if err != nil {
		return nil, err
	}
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	activeConversation, err := field(ctx, x.msg.GetActiveConversation,
		&msg.GetActiveConversationReq{ConversationIDs: conversationIDs}, (*msg.GetActiveConversationResp).GetConversations)
	if err != nil {
		return nil, err
	}
	if len(activeConversation) == 0 {
		return nil, nil
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
	sortList := sortConversations.Top(limitGetActiveConversation)
	conversations, err := field(ctx, x.conv.GetConversations,
		&conversation.GetConversationsReq{
			OwnerUserID: opUserID,
			ConversationIDs: datautil.Slice(sortList, func(c *msg.ActiveConversation) string {
				return c.ConversationID
			})}, (*conversation.GetConversationsResp).GetConversations)
	if err != nil {
		return nil, err
	}
	//readSeq, err := field(ctx, x.msg.GetHasReadSeqs,
	//	&msg.GetHasReadSeqsReq{UserID: opUserID, ConversationIDs: conversationIDs}, (*msg.SeqsInfoResp).GetMaxSeqs)
	//if err != nil {
	//	return nil, err
	//}
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
			//MaxSeq:       c.MaxSeq,
			//MaxSeqTime:   c.LastTime,
			//ReadSeq:      readSeq[c.ConversationID],
		})
	}
	return resp, nil
}

func (x *JSSdk) getConversations(ctx *gin.Context) ([]ConversationMsg, error) {
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
		return nil, nil
	}
	maxSeqs, err := field(ctx, x.msg.GetMaxSeqs,
		&msg.GetMaxSeqsReq{ConversationIDs: datautil.Slice(conversations, func(c *conversation.Conversation) string {
			return c.ConversationID
		})}, (*msg.SeqsInfoResp).GetMaxSeqs)
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
		})
	}
	return resp, nil
}

type ConversationMsg struct {
	Conversation *conversation.Conversation `json:"conversation"`
	LastMsg      *sdkws.MsgData             `json:"lastMsg"`
	//ReadSeq      int64                      `json:"readSeq"`
	//MaxSeq       int64                      `json:"maxSeq"`
	//MaxSeqTime   int64                      `json:"maxSeqTime"`
}

type sortActiveConversations struct {
	Conversation          []*msg.ActiveConversation
	PinnedConversationIDs map[string]struct{}
}

func (s sortActiveConversations) Top(limit int) []*msg.ActiveConversation {
	if limit > 0 && len(s.Conversation) > limit {
		return s.Conversation[:limit]
	}
	return s.Conversation
}

func (s sortActiveConversations) Len() int {
	return len(s.Conversation)
}

func (s sortActiveConversations) Less(i, j int) bool {
	iv, jv := s.Conversation[i], s.Conversation[j]
	_, ip := s.PinnedConversationIDs[iv.ConversationID]
	_, jp := s.PinnedConversationIDs[jv.ConversationID]
	if ip != jp {
		return ip
	}
	return iv.LastTime > jv.LastTime
}

func (s sortActiveConversations) Swap(i, j int) {
	s.Conversation[i], s.Conversation[j] = s.Conversation[j], s.Conversation[i]
}

func field[A, B, C any](ctx context.Context, fn func(ctx context.Context, req *A, opts ...grpc.CallOption) (*B, error), req *A, get func(*B) C) (C, error) {
	resp, err := fn(ctx, req)
	if err != nil {
		var c C
		return c, err
	}
	return get(resp), nil
}

func call[R any](c *gin.Context, fn func(ctx *gin.Context) (R, error)) {
	resp, err := fn(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

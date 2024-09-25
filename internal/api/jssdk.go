package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
	"sort"
)

const limitGetActiveConversation = 100

type JSSdk struct {
	msg  msg.MsgClient
	conv conversation.ConversationClient
}

func field[A, B, C any](ctx context.Context, fn func(ctx context.Context, req *A, opts ...grpc.CallOption) (*B, error), req *A, get func(*B) C) (C, error) {
	resp, err := fn(ctx, req)
	if err != nil {
		var c C
		return c, err
	}
	return get(resp), nil
}

func (x *JSSdk) GetActiveConversation(ctx *gin.Context) ([]ConversationMsg, error) {
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
		// todo get pinned conversation ids
	}
	sort.Sort(&sortConversations)
	sortList := sortConversations.Top(limitGetActiveConversation)
	conversations, err := field(ctx, x.conv.GetConversations,
		&conversation.GetConversationsReq{ConversationIDs: datautil.Slice(sortList, func(c *msg.ActiveConversation) string {
			return c.ConversationID
		})}, (*conversation.GetConversationsResp).GetConversations)
	if err != nil {
		return nil, err
	}
	readSeq, err := field(ctx, x.msg.GetHasReadSeqs,
		&msg.GetHasReadSeqsReq{UserID: opUserID, ConversationIDs: conversationIDs}, (*msg.SeqsInfoResp).GetMaxSeqs)
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
		msgList, ok := msgs[c.ConversationID]
		if ok {
			continue
		}
		var lastMsg *sdkws.MsgData
		if len(msgList.Msgs) > 0 {
			lastMsg = msgList.Msgs[0]
		}
		resp = append(resp, ConversationMsg{
			Conversation: conv,
			LastMsg:      lastMsg,
			MaxSeq:       c.MaxSeq,
			MaxSeqTime:   c.LastTime,
			ReadSeq:      readSeq[c.ConversationID],
		})
	}
	return resp, nil
}

type ConversationMsg struct {
	Conversation *conversation.Conversation `json:"conversation"`
	LastMsg      *sdkws.MsgData             `json:"lastMsg"`
	ReadSeq      int64                      `json:"readSeq"`
	MaxSeq       int64                      `json:"maxSeq"`
	MaxSeqTime   int64                      `json:"maxSeqTime"`
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

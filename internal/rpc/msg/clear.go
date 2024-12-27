package msg

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/log"
	"strings"
)

// DestructMsgs hard delete in Database.
func (m *msgServer) DestructMsgs(ctx context.Context, req *msg.DestructMsgsReq) (*msg.DestructMsgsResp, error) {
	if err := authverify.CheckAdmin(ctx, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	docs, err := m.MsgDatabase.GetRandBeforeMsg(ctx, req.Timestamp, int(req.Limit))
	if err != nil {
		return nil, err
	}
	for i, doc := range docs {
		if err := m.MsgDatabase.DeleteDoc(ctx, doc.DocID); err != nil {
			return nil, err
		}
		log.ZDebug(ctx, "DestructMsgs delete doc", "index", i, "docID", doc.DocID)
		index := strings.LastIndex(doc.DocID, ":")
		if index < 0 {
			continue
		}
		var minSeq int64
		for _, model := range doc.Msg {
			if model.Msg == nil {
				continue
			}
			if model.Msg.Seq > minSeq {
				minSeq = model.Msg.Seq
			}
		}
		if minSeq <= 0 {
			continue
		}
		conversationID := doc.DocID[:index]
		if conversationID == "" {
			continue
		}
		minSeq++
		if err := m.MsgDatabase.SetMinSeq(ctx, conversationID, minSeq); err != nil {
			return nil, err
		}
		log.ZDebug(ctx, "DestructMsgs delete doc set min seq", "index", i, "docID", doc.DocID, "conversationID", conversationID, "setMinSeq", minSeq)
	}
	return &msg.DestructMsgsResp{Count: int32(len(docs))}, nil
}

func (m *msgServer) GetLastMessageSeqByTime(ctx context.Context, req *msg.GetLastMessageSeqByTimeReq) (*msg.GetLastMessageSeqByTimeResp, error) {
	seq, err := m.MsgDatabase.GetLastMessageSeqByTime(ctx, req.ConversationID, req.Time)
	if err != nil {
		return nil, err
	}
	return &msg.GetLastMessageSeqByTimeResp{Seq: seq}, nil
}

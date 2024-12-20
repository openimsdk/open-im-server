package msg

import (
	"context"
	"strings"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	pbconv "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/idutil"
	"github.com/openimsdk/tools/utils/stringutil"
	"golang.org/x/sync/errgroup"
)

// hard delete in Database.
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

// soft delete for user self
func (m *msgServer) ClearMsg(ctx context.Context, req *msg.ClearMsgReq) (*msg.ClearMsgResp, error) {
	temp := convert.ConversationsPb2DB(req.Conversations)

	batchNum := 100

	errg, _ := errgroup.WithContext(ctx)
	errg.SetLimit(100)

	for i := 0; i < len(temp); i += batchNum {
		batch := temp[i:min(i+batchNum, len(temp))]

		errg.Go(func() error {
			for _, conversation := range batch {
				handleCtx := mcontext.NewCtx(stringutil.GetSelfFuncName() + "-" + idutil.OperationIDGenerator() + "-" + conversation.ConversationID + "-" + conversation.OwnerUserID)
				log.ZDebug(handleCtx, "User MsgsDestruct",
					"conversationID", conversation.ConversationID,
					"ownerUserID", conversation.OwnerUserID,
					"msgDestructTime", conversation.MsgDestructTime,
					"lastMsgDestructTime", conversation.LatestMsgDestructTime)

				seqs, err := m.MsgDatabase.ClearUserMsgs(handleCtx, conversation.OwnerUserID, conversation.ConversationID, conversation.MsgDestructTime, conversation.LatestMsgDestructTime)
				if err != nil {
					log.ZError(handleCtx, "user msg destruct failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
					continue
				}

				if len(seqs) > 0 {
					minseq := datautil.Max(seqs...)

					// update
					if err := pbconv.UpdateConversationCaller.Execute(ctx, &pbconv.UpdateConversationReq{
						ConversationID:        conversation.ConversationID,
						UserIDs:               []string{conversation.OwnerUserID},
						MinSeq:                wrapperspb.Int64(minseq),
						LatestMsgDestructTime: wrapperspb.Int64(time.Now().UnixMilli()),
					}); err != nil {
						log.ZError(handleCtx, "updateUsersConversationField failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
						continue
					}

					if err := pbconv.SetConversationMinSeqCaller.Execute(ctx, &pbconv.SetConversationMinSeqReq{
						ConversationID: conversation.ConversationID,
						OwnerUserID:    []string{conversation.OwnerUserID},
						MinSeq:         minseq,
					}); err != nil {
						return err
					}

					// if you need Notify SDK client userseq is update.
					// m.msgNotificationSender.UserDeleteMsgsNotification(handleCtx, conversation.OwnerUserID, conversation.ConversationID, seqs)
				}
			}
			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *msgServer) GetLastMessageSeqByTime(ctx context.Context, req *msg.GetLastMessageSeqByTimeReq) (*msg.GetLastMessageSeqByTimeResp, error) {
	seq, err := m.MsgDatabase.GetLastMessageSeqByTime(ctx, req.ConversationID, req.Time)
	if err != nil {
		return nil, err
	}
	return &msg.GetLastMessageSeqByTimeResp{Seq: seq}, nil
}

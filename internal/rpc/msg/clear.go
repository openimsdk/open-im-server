package msg

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/conversation"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/idutil"
	"github.com/openimsdk/tools/utils/stringutil"
)

// hard delete in Database.
func (m *msgServer) ClearMsg(ctx context.Context, req *msg.ClearMsgReq) (_ *msg.ClearMsgResp, err error) {
	if err := authverify.CheckAdmin(ctx, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	if req.Timestamp > time.Now().UnixMilli() {
		return nil, errs.ErrArgs.WrapMsg("request millisecond timestamp error")
	}
	var (
		docNum int
		msgNum int
		start  = time.Now()
	)
	clearMsg := func(ctx context.Context) (bool, error) {
		conversationSeqs := make(map[string]struct{})

		// update latest msg destruct time in conversation DB.
		defer func() {
			req := &conversation.UpdateConversationReq{
				LatestMsgDestructTime: wrapperspb.Int64(time.Now().UnixMilli()),
			}
			for conversationID := range conversationSeqs {
				req.ConversationID = conversationID
				if err := m.Conversation.UpdateConversation(ctx, req); err != nil {
					log.ZError(ctx, "update conversation max seq failed", err, "conversationID", conversationID, "msgDestructTime", req.MsgDestructTime)
				}
			}
		}()

		msgs, err := m.MsgDatabase.GetBeforeMsg(ctx, req.Timestamp, 100)
		if err != nil {
			return false, err
		}
		if len(msgs) == 0 {
			return false, nil
		}
		for _, msg := range msgs {
			index, err := m.MsgDatabase.DeleteDocMsgBefore(ctx, req.Timestamp, msg)
			if err != nil {
				return false, err
			}
			if len(index) == 0 {
				return false, errs.ErrInternalServer.WrapMsg("delete doc msg failed")
			}
			docNum++
			msgNum += len(index)
			conversationID := msg.DocID[:strings.LastIndex(msg.DocID, ":")]
			if _, ok := conversationSeqs[conversationID]; !ok {
				conversationSeqs[conversationID] = struct{}{}
			}
		}
		return true, nil
	}

	for {
		keep, err := clearMsg(ctx)
		if err != nil {
			log.ZError(ctx, "clear msg failed", err, "docNum", docNum, "msgNum", msgNum, "cost", time.Since(start))
			return nil, err
		}
		if !keep {
			log.ZInfo(ctx, "clear msg success", "docNum", docNum, "msgNum", msgNum, "cost", time.Since(start))
			break
		}
		log.ZInfo(ctx, "clearing message", "docNum", docNum, "msgNum", msgNum, "cost", time.Since(start))
	}
	return &msg.ClearMsgResp{}, nil
}

// soft delete for self
func (m *msgServer) DestructMsgs(ctx context.Context, req *msg.DestructMsgsReq) (_ *msg.DestructMsgsResp, err error) {
	temp := convert.ConversationsPb2DB(req.Conversations)

	batchNum := 100
	var wg sync.WaitGroup
	wg.Add((len(temp) + batchNum - 1) / batchNum)

	for i := 0; i < len(temp); i += batchNum {
		batch := temp[i:min(i+batchNum, len(temp))]

		go func(batch []*model.Conversation) {
			defer wg.Done()

			for _, conversation := range temp {
				handleCtx := mcontext.NewCtx(stringutil.GetSelfFuncName() + "-" + idutil.OperationIDGenerator() + "-" + conversation.ConversationID + "-" + conversation.OwnerUserID)
				log.ZDebug(handleCtx, "User MsgsDestruct",
					"conversationID", conversation.ConversationID,
					"ownerUserID", conversation.OwnerUserID,
					"msgDestructTime", conversation.MsgDestructTime,
					"lastMsgDestructTime", conversation.LatestMsgDestructTime)

				seqs, err := m.MsgDatabase.UserMsgsDestruct(handleCtx, conversation.OwnerUserID, conversation.ConversationID, conversation.MsgDestructTime, conversation.LatestMsgDestructTime)
				if err != nil {
					log.ZError(handleCtx, "user msg destruct failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
					continue
				}

				if len(seqs) > 0 {
					if err := m.Conversation.UpdateConversation(handleCtx,
						&pbconversation.UpdateConversationReq{
							UserIDs:               []string{conversation.OwnerUserID},
							ConversationID:        conversation.ConversationID,
							LatestMsgDestructTime: wrapperspb.Int64(time.Now().UnixMilli())}); err != nil {
						log.ZError(handleCtx, "updateUsersConversationField failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
						continue
					}

					// if you need Notify SDK client userseq is update.
					// m.msgNotificationSender.UserDeleteMsgsNotification(handleCtx, conversation.OwnerUserID, conversation.ConversationID, seqs)
				}
			}
		}(batch)

	}
	wg.Wait()

	return nil, nil
}

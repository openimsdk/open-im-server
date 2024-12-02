package msg

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/idutil"
	"github.com/openimsdk/tools/utils/stringutil"
	"golang.org/x/sync/errgroup"
)

// hard delete in Database.
func (m *msgServer) DestructMsgs(ctx context.Context, req *msg.DestructMsgsReq) (_ *msg.DestructMsgsResp, err error) {
	if err := authverify.CheckAdmin(ctx, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	if req.Timestamp > time.Now().UnixMilli() {
		return nil, errs.ErrArgs.WrapMsg("request millisecond timestamp error")
	}
	var (
		docNum   int
		msgNum   int
		start    = time.Now()
		getLimit = 5000
	)

	destructMsg := func(ctx context.Context) (bool, error) {
		docIDs, err := m.MsgDatabase.GetDocIDs(ctx)
		if err != nil {
			return false, err
		}

		msgs, err := m.MsgDatabase.GetBeforeMsg(ctx, req.Timestamp, docIDs, getLimit)
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
		}

		return true, nil
	}

	_, err = destructMsg(ctx)
	if err != nil {
		log.ZError(ctx, "clear msg failed", err, "docNum", docNum, "msgNum", msgNum, "cost", time.Since(start))
		return nil, err
	}

	log.ZDebug(ctx, "clearing message", "docNum", docNum, "msgNum", msgNum, "cost", time.Since(start))

	return &msg.DestructMsgsResp{}, nil
}

// soft delete for user self
func (m *msgServer) ClearMsg(ctx context.Context, req *msg.ClearMsgReq) (_ *msg.ClearMsgResp, err error) {
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
					if err := m.Conversation.UpdateConversation(handleCtx,
						&pbconversation.UpdateConversationReq{
							UserIDs:               []string{conversation.OwnerUserID},
							ConversationID:        conversation.ConversationID,
							LatestMsgDestructTime: wrapperspb.Int64(time.Now().UnixMilli()),
							MinSeq:                wrapperspb.Int64(minseq),
						}); err != nil {
						log.ZError(handleCtx, "updateUsersConversationField failed", err, "conversationID", conversation.ConversationID, "ownerUserID", conversation.OwnerUserID)
						continue
					}

					if err := m.Conversation.SetConversationMinSeq(handleCtx, []string{conversation.OwnerUserID}, conversation.ConversationID, minseq); err != nil {
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

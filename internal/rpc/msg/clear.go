package msg

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"strings"
	"time"
)

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
		defer func() {
			req := &conversation.UpdateConversationReq{
				MsgDestructTime: wrapperspb.Int64(time.Now().UnixMilli()),
			}
			for conversationID := range conversationSeqs {
				req.ConversationID = conversationID
				if err := m.Conversation.UpdateConversations(ctx, req); err != nil {
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

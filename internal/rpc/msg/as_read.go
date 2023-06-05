package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
)

func (m *msgServer) MarkMsgsAsRead(ctx context.Context, req *msg.MarkMsgsAsReadReq) (resp *msg.MarkMsgsAsReadResp, err error) {
	conversations, err := m.Conversation.GetConversationsByConversationID(ctx, []string{req.ConversationID})
	if err != nil {
		return
	}
	var recvID string
	if conversations[0].ConversationType == constant.SingleChatType || conversations[0].ConversationType == constant.NotificationChatType {
		if req.UserID == conversations[0].OwnerUserID {
			recvID = conversations[0].UserID
		} else {
			recvID = conversations[0].OwnerUserID
		}
	} else if conversations[0].ConversationType == constant.SuperGroupChatType {
		recvID = conversations[0].GroupID
	}
	err = m.MsgDatabase.MarkSingleChatMsgsAsRead(ctx, req.ConversationID, req.UserID, req.Seqs)
	if err != nil {
		return
	}
	if err = m.sendMarkAsReadNotification(ctx, req.ConversationID, req.UserID, recvID, req.Seqs); err != nil {
		return
	}
	return &msg.MarkMsgsAsReadResp{}, nil
}

func (m *msgServer) sendMarkAsReadNotification(ctx context.Context, conversationID string, sendID, recvID string, seqs []int64) error {
	// tips := &sdkws.MarkAsReadTips{
	// 	MarkAsReadUserID: sendID,
	// 	ConversationID:   conversationID,
	// 	Seqs:             seqs,
	// }
	// m.notificationSender.NotificationWithSesstionType(ctx)
	return nil
}

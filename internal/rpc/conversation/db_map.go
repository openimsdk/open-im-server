package conversation

import (
	"context"

	dbModel "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/conversation"
)

func UpdateConversationsMap(ctx context.Context, req *conversation.SetConversationsReq) (m map[string]any, conversation dbModel.Conversation, err error) {
	m = make(map[string]any)

	conversation.ConversationID = req.Conversation.ConversationID
	conversation.ConversationType = req.Conversation.ConversationType
	conversation.UserID = req.Conversation.UserID
	conversation.GroupID = req.Conversation.GroupID

	if req.Conversation.RecvMsgOpt != nil {
		conversation.RecvMsgOpt = req.Conversation.RecvMsgOpt.Value
		m["recv_msg_opt"] = req.Conversation.RecvMsgOpt.Value
	}

	if req.Conversation.AttachedInfo != nil {
		conversation.AttachedInfo = req.Conversation.AttachedInfo.Value
		m["attached_info"] = req.Conversation.AttachedInfo.Value
	}

	if req.Conversation.Ex != nil {
		conversation.Ex = req.Conversation.Ex.Value
		m["ex"] = req.Conversation.Ex.Value
	}
	if req.Conversation.IsPinned != nil {
		conversation.IsPinned = req.Conversation.IsPinned.Value
		m["is_pinned"] = req.Conversation.IsPinned.Value
	}
	if req.Conversation.GroupAtType != nil {
		conversation.GroupAtType = req.Conversation.GroupAtType.Value
		m["group_at_type"] = req.Conversation.GroupAtType.Value
	}
	if req.Conversation.MsgDestructTime != nil {
		conversation.MsgDestructTime = req.Conversation.MsgDestructTime.Value
		m["msg_destruct_time"] = req.Conversation.MsgDestructTime.Value
	}
	if req.Conversation.IsMsgDestruct != nil {
		conversation.IsMsgDestruct = req.Conversation.IsMsgDestruct.Value
		m["is_msg_destruct"] = req.Conversation.IsMsgDestruct.Value
	}
	if req.Conversation.BurnDuration != nil {
		conversation.BurnDuration = req.Conversation.BurnDuration.Value
		m["burn_duration"] = req.Conversation.BurnDuration.Value
	}

	return m, conversation, nil
}

func UserUpdateCheckMap(ctx context.Context, userID string, req *conversation.ConversationReq, conversation *dbModel.Conversation) (unequal bool) {
	unequal = false

	if req.RecvMsgOpt != nil && conversation.RecvMsgOpt != req.RecvMsgOpt.Value {
		unequal = true
	}
	if req.AttachedInfo != nil && conversation.AttachedInfo != req.AttachedInfo.Value {
		unequal = true
	}
	if req.Ex != nil && conversation.Ex != req.Ex.Value {
		unequal = true
	}
	if req.IsPinned != nil && conversation.IsPinned != req.IsPinned.Value {
		unequal = true
	}
	if req.GroupAtType != nil && conversation.GroupAtType != req.GroupAtType.Value {
		unequal = true
	}
	if req.MsgDestructTime != nil && conversation.MsgDestructTime != req.MsgDestructTime.Value {
		unequal = true
	}
	if req.IsMsgDestruct != nil && conversation.IsMsgDestruct != req.IsMsgDestruct.Value {
		unequal = true
	}
	if req.BurnDuration != nil && conversation.BurnDuration != req.BurnDuration.Value {
		unequal = true
	}

	return unequal
}

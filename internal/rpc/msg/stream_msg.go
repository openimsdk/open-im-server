package msg

import (
	"context"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"time"
)

const StreamDeadlineTime = time.Second * 60 * 10

func (m *msgServer) handlerStreamMsg(ctx context.Context, msgData *sdkws.MsgData) error {
	now := time.Now()
	val := &model.StreamMsg{
		ClientMsgID:    msgData.ClientMsgID,
		ConversationID: msgprocessor.GetConversationIDByMsg(msgData),
		UserID:         msgData.SendID,
		CreateTime:     now,
		DeadlineTime:   now.Add(StreamDeadlineTime),
	}
	return m.StreamMsgDatabase.CreateStreamMsg(ctx, val)
}

func (m *msgServer) getStreamMsg(ctx context.Context, clientMsgID string) (*model.StreamMsg, error) {
	res, err := m.StreamMsgDatabase.GetStreamMsg(ctx, clientMsgID)
	if err != nil {
		return nil, err
	}
	if !res.End && res.DeadlineTime.Before(time.Now()) {
		res.End = true
	}
	return res, nil
}

func (m *msgServer) AppendStreamMsg(ctx context.Context, req *msg.AppendStreamMsgReq) (*msg.AppendStreamMsgResp, error) {
	res, err := m.getStreamMsg(ctx, req.ClientMsgID)
	if err != nil {
		return nil, err
	}
	if res.End {
		return nil, errs.ErrNoPermission.WrapMsg("stream msg is end")
	}
	if len(res.Packets) < int(req.StartIndex) {
		return nil, errs.ErrNoPermission.WrapMsg("start index is invalid")
	}
	if val := len(res.Packets) - int(req.StartIndex); val > 0 {
		exist := res.Packets[int(req.StartIndex):]
		for i, s := range exist {
			if len(req.Packets) == 0 {
				break
			}
			if s != req.Packets[i] {
				return nil, errs.ErrNoPermission.WrapMsg(fmt.Sprintf("packet %d has been written and is inconsistent", i))
			}
			req.StartIndex++
			req.Packets = req.Packets[1:]
		}
	}
	if len(req.Packets) == 0 && res.End == req.End {
		return &msg.AppendStreamMsgResp{}, nil
	}
	if err := m.StreamMsgDatabase.AppendStreamMsg(ctx, req.ClientMsgID, int(req.StartIndex), req.Packets, req.End); err != nil {
		return nil, err
	}
	conversation, err := m.Conversation.GetConversation(ctx, res.UserID, res.ConversationID)
	if err != nil {
		return nil, err
	}
	tips := &sdkws.StreamMsgTips{
		ClientMsgID: res.ClientMsgID,
		StartIndex:  req.StartIndex,
		Packets:     req.Packets,
		End:         req.End,
	}
	var (
		recvID      string
		sessionType int32
	)
	if conversation.GroupID == "" {
		sessionType = constant.SingleChatType
		recvID = conversation.UserID
	} else {
		sessionType = constant.ReadGroupChatType
		recvID = conversation.GroupID
	}
	m.msgNotificationSender.StreamMsgNotification(ctx, res.UserID, recvID, sessionType, tips)
	return &msg.AppendStreamMsgResp{}, nil
}

func (m *msgServer) GetStreamMsg(ctx context.Context, req *msg.GetStreamMsgReq) (*msg.GetStreamMsgResp, error) {
	res, err := m.getStreamMsg(ctx, req.ClientMsgID)
	if err != nil {
		return nil, err
	}
	return &msg.GetStreamMsgResp{
		ClientMsgID:    res.ClientMsgID,
		ConversationID: res.ConversationID,
		UserID:         res.UserID,
		Packets:        res.Packets,
		End:            res.End,
		CreateTime:     res.CreateTime.UnixMilli(),
		DeadlineTime:   res.DeadlineTime.UnixMilli(),
	}, nil
}

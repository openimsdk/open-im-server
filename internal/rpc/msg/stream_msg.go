package msg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

const (
	StreamTimeoutEnd            = time.Minute * 10
	StreamTimeoutEndMillisecond = int64(StreamTimeoutEnd / time.Millisecond)
)

func (m *msgServer) createStreamMsgHandler(ctx context.Context, msgData *sdkws.MsgData) error {
	var elem apistruct.StreamMsgElem
	if err := json.Unmarshal(msgData.Content, &elem); err != nil {
		return errs.ErrArgs.WrapMsg("stream msg content is invalid", "content", string(msgData.Content))
	}
	conversationID := msgprocessor.GetConversationIDByMsg(msgData)
	if _, err := m.StreamMsgDatabase.GetStreamMsg(ctx, conversationID, msgData.ClientMsgID); err != nil {
		if !errs.ErrRecordNotFound.Is(err) {
			return err
		}
	}
	streamMsg := &cache.StreamMsg{
		SendUserID:    msgData.SendID,
		RecvID:        msgData.RecvID,
		SessionType:   msgData.SessionType,
		UpdateTime:    time.Now().UnixMilli(),
		StreamType:    elem.Type,
		StreamContent: elem.Content,
	}
	switch msgData.SessionType {
	case constant.ReadGroupChatType, constant.WriteGroupChatType:
		streamMsg.RecvID = msgData.GroupID
	}
	return m.StreamMsgDatabase.CreateStreamMsg(ctx, msgprocessor.GetConversationIDByMsg(msgData), msgData.ClientMsgID, streamMsg)
}

func (m *msgServer) AppendStreamMsg(ctx context.Context, req *msg.AppendStreamMsgReq) (*msg.AppendStreamMsgResp, error) {
	end, err := m.StreamMsgDatabase.GetStreamMsgEnd(ctx, req.ConversationID, req.ClientMsgID)
	if err != nil {
		return nil, err
	}
	if end {
		return nil, errs.ErrNoPermission.WrapMsg("stream msg is end")
	}
	res, err := m.StreamMsgDatabase.AppendStreamMsg(ctx, req.ConversationID, req.ClientMsgID, int(req.StartIndex), req.Packets, req.End, req.End)
	if err != nil {
		return nil, err
	}
	tips := &sdkws.StreamMsgTips{
		ConversationID: req.ConversationID,
		ClientMsgID:    req.ClientMsgID,
		StartIndex:     req.StartIndex,
		Packets:        req.Packets,
		End:            req.End,
	}
	m.msgNotificationSender.StreamMsgNotification(ctx, res.SendUserID, res.RecvID, res.SessionType, tips)
	if req.End {
		m.modifyStreamMessage(ctx, req.ConversationID, req.ClientMsgID, res)
	}
	return &msg.AppendStreamMsgResp{}, nil
}

func (m *msgServer) modifyStreamMessage(ctx context.Context, conversationID string, clientMsgID string, res *cache.StreamMsg) {
	packets := make([]string, 0, len(res.Packets))
	for i := int64(0); ; i++ {
		data, ok := res.Packets[i]
		if !ok {
			break
		}
		packets = append(packets, data)
	}
	content, err := json.Marshal(&apistruct.StreamMsgElem{
		Type:     res.StreamType,
		Content:  res.StreamContent,
		Packets:  packets,
		End:      res.End,
		Deadline: res.UpdateTime,
	})
	if err != nil {
		log.ZError(ctx, "modifyStreamMessage json.Marshal", err, "conversationID", conversationID, "clientMsgID", clientMsgID)
		return
	}
	req := &msg.ModifyMessageReq{
		ConversationID: conversationID,
		NewContent:     string(content),
	}
	modifyMessage := func() error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		if req.Seq == 0 {
			req.Seq, err = m.MsgDatabase.GetMessageSeq(ctx, conversationID, clientMsgID)
			if err != nil {
				return err
			}
		}
		if _, err := m.ModifyMessage(ctx, req); err != nil {
			return err
		}
		return nil
	}

	if err := modifyMessage(); err != nil {
		log.ZError(ctx, "sync modifyStreamMessage", err, "conversationID", conversationID, "content", string(content))
		ctx = context.WithoutCancel(ctx)
		go func() {
			for i := 1; i <= 10; i++ {
				if err := modifyMessage(); err == nil {
					log.ZDebug(ctx, "async modifyStreamMessage success", "conversationID", conversationID, "content", string(content), "count", i)
					return
				} else {
					log.ZError(ctx, "modifyStreamMessage", err, "conversationID", conversationID, "content", string(content), "count", i)
					time.Sleep(time.Second * time.Duration(i))
				}
			}
		}()
	}
}

func (m *msgServer) GetStreamMsg(ctx context.Context, req *msg.GetStreamMsgReq) (*msg.GetStreamMsgResp, error) {
	value, err := m.StreamMsgDatabase.GetStreamMsg(ctx, req.ConversationID, req.ClientMsgID)
	if err == nil {
		resp := msg.GetStreamMsgResp{
			UserID:  value.SendUserID,
			Packets: make([]string, 0, len(value.Packets)),
			End:     value.End,
		}
		for i := int64(0); ; i++ {
			data, ok := value.Packets[i]
			if !ok {
				break
			}
			resp.Packets = append(resp.Packets, data)
		}
		if resp.End {
			resp.DeadlineTime = value.UpdateTime
		} else {
			if now := time.Now().UnixMilli(); now-value.UpdateTime >= StreamTimeoutEndMillisecond {
				resp.DeadlineTime = now + StreamTimeoutEndMillisecond
				resp.End = true
			}
		}
		return &resp, nil
	} else if !errs.ErrRecordNotFound.Is(err) {
		return nil, err
	}
	if req.Seq <= 0 || errs.ErrRecordNotFound.Is(err) == false {
		return nil, err
	}
	msgs, err := m.MsgDatabase.GetMessageBySeqs(ctx, req.ConversationID, mcontext.GetOpUserID(ctx), []int64{req.Seq})
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 || msgs[0] == nil {
		return nil, errs.ErrRecordNotFound.WrapMsg("stream message not found")
	}
	msgData := msgs[0]
	if msgData.ClientMsgID != req.ClientMsgID {
		return nil, errs.ErrRecordNotFound.WrapMsg("stream message id not match")
	}
	if msgData.ContentType != constant.Stream {
		return nil, errs.ErrNoPermission.WrapMsg("stream message content type not match")
	}
	var elem apistruct.StreamMsgElem
	if len(msgData.Content) > 0 {
		if err := json.Unmarshal(msgData.Content, &elem); err != nil {
			log.ZError(ctx, "stream msg unmarshal", err, "content", string(msgData.Content), "conversationID", req.ConversationID, "seq", req.Seq)
		}
	}
	resp := &msg.GetStreamMsgResp{
		UserID:       msgData.SendID,
		Packets:      elem.Packets,
		End:          elem.End,
		DeadlineTime: elem.Deadline,
	}
	if !resp.End {
		resp.End = true
		resp.DeadlineTime = msgData.SendTime + StreamTimeoutEndMillisecond
	}
	return resp, nil
}

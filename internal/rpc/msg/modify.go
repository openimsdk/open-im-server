package msg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/constant"
	msgpb "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

func (m *msgServer) getModifyRawMessage(ctx context.Context, req *msgpb.ModifyMessageReq) (*model.MsgDataModel, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	msgs, err := m.MsgDatabase.GetMessageBySeqsDB(ctx, req.ConversationID, opUserID, []int64{req.Seq})
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg("msg seq not found")
	}
	val := msgs[0]
	if val == nil || val.Msg == nil || val.Msg.Status == constant.MsgStatusHasDeleted {
		return nil, servererrs.ErrRecordNotFound.WrapMsg("msg already delete")
	}
	if val.Revoke != nil {
		return nil, servererrs.ErrMsgAlreadyRevoke.WrapMsg("msg already revoke")
	}
	msgData := val.Msg
	if req.OldContent != "" {
		if req.OldContent != msgData.Content {
			return nil, servererrs.ErrArgs.WrapMsg("old msg content not match")
		}
	}
	if req.NewContent == msgData.Content {
		return nil, errs.ErrArgs.WrapMsg("new content same as old content")
	}
	if datautil.Contain(opUserID, m.config.Share.IMAdminUser.UserIDs...) {
		return msgData, nil
	}
	isGroup := msgprocessor.IsGroupConversationID(req.ConversationID)
	if !isGroup {
		if msgData.SendID != opUserID {
			return nil, servererrs.ErrNoPermission.WrapMsg("no permission")
		}
		return msgData, nil
	}
	groupID := msgData.GroupID
	if groupID == "" {
		groupID = msgData.RecvID
	}
	groupInfo, err := m.GroupLocalCache.GetGroupInfo(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.Wrap()
	}
	var memberUserIDs []string
	if msgData.SendID == opUserID {
		memberUserIDs = []string{opUserID}
	} else {
		memberUserIDs = []string{opUserID, msgData.SendID}
	}
	members, err := m.GroupLocalCache.GetGroupMemberInfoMap(ctx, groupID, memberUserIDs)
	if err != nil {
		return nil, err
	}
	opMember, ok := members[opUserID]
	if !ok {
		return nil, servererrs.ErrNoPermission.WrapMsg("opUser no in group")
	}
	if msgData.SendID == opUserID {
		return msgData, nil
	}
	if opMember.RoleLevel <= constant.GroupOrdinaryUsers {
		return nil, errs.ErrNoPermission.WrapMsg("no permission update other user msg")
	}
	var sendRoleLevel int32
	if sendMember, ok := members[msgData.SendID]; ok {
		sendRoleLevel = sendMember.RoleLevel
	}
	if sendRoleLevel >= opMember.RoleLevel {
		return nil, errs.ErrNoPermission.WrapMsg("no permission update other user msg")
	}
	return msgData, nil
}

func (m *msgServer) ModifyMessage(ctx context.Context, req *msgpb.ModifyMessageReq) (*msgpb.ModifyMessageResp, error) {
	lockKey := fmt.Sprintf("MODIFYMESSAGE:%s:%d", req.ConversationID, req.Seq)
	lockValue, err := m.lock.Lock(ctx, lockKey, time.Second*30)
	if err != nil {
		return nil, err
	}
	defer m.lock.Unlock(ctx, lockKey, lockValue)
	msg, err := m.getModifyRawMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	var attachedInfo map[string]json.RawMessage
	if msg.AttachedInfo != "" && msg.AttachedInfo != "null" && msg.AttachedInfo != "{}" {
		if err = json.Unmarshal([]byte(msg.AttachedInfo), &attachedInfo); err != nil {
			log.ZWarn(ctx, "json.Unmarshal", err, "attachedInfo", msg.AttachedInfo)
		}
	}
	if attachedInfo == nil {
		attachedInfo = make(map[string]json.RawMessage)
	}
	const modifyAttachedKey = "lastModified"
	type LastModified struct {
		UserID        string `json:"userID"`        // last modified user ID
		ModifiedTime  int64  `json:"modifiedTime"`  // last modified time
		ModifiedCount int64  `json:"modifiedCount"` // last modified count
	}
	var modifyValue LastModified
	if val := attachedInfo[modifyAttachedKey]; len(val) > 0 {
		if err = json.Unmarshal(val, &modifyValue); err != nil {
			return nil, errs.WrapMsg(err, "json.Unmarshal modifyValue", "val", val)
		}
		if modifyValue.ModifiedCount < 1 {
			modifyValue.ModifiedCount = 1
		}
	}
	modifyValue.ModifiedCount++
	modifyValue.ModifiedTime = time.Now().UnixMilli()
	modifyValue.UserID = mcontext.GetOpUserID(ctx)
	modifyVal, err := json.Marshal(&modifyValue)
	if err != nil {
		return nil, err
	}
	attachedInfo[modifyAttachedKey] = modifyVal
	attached, err := json.Marshal(attachedInfo)
	if err != nil {
		return nil, errs.ErrInternalServer.WrapMsg("json.Marshal attachedInfo", "attachedInfo", attachedInfo)
	}
	msg.Content = req.NewContent
	msg.AttachedInfo = string(attached)
	if err := m.MsgDatabase.UpdateMsg(ctx, req.ConversationID, msg); err != nil {
		return nil, err
	}
	tips := &sdkws.ModifyMsgTips{
		ConversationID: req.ConversationID,
		Seq:            req.Seq,
		ClientMsgID:    msg.ClientMsgID,
		NewContent:     req.NewContent,
		ModifiedTime:   modifyValue.ModifiedTime,
		ModifiedCount:  modifyValue.ModifiedCount,
		UserID:         modifyValue.UserID,
	}
	recvID := msg.GroupID
	if recvID == "" {
		recvID = msg.RecvID
	}
	m.notificationSender.NotificationWithSessionType(ctx, msg.SendID, recvID, constant.ModifyMessageNotification, msg.SessionType, tips)
	return &msgpb.ModifyMessageResp{
		ModifiedTime:  modifyValue.ModifiedTime,
		ModifiedCount: modifyValue.ModifiedCount,
	}, nil
}

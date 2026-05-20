// Copyright © 2026 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package group

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/protocol/constant"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

// 群置顶消息相关 RPC 实现：
// - 自动滚动保留最近 N 条置顶消息（N=model.GroupPinnedMsgMaxKeep，默认为 3）
// - 置顶时记录消息 seq，并按 seq 做完整内容快照存档
// - 返回置顶列表与群通知均按各用户会话 minSeq/maxSeq 过滤可见置顶
// - 每条置顶记录拥有唯一 pinID，作为 unpin 时的精准删除凭据
// - 权限：默认全员可置顶；当 group.AllowPinMsg=1 时，仅群主/管理员可置顶或取消置顶

const (
	groupPinnedActionPin   = int32(1)
	groupPinnedActionUnpin = int32(2)
)

// PinGroupMessage 群聊中置顶单条消息
func (s *groupServer) PinGroupMessage(ctx context.Context, req *pbgroup.PinGroupMessageReq) (*pbgroup.PinGroupMessageResp, error) {
	if req.GroupID == "" {
		return nil, errs.ErrArgs.WrapMsg("groupID empty")
	}
	if req.Seq <= 0 {
		return nil, errs.ErrArgs.WrapMsg("seq must be positive")
	}

	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.Wrap()
	}

	if err := s.checkPinPermission(ctx, group); err != nil {
		return nil, err
	}

	conversationID := msgprocessor.GetConversationIDBySessionType(constant.ReadGroupChatType, req.GroupID)
	msgData, err := s.msgClient.GetSingleMsgBySeq(ctx, conversationID, req.Seq)
	if err != nil {
		return nil, err
	}
	if msgData == nil {
		return nil, servererrs.ErrRecordNotFound.WrapMsg("message not found by seq")
	}
	if msgData.GroupID != "" && msgData.GroupID != req.GroupID {
		return nil, errs.ErrArgs.WrapMsg("seq does not belong to this group")
	}
	if msgData.Status >= constant.MsgStatusHasDeleted {
		return nil, servererrs.ErrRecordNotFound.WrapMsg("message has been deleted")
	}

	pin := buildPinSnapshot(req.GroupID, conversationID, mcontext.GetOpUserID(ctx), msgData)

	pinnedList, err := s.db.PinGroupMessage(ctx, req.GroupID, pin)
	if err != nil {
		return nil, err
	}

	pbPinned := pinnedMsgDB2PB(pin)
	pbListAll := pinnedListDB2PB(pinnedList)
	visibleList, err := s.filterPinnedMsgsByUserSeq(ctx, req.GroupID, mcontext.GetOpUserID(ctx), pinnedList)
	if err != nil {
		return nil, err
	}

	s.notification.GroupMessagePinnedNotification(ctx, req.GroupID, groupPinnedActionPin, pbPinned, pbListAll)

	return &pbgroup.PinGroupMessageResp{
		PinnedMsg:  pbPinned,
		PinnedList: pinnedListDB2PB(visibleList),
	}, nil
}

// UnpinGroupMessage 群聊中取消置顶单条消息（pinID 优先；为空则按 seq）
func (s *groupServer) UnpinGroupMessage(ctx context.Context, req *pbgroup.UnpinGroupMessageReq) (*pbgroup.UnpinGroupMessageResp, error) {
	if req.GroupID == "" {
		return nil, errs.ErrArgs.WrapMsg("groupID empty")
	}
	if req.PinID == "" && req.Seq <= 0 {
		return nil, errs.ErrArgs.WrapMsg("either pinID or seq must be provided")
	}

	group, err := s.db.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, servererrs.ErrDismissedAlready.Wrap()
	}
	if err := s.checkPinPermission(ctx, group); err != nil {
		return nil, err
	}

	current, err := s.db.GetGroupPinnedMessages(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	var target *model.GroupPinnedMessage
	for _, m := range current {
		if req.PinID != "" {
			if m.PinID == req.PinID {
				target = m
				break
			}
		} else if m.Seq == req.Seq {
			target = m
			break
		}
	}
	if target == nil {
		return nil, servererrs.ErrRecordNotFound.WrapMsg("pinned message not found")
	}

	pinnedList, err := s.db.UnpinGroupMessage(ctx, req.GroupID, req.PinID, req.Seq)
	if err != nil {
		return nil, err
	}

	pbPinned := pinnedMsgDB2PB(target)
	pbListAll := pinnedListDB2PB(pinnedList)
	visibleList, err := s.filterPinnedMsgsByUserSeq(ctx, req.GroupID, mcontext.GetOpUserID(ctx), pinnedList)
	if err != nil {
		return nil, err
	}

	s.notification.GroupMessagePinnedNotification(ctx, req.GroupID, groupPinnedActionUnpin, pbPinned, pbListAll)

	return &pbgroup.UnpinGroupMessageResp{PinnedList: pinnedListDB2PB(visibleList)}, nil
}

// GetGroupPinnedMessages 获取群置顶消息列表
func (s *groupServer) GetGroupPinnedMessages(ctx context.Context, req *pbgroup.GetGroupPinnedMessagesReq) (*pbgroup.GetGroupPinnedMessagesResp, error) {
	if req.GroupID == "" {
		return nil, errs.ErrArgs.WrapMsg("groupID empty")
	}
	if err := s.checkAdminOrInGroup(ctx, req.GroupID); err != nil {
		return nil, err
	}
	pinnedList, err := s.db.GetGroupPinnedMessages(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	userID := mcontext.GetOpUserID(ctx)
	visibleList, err := s.filterPinnedMsgsByUserSeq(ctx, req.GroupID, userID, pinnedList)
	if err != nil {
		return nil, err
	}
	return &pbgroup.GetGroupPinnedMessagesResp{
		PinnedList: pinnedListDB2PB(visibleList),
	}, nil
}

// filterPinnedMsgsByUserSeq 按用户在该群会话的 minSeq/maxSeq 过滤置顶消息。
// 与拉取群消息可见范围一致：新成员（minSeq 被抬高）看不到入群前的置顶；被踢成员（maxSeq 受限）看不到踢出后的置顶。
func (s *groupServer) filterPinnedMsgsByUserSeq(ctx context.Context, groupID, userID string, list []*model.GroupPinnedMessage) ([]*model.GroupPinnedMessage, error) {
	if len(list) == 0 {
		return nil, nil
	}
	if authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		return list, nil
	}
	if userID == "" {
		return nil, errs.ErrNoPermission.WrapMsg("op user id empty")
	}
	conversationID := msgprocessor.GetConversationIDBySessionType(constant.ReadGroupChatType, groupID)
	conv, err := s.conversationClient.GetConversation(ctx, conversationID, userID)
	if err != nil {
		if errs.ErrRecordNotFound.Is(err) {
			return nil, nil
		}
		return nil, err
	}
	out := make([]*model.GroupPinnedMessage, 0, len(list))
	for _, m := range list {
		if m == nil || !isPinnedSeqVisible(m.Seq, conv.MinSeq, conv.MaxSeq) {
			continue
		}
		out = append(out, m)
	}
	return out, nil
}

func isPinnedSeqVisible(seq, minSeq, maxSeq int64) bool {
	if seq <= 0 {
		return false
	}
	if minSeq > 0 && seq < minSeq {
		return false
	}
	if maxSeq > 0 && seq > maxSeq {
		return false
	}
	return true
}

func filterPinnedListPB(list []*sdkws.GroupPinnedMsgInfo, minSeq, maxSeq int64) []*sdkws.GroupPinnedMsgInfo {
	if len(list) == 0 {
		return nil
	}
	out := make([]*sdkws.GroupPinnedMsgInfo, 0, len(list))
	for _, m := range list {
		if m == nil || !isPinnedSeqVisible(m.Seq, minSeq, maxSeq) {
			continue
		}
		out = append(out, m)
	}
	return out
}

func pinnedMsgPBVisibleToUser(pinned *sdkws.GroupPinnedMsgInfo, minSeq, maxSeq int64) *sdkws.GroupPinnedMsgInfo {
	if pinned == nil || !isPinnedSeqVisible(pinned.Seq, minSeq, maxSeq) {
		return nil
	}
	return pinned
}

// checkPinPermission 校验当前操作者是否具备群消息置顶权限
func (s *groupServer) checkPinPermission(ctx context.Context, group *model.Group) error {
	if authverify.IsAppManagerUid(ctx, s.config.Share.IMAdminUserID) {
		return nil
	}
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		return errs.ErrNoPermission.WrapMsg("op user id empty")
	}
	member, err := s.db.TakeGroupMember(ctx, group.GroupID, opUserID)
	if err != nil {
		return err
	}
	isOwnerOrAdmin := member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin
	if group.AllowPinMsg == model.GroupPermAdminOnly && !isOwnerOrAdmin {
		return errs.ErrNoPermission.WrapMsg("only owner or admin can pin/unpin group message")
	}
	return nil
}

// buildPinSnapshot 把 sdkws.MsgData 完整快照成 GroupPinnedMessage
// PinID 在 mgo 层 Pin 时若为空会自动生成；这里留空交由存储层处理
func buildPinSnapshot(groupID, conversationID, opUserID string, m *sdkws.MsgData) *model.GroupPinnedMessage {
	pin := &model.GroupPinnedMessage{
		GroupID:          groupID,
		ConversationID:   conversationID,
		Seq:              m.Seq,
		ServerMsgID:      m.ServerMsgID,
		ClientMsgID:      m.ClientMsgID,
		SendID:           m.SendID,
		RecvID:           m.RecvID,
		SenderPlatformID: m.SenderPlatformID,
		SenderNickname:   m.SenderNickname,
		SenderFaceURL:    m.SenderFaceURL,
		SessionType:      m.SessionType,
		MsgFrom:          m.MsgFrom,
		ContentType:      m.ContentType,
		Content:          string(m.Content),
		AtUserIDList:     append([]string(nil), m.AtUserIDList...),
		Options:          copyOptions(m.Options),
		AttachedInfo:     m.AttachedInfo,
		Ex:               m.Ex,
		SendTime:         m.SendTime,
		CreateTime:       m.CreateTime,
		Status:           m.Status,
		PinUserID:        opUserID,
		PinTime:          time.Now().UnixMilli(),
	}
	if m.OfflinePushInfo != nil {
		pin.OfflinePush = &model.GroupPinnedOfflinePush{
			Title:         m.OfflinePushInfo.Title,
			Desc:          m.OfflinePushInfo.Desc,
			Ex:            m.OfflinePushInfo.Ex,
			IOSPushSound:  m.OfflinePushInfo.IOSPushSound,
			IOSBadgeCount: m.OfflinePushInfo.IOSBadgeCount,
			SignalInfo:    m.OfflinePushInfo.SignalInfo,
		}
	}
	return pin
}

func copyOptions(src map[string]bool) map[string]bool {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]bool, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func pinnedMsgDB2PB(m *model.GroupPinnedMessage) *sdkws.GroupPinnedMsgInfo {
	if m == nil {
		return nil
	}
	return &sdkws.GroupPinnedMsgInfo{
		PinID:            m.PinID,
		GroupID:          m.GroupID,
		ConversationID:   m.ConversationID,
		Seq:              m.Seq,
		ServerMsgID:      m.ServerMsgID,
		ClientMsgID:      m.ClientMsgID,
		SendID:           m.SendID,
		RecvID:           m.RecvID,
		SenderPlatformID: m.SenderPlatformID,
		SenderNickname:   m.SenderNickname,
		SenderFaceURL:    m.SenderFaceURL,
		SessionType:      m.SessionType,
		MsgFrom:          m.MsgFrom,
		ContentType:      m.ContentType,
		Content:          m.Content,
		AtUserIDList:     append([]string(nil), m.AtUserIDList...),
		Options:          copyOptions(m.Options),
		AttachedInfo:     m.AttachedInfo,
		Ex:               m.Ex,
		SendTime:         m.SendTime,
		CreateTime:       m.CreateTime,
		Status:           m.Status,
		PinUserID:        m.PinUserID,
		PinTime:          m.PinTime,
	}
}

func pinnedListDB2PB(list []*model.GroupPinnedMessage) []*sdkws.GroupPinnedMsgInfo {
	if len(list) == 0 {
		return nil
	}
	result := make([]*sdkws.GroupPinnedMsgInfo, 0, len(list))
	for _, m := range list {
		result = append(result, pinnedMsgDB2PB(m))
	}
	return result
}

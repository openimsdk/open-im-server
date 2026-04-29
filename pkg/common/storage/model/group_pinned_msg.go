// Copyright © 2026 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package model

// GroupPinnedMsgMaxKeep 群置顶消息最多保留的条数（最新置顶的在最前）
const GroupPinnedMsgMaxKeep = 3

// GroupPinnedOfflinePush 离线推送信息快照
type GroupPinnedOfflinePush struct {
	Title         string `bson:"title"`
	Desc          string `bson:"desc"`
	Ex            string `bson:"ex"`
	IOSPushSound  string `bson:"ios_push_sound"`
	IOSBadgeCount bool   `bson:"ios_badge_count"`
	SignalInfo    string `bson:"signal_info"`
}

// GroupPinnedMessage 一条群置顶消息的完整内容快照
// 置顶时把消息整体快照入库，避免后续消息删除/撤回影响已置顶展示
type GroupPinnedMessage struct {
	// PinID 全局唯一 id，用于精准取消置顶（生产由 mongo ObjectID().Hex() 生成）
	PinID string `bson:"pin_id"`

	// 会话 / 群信息
	ConversationID string `bson:"conversation_id"`
	GroupID        string `bson:"group_id"`

	// 消息标识
	Seq         int64  `bson:"seq"`
	ServerMsgID string `bson:"server_msg_id"`
	ClientMsgID string `bson:"client_msg_id"`

	// 发送方信息
	SendID           string `bson:"send_id"`
	RecvID           string `bson:"recv_id"`
	SenderPlatformID int32  `bson:"sender_platform_id"`
	SenderNickname   string `bson:"sender_nickname"`
	SenderFaceURL    string `bson:"sender_face_url"`

	// 消息内容快照
	SessionType  int32           `bson:"session_type"`
	MsgFrom      int32           `bson:"msg_from"`
	ContentType  int32           `bson:"content_type"`
	Content      string          `bson:"content"`
	AtUserIDList []string        `bson:"at_user_id_list"`
	Options      map[string]bool `bson:"options"`
	AttachedInfo string          `bson:"attached_info"`
	Ex           string          `bson:"ex"`

	OfflinePush *GroupPinnedOfflinePush `bson:"offline_push"`

	// 时间
	SendTime   int64 `bson:"send_time"`
	CreateTime int64 `bson:"create_time"`
	Status     int32 `bson:"status"`

	// 操作人 & 时间
	PinUserID string `bson:"pin_user_id"`
	PinTime   int64  `bson:"pin_time"`
}

// GroupPinnedMsg 一个群的置顶消息文档，按 group_id 唯一
type GroupPinnedMsg struct {
	GroupID    string                `bson:"group_id"`
	PinnedMsgs []*GroupPinnedMessage `bson:"pinned_msgs"`
}

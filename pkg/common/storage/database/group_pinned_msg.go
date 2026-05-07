// Copyright © 2026 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// GroupPinnedMsg 群置顶消息的存储抽象
type GroupPinnedMsg interface {
	// Pin 置顶一条消息：若 PinID 为空会自动生成；自动滚动保留最近 N 条
	Pin(ctx context.Context, groupID string, msg *model.GroupPinnedMessage) ([]*model.GroupPinnedMessage, error)
	// Unpin 取消置顶；pinID 非空时按 pinID 精确删除，否则按 seq 删除
	Unpin(ctx context.Context, groupID string, pinID string, seq int64) ([]*model.GroupPinnedMessage, error)
	// Get 获取群置顶消息列表（最新的在前）
	Get(ctx context.Context, groupID string) ([]*model.GroupPinnedMessage, error)
}

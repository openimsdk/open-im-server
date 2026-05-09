// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// PhoneSN 手机号 is_snd 持久化
type PhoneSN interface {
	// GetByPhone 按手机号查询；无记录时返回 (nil, nil)
	GetByPhone(ctx context.Context, phone string) (*model.PhoneSNInfo, error)
	// Upsert 写入或更新 is_snd 与 user_id
	Upsert(ctx context.Context, phone string, userID int64, isSnd bool) error
	// DeleteByPhone 按手机号删除记录；记录不存在时不报错
	DeleteByPhone(ctx context.Context, phone string) error
}

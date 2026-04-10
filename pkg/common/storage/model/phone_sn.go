// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package model

// PhoneSNInfo 手机号与 is_snd、关联 user_id（每条以 phone 唯一）
type PhoneSNInfo struct {
	Phone      string `bson:"phone"`
	UserID     int64  `bson:"user_id"`
	IsSnd      bool   `bson:"is_snd"`
	UpdateTime int64  `bson:"update_time"`
}

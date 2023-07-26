// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"github.com/go-playground/validator/v10"

	"github.com/OpenIMSDK/protocol/constant"
)

func RequiredIf(fl validator.FieldLevel) bool {
	sessionType := fl.Parent().FieldByName("SessionType").Int()
	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		if fl.FieldName() == "RecvID" {
			return fl.Field().String() != ""
		}
	case constant.GroupChatType, constant.SuperGroupChatType:
		if fl.FieldName() == "GroupID" {
			return fl.Field().String() != ""
		}
	default:
		return true
	}
	return true
}

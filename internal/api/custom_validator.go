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

	"github.com/openimsdk/protocol/constant"
)

// RequiredIf validates if the specified field is required based on the session type.
func RequiredIf(fl validator.FieldLevel) bool {
	sessionType := fl.Parent().FieldByName("SessionType").Int()

	switch sessionType {
	case constant.SingleChatType, constant.NotificationChatType:
		return fl.FieldName() != "RecvID" || fl.Field().String() != ""
	case constant.WriteGroupChatType, constant.ReadGroupChatType:
		return fl.FieldName() != "GroupID" || fl.Field().String() != ""
	default:
		return true
	}
}

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

package sdkws

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

func (x *MsgData) Check() error {
	if x.SendID == "" {
		return errs.ErrArgs.Wrap("sendID is empty")
	}
	if x.Content == nil {
		return errs.ErrArgs.Wrap("content is empty")
	}
	if x.ContentType <= constant.ContentTypeBegin || x.ContentType >= constant.NotificationEnd {
		return errs.ErrArgs.Wrap("content type is invalid")
	}
	if x.SessionType < constant.SingleChatType || x.SessionType > constant.NotificationChatType {
		return errs.ErrArgs.Wrap("sessionType is invalid")
	}
	if x.SessionType == constant.SingleChatType || x.SessionType == constant.NotificationChatType {
		if x.RecvID == "" {
			return errs.ErrArgs.Wrap("recvID is empty")
		}
	}
	if x.SessionType == constant.GroupChatType || x.SessionType == constant.SuperGroupChatType {
		if x.GroupID == "" {
			return errs.ErrArgs.Wrap("GroupID is empty")
		}
	}
	return nil
}

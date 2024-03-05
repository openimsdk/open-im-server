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

package msg

import (
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func isMessageHasReadEnabled(msgData *sdkws.MsgData, config *config.GlobalConfig) bool {
	switch {
	case msgData.ContentType == constant.HasReadReceipt && msgData.SessionType == constant.SingleChatType:
		if config.SingleMessageHasReadReceiptEnable {
			return true
		} else {
			return false
		}
	case msgData.ContentType == constant.HasReadReceipt && msgData.SessionType == constant.SuperGroupChatType:
		if config.GroupMessageHasReadReceiptEnable {
			return true
		} else {
			return false
		}
	}
	return true
}

func IsNotFound(err error) bool {
	switch utils.Unwrap(err) {
	case redis.Nil, mongo.ErrNoDocuments:
		return true
	default:
		return false
	}
}

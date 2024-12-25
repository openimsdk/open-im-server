// Copyright © 2024 OpenIM. All rights reserved.
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

package cachekey

import (
	"strconv"
)

const (
	messageCache      = "MESSAGE_CACHE:"
	sendMsgFailedFlag = "SEND_MSG_FAILED_FLAG:"
	messageCacheV2    = "MESSAGE_CACHE_V2:"
)

func GetMessageCacheKey(conversationID string, seq int64) string {
	return messageCache + conversationID + "_" + strconv.Itoa(int(seq))
}

func GetMessageCacheKeyV2(conversationID string, seq int64) string {
	return messageCacheV2 + conversationID + "_" + strconv.Itoa(int(seq))
}

func GetSendMsgKey(id string) string {
	return sendMsgFailedFlag + id
}

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
	getuiToken              = "GETUI_TOKEN"
	getuiTaskID             = "GETUI_TASK_ID"
	fmcToken                = "FCM_TOKEN:"
	userBadgeUnreadCountSum = "USER_BADGE_UNREAD_COUNT_SUM:"
)

func GetFcmAccountTokenKey(account string, platformID int) string {
	return fmcToken + account + ":" + strconv.Itoa(platformID)
}

func GetUserBadgeUnreadCountSumKey(userID string) string {
	return userBadgeUnreadCountSum + userID
}

func GetGetuiTokenKey() string {
	return getuiToken
}
func GetGetuiTaskIDKey() string {
	return getuiTaskID
}

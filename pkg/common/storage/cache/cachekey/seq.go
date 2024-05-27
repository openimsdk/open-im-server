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

const (
	maxSeq                 = "MAX_SEQ:"
	minSeq                 = "MIN_SEQ:"
	conversationUserMinSeq = "CON_USER_MIN_SEQ:"
	hasReadSeq             = "HAS_READ_SEQ:"
)

func GetMaxSeqKey(conversationID string) string {
	return maxSeq + conversationID
}

func GetMinSeqKey(conversationID string) string {
	return minSeq + conversationID
}

func GetHasReadSeqKey(conversationID string, userID string) string {
	return hasReadSeq + userID + ":" + conversationID
}

func GetConversationUserMinSeqKey(conversationID, userID string) string {
	return conversationUserMinSeq + conversationID + "u:" + userID
}

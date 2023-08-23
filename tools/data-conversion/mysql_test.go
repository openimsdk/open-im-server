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

package data_conversion

import "testing"

// pass
func TestUserConversion(t *testing.T) {
	UserConversion()
}

// pass
func TestFriendConversion(t *testing.T) {
	FriendConversion()
}

// pass
func TestGroupConversion(t *testing.T) {
	GroupConversion()
	GroupMemberConversion()
}

// pass
func TestBlacksConversion(t *testing.T) {
	BlacksConversion()
}

// pass
func TestRequestConversion(t *testing.T) {
	RequestConversion()
}

// pass
func TestChatLogsConversion(t *testing.T) {
	// If the printed result is too long, the console will not display it, but it can run normally
	ChatLogsConversion()
}

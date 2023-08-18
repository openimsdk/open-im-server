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

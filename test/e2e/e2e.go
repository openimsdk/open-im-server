package e2e

import (
	"testing"

	gettoken "github.com/openimsdk/open-im-server/v3/test/e2e/api/token"
	"github.com/openimsdk/open-im-server/v3/test/e2e/api/user"
)

// RunE2ETests checks configuration parameters (specified through flags) and then runs
// E2E tests using the Ginkgo runner.
// If a "report directory" is specified, one or more JUnit test reports will be
// generated in this directory, and cluster logs will also be saved.
// This function is called on each Ginkgo node in parallel mode.
func RunE2ETests(t *testing.T) {

	// Example usage of new functions
	token, _ := gettoken.GetUserToken("openIM123456")

	// Example of getting user info
	_ = user.GetUsersInfo(token, []string{"user1", "user2"})

	// Example of updating user info
	_ = user.UpdateUserInfo(token, "user1", "NewNickname", "https://github.com/openimsdk/open-im-server/blob/main/assets/logo/openim-logo.png")

	// Example of getting users' online status
	_ = user.GetUsersOnlineStatus(token, []string{"user1", "user2"})

	// Example of forcing a logout
	_ = user.ForceLogout(token, "4950983283", 2)

	// Example of checking user account
	_ = user.CheckUserAccount(token, []string{"openIM123456", "anotherUserID"})

	// Example of getting users
	_ = user.GetUsers(token, 1, 100)
}

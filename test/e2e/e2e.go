package e2e

import (
	"testing"

	gettoken "github.com/openimsdk/open-im-server/v3/test/e2e/api/token"
	"github.com/openimsdk/open-im-server/v3/test/e2e/api/user"
)

func RunE2ETests(t *testing.T) {
	token, _ := gettoken.GetUserToken("openIM123456")
	_ = user.GetUsersInfo(token, []string{"user1", "user2"})
	_ = user.UpdateUserInfo(token, "user1", "NewNickname", "https://github.com/openimsdk/open-im-server/blob/main/assets/logo/openim-logo.png")
	_ = user.GetUsersOnlineStatus(token, []string{"user1", "user2"})
	_ = user.ForceLogout(token, "4950983283", 2)
	_ = user.CheckUserAccount(token, []string{"openIM123456", "anotherUserID"})
	_ = user.GetUsers(token, 1, 100)
}

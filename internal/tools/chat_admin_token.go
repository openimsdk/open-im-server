package tools

import (
	"context"

	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/tools/log"
)

// fetchChatAdminToken 通过 IM auth-rpc GetAdminToken 获取管理员 token。
// 使用 config.Share.Secret 和第一个 IMAdminUserID 作为凭据。
func (c *cronServer) fetchChatAdminToken(ctx context.Context) (string, error) {
	userID := c.config.Share.IMAdminUserID[0]
	resp, err := c.authClient.GetAdminToken(ctx, &auth.GetAdminTokenReq{
		Secret: c.config.Share.Secret,
		UserID: userID,
	})
	if err != nil {
		return "", err
	}
	log.ZDebug(ctx, "fetchChatAdminToken: ok", "userID", userID)
	return resp.Token, nil
}

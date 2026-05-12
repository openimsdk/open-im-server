package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

const deleteExpiredUserBatchLimit = 100

// chatHTTPClient 带超时，防止 chat 服务无响应时 cron worker 永久挂起。
var chatHTTPClient = &http.Client{Timeout: 3 * time.Second}

// deleteExpiredOfflineUsers 是 cron "@hourly" 触发的入口。
// 批量查询离线时长超过 delete_account_interval 的用户并依次调用 chat /account/del 删除。
func (c *cronServer) deleteExpiredOfflineUsers() {
	now := time.Now()
	operationID := fmt.Sprintf("cron_del_expired_user_%d_%d", os.Getpid(), now.UnixMilli())
	ctx := mcontext.SetOperationID(c.ctx, operationID)
	log.ZInfo(ctx, "deleteExpiredOfflineUsers: start", "time", now)

	users, err := c.userOfflineRecordDB.FindExpiredUsers(ctx, now, deleteExpiredUserBatchLimit)
	if err != nil {
		log.ZError(ctx, "deleteExpiredOfflineUsers: FindExpiredUsers failed", err)
		return
	}
	if len(users) == 0 {
		log.ZDebug(ctx, "deleteExpiredOfflineUsers: no expired users found")
		return
	}
	log.ZInfo(ctx, "deleteExpiredOfflineUsers: found expired users", "count", len(users))

	adminToken, err := c.fetchChatAdminToken(ctx)
	if err != nil {
		log.ZError(ctx, "deleteExpiredOfflineUsers: fetchChatAdminToken failed", err)
		return
	}

	for i, u := range users {
		subCtx := mcontext.SetOperationID(c.ctx, fmt.Sprintf("%s_%d", operationID, i))
		c.deleteExpiredUser(subCtx, adminToken, u.UserID)
	}
	log.ZInfo(ctx, "deleteExpiredOfflineUsers: done", "count", len(users), "elapsed", time.Since(now))
}

// deleteExpiredUser 通过 chat HTTP API POST /account/del 删除单个过期用户。
// chat 服务端会处理：强制登出、删除好友/群组关系、清理 chat 账号数据等。
// adminToken 为当次批次开始时通过 admin-api /account/login 获取的管理员 token。
func (c *cronServer) deleteExpiredUser(ctx context.Context, adminToken, userID string) {
	log.ZInfo(ctx, "deleteExpiredUser: start", "userID", userID)

	operationID := mcontext.GetOperationID(ctx)

	body, _ := json.Marshal(map[string]any{"userIDs": []string{userID}})
	url := c.chatAPIAddress + "/account/del"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.ZError(ctx, "deleteExpiredUser: build request failed", err, "userID", userID)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", adminToken)
	req.Header.Set("operationID", operationID)

	resp, err := chatHTTPClient.Do(req)
	if err != nil {
		log.ZError(ctx, "deleteExpiredUser: HTTP call failed", err, "userID", userID, "url", url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result map[string]any
		_ = json.NewDecoder(resp.Body).Decode(&result)
		log.ZError(ctx, "deleteExpiredUser: chat API returned error",
			fmt.Errorf("status %d", resp.StatusCode),
			"userID", userID, "response", result)
		return
	}

	// chat /account/del 已处理好友/群组/IM用户删除；仅清理 user_offline_record 防止重复触发
	if err := c.userOfflineRecordDB.Delete(ctx, userID); err != nil {
		log.ZWarn(ctx, "deleteExpiredUser: Delete offline record failed", err, "userID", userID)
	}

	log.ZInfo(ctx, "deleteExpiredUser: done", "userID", userID)
}

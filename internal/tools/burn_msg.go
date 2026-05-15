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

package tools

import (
	"fmt"
	"os"
	"time"

	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

// clearBurnExpiredMsgs 阅后即焚 cron 入口：循环调用 conversation 服务的
// ClearBurnExpiredMsgs，每次至多处理 burnLimit 个 (user, conversation) 分组，
// 直至本轮没有新的过期分组或达到防御性的最大循环次数。
func (c *cronServer) clearBurnExpiredMsgs() {
	now := time.Now()
	operationID := fmt.Sprintf("cron_burn_msg_%d_%d", os.Getpid(), now.UnixMilli())
	ctx := mcontext.SetOperationID(c.ctx, operationID)
	log.ZDebug(ctx, "clear burn expired msgs cron start")
	const (
		maxLoop   = 10000
		burnLimit = 100
	)
	var count int
	for i := 1; i <= maxLoop; i++ {
		resp, err := c.conversationClient.ClearBurnExpiredMsgs(ctx, &pbconversation.ClearBurnExpiredMsgsReq{
			Timestamp: now.UnixMilli(),
			Limit:     burnLimit,
		})
		if err != nil {
			log.ZError(ctx, "ClearBurnExpiredMsgs failed.", err)
			return
		}
		count += int(resp.Count)
		if resp.Count < burnLimit {
			break
		}
	}
	log.ZDebug(ctx, "clear burn expired msgs cron completed", "cost", time.Since(now), "count", count)
}

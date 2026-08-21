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

package webhook

import (
	"context"
	"testing"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/mq/memamq"
)

func TestAsyncPost_whenQueueFull_returnsWithoutWaiting(t *testing.T) {
	queue := memamq.NewMemoryQueue(1, 1)
	taskStarted := make(chan struct{})
	release := make(chan struct{})
	t.Cleanup(func() {
		close(release)
		queue.Stop()
	})

	if err := queue.Push(func() {
		close(taskStarted)
		<-release
	}); err != nil {
		t.Fatal(err)
	}
	<-taskStarted
	if err := queue.Push(func() {}); err != nil {
		t.Fatal(err)
	}

	client := NewWebhookClient("http://example.com", queue)
	started := time.Now()
	client.AsyncPost(
		context.Background(),
		"command",
		&callbackstruct.CommonCallbackReq{CallbackCommand: "command"},
		&callbackstruct.CommonCallbackResp{},
		&config.AfterConfig{Enable: true},
	)

	if elapsed := time.Since(started); elapsed >= time.Second {
		t.Fatalf("AsyncPost waited %s for queue capacity", elapsed)
	}
}

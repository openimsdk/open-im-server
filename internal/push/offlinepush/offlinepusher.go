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

package offlinepush

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/dummy"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/fcm"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/getui"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/jpush"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/options"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"strings"
)

const (
	geTUI    = "getui"
	firebase = "fcm"
	jPush    = "jpush"
)

// OfflinePusher Offline Pusher.
type OfflinePusher interface {
	Push(ctx context.Context, userIDs []string, title, content string, opts *options.Opts) error
}

func NewOfflinePusher(pushConf *config.Push, cache cache.ThirdCache, fcmConfigPath string) (OfflinePusher, error) {
	var offlinePusher OfflinePusher
	pushConf.Enable = strings.ToLower(pushConf.Enable)
	switch pushConf.Enable {
	case geTUI:
		offlinePusher = getui.NewClient(pushConf, cache)
	case firebase:
		return fcm.NewClient(pushConf, cache, fcmConfigPath)
	case jPush:
		offlinePusher = jpush.NewClient(pushConf)
	default:
		offlinePusher = dummy.NewClient()
	}
	return offlinePusher, nil
}

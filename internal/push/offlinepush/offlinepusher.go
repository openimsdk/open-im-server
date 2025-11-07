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

package push

import (
	"OpenIM/internal/push/fcm"
	"OpenIM/internal/push/getui"
	"OpenIM/internal/push/jpush"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/db/cache"
	"context"
)

type OfflinePusher interface {
	Push(ctx context.Context, userIDs []string, title, content string, opts *Opts) error
}

func NewOfflinePusher(cache cache.Model) OfflinePusher {
	var offlinePusher OfflinePusher
	if config.Config.Push.Getui.Enable {
		offlinePusher = getui.NewClient(cache)
	}
	if config.Config.Push.Fcm.Enable {
		offlinePusher = fcm.NewClient(cache)
	}
	if config.Config.Push.Jpns.Enable {
		offlinePusher = jpush.NewClient()
	}
	return offlinePusher
}

type Opts struct {
	Signal        *Signal
	IOSPushSound  string
	IOSBadgeCount bool
	Ex            string
}

type Signal struct {
	ClientMsgID string
}

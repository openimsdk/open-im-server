package dummy

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/options"
	"github.com/openimsdk/tools/log"
	"sync/atomic"
)

func NewClient() *Dummy {
	return &Dummy{}
}

type Dummy struct {
	v atomic.Bool
}

func (d *Dummy) Push(ctx context.Context, userIDs []string, title, content string, opts *options.Opts) error {
	if d.v.CompareAndSwap(false, true) {
		log.ZWarn(ctx, "dummy push", nil, "ps", "the offline push is not configured. to configure it, please go to config/openim-push.yml")
	}
	return nil
}

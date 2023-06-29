package offlinepush

import (
	"context"
)

type OfflinePusher interface {
	Push(ctx context.Context, userIDs []string, title, content string, opts *Opts) error
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

package dummy

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush"
)

func NewClient() *Dummy {
	return &Dummy{}
}

type Dummy struct {
}

func (d *Dummy) Push(ctx context.Context, userIDs []string, title, content string, opts *offlinepush.Opts) error {
	return nil
}

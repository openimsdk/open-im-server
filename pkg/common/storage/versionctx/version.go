package versionctx

import (
	"context"
	"sync"

	tablerelation "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type Collection struct {
	Name string
	Doc  *tablerelation.VersionLog
}

type versionKey struct{}

func WithVersionLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, versionKey{}, &VersionLog{})
}

func GetVersionLog(ctx context.Context) *VersionLog {
	if v, ok := ctx.Value(versionKey{}).(*VersionLog); ok {
		return v
	}
	return nil
}

type VersionLog struct {
	lock sync.Mutex
	data []Collection
}

func (v *VersionLog) Append(data ...Collection) {
	if v == nil || len(data) == 0 {
		return
	}
	v.lock.Lock()
	defer v.lock.Unlock()
	v.data = append(v.data, data...)
}

func (v *VersionLog) Get() []Collection {
	if v == nil {
		return nil
	}
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.data
}

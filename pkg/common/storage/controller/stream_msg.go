package controller

import "github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"

type StreamMsgDatabase interface {
	cache.StreamMsgCache
}

func NewStreamMsgDatabase(db cache.StreamMsgCache) StreamMsgDatabase {
	return &streamMsgDatabase{db}
}

type streamMsgDatabase struct {
	cache.StreamMsgCache
}

package mq

import "github.com/golang/protobuf/proto"

type Producer interface {
	SendMessage(m proto.Message, key ...string) (int32, int64, error)
}

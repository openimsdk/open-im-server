package msg

import (
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsNotFound(err error) bool {
	switch errs.Unwrap(err) {
	case redis.Nil, mongo.ErrNoDocuments:
		return true
	default:
		return false
	}
}

type activeConversations []*msg.ActiveConversation

func (s activeConversations) Len() int {
	return len(s)
}

func (s activeConversations) Less(i, j int) bool {
	return s[i].LastTime > s[j].LastTime
}

func (s activeConversations) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

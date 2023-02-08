package cache

import (
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"golang.org/x/tools/go/ssa/testdata/src/strconv"
	"time"
)

type DBFun func() (string, error)

type ConversationCache interface {
	GetUserConversationIDListFromCache(userID string, fn DBFun) ([]string, error)
	DelUserConversationIDListFromCache(userID string) error
	GetConversationFromCache(ownerUserID, conversationID string, fn DBFun) (*table.ConversationModel, error)
	GetConversationsFromCache(ownerUserID string, conversationIDList []string, fn DBFun) ([]*table.ConversationModel, error)
	GetUserAllConversationList(ownerUserID string, fn DBFun) ([]*table.ConversationModel, error)
	DelConversationFromCache(ownerUserID, conversationID string) error
}
type ConversationRedis struct {
	rcClient *rockscache.Client
}

func NewConversationRedis(rcClient *rockscache.Client) *ConversationRedis {
	return &ConversationRedis{rcClient: rcClient}
}

func NewConversationCache(rdb redis.UniversalClient, conversationDB *relation.ConversationGorm, options rockscache.Options) *ConversationCache {
	return &ConversationCache{conversationDB: conversationDB, expireTime: conversationExpireTime, rcClient: rockscache.NewClient(rdb, options)}
}

func (c *ConversationCache) getConversationKey(ownerUserID, conversationID string) string {
	return conversationKey + ownerUserID + ":" + conversationID
}

func (c *ConversationCache) getConversationIDsKey(ownerUserID string) string {
	return conversationIDsKey + ownerUserID
}

func (c *ConversationCache) getRecvMsgOptKey(ownerUserID, conversationID string) string {
	return recvMsgOptKey + ownerUserID + ":" + conversationID
}

func (c *ConversationCache) getSuperGroupRecvNotNotifyUserIDsKey(groupID string) string {
	return superGroupRecvMsgNotNotifyUserIDsKey + groupID
}

func (c *ConversationCache) GetUserConversationIDs(ctx context.Context, ownerUserID string, f func(userID string) ([]string, error)) (conversationIDs []string, err error) {
	//getConversationIDs := func() (string, error) {
	//	conversationIDs, err := relation.GetConversationIDsByUserID(ownerUserID)
	//	if err != nil {
	//		return "", err
	//	}
	//	bytes, err := json.Marshal(conversationIDs)
	//	if err != nil {
	//		return "", utils.Wrap(err, "")
	//	}
	//	return string(bytes), nil
	//}
	//defer func() {
	//	tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationIDs", conversationIDs)
	//}()
	//conversationIDsStr, err := c.rcClient.Fetch(c.getConversationIDsKey(ownerUserID), time.Second*30*60, getConversationIDs)
	//err = json.Unmarshal([]byte(conversationIDsStr), &conversationIDs)
	//if err != nil {
	//	return nil, utils.Wrap(err, "")
	//}
	//return conversationIDs, nil
	return GetCache(c.rcClient, c.getConversationIDsKey(ownerUserID), time.Second*30*60, func() ([]string, error) {
		return f(ownerUserID)
	})
}

func (c *ConversationCache) GetUserConversationIDs1(ctx context.Context, ownerUserID string) (conversationIDs []string, err error) {
	//getConversationIDs := func() (string, error) {
	//	conversationIDs, err := relation.GetConversationIDsByUserID(ownerUserID)
	//	if err != nil {
	//		return "", err
	//	}
	//	bytes, err := json.Marshal(conversationIDs)
	//	if err != nil {
	//		return "", utils.Wrap(err, "")
	//	}
	//	return string(bytes), nil
	//}
	//defer func() {
	//	tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationIDs", conversationIDs)
	//}()
	//conversationIDsStr, err := c.rcClient.Fetch(c.getConversationIDsKey(ownerUserID), time.Second*30*60, getConversationIDs)
	//err = json.Unmarshal([]byte(conversationIDsStr), &conversationIDs)
	//if err != nil {
	//	return nil, utils.Wrap(err, "")
	//}
	//return conversationIDs, nil
	return GetCache1[[]string](c.rcClient, c.getConversationIDsKey(ownerUserID), time.Second*30*60, fn)
}

func GetCache1[T any](rcClient *rockscache.Client, key string, expire time.Duration, fn func() (any, error)) (T, error) {
	v, err := rcClient.Fetch(key, expire, func() (string, error) {
		v, err := fn()
		if err != nil {
			return "", err
		}
		bs, err := json.Marshal(v)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bs), nil
	})
	var t T
	if err != nil {
		return t, err
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		return t, utils.Wrap(err, "")
	}
	return t, nil
}

func GetCache[T any](rcClient *rockscache.Client, key string, expire time.Duration, fn func() (T, error)) (T, error) {
	v, err := rcClient.Fetch(key, expire, func() (string, error) {
		v, err := fn()
		if err != nil {
			return "", err
		}
		bs, err := json.Marshal(v)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bs), nil
	})
	var t T
	if err != nil {
		return t, err
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		return t, utils.Wrap(err, "")
	}
	return t, nil
}

func (c *ConversationCache) DelUserConversationIDs(ctx context.Context, ownerUserID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID)
	}()
	return utils.Wrap(c.rcClient.TagAsDeleted(c.getConversationIDsKey(ownerUserID)), "DelUserConversationIDs err")
}

func (c *ConversationCache) GetConversation(ctx context.Context, ownerUserID, conversationID string) (conversation *relationTb.ConversationModel, err error) {
	getConversation := func() (string, error) {
		conversation, err := relation.GetConversation(ownerUserID, conversationID)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(conversation)
		if err != nil {
			return "", utils.Wrap(err, "conversation Marshal failed")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationID", conversationID, "conversation", *conversation)
	}()
	conversationStr, err := c.rcClient.Fetch(c.getConversationKey(ownerUserID, conversationID), c.expireTime, getConversation)
	if err != nil {
		return nil, err
	}
	conversation = &relationTb.ConversationModel{}
	err = json.Unmarshal([]byte(conversationStr), &conversation)
	return conversation, utils.Wrap(err, "Unmarshal failed")
}

func (c *ConversationCache) DelConversation(ctx context.Context, ownerUserID, conversationID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationID", conversationID)
	}()
	return utils.Wrap(c.rcClient.TagAsDeleted(c.getConversationKey(ownerUserID, conversationID)), "DelConversation err")
}

func (c *ConversationCache) GetConversations(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []relationTb.ConversationModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationIDs", conversationIDs, "conversations", conversations)
	}()
	for _, conversationID := range conversationIDs {
		conversation, err := c.GetConversation(ctx, ownerUserID, conversationID)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, *conversation)
	}
	return conversations, nil
}

func (c *ConversationCache) GetUserAllConversations(ctx context.Context, ownerUserID string) (conversations []relationTb.ConversationModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversations", conversations)
	}()
	IDs, err := c.GetUserConversationIDs(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}
	var conversationIDs []relationTb.ConversationModel
	for _, conversationID := range IDs {
		conversation, err := c.GetConversation(ctx, ownerUserID, conversationID)
		if err != nil {
			return nil, err
		}
		conversationIDs = append(conversationIDs, *conversation)
	}
	return conversationIDs, nil
}

func (c *ConversationCache) GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error) {
	getConversation := func() (string, error) {
		conversation, err := relation.GetConversation(ownerUserID, conversationID)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(conversation.RecvMsgOpt)), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationID", conversationID, "opt", opt)
	}()
	optStr, err := c.rcClient.Fetch(c.getConversationKey(ownerUserID, conversationID), c.expireTime, getConversation)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(optStr)
}

func (c *ConversationCache) DelUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) error {
	return utils.Wrap(c.rcClient.TagAsDeleted(c.getConversationKey(ownerUserID, conversationID)), "DelUserRecvMsgOpt failed")
}

func (c *ConversationCache) GetSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) (userIDs []string, err error) {
	return nil, nil
}

func (c *ConversationCache) DelSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) (err error) {
	return nil
}

func (c *ConversationCache) GetSuperGroupRecvMsgNotNotifyUserIDsHash(ctx context.Context, groupID string) (hash uint32, err error) {
	return
}

func (c *ConversationCache) DelSuperGroupRecvMsgNotNotifyUserIDsHash(ctx context.Context, groupID string) {
	return
}

package cache

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	conversationKey                      = "CONVERSATION:"
	conversationIDsKey                   = "CONVERSATION_IDS:"
	recvMsgOptKey                        = "RECV_MSG_OPT:"
	superGroupRecvMsgNotNotifyUserIDsKey = "SUPER_GROUP_RECV_MSG_NOT_NOTIFY_USER_IDS:"
	conversationExpireTime               = time.Second * 60 * 60 * 12
)

type FuncDB func() (string, error)

// arg fn will exec when no data in cache
type ConversationCache interface {
	// get user's conversationIDs from cache
	GetUserConversationIDs(ctx context.Context, userID string, fn FuncDB) ([]string, error)
	// del user's conversationIDs from cache, call when a user add or reduce a conversation
	DelUserConversationIDs(ctx context.Context, userID string) error
	DelUsersConversationIDs(ctx context.Context, userIDList []string) error
	// get one conversation from cache
	GetConversation(ctx context.Context, ownerUserID, conversationID string, fn FuncDB) (*relationTb.ConversationModel, error)
	// get one conversation from cache
	GetConversations(ctx context.Context, ownerUserID string, conversationIDs []string, fn FuncDB) ([]*relationTb.ConversationModel, error)
	// get one user's all conversations from cache
	GetUserAllConversations(ctx context.Context, ownerUserID string, fn FuncDB) ([]*relationTb.ConversationModel, error)
	// del one conversation from cache, call when one user's conversation Info changed
	DelConversation(ctx context.Context, ownerUserID, conversationID string) error
	DelUserConversations(ctx context.Context, ownerUserID string, conversationIDList []string) error
	DelUsersConversation(ctx context.Context, ownerUserIDList []string, conversationID string) error
	// get user conversation recv msg from cache
	GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string, fn func(ctx context.Context, ownerUserID, conversationID string) (opt int, err error)) (opt int, err error)
	// del user recv msg opt from cache, call when user's conversation recv msg opt changed
	DelUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) error
	// get one super group recv msg but do not notification userID list
	GetSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string, fn func(ctx context.Context, groupID string) (userIDs []string, err error)) (userIDs []string, err error)
	// del one super group recv msg but do not notification userID list, call it when this list changed
	DelSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) error
	//GetSuperGroupRecvMsgNotNotifyUserIDsHash(ctx context.Context, groupID string) (hash uint32, err error)
	//DelSuperGroupRecvMsgNotNotifyUserIDsHash(ctx context.Context, groupID string)
}

func NewConversationRedis(rdb redis.UniversalClient, opts rockscache.Options) ConversationCache {
	return &ConversationRedis{rcClient: rockscache.NewClient(rdb, opts)}
}

type ConversationRedis struct {
	rcClient   *rockscache.Client
	expireTime time.Duration
}

func (c *ConversationRedis) GetUserConversationIDs(ctx context.Context, userID string, fn FuncDB) ([]string, error) {
	return nil, nil
}

func (c *ConversationRedis) GetConversation(ctx context.Context, ownerUserID, conversationID string, fn FuncDB) (*relationTb.ConversationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConversationRedis) GetConversations(ctx context.Context, ownerUserID string, conversationIDs []string, fn FuncDB) ([]*relationTb.ConversationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConversationRedis) GetUserAllConversations(ctx context.Context, ownerUserID string, fn FuncDB) ([]*relationTb.ConversationModel, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConversationRedis) DelUserConversations(ctx context.Context, ownerUserID string, conversationIDList []string) error {
	//TODO implement me
	panic("implement me")
}

func (c *ConversationRedis) GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string, fn func(ctx context.Context, ownerUserID string, conversationID string) (opt int, err error)) (opt int, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConversationRedis) GetSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string, fn func(ctx context.Context, groupID string) (userIDs []string, err error)) (userIDs []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *ConversationRedis) DelUsersConversationIDs(ctx context.Context, userIDList []string) error {
	panic("implement me")
}

func (c *ConversationRedis) DelUsersConversation(ctx context.Context, ownerUserIDList []string, conversationID string) error {
	panic("implement me")
}

func NewNewConversationRedis(rdb redis.UniversalClient, conversationDB *relation.ConversationGorm, options rockscache.Options) *ConversationRedis {
	return &ConversationRedis{rcClient: rockscache.NewClient(rdb, options)}
}

func (c *ConversationRedis) getConversationKey(ownerUserID, conversationID string) string {
	return conversationKey + ownerUserID + ":" + conversationID
}

func (c *ConversationRedis) getConversationIDsKey(ownerUserID string) string {
	return conversationIDsKey + ownerUserID
}

func (c *ConversationRedis) getRecvMsgOptKey(ownerUserID, conversationID string) string {
	return recvMsgOptKey + ownerUserID + ":" + conversationID
}

func (c *ConversationRedis) getSuperGroupRecvNotNotifyUserIDsKey(groupID string) string {
	return superGroupRecvMsgNotNotifyUserIDsKey + groupID
}

//func (c *ConversationRedis) GetUserConversationIDs(ctx context.Context, ownerUserID string) (conversationIDs []string, err error) {
//	//getConversationIDs := func() (string, error) {
//	//	conversationIDs, err := relation.GetConversationIDsByUserID(ownerUserID)
//	//	if err != nil {
//	//		return "", err
//	//	}
//	//	bytes, err := json.Marshal(conversationIDs)
//	//	if err != nil {
//	//		return "", utils.Wrap(err, "")
//	//	}
//	//	return string(bytes), nil
//	//}
//	//defer func() {
//	//	mcontext.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationIDs", conversationIDs)
//	//}()
//	//conversationIDsStr, err := c.rcClient.Fetch(c.getConversationIDsKey(ownerUserID), time.Second*30*60, getConversationIDs)
//	//err = json.Unmarshal([]byte(conversationIDsStr), &conversationIDs)
//	//if err != nil {
//	//	return nil, utils.Wrap(err, "")
//	//}
//	//return conversationIDs, nil
//	return GetCache(ctx, c.rcClient, c.getConversationIDsKey(ownerUserID), conversationExpireTime, func(ctx context.Context) ([]string, error) {
//		panic("implement me")
//	})
//}

func (c *ConversationRedis) GetUserConversationIDs1(ctx context.Context, ownerUserID string) (conversationIDs []string, err error) {
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
	//	mcontext.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationIDs", conversationIDs)
	//}()
	//conversationIDsStr, err := c.rcClient.Fetch(c.getConversationIDsKey(ownerUserID), time.Second*30*60, getConversationIDs)
	//err = json.Unmarshal([]byte(conversationIDsStr), &conversationIDs)
	//if err != nil {
	//	return nil, utils.Wrap(err, "")
	//}
	//return conversationIDs, nil
	//return GetCache1[[]string](c.rcClient, c.getConversationIDsKey(ownerUserID), conversationExpireTime, fn)

	return GetCache(ctx, c.rcClient, c.getConversationIDsKey(ownerUserID), conversationExpireTime, func(ctx context.Context) ([]string, error) {
		panic("")
	})
}

//func GetCache1[T any](rcClient *rockscache.Client, key string, expire time.Duration, fn func() (any, error)) (T, error) {
//	v, err := rcClient.Fetch(key, expire, func() (string, error) {
//		v, err := fn()
//		if err != nil {
//			return "", err
//		}
//		bs, err := json.Marshal(v)
//		if err != nil {
//			return "", utils.Wrap(err, "")
//		}
//		return string(bs), nil
//	})
//	var t T
//	if err != nil {
//		return t, err
//	}
//	err = json.Unmarshal([]byte(v), &t)
//	if err != nil {
//		return t, utils.Wrap(err, "")
//	}
//	return t, nil
//}

func (c *ConversationRedis) DelUserConversationIDs(ctx context.Context, ownerUserID string) (err error) {
	return utils.Wrap(c.rcClient.TagAsDeleted(c.getConversationIDsKey(ownerUserID)), "DelUserConversationIDs err")
}

//func (c *ConversationRedis) GetConversation(ctx context.Context, ownerUserID, conversationID string) (conversation *relationTb.ConversationModel, err error) {
//	return GetCache(ctx, c.rcClient, c.getConversationKey(ownerUserID, conversationID), c.expireTime, func(ctx context.Context) (*relationTb.ConversationModel, error) {
//		panic("implement me")
//	})
//}

func (c *ConversationRedis) DelConversation(ctx context.Context, ownerUserID, conversationID string) (err error) {
	return utils.Wrap(c.rcClient.TagAsDeleted(c.getConversationKey(ownerUserID, conversationID)), "DelConversation err")
}

//func (c *ConversationRedis) GetConversations(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []relationTb.ConversationModel, err error) {
//	defer func() {
//		mcontext.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationIDs", conversationIDs, "conversations", conversations)
//	}()
//	for _, conversationID := range conversationIDs {
//		conversation, err := c.GetConversation(ctx, ownerUserID, conversationID)
//		if err != nil {
//			return nil, err
//		}
//		conversations = append(conversations, *conversation)
//	}
//	return conversations, nil
//}

//func (c *ConversationRedis) GetUserAllConversations(ctx context.Context, ownerUserID string) (conversations []relationTb.ConversationModel, err error) {
//	defer func() {
//		mcontext.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversations", conversations)
//	}()
//	IDs, err := c.GetUserConversationIDs(ctx, ownerUserID)
//	if err != nil {
//		return nil, err
//	}
//	var conversationIDs []relationTb.ConversationModel
//	for _, conversationID := range IDs {
//		conversation, err := c.GetConversation(ctx, ownerUserID, conversationID)
//		if err != nil {
//			return nil, err
//		}
//		conversationIDs = append(conversationIDs, *conversation)
//	}
//	return conversationIDs, nil
//}

//func (c *ConversationRedis) GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error) {
//	//getConversation := func() (string, error) {
//	//	conversation, err := relation.GetConversation(ownerUserID, conversationID)
//	//	if err != nil {
//	//		return "", err
//	//	}
//	//	return strconv.Itoa(int(conversation.RecvMsgOpt)), nil
//	//}
//	//defer func() {
//	//	mcontext.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "conversationID", conversationID, "opt", opt)
//	//}()
//	//optStr, err := c.rcClient.Fetch(c.getConversationKey(ownerUserID, conversationID), c.expireTime, getConversation)
//	//if err != nil {
//	//	return 0, err
//	//}
//	//return strconv.Atoi(optStr)
//	// panic("implement me")
//	return GetCache(ctx, c.rcClient, c.getConversationKey(ownerUserID, conversationID), c.expireTime, func(ctx context.Context) (int, error) {
//		panic("implement me")
//	})
//}

func (c *ConversationRedis) DelUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) error {
	return utils.Wrap(c.rcClient.TagAsDeleted(c.getConversationKey(ownerUserID, conversationID)), "DelUserRecvMsgOpt failed")
}

func (c *ConversationRedis) DelSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) (err error) {
	panic("implement me")
}

func (c *ConversationRedis) GetSuperGroupRecvMsgNotNotifyUserIDsHash(ctx context.Context, groupID string) (hash uint32, err error) {
	panic("implement me")
}

func (c *ConversationRedis) DelSuperGroupRecvMsgNotNotifyUserIDsHash(ctx context.Context, groupID string) {
	panic("implement me")
}

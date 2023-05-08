package controller

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
)

// for mongoDB
type ExtendMsgDatabase interface {
	CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error
	GetAllExtendMsgSet(ctx context.Context, ID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error)
	GetExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error)
	InsertExtendMsg(ctx context.Context, conversationID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error
	InsertOrUpdateReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*unRelationTb.KeyValueModel) error
	DeleteReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*unRelationTb.KeyValueModel) error
	GetExtendMsg(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error)
}

type extendMsgDatabase struct {
	database unRelationTb.ExtendMsgSetModelInterface
	cache    cache.ExtendMsgSetCache
	ctxTx    tx.CtxTx
}

func NewExtendMsgDatabase(extendMsgModel unRelationTb.ExtendMsgSetModelInterface, cache cache.ExtendMsgSetCache, ctxTx tx.CtxTx) ExtendMsgDatabase {
	return &extendMsgDatabase{database: extendMsgModel, cache: cache, ctxTx: ctxTx}
}

func (e *extendMsgDatabase) CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error {
	return e.database.CreateExtendMsgSet(ctx, set)
}

func (e *extendMsgDatabase) GetAllExtendMsgSet(ctx context.Context, conversationID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error) {
	return e.database.GetAllExtendMsgSet(ctx, conversationID, opts)
}

func (e *extendMsgDatabase) GetExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error) {
	return e.database.GetExtendMsgSet(ctx, conversationID, sessionType, maxMsgUpdateTime)
}

func (e *extendMsgDatabase) InsertExtendMsg(ctx context.Context, conversationID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error {
	return e.database.InsertExtendMsg(ctx, conversationID, sessionType, msg)
}

func (e *extendMsgDatabase) InsertOrUpdateReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*unRelationTb.KeyValueModel) error {
	return e.database.InsertOrUpdateReactionExtendMsgSet(ctx, conversationID, sessionType, clientMsgID, msgFirstModifyTime, reactionExtensionList)
}

func (e *extendMsgDatabase) DeleteReactionExtendMsgSet(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*unRelationTb.KeyValueModel) error {
	return e.database.DeleteReactionExtendMsgSet(ctx, conversationID, sessionType, clientMsgID, msgFirstModifyTime, reactionExtensionList)
}

func (e *extendMsgDatabase) GetExtendMsg(ctx context.Context, conversationID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error) {
	return e.cache.GetExtendMsg(ctx, conversationID, sessionType, clientMsgID, maxMsgUpdateTime)
}

package controller

import (
	unRelationTb "OpenIM/pkg/common/db/table/unrelation"
	"OpenIM/pkg/proto/sdkws"
	"context"
)

type ExtendMsgDatabase interface {
	CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error
	GetAllExtendMsgSet(ctx context.Context, ID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error)
	GetExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error)
	InsertExtendMsg(ctx context.Context, sourceID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error
	InsertOrUpdateReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error
	DeleteReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error
	GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error)
}

type extendMsgDatabase struct {
	db unRelationTb.ExtendMsgSetModelInterface
}

func NewExtendMsgDatabase() ExtendMsgDatabase {
	return &extendMsgDatabase{}
}

func (e *extendMsgDatabase) CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error {
	return e.db.CreateExtendMsgSet(ctx, set)
}

func (e *extendMsgDatabase) GetAllExtendMsgSet(ctx context.Context, sourceID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error) {
	return e.db.GetAllExtendMsgSet(ctx, sourceID, opts)
}

func (e *extendMsgDatabase) GetExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error) {
	return e.db.GetExtendMsgSet(ctx, sourceID, sessionType, maxMsgUpdateTime)
}

func (e *extendMsgDatabase) InsertExtendMsg(ctx context.Context, sourceID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error {
	return e.db.InsertExtendMsg(ctx, sourceID, sessionType, msg)
}

func (e *extendMsgDatabase) InsertOrUpdateReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error {
	return e.db.InsertOrUpdateReactionExtendMsgSet(ctx, sourceID, sessionType, clientMsgID, msgFirstModifyTime, reactionExtensionList)
}
func (e *extendMsgDatabase) DeleteReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error {
	return e.db.DeleteReactionExtendMsgSet(ctx, sourceID, sessionType, clientMsgID, msgFirstModifyTime, reactionExtensionList)
}

func (e *extendMsgDatabase) GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error) {
	return e.db.TakeExtendMsg(ctx, sourceID, sessionType, clientMsgID, maxMsgUpdateTime)
}

package controller

import (
	unRelationTb "OpenIM/pkg/common/db/table/unrelation"
	"OpenIM/pkg/proto/sdkws"
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExtendMsgInterface interface {
	CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error
	GetAllExtendMsgSet(ctx context.Context, ID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error)
	GetExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error)
	InsertExtendMsg(ctx context.Context, sourceID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error
	InsertOrUpdateReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error
	DeleteReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error
	GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error)
}

type ExtendMsgController struct {
	database ExtendMsgDatabase
}

func NewExtendMsgController(mgo *mongo.Client, rdb redis.UniversalClient) *ExtendMsgController {
	return &ExtendMsgController{}
}

func (e *ExtendMsgController) CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error {
	return e.database.CreateExtendMsgSet(ctx, set)
}

func (e *ExtendMsgController) GetAllExtendMsgSet(ctx context.Context, ID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error) {
	return e.GetAllExtendMsgSet(ctx, ID, opts)
}

func (e *ExtendMsgController) GetExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error) {
	return e.GetExtendMsgSet(ctx, sourceID, sessionType, maxMsgUpdateTime)
}

func (e *ExtendMsgController) InsertExtendMsg(ctx context.Context, sourceID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error {
	return e.InsertExtendMsg(ctx, sourceID, sessionType, msg)
}

func (e *ExtendMsgController) InsertOrUpdateReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error {
	return e.InsertOrUpdateReactionExtendMsgSet(ctx, sourceID, sessionType, clientMsgID, msgFirstModifyTime, reactionExtensionList)
}
func (e *ExtendMsgController) DeleteReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error {
	return e.DeleteReactionExtendMsgSet(ctx, sourceID, sessionType, clientMsgID, msgFirstModifyTime, reactionExtensionList)
}

func (e *ExtendMsgController) GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error) {
	return e.GetExtendMsg(ctx, sourceID, sessionType, clientMsgID, maxMsgUpdateTime)
}

type ExtendMsgDatabaseInterface interface {
	CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error
	GetAllExtendMsgSet(ctx context.Context, ID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error)
	GetExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error)
	InsertExtendMsg(ctx context.Context, sourceID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error
	InsertOrUpdateReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error
	DeleteReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error
	GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error)
}

type ExtendMsgDatabase struct {
	model unRelationTb.ExtendMsgSetModelInterface
}

func NewExtendMsgDatabase() ExtendMsgDatabaseInterface {
	return &ExtendMsgDatabase{}
}

func (e *ExtendMsgDatabase) CreateExtendMsgSet(ctx context.Context, set *unRelationTb.ExtendMsgSetModel) error {
	return e.model.CreateExtendMsgSet(ctx, set)
}

func (e *ExtendMsgDatabase) GetAllExtendMsgSet(ctx context.Context, sourceID string, opts *unRelationTb.GetAllExtendMsgSetOpts) (sets []*unRelationTb.ExtendMsgSetModel, err error) {
	return e.model.GetAllExtendMsgSet(ctx, sourceID, opts)
}

func (e *ExtendMsgDatabase) GetExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, maxMsgUpdateTime int64) (*unRelationTb.ExtendMsgSetModel, error) {
	return e.model.GetExtendMsgSet(ctx, sourceID, sessionType, maxMsgUpdateTime)
}

func (e *ExtendMsgDatabase) InsertExtendMsg(ctx context.Context, sourceID string, sessionType int32, msg *unRelationTb.ExtendMsgModel) error {
	return e.model.InsertExtendMsg(ctx, sourceID, sessionType, msg)
}

func (e *ExtendMsgDatabase) InsertOrUpdateReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error {
	return e.model.InsertOrUpdateReactionExtendMsgSet(ctx, sourceID, sessionType, clientMsgID, msgFirstModifyTime, reactionExtensionList)
}
func (e *ExtendMsgDatabase) DeleteReactionExtendMsgSet(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, msgFirstModifyTime int64, reactionExtensionList map[string]*sdkws.KeyValue) error {
	return e.model.DeleteReactionExtendMsgSet(ctx, sourceID, sessionType, clientMsgID, msgFirstModifyTime, reactionExtensionList)
}

func (e *ExtendMsgDatabase) GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, maxMsgUpdateTime int64) (extendMsg *unRelationTb.ExtendMsgModel, err error) {
	return e.model.TakeExtendMsg(ctx, sourceID, sessionType, clientMsgID, maxMsgUpdateTime)
}

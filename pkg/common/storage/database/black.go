package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

type Black interface {
	Create(ctx context.Context, blacks []*model.Black) (err error)
	Delete(ctx context.Context, blacks []*model.Black) (err error)
	Find(ctx context.Context, blacks []*model.Black) (blackList []*model.Black, err error)
	Take(ctx context.Context, ownerUserID, blockUserID string) (black *model.Black, err error)
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*model.Black, err error)
	FindOwnerBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*model.Black, err error)
	FindBlackUserIDs(ctx context.Context, ownerUserID string) (blackUserIDs []string, err error)
}

var (
	_ Black = (*mgoImpl)(nil)
	_ Black = (*redisImpl)(nil)
)

type mgoImpl struct {
}

func (m *mgoImpl) Create(ctx context.Context, blacks []*model.Black) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mgoImpl) Delete(ctx context.Context, blacks []*model.Black) (err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mgoImpl) Find(ctx context.Context, blacks []*model.Black) (blackList []*model.Black, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mgoImpl) Take(ctx context.Context, ownerUserID, blockUserID string) (black *model.Black, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mgoImpl) FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*model.Black, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mgoImpl) FindOwnerBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*model.Black, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *mgoImpl) FindBlackUserIDs(ctx context.Context, ownerUserID string) (blackUserIDs []string, err error) {
	//TODO implement me
	panic("implement me")
}

type redisImpl struct {
}

func (r *redisImpl) Create(ctx context.Context, blacks []*model.Black) (err error) {

	//TODO implement me
	panic("implement me")
}

func (r *redisImpl) Delete(ctx context.Context, blacks []*model.Black) (err error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisImpl) Find(ctx context.Context, blacks []*model.Black) (blackList []*model.Black, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisImpl) Take(ctx context.Context, ownerUserID, blockUserID string) (black *model.Black, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisImpl) FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*model.Black, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisImpl) FindOwnerBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*model.Black, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisImpl) FindBlackUserIDs(ctx context.Context, ownerUserID string) (blackUserIDs []string, err error) {
	//TODO implement me
	panic("implement me")
}

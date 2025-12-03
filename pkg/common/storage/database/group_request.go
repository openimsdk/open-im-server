package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

type GroupRequest interface {
	Create(ctx context.Context, groupRequests []*model.GroupRequest) (err error)
	Delete(ctx context.Context, groupID string, userID string) (err error)
	UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error)
	Take(ctx context.Context, groupID string, userID string) (groupRequest *model.GroupRequest, err error)
	FindGroupRequests(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupRequest, error)
	Page(ctx context.Context, userID string, groupIDs []string, handleResults []int, pagination pagination.Pagination) (total int64, groups []*model.GroupRequest, err error)
	PageGroup(ctx context.Context, groupIDs []string, handleResults []int, pagination pagination.Pagination) (total int64, groups []*model.GroupRequest, err error)
	GetUnhandledCount(ctx context.Context, groupIDs []string, ts int64) (int64, error)
}

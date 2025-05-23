// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

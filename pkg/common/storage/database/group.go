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
	"time"

	"github.com/openimsdk/tools/db/pagination"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type Group interface {
	Create(ctx context.Context, groups []*model.Group) (err error)
	UpdateMap(ctx context.Context, groupID string, args map[string]any) (err error)
	UpdateStatus(ctx context.Context, groupID string, status int32) (err error)
	Find(ctx context.Context, groupIDs []string) (groups []*model.Group, err error)
	Take(ctx context.Context, groupID string) (group *model.Group, err error)
	Search(ctx context.Context, keyword string, pagination pagination.Pagination) (total int64, groups []*model.Group, err error)
	// Get Group total quantity
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// Get Group total quantity every day
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)

	FindJoinSortGroupID(ctx context.Context, groupIDs []string) ([]string, error)

	SearchJoin(ctx context.Context, groupIDs []string, keyword string, pagination pagination.Pagination) (int64, []*model.Group, error)
}

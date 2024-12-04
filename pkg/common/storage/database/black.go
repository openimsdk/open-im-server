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

	"github.com/openimsdk/tools/db/pagination"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
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

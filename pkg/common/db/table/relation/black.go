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

package relation

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/pagination"
)

type BlackModel struct {
	OwnerUserID    string    `bson:"owner_user_id"`
	BlockUserID    string    `bson:"block_user_id"`
	CreateTime     time.Time `bson:"create_time"`
	AddSource      int32     `bson:"add_source"`
	OperatorUserID string    `bson:"operator_user_id"`
	Ex             string    `bson:"ex"`
}

type BlackModelInterface interface {
	Create(ctx context.Context, blacks []*BlackModel) (err error)
	Delete(ctx context.Context, blacks []*BlackModel) (err error)
	// UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]any) (err error)
	// Update(ctx context.Context, blacks []*BlackModel) (err error)
	Find(ctx context.Context, blacks []*BlackModel) (blackList []*BlackModel, err error)
	Take(ctx context.Context, ownerUserID, blockUserID string) (black *BlackModel, err error)
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pagination pagination.Pagination) (total int64, blacks []*BlackModel, err error)
	FindOwnerBlackInfos(ctx context.Context, ownerUserID string, userIDs []string) (blacks []*BlackModel, err error)
	FindBlackUserIDs(ctx context.Context, ownerUserID string) (blackUserIDs []string, err error)
}

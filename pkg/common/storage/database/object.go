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

type ObjectInfo interface {
	SetObject(ctx context.Context, obj *model.Object) error
	Take(ctx context.Context, engine string, name string) (*model.Object, error)
	Delete(ctx context.Context, engine string, name string) error
	FindNeedDeleteObjectByDB(ctx context.Context, duration time.Time, needDelType []string, pagination pagination.Pagination) (total int64, objects []*model.Object, err error)
	FindModelsByKey(ctx context.Context, key string) (objects []*model.Object, err error)
}

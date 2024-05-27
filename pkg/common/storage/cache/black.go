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

package cache

import (
	"context"
)

type BlackCache interface {
	BatchDeleter
	CloneBlackCache() BlackCache
	GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error)
	// del user's blackIDs msgCache, exec when a user's black list changed
	DelBlackIDs(ctx context.Context, userID string) BlackCache
}

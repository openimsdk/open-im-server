// Copyright © 2024 OpenIM. All rights reserved.
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

// CallRecordDatabase defines storage operations for the call record table.
type CallRecordDatabase interface {
	// CreateCallRecord writes a new call record entry.
	CreateCallRecord(ctx context.Context, record *model.CallRecord) error
	// SearchCallRecords returns paginated call records involving userID,
	// optionally filtered by status, time range and a keyword (fuzzy match on InviterUserNickname).
	SearchCallRecords(ctx context.Context, userID string, status int32, startTime, endTime int64, keyword string, pg pagination.Pagination) (int64, []*model.CallRecord, error)
	// DeleteCallRecords removes call records by their SIDs.
	DeleteCallRecords(ctx context.Context, sids []string) error
}

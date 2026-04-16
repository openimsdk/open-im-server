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
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

type SpamReport interface {
	// Create inserts a new spam report record.
	Create(ctx context.Context, report *model.SpamReport) error
	// Find queries spam reports with optional filters, returns total count and records.
	Find(ctx context.Context, status int32, reportedUserID, reporterUserID string,
		start, end time.Time, pagination pagination.Pagination) (int64, []*model.SpamReport, error)
	// UpdateStatus updates the handling status of a spam report.
	UpdateStatus(ctx context.Context, reportID string, status int32, handlerUserID string, handleTime time.Time) error
	// Get retrieves a single spam report by its reportID.
	Get(ctx context.Context, reportID string) (*model.SpamReport, error)
}

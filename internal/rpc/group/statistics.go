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

package group

import (
	"context"
	"time"

	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/tools/errs"
)

func (s *groupServer) GroupCreateCount(ctx context.Context, req *group.GroupCreateCountReq) (*group.GroupCreateCountResp, error) {
	if req.Start > req.End {
		return nil, errs.ErrArgs.Wrap("start > end")
	}
	total, err := s.db.CountTotal(ctx, nil)
	if err != nil {
		return nil, err
	}
	start := time.UnixMilli(req.Start)
	before, err := s.db.CountTotal(ctx, &start)
	if err != nil {
		return nil, err
	}
	count, err := s.db.CountRangeEverydayTotal(ctx, start, time.UnixMilli(req.End))
	if err != nil {
		return nil, err
	}
	return &group.GroupCreateCountResp{Total: total, Before: before, Count: count}, nil
}

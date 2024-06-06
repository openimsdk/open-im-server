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

package friend

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/mcontext"
)

func (s *friendServer) GetPaginationBlacks(ctx context.Context, req *relation.GetPaginationBlacksReq) (resp *relation.GetPaginationBlacksResp, err error) {
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	total, blacks, err := s.blackDatabase.FindOwnerBlacks(ctx, req.UserID, req.Pagination)
	if err != nil {
		return nil, err
	}
	resp = &relation.GetPaginationBlacksResp{}
	resp.Blacks, err = convert.BlackDB2Pb(ctx, blacks, s.userRpcClient.GetUsersInfoMap)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

func (s *friendServer) IsBlack(ctx context.Context, req *relation.IsBlackReq) (*relation.IsBlackResp, error) {
	in1, in2, err := s.blackDatabase.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	resp := &relation.IsBlackResp{}
	resp.InUser1Blacks = in1
	resp.InUser2Blacks = in2
	return resp, nil
}

func (s *friendServer) RemoveBlack(ctx context.Context, req *relation.RemoveBlackReq) (*relation.RemoveBlackResp, error) {
	if err := s.userRpcClient.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	if err := s.blackDatabase.Delete(ctx, []*model.Black{{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID}}); err != nil {
		return nil, err
	}

	s.notificationSender.BlackDeletedNotification(ctx, req)

	return &relation.RemoveBlackResp{}, nil
}

func (s *friendServer) AddBlack(ctx context.Context, req *relation.AddBlackReq) (*relation.AddBlackResp, error) {
	if err := authverify.CheckAccessV3(ctx, req.OwnerUserID, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	_, err := s.userRpcClient.GetUsersInfo(ctx, []string{req.OwnerUserID, req.BlackUserID})
	if err != nil {
		return nil, err
	}
	black := model.Black{
		OwnerUserID:    req.OwnerUserID,
		BlockUserID:    req.BlackUserID,
		OperatorUserID: mcontext.GetOpUserID(ctx),
		CreateTime:     time.Now(),
		Ex:             req.Ex,
	}

	if err := s.blackDatabase.Create(ctx, []*model.Black{&black}); err != nil {
		return nil, err
	}
	s.notificationSender.BlackAddedNotification(ctx, req)
	return &relation.AddBlackResp{}, nil
}

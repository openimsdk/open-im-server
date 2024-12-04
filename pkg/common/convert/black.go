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

package convert

import (
	"context"

	"github.com/openimsdk/protocol/sdkws"
	sdk "github.com/openimsdk/protocol/sdkws"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func BlackDB2Pb(ctx context.Context, blackDBs []*model.Black, f func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error)) (blackPbs []*sdk.BlackInfo, err error) {
	if len(blackDBs) == 0 {
		return nil, nil
	}
	var userIDs []string
	for _, blackDB := range blackDBs {
		userIDs = append(userIDs, blackDB.BlockUserID)
	}
	userInfos, err := f(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	for _, blackDB := range blackDBs {
		blackPb := &sdk.BlackInfo{
			OwnerUserID:    blackDB.OwnerUserID,
			CreateTime:     blackDB.CreateTime.Unix(),
			AddSource:      blackDB.AddSource,
			Ex:             blackDB.Ex,
			OperatorUserID: blackDB.OperatorUserID,
			BlackUserInfo: &sdkws.PublicUserInfo{
				UserID:   userInfos[blackDB.BlockUserID].UserID,
				Nickname: userInfos[blackDB.BlockUserID].Nickname,
				FaceURL:  userInfos[blackDB.BlockUserID].FaceURL,
				Ex:       userInfos[blackDB.BlockUserID].Ex,
			},
		}
		blackPbs = append(blackPbs, blackPb)
	}
	return blackPbs, nil
}

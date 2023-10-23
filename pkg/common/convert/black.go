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

	"github.com/OpenIMSDK/protocol/sdkws"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

func BlackDB2Pb(
	ctx context.Context,
	blackDBs []*relation.BlackModel,
	f func(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error),
) (blackPbs []*sdkws.BlackInfo, err error) {
	if len(blackDBs) == 0 {
		return nil, nil
	}
	userIDs := make([]string, 0, len(blackDBs))
	for _, blackDB := range blackDBs {
		userIDs = append(userIDs, blackDB.BlockUserID)
	}
	userInfos, err := f(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	for _, blackDB := range blackDBs {
		blackPb := &sdkws.BlackInfo{
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

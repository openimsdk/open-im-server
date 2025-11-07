package convert

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/openimsdk/protocol/sdkws"
	sdk "github.com/openimsdk/protocol/sdkws"
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

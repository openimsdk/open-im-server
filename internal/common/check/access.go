package check

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
)

func (u *UserCheck) Access(ctx context.Context, ownerUserID string) (err error) {
	_, err = u.GetUserInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}
	return tokenverify.CheckAccessV3(ctx, ownerUserID)
}

package check

import (
	"OpenIM/pkg/common/tokenverify"
	"context"
)

func (u *UserCheck) Access(ctx context.Context, ownerUserID string) (err error) {
	_, err = u.GetUsersInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}
	return tokenverify.CheckAccessV3(ctx, ownerUserID)
}

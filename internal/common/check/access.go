package check

import (
	"Open_IM/pkg/common/tokenverify"
	"context"
)

func Access(ctx context.Context, ownerUserID string) (err error) {
	_, err = GetUsersInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}
	return tokenverify.CheckAccessV3(ctx, ownerUserID)
}

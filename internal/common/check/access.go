package check

import (
	"Open_IM/pkg/common/token_verify"
	"context"
)

func Access(ctx context.Context, ownerUserID string) (err error) {
	_, err = GetUsersInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}
	return token_verify.CheckAccessV3(ctx, ownerUserID)
}

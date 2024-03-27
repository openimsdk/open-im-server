package checks

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type ServiceChecker interface {
	Check(ctx context.Context, config *config.GlobalConfig) error
}

func CheckServices(ctx context.Context, cfg *config.GlobalConfig, checkers []ServiceChecker) error {
	for _, checker := range checkers {
		if err := checker.Check(ctx, cfg); err != nil {
			return err
		}
	}
	return nil
}

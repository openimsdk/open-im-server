package user

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbuser "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/utils/datautil"
)

func (s *userServer) GetUserClientConfig(ctx context.Context, req *pbuser.GetUserClientConfigReq) (*pbuser.GetUserClientConfigResp, error) {
	if req.UserID != "" {
		if err := authverify.CheckAccessV3(ctx, req.UserID, s.config.Share.IMAdminUserID); err != nil {
			return nil, err
		}
		if _, err := s.db.GetUserByID(ctx, req.UserID); err != nil {
			return nil, err
		}
	}
	res, err := s.clientConfig.GetUserConfig(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetUserClientConfigResp{Configs: res}, nil
}

func (s *userServer) SetUserClientConfig(ctx context.Context, req *pbuser.SetUserClientConfigReq) (*pbuser.SetUserClientConfigResp, error) {
	if err := authverify.CheckAdmin(ctx, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	if req.UserID != "" {
		if _, err := s.db.GetUserByID(ctx, req.UserID); err != nil {
			return nil, err
		}
	}
	if err := s.clientConfig.SetUserConfig(ctx, req.UserID, req.Configs); err != nil {
		return nil, err
	}
	return &pbuser.SetUserClientConfigResp{}, nil
}

func (s *userServer) DelUserClientConfig(ctx context.Context, req *pbuser.DelUserClientConfigReq) (*pbuser.DelUserClientConfigResp, error) {
	if err := authverify.CheckAdmin(ctx, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	if err := s.clientConfig.DelUserConfig(ctx, req.UserID, req.Keys); err != nil {
		return nil, err
	}
	return &pbuser.DelUserClientConfigResp{}, nil
}

func (s *userServer) PageUserClientConfig(ctx context.Context, req *pbuser.PageUserClientConfigReq) (*pbuser.PageUserClientConfigResp, error) {
	if err := authverify.CheckAdmin(ctx, s.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	total, res, err := s.clientConfig.GetUserConfigPage(ctx, req.UserID, req.Key, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &pbuser.PageUserClientConfigResp{
		Total: total,
		Configs: datautil.Slice(res, func(e *model.ClientConfig) *pbuser.ClientConfig {
			return &pbuser.ClientConfig{
				UserID: e.UserID,
				Key:    e.Key,
				Value:  e.Value,
			}
		}),
	}, nil
}

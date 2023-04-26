package rpcclient

import (
	"context"
	"fmt"

	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"google.golang.org/grpc"
)

type MetaClient struct {
	// contains filtered or unexported fields
	client          discoveryregistry.SvcDiscoveryRegistry
	rpcRegisterName string
	getUsersInfo    func(ctx context.Context, userIDs []string) ([]CommonUser, error)
}

func NewMetaClient(client discoveryregistry.SvcDiscoveryRegistry, rpcRegisterName string, opts ...MetaClientOptions) *MetaClient {
	c := &MetaClient{
		client:          client,
		rpcRegisterName: rpcRegisterName,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type MetaClientOptions func(*MetaClient)

func WithDBFunc(fn func(ctx context.Context, userIDs []string) (users []*relationTb.UserModel, err error)) MetaClientOptions {
	return func(s *MetaClient) {
		f := func(ctx context.Context, userIDs []string) (result []CommonUser, err error) {
			users, err := fn(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			for _, user := range users {
				result = append(result, user)
			}
			return result, nil
		}
		s.getUsersInfo = f
	}
}

func WithRpcFunc(fn func(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error)) MetaClientOptions {
	return func(s *MetaClient) {
		f := func(ctx context.Context, userIDs []string) (result []CommonUser, err error) {
			users, err := fn(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			for _, user := range users {
				result = append(result, user)
			}
			return result, err
		}
		s.getUsersInfo = f
	}
}

func (m *MetaClient) getFaceURLAndName(userID string) (faceURL, nickname string, err error) {
	users, err := m.getUsersInfo(context.Background(), []string{userID})
	if err != nil {
		return "", "", err
	}
	if len(users) == 0 {
		return "", "", errs.ErrRecordNotFound.Wrap(fmt.Sprintf("notification user %s not found", userID))
	}
	return users[0].GetFaceURL(), users[0].GetNickname(), nil
}

func (m *MetaClient) getConn() (*grpc.ClientConn, error) {
	return m.client.GetConn(m.rpcRegisterName)
}

type CommonUser interface {
	GetNickname() string
	GetFaceURL() string
	GetUserID() string
	GetEx() string
}

type CommonGroup interface {
	GetNickname() string
	GetFaceURL() string
	GetGroupID() string
	GetEx() string
}

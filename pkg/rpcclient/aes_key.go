package rpcclient

import (
	"context"
	aesKey "github.com/OpenIMSDK/protocol/aeskey"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"google.golang.org/grpc"
)

type AesKey struct {
	conn   grpc.ClientConnInterface
	Client aesKey.AesKeyClient
	Discov discoveryregistry.SvcDiscoveryRegistry
}

func NewAesKey(discov discoveryregistry.SvcDiscoveryRegistry) *AesKey {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImAesKeyName)
	if err != nil {
		panic(err)
	}
	client := aesKey.NewAesKeyClient(conn)
	return &AesKey{Discov: discov, Client: client, conn: conn}
}

type AesKeyRpcClient AesKey

func NewAesKeyRpcClient(discov discoveryregistry.SvcDiscoveryRegistry) AesKeyRpcClient {
	return AesKeyRpcClient(*NewAesKey(discov))
}

func (a *AesKeyRpcClient) AcquireAesKey(ctx context.Context, conversationType int32, userID, friendUserID, groupID string) (*aesKey.AcquireAesKeyResp, error) {
	return a.Client.AcquireAesKey(ctx, &aesKey.AcquireAesKeyReq{ConversationType: conversationType, OwnerUserID: userID, FriendUserID: friendUserID, GroupID: groupID})
}

func (a *AesKeyRpcClient) AcquireAesKeys(ctx context.Context, userID string) (*aesKey.AcquireAesKeysResp, error) {
	return a.Client.AcquireAesKeys(ctx, &aesKey.AcquireAesKeysReq{UserID: userID})
}

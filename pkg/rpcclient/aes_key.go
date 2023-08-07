package rpcclient

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/aes_key"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"google.golang.org/grpc"
)

type AesKey struct {
	conn   grpc.ClientConnInterface
	Client aes_key.AesKeyClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func NewAesKey(discov discoveryregistry.SvcDiscoveryRegistry) *AesKey {
	conn, err := discov.GetConn(context.Background(), "aesKey")
	if err != nil {
		panic(err)
	}
	client := aes_key.NewAesKeyClient(conn)
	return &AesKey{discov: discov, conn: conn, Client: client}
}

type AesKeyRpcClient AesKey

func NewAesKeyRpcClient(discov discoveryregistry.SvcDiscoveryRegistry) AesKeyRpcClient {
	return AesKeyRpcClient(*NewAesKey(discov))
}
func (a *AesKeyRpcClient) GetKey(ctx context.Context, sId string, sType int32) (*aes_key.GetAesKeyResp, error) {
	key, err := a.Client.GetAesKey(ctx, &aes_key.GetAesKeyReq{SId: sId, SType: sType})
	if err != nil {
		return nil, err
	}
	return key, nil
}
func (a *AesKeyRpcClient) GetAllKey(ctx context.Context, uId string) (*aes_key.GetAllAesKeyResp, error) {
	key, err := a.Client.GetAllAesKey(ctx, &aes_key.GetAllAesKeyReq{UId: uId})
	if err != nil {
		return nil, err
	}
	return key, nil
}

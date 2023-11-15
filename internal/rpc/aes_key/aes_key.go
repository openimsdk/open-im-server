package aes_key

import (
	"context"
	key "github.com/OpenIMSDK/protocol/aeskey"
	registry "github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/relation"
	tablerelation "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"google.golang.org/grpc"
)

type aesKeyServer struct {
	aesKeyDatabase controller.AesKeyDatabase
	RegisterCenter registry.SvcDiscoveryRegistry
}

func Start(client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&tablerelation.AesKeyModel{}); err != nil {
		return err
	}
	gorm := relation.NewAesKeyGorm(db)
	key.RegisterAesKeyServer(server, &aesKeyServer{
		aesKeyDatabase: controller.NewAesKeyDatabase(gorm),
		RegisterCenter: client,
	})
	return nil
}
func (a *aesKeyServer) AcquireAesKey(ctx context.Context, req *key.AcquireAesKeyReq) (*key.AcquireAesKeyResp, error) {
	aesKey, err := a.aesKeyDatabase.AcquireAesKey(ctx, req.ConversationType, req.OwnerUserID, req.FriendUserID, req.GroupID)
	if err != nil {
		return nil, err
	}
	resp := key.AcquireAesKeyResp{}
	utils.CopyStructFields(resp.AesKey, &aesKey)
	return &resp, nil
}

func (a *aesKeyServer) AcquireAesKeys(ctx context.Context, req *key.AcquireAesKeysReq) (*key.AcquireAesKeysResp, error) {
	keysm, err := a.aesKeyDatabase.AcquireAesKeys(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	var keys []*key.AesKey
	resp := key.AcquireAesKeysResp{AesKeys: keys}
	utils.CopyStructFields(&resp.AesKeys, &keysm)
	return &resp, nil
}

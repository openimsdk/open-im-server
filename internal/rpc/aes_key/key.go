package aes_key

import (
	"context"
	utils "github.com/OpenIMSDK/Open-IM-Server/pkg/aes_utils"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/aes_key"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"google.golang.org/grpc"
)

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&relationTb.AesKeyModel{}); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	//redis, err := cache.NewRedis()
	//if err != nil {
	//	return err
	//}
	//friendDB := relation.NewFriendGorm(db)
	userRpcClient := rpcclient.NewAesKeyRpcClient(client)
	conversationRpcClient := rpcclient.NewConversationRpcClient(client)
	aesKeydatabase := controller.NewAesKeyDatabase(relation.NewAesKeyGorm(db))
	aes_key.RegisterAesKeyServer(server, &keyServer{
		AesKeyDatabase: aesKeydatabase,
		//FriendDatabase: controller.NewFriendDatabase(
		//	friendDB,
		//	relation.NewFriendRequestGorm(db),
		//	cache.NewFriendCacheRedis(redis, friendDB, cache.GetDefaultOpt()),
		//	tx.NewGorm(db),
		//),
		//GroupDatabase: nil,
		User:                  userRpcClient,
		conversationRpcClient: conversationRpcClient,
	})
	return nil
}

type keyServer struct {
	AesKeyDatabase        controller.AesKeyDatabase
	FriendDatabase        controller.FriendDatabase
	GroupDatabase         controller.GroupDatabase
	User                  rpcclient.AesKeyRpcClient
	conversationRpcClient rpcclient.ConversationRpcClient
}

func (k keyServer) GetAesKey(ctx context.Context, req *aes_key.GetAesKeyReq) (*aes_key.GetAesKeyResp, error) {
	key, err := k.AesKeyDatabase.GetAesKey(ctx, req.UId, req.SId, req.SType)
	if err != nil {
		aesKey := utils.GenerateAesKey(utils.Get2StringHash(req.UId, req.SId))
		model := relationTb.AesKeyModel{
			UserID:           req.UId,
			ConversationID:   req.SId,
			AesKey:           aesKey,
			ConversationType: req.SType,
		}
		err1 := k.AesKeyDatabase.InstallAesKey(ctx, model)
		if err1 != nil {
			return nil, err
		}
		model.UserID = req.SId
		model.ConversationID = req.UId
		err1 = k.AesKeyDatabase.InstallAesKey(ctx, model)
		if err1 != nil {
			return nil, err
		}
		return &aes_key.GetAesKeyResp{Key: aesKey}, nil
	}
	return &aes_key.GetAesKeyResp{Key: key.AesKey}, nil
}

func (k keyServer) GetAllAesKey(ctx context.Context, req *aes_key.GetAllAesKeyReq) (*aes_key.GetAllAesKeyResp, error) {
	key, err := k.AesKeyDatabase.GetAllAesKey(ctx, req.UId)
	if err != nil {
		return &aes_key.GetAllAesKeyResp{}, err
	}
	resp := &aes_key.GetAllAesKeyResp{Keys: make(map[string]string)}
	for i := range key {
		resp.Keys[key[i].ConversationID] = key[i].AesKey
	}
	return resp, nil
}

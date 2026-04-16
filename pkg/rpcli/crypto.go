package rpcli

import (
	"context"

	pbcrypto "github.com/openimsdk/protocol/crypto"
	"github.com/openimsdk/tools/log"
	"google.golang.org/grpc"
)

func NewCryptoClient(cc grpc.ClientConnInterface) *CryptoClient {
	return &CryptoClient{pbcrypto.NewCryptoServiceClient(cc)}
}

type CryptoClient struct {
	pbcrypto.CryptoServiceClient
}

func (x *CryptoClient) BumpGroupKeyVersion(ctx context.Context, groupID, operatorUserID, eventType string) {
	log.ZDebug(ctx, "BumpGroupKeyVersion start", "groupID", groupID, "operatorUserID", operatorUserID, "eventType", eventType)
	resp, err := x.CryptoServiceClient.BumpGroupKeyVersion(ctx, &pbcrypto.BumpGroupKeyVersionReq{
		GroupID:        groupID,
		OperatorUserID: operatorUserID,
		EventType:      eventType,
	})
	if err != nil {
		log.ZError(ctx, "BumpGroupKeyVersion failed", err,
			"groupID", groupID,
			"operatorUserID", operatorUserID,
			"eventType", eventType,
		)
		return
	}
	log.ZDebug(ctx, "BumpGroupKeyVersion success", "groupID", groupID, "newVersion", resp.GroupKeyVersion)
}

func (x *CryptoClient) InitGroupKeyVersion(ctx context.Context, groupID string) {
	log.ZDebug(ctx, "InitGroupKeyVersion start", "groupID", groupID, "eventType", "group_created")
	_, err := x.CryptoServiceClient.BumpGroupKeyVersion(ctx, &pbcrypto.BumpGroupKeyVersionReq{
		GroupID:   groupID,
		EventType: "group_created",
	})
	if err != nil {
		log.ZError(ctx, "InitGroupKeyVersion failed", err,
			"groupID", groupID,
			"eventType", "group_created",
		)
		return
	}
	log.ZDebug(ctx, "InitGroupKeyVersion success", "groupID", groupID)
}

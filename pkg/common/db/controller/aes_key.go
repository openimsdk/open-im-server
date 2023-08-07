package controller

import (
	"context"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
)

type aesKeyDatabase struct {
	aesKeyDB relationTb.AesKeyModelInterface
}

type AesKeyDatabase interface {
	//生成AesKey
	InstallAesKey(ctx context.Context, aesKey relationTb.AesKeyModel) error
	GetAesKey(ctx context.Context, userId, cid string, cType int32) (aesKey *relationTb.AesKeyModel, err error)
	GetAllAesKey(ctx context.Context, userId string) (aesKey []*relationTb.AesKeyModel, err error)
}

func NewAesKeyDatabase(aesKey relationTb.AesKeyModelInterface) AesKeyDatabase {
	return &aesKeyDatabase{
		aesKeyDB: aesKey,
	}
}

func (a *aesKeyDatabase) InstallAesKey(ctx context.Context, aesKey relationTb.AesKeyModel) error {
	err := a.aesKeyDB.Install(ctx, aesKey)
	if err != nil {
		return err
	}
	return nil
}
func (a *aesKeyDatabase) GetAesKey(ctx context.Context, userId, cid string, cType int32) (aesKey *relationTb.AesKeyModel, err error) {
	return a.aesKeyDB.GetAesKey(ctx, userId, cid, cType)
}

func (a *aesKeyDatabase) GetAllAesKey(ctx context.Context, userId string) (aesKey []*relationTb.AesKeyModel, err error) {
	return a.aesKeyDB.GetAllAesKey(ctx, userId)
}

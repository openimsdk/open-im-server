package controller

import (
	"context"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
)

type KeyDatabase interface {
	DBGetKey(ctx context.Context, cID string) (key relationTb.KeyModel, err error)
	DBInstallKey(ctx context.Context, key relationTb.KeyModel) (err error)
}
type keyDatabase struct {
	keyDB relationTb.KeyModelInterface
	//cache cache.KeyCache
}

func NewKeyDatabase(key relationTb.KeyModelInterface) KeyDatabase {
	return &keyDatabase{
		keyDB: key,
	}
}
func (k *keyDatabase) DBGetKey(ctx context.Context, cID string) (key relationTb.KeyModel, err error) {
	return k.keyDB.GetKey(ctx, cID)
}
func (k *keyDatabase) DBInstallKey(ctx context.Context, key relationTb.KeyModel) (err error) {
	return k.keyDB.InstallKey(ctx, key)
}

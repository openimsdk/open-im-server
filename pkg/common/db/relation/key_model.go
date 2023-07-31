package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

var _ relation.KeyModelInterface = (*KeyGorm)(nil)

type KeyGorm struct {
	*MetaDB
}

func NewKeyDB(db *gorm.DB) relation.KeyModelInterface {
	return &KeyGorm{NewMetaDB(db, &relation.KeyModel{})}
}
func (k KeyGorm) InstallKey(ctx context.Context, key relation.KeyModel) (err error) {
	return utils.Wrap(k.DB.Create(&key).Error, "")
}

func (k KeyGorm) GetKey(ctx context.Context, cid string) (key relation.KeyModel, err error) {
	err = k.db(ctx).Where("cid=?", cid).Find(&key).Error
	return key, err
}

package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/tools/utils"
	"gorm.io/gorm"
)

type AesKeyGorm struct {
	*MetaDB
}

func NewAesKeyGorm(db *gorm.DB) relation.AesKeyModelInterface {
	return &AesKeyGorm{NewMetaDB(db, &relation.AesKeyModel{})}
}
func (a AesKeyGorm) Install(ctx context.Context, aesKey relation.AesKeyModel) (err error) {
	return utils.Wrap(a.db(ctx).Create(&aesKey).Error, "")
}
func (a AesKeyGorm) GetAesKey(ctx context.Context, userId, cid string, cType int32) (aesKey *relation.AesKeyModel, err error) {
	return aesKey, utils.Wrap(a.db(ctx).Where("user_id = ? and conversation_id=? and conversation_type =?", userId, cid, cType).Find(&aesKey).Error, "")
}

func (a AesKeyGorm) GetAllAesKey(ctx context.Context, userId string) (aesKey []*relation.AesKeyModel, err error) {
	return aesKey, utils.Wrap(a.db(ctx).Where("user_id = ?", userId).Find(aesKey).Error, "")
}

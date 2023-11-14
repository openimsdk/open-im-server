package relation

import (
	"context"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"gorm.io/gorm"
)

type AesKeyGorm struct {
	*MetaDB
}

func NewAesKeyGorm(db *gorm.DB) *AesKeyGorm {
	return &AesKeyGorm{NewMetaDB(db, &relation.AesKeyModel{})}
}

func (a *AesKeyGorm) Installs(ctx context.Context, keys []*relation.AesKeyModel) (err error) {
	return utils.Wrap(a.db(ctx).Create(&keys).Error, "")
}

func (a *AesKeyGorm) GetAesKey(tx context.Context, KeyConversationsID string) (key *relation.AesKeyModel, err error) {
	key = &relation.AesKeyModel{}
	return key, utils.Wrap(a.db(tx).Where("key_conversations_id = ? ", KeyConversationsID).Take(key).Error, "")
}

func (a *AesKeyGorm) GetAllAesKey(tx context.Context, UserID string) (keys []*relation.AesKeyModel, err error) {
	return keys, utils.Wrap(a.db(tx).Where("owner_user_id = ? or friend_user_id = ? ", UserID, UserID).Take(keys).Error, "")
}

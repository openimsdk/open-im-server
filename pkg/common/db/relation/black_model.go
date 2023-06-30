package relation

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

type BlackGorm struct {
	*MetaDB
}

func NewBlackGorm(db *gorm.DB) relation.BlackModelInterface {
	return &BlackGorm{NewMetaDB(db, &relation.BlackModel{})}
}

func (b *BlackGorm) Create(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	return utils.Wrap(b.db(ctx).Create(&blacks).Error, "")
}

func (b *BlackGorm) Delete(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	return utils.Wrap(b.db(ctx).Delete(blacks).Error, "")
}

func (b *BlackGorm) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error) {
	return utils.Wrap(b.db(ctx).Where("block_user_id = ? and block_user_id = ?", ownerUserID, blockUserID).Updates(args).Error, "")
}

func (b *BlackGorm) Update(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	return utils.Wrap(b.db(ctx).Updates(&blacks).Error, "")
}

func (b *BlackGorm) Find(ctx context.Context, blacks []*relation.BlackModel) (blackList []*relation.BlackModel, err error) {
	var where [][]interface{}
	for _, black := range blacks {
		where = append(where, []interface{}{black.OwnerUserID, black.BlockUserID})
	}
	return blackList, utils.Wrap(b.db(ctx).Where("(owner_user_id, block_user_id) in ?", where).Find(&blackList).Error, "")
}

func (b *BlackGorm) Take(ctx context.Context, ownerUserID, blockUserID string) (black *relation.BlackModel, err error) {
	black = &relation.BlackModel{}
	return black, utils.Wrap(b.db(ctx).Where("owner_user_id = ? and block_user_id = ?", ownerUserID, blockUserID).Take(black).Error, "")
}

func (b *BlackGorm) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*relation.BlackModel, total int64, err error) {
	err = b.db(ctx).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	totalUint32, blacks, err := ormutil.GormPage[relation.BlackModel](b.db(ctx).Where("owner_user_id = ?", ownerUserID), pageNumber, showNumber)
	total = int64(totalUint32)
	return
}

func (b *BlackGorm) FindBlackUserIDs(ctx context.Context, ownerUserID string) (blackUserIDs []string, err error) {
	return blackUserIDs, utils.Wrap(b.db(ctx).Where("owner_user_id = ?", ownerUserID).Pluck("block_user_id", &blackUserIDs).Error, "")
}

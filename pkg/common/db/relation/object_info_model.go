package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
	"time"
)

func NewObjectInfo(db *gorm.DB) relation.ObjectInfoModelInterface {
	return &ObjectInfoGorm{
		DB: db,
	}
}

type ObjectInfoGorm struct {
	DB *gorm.DB
}

func (o *ObjectInfoGorm) NewTx(tx any) relation.ObjectInfoModelInterface {
	return &ObjectInfoGorm{
		DB: tx.(*gorm.DB),
	}
}

func (o *ObjectInfoGorm) SetObject(ctx context.Context, obj *relation.ObjectInfoModel) (err error) {
	return utils.Wrap1(o.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("name = ?", obj.Name).Delete(&relation.ObjectInfoModel{}).Error; err != nil {
			return err
		}
		return tx.Create(obj).Error
	}))
}

func (o *ObjectInfoGorm) Take(ctx context.Context, name string) (info *relation.ObjectInfoModel, err error) {
	info = &relation.ObjectInfoModel{}
	return info, utils.Wrap1(o.DB.Where("name = ?", name).Take(info).Error)
}

func (o *ObjectInfoGorm) DeleteExpiration(ctx context.Context, expiration time.Time) (err error) {
	return utils.Wrap1(o.DB.Where("expiration_time IS NOT NULL AND expiration_time <= ?", expiration).Delete(&relation.ObjectInfoModel{}).Error)
}

package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type ObjectInfoGorm struct {
	*MetaDB
}

func NewObjectInfo(db *gorm.DB) relation.ObjectInfoModelInterface {
	return &ObjectInfoGorm{
		NewMetaDB(db, &relation.ObjectInfoModel{}),
	}
}

func (o *ObjectInfoGorm) NewTx(tx any) relation.ObjectInfoModelInterface {
	return &ObjectInfoGorm{
		NewMetaDB(tx.(*gorm.DB), &relation.ObjectInfoModel{}),
	}
}

func (o *ObjectInfoGorm) SetObject(ctx context.Context, obj *relation.ObjectInfoModel) (err error) {
	if err := o.DB.WithContext(ctx).Where("name = ?", obj.Name).Delete(&relation.ObjectInfoModel{}).Error; err != nil {
		return errs.Wrap(err)
	}
	return errs.Wrap(o.DB.WithContext(ctx).Create(obj).Error)
	//return errs.Wrap(o.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
	//	if err := tx.Where("name = ?", obj.Name).Delete(&relation.ObjectInfoModel{}).Error; err != nil {
	//		return errs.Wrap(err)
	//	}
	//	return errs.Wrap(tx.Create(obj).Error)
	//}))
}

func (o *ObjectInfoGorm) Take(ctx context.Context, name string) (info *relation.ObjectInfoModel, err error) {
	info = &relation.ObjectInfoModel{}
	return info, utils.Wrap1(o.DB.WithContext(ctx).Where("name = ?", name).Take(info).Error)
}

func (o *ObjectInfoGorm) DeleteExpiration(ctx context.Context, expiration time.Time) (err error) {
	return utils.Wrap1(o.DB.WithContext(ctx).Where("expiration_time IS NOT NULL AND expiration_time <= ?", expiration).Delete(&relation.ObjectInfoModel{}).Error)
}

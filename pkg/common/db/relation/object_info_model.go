package relation

import (
	"OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"gorm.io/gorm"
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
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "objectInfo", obj)
	}()
	return utils.Wrap1(o.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("name = ?", obj.Name).Delete(&relation.ObjectInfoModel{}).Error; err != nil {
			return err
		}
		return tx.Create(obj).Error
	}))
}

func (o *ObjectInfoGorm) Take(ctx context.Context, name string) (info *relation.ObjectInfoModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "name", name, "info", info)
	}()
	info = &relation.ObjectInfoModel{}
	return info, utils.Wrap1(o.DB.Where("name = ?", name).Take(info).Error)
}

package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"gorm.io/gorm"
)

type ObjectInfoGorm struct {
	*MetaDB
}

func NewObjectInfo(db *gorm.DB) relation.ObjectInfoModelInterface {
	return &ObjectInfoGorm{
		NewMetaDB(db, &relation.ObjectModel{}),
	}
}

func (o *ObjectInfoGorm) NewTx(tx any) relation.ObjectInfoModelInterface {
	return &ObjectInfoGorm{
		NewMetaDB(tx.(*gorm.DB), &relation.ObjectModel{}),
	}
}

func (o *ObjectInfoGorm) SetObject(ctx context.Context, obj *relation.ObjectModel) (err error) {
	if err := o.DB.WithContext(ctx).Where("name = ?", obj.Name).FirstOrCreate(obj).Error; err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (o *ObjectInfoGorm) Take(ctx context.Context, name string) (info *relation.ObjectModel, err error) {
	info = &relation.ObjectModel{}
	return info, errs.Wrap(o.DB.WithContext(ctx).Where("name = ?", name).Take(info).Error)
}

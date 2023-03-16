package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type ObjectPutGorm struct {
	*MetaDB
}

func NewObjectPut(db *gorm.DB) relation.ObjectPutModelInterface {
	return &ObjectPutGorm{
		NewMetaDB(db, &relation.ObjectPutModel{}),
	}
}

func (o *ObjectPutGorm) NewTx(tx any) relation.ObjectPutModelInterface {
	return &ObjectPutGorm{
		NewMetaDB(tx.(*gorm.DB), &relation.ObjectPutModel{}),
	}
}

func (o *ObjectPutGorm) Create(ctx context.Context, m []*relation.ObjectPutModel) (err error) {
	return utils.Wrap1(o.DB.Create(m).Error)
}

func (o *ObjectPutGorm) Take(ctx context.Context, putID string) (put *relation.ObjectPutModel, err error) {
	put = &relation.ObjectPutModel{}
	return put, utils.Wrap1(o.DB.Where("put_id = ?", putID).Take(put).Error)
}

func (o *ObjectPutGorm) SetCompleted(ctx context.Context, putID string) (err error) {
	return utils.Wrap1(o.DB.Model(&relation.ObjectPutModel{}).Where("put_id = ?", putID).Update("complete", true).Error)
}

func (o *ObjectPutGorm) FindExpirationPut(ctx context.Context, expirationTime time.Time, num int) (list []*relation.ObjectPutModel, err error) {
	err = o.DB.Where("effective_time <= ?", expirationTime).Limit(num).Find(&list).Error
	return list, utils.Wrap1(err)
}

func (o *ObjectPutGorm) DelPut(ctx context.Context, ids []string) (err error) {
	return utils.Wrap1(o.DB.Where("put_id IN ?", ids).Delete(&relation.ObjectPutModel{}).Error)
}

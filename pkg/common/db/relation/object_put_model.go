package relation

import (
	"OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

func NewObjectPut(db *gorm.DB) relation.ObjectPutModelInterface {
	return &ObjectPutGorm{
		DB: db,
	}
}

type ObjectPutGorm struct {
	DB *gorm.DB
}

func (o *ObjectPutGorm) NewTx(tx any) relation.ObjectPutModelInterface {
	return &ObjectPutGorm{
		DB: tx.(*gorm.DB),
	}
}

func (o *ObjectPutGorm) Create(ctx context.Context, m []*relation.ObjectPutModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "objectPut", m)
	}()
	return utils.Wrap1(o.DB.Create(m).Error)
}

func (o *ObjectPutGorm) Take(ctx context.Context, putID string) (put *relation.ObjectPutModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "putID", putID, "put", put)
	}()
	put = &relation.ObjectPutModel{}
	return put, utils.Wrap1(o.DB.Where("put_id = ?", putID).Take(put).Error)
}

func (o *ObjectPutGorm) SetCompleted(ctx context.Context, putID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "putID", putID)
	}()
	return utils.Wrap1(o.DB.Model(&relation.ObjectPutModel{}).Where("put_id = ?", putID).Update("complete", true).Error)
}

func (o *ObjectPutGorm) FindExpirationPut(ctx context.Context, expirationTime time.Time, num int) (list []*relation.ObjectPutModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "expirationTime", expirationTime, "num", num, "list", list)
	}()
	err = o.DB.Where("effective_time <= ?", expirationTime).Limit(num).Find(&list).Error
	return list, utils.Wrap1(err)
}

func (o *ObjectPutGorm) DelPut(ctx context.Context, ids []string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ids", ids)
	}()
	return utils.Wrap1(o.DB.Where("put_id IN ?", ids).Delete(&relation.ObjectPutModel{}).Error)
}

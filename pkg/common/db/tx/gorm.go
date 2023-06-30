package tx

import (
	"gorm.io/gorm"
)

func NewGorm(db *gorm.DB) Tx {
	return &_Gorm{tx: db}
}

type _Gorm struct {
	tx *gorm.DB
}

func (g *_Gorm) Transaction(fn func(tx any) error) error {
	return g.tx.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

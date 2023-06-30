package relation

import (
	"context"

	"gorm.io/gorm"
)

type MetaDB struct {
	DB    *gorm.DB
	table any
}

func NewMetaDB(db *gorm.DB, table any) *MetaDB {
	return &MetaDB{
		DB:    db,
		table: table,
	}
}

func (g *MetaDB) db(ctx context.Context) *gorm.DB {
	return g.DB.WithContext(ctx).Model(g.table)
}

package relation

import (
	"context"
	"gorm.io/gorm"
)

type MetaDB struct {
	DB    *gorm.DB
	table interface{}
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

func (g *MetaDB) page(pageNumber, showNumber int32) *gorm.DB {
	return g.DB.Limit(int(showNumber)).Offset(int(pageNumber*showNumber-1))
}

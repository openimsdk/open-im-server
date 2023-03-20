package ormutil

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

func GormPage[E any](db *gorm.DB, pageNumber, showNumber int32) (uint32, []*E, error) {
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, nil, utils.Wrap(err, "")
	}
	var es []*E
	if err := db.Limit(int(showNumber)).Offset(int(pageNumber * showNumber)).Find(&es).Error; err != nil {
		return 0, nil, utils.Wrap(err, "")
	}
	return uint32(count), es, nil
}

func GormSearch[E any](db *gorm.DB, fields []string, value string, pageNumber, showNumber int32) (uint32, []*E, error) {
	if len(fields) > 0 && value != "" {
		value = "%" + value + "%"
		if len(fields) == 1 {
			db = db.Where(fields[0]+" like ?", value)
		} else {
			t := db
			for _, field := range fields {
				t = t.Or(field+" like ?", value)
			}
			db = db.Where(t)
		}
	}
	return GormPage[E](db, pageNumber, showNumber)
}

func GormIn[E any](db **gorm.DB, field string, es []E) {
	if len(es) == 0 {
		return
	}
	*db = (*db).Where(field+" in (?)", es)
}

func MapCount(db *gorm.DB, field string) (map[string]uint32, error) {
	var items []struct {
		ID    string `gorm:"column:id"`
		Count uint32 `gorm:"column:count"`
	}
	if err := db.Select(field + " as id, count(1) as count").Group(field).Find(&items).Error; err != nil {
		return nil, err
	}
	m := make(map[string]uint32)
	for _, item := range items {
		m[item.ID] = item.Count
	}
	return m, nil
}

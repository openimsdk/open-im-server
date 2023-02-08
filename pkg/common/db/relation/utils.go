package relation

import (
	"Open_IM/pkg/utils"
	"gorm.io/gorm"
)

func gormPage[E any](db *gorm.DB, pageNumber, showNumber int32) (int32, []*E, error) {
	var count int64
	if err := db.Model(new(E)).Count(&count).Error; err != nil {
		return 0, nil, utils.Wrap(err, "")
	}
	var es []*E
	if err := db.Limit(int(showNumber)).Offset(int(pageNumber * showNumber)).Find(&es).Error; err != nil {
		return 0, nil, utils.Wrap(err, "")
	}
	return int32(count), es, nil
}

func gormSearch[E any](db *gorm.DB, fields []string, value string, pageNumber, showNumber int32) (int32, []*E, error) {
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
	return gormPage[E](db, pageNumber, showNumber)
}

func gormIn[E any](db **gorm.DB, field string, es []E) {
	if len(es) == 0 {
		return
	}
	*db = (*db).Where(field+" in (?)", es)
}

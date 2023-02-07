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

func gormSearch[E any](db *gorm.DB, field string, value string, pageNumber, showNumber int32) (int32, []*E, error) {
	if field != "" && value != "" {
		db = db.Where(field+" like ?", "%"+value+"%")
	}
	return gormPage[E](db, pageNumber, showNumber)
}

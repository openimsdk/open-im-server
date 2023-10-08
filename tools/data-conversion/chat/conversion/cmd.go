package conversion

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func FindAndInsert[V2 any, V3 schema.Tabler](v2db *gorm.DB, v3db *gorm.DB, fn func(V2) V3) (string, error) {
	var t V3
	name := t.TableName()
	if err := v3db.AutoMigrate(&t); err != nil {
		return name, fmt.Errorf("auto migrate v3 %s failed %w", name, err)
	}
	const size = 100
	for i := 0; ; i++ {
		var v2s []V2
		if err := v2db.Offset(i * size).Limit(size).Find(&v2s).Error; err != nil {
			return name, fmt.Errorf("find v2 %s failed %w", name, err)
		}
		if len(v2s) == 0 {
			return name, nil
		}
		v3s := make([]V3, 0, len(v2s))
		for _, v := range v2s {
			v3s = append(v3s, fn(v))
		}
		if err := v3db.Create(&v3s).Error; err != nil {
			return name, fmt.Errorf("insert v3 %s failed %w", name, err)
		}
	}
}

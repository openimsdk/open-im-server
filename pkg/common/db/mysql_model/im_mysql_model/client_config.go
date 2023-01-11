package im_mysql_model

import (
	"gorm.io/gorm"
)

var InitConfigDB *gorm.DB

func SetClientInitConfig(m map[string]interface{}) error {
	result := InitConfigDB.Model(&ClientInitConfig{}).Where("1=1").Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := InitConfigDB.Model(&ClientInitConfig{}).Create(m).Error
		return err
	}

	return nil
}

func GetClientInitConfig() (ClientInitConfig, error) {
	var config ClientInitConfig
	err := InitConfigDB.Model(&ClientInitConfig{}).First(&config).Error
	return config, err
}

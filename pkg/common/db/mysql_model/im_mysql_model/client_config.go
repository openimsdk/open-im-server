package im_mysql_model

import "Open_IM/pkg/common/db"

func SetClientInitConfig(m map[string]interface{}) error {
	result := db.DB.MysqlDB.DefaultGormDB().Model(&ClientInitConfig{}).Where("1=1").Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := db.DB.MysqlDB.DefaultGormDB().Model(&ClientInitConfig{}).Create(m).Error
		return err
	}

	return nil
}

func GetClientInitConfig() (ClientInitConfig, error) {
	var config ClientInitConfig
	err := db.DB.MysqlDB.DefaultGormDB().Model(&ClientInitConfig{}).First(&config).Error
	return config, err
}

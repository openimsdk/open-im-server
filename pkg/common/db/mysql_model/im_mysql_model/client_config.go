package im_mysql_model

import "Open_IM/pkg/common/db"

func SetClientInitConfig(m map[string]interface{}) error {
	result := db.DB.MysqlDB.DefaultGormDB().Model(&db.ClientInitConfig{}).Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 1 {
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.ClientInitConfig{}).Create(m).Error
		return err
	}

	return nil
}

func GetClientInitConfig() (db.ClientInitConfig, error) {
	var config db.ClientInitConfig
	err := db.DB.MysqlDB.DefaultGormDB().Model((&db.ClientInitConfig{})).First(&config).Error
	return config, err
}

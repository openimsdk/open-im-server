package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func UpdateAppVersion(appType int, version string, forceUpdate bool, fileName, yamlName string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	app := db.AppVersion{
		Version:    version,
		Type:       appType,
		UpdateTime: int(time.Now().Unix()),
		FileName:   fileName,
		YamlName:   yamlName,
	}
	result := dbConn.Model(db.AppVersion{}).Where("type = ?", appType).Updates(&app).Update(map[string]interface{}{"force_update": forceUpdate})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := dbConn.Create(&app).Error
		return err
	}
	return nil
}

func GetNewestVersion(appType int) (*db.AppVersion, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	app := db.AppVersion{}
	if err != nil {
		return &app, err
	}
	dbConn.LogMode(true)
	return &app, dbConn.Model(db.AppVersion{}).First(&app, appType).Error
}

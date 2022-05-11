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
	updateTime := int(time.Now().Unix())
	app := db.AppVersion{
		Version:     version,
		Type:        appType,
		UpdateTime:  updateTime,
		FileName:    fileName,
		YamlName:    yamlName,
		ForceUpdate: forceUpdate,
	}
	result := dbConn.Model(db.AppVersion{}).Where("type = ?", appType).Update(map[string]interface{}{"force_update": forceUpdate,
		"version": version, "update_time": int(time.Now().Unix()), "file_name": fileName, "yaml_name": yamlName, "type": appType})
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

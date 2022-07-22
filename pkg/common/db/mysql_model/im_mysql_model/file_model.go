package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func UpdateAppVersion(appType int, version string, forceUpdate bool, fileName, yamlName, updateLog string) error {
	updateTime := int(time.Now().Unix())
	app := db.AppVersion{
		Version:     version,
		Type:        appType,
		UpdateTime:  updateTime,
		FileName:    fileName,
		YamlName:    yamlName,
		ForceUpdate: forceUpdate,
		UpdateLog:   updateLog,
	}
	result := db.DB.MysqlDB.DefaultGormDB().Model(db.AppVersion{}).Where("type = ?", appType).Updates(map[string]interface{}{"force_update": forceUpdate,
		"version": version, "update_time": int(time.Now().Unix()), "file_name": fileName, "yaml_name": yamlName, "type": appType, "update_log": updateLog})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := db.DB.MysqlDB.DefaultGormDB().Create(&app).Error
		return err
	}
	return nil
}

func GetNewestVersion(appType int) (*db.AppVersion, error) {
	app := db.AppVersion{}
	return &app, db.DB.MysqlDB.DefaultGormDB().Model(db.AppVersion{}).First(&app, appType).Error
}

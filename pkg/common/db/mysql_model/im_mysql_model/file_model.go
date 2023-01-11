package im_mysql_model

import (
	"gorm.io/gorm"
	"time"
)

var AppDB *gorm.DB

func UpdateAppVersion(appType int, version string, forceUpdate bool, fileName, yamlName, updateLog string) error {
	updateTime := int(time.Now().Unix())
	app := AppVersion{
		Version:     version,
		Type:        appType,
		UpdateTime:  updateTime,
		FileName:    fileName,
		YamlName:    yamlName,
		ForceUpdate: forceUpdate,
		UpdateLog:   updateLog,
	}
	result := AppDB.Model(AppVersion{}).Where("type = ?", appType).Updates(map[string]interface{}{"force_update": forceUpdate,
		"version": version, "update_time": int(time.Now().Unix()), "file_name": fileName, "yaml_name": yamlName, "type": appType, "update_log": updateLog})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := AppDB.Create(&app).Error
		return err
	}
	return nil
}

func GetNewestVersion(appType int) (*AppVersion, error) {
	app := AppVersion{}
	return &app, AppDB.Model(AppVersion{}).First(&app, appType).Error
}

package mysql

import (
	"gorm.io/gorm"
	"time"
)

var AppDB *gorm.DB

type AppVersion struct {
	Version     string `gorm:"column:version;size:64" json:"version"`
	Type        int    `gorm:"column:type;primary_key" json:"type"`
	UpdateTime  int    `gorm:"column:update_time" json:"update_time"`
	ForceUpdate bool   `gorm:"column:force_update" json:"force_update"`
	FileName    string `gorm:"column:file_name" json:"file_name"`
	YamlName    string `gorm:"column:yaml_name" json:"yaml_name"`
	UpdateLog   string `gorm:"column:update_log" json:"update_log"`
}

func (AppVersion) TableName() string {
	return "app_version"
}

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

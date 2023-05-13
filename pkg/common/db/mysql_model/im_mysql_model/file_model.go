// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

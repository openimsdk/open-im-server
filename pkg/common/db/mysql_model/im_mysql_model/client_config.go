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

import "Open_IM/pkg/common/db"

func SetClientInitConfig(m map[string]interface{}) error {
	result := db.DB.MysqlDB.DefaultGormDB().Model(&db.ClientInitConfig{}).Where("1=1").Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		err := db.DB.MysqlDB.DefaultGormDB().Model(&db.ClientInitConfig{}).Create(m).Error
		return err
	}

	return nil
}

func GetClientInitConfig() (db.ClientInitConfig, error) {
	var config db.ClientInitConfig
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.ClientInitConfig{}).First(&config).Error
	return config, err
}

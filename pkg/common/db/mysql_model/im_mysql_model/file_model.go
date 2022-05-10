package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"time"
)

func UpdateAppVersion(appType int, version string, forceUpdate bool) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	app := db.AppVersion{
		Version:     version,
		Type:        appType,
		UpdateTime:  int(time.Now().Unix()),
		ForceUpdate: forceUpdate,
	}
	result := dbConn.Model(db.AppVersion{}).Where("app_type = ?", appType).Updates(&app)
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
	if err != nil {
		return nil, err
	}
	dbConn.LogMode(true)
	app := db.AppVersion{}
	return &app, dbConn.Model(db.AppVersion{}).First(&app, appType).Error
}

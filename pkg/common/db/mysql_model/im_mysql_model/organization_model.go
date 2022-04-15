package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"errors"
	"time"
)

func CreateDepartment(department *db.Department) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	department.CreateTime = time.Now()
	return dbConn.Table("departments").Create(department).Error
}

func GetDepartment(departmentID string) (error, *db.Department) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err, nil
	}
	var department db.Department
	err = dbConn.Table("departments").Where("department_id=?", departmentID).Find(&department).Error
	return err, &department
}

func UpdateDepartment(department *db.Department, args map[string]interface{}) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	if err = dbConn.Table("departments").Where("department_id=?", department.DepartmentID).Updates(department).Error; err != nil {
		return err
	}
	if args != nil {
		return dbConn.Table("departments").Where("department_id=?", department.DepartmentID).Updates(args).Error
	}
	return nil
}

func GetSubDepartmentList(departmentID string) (error, []db.Department) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err, nil
	}
	var departmentList []db.Department
	err = dbConn.Table("departments").Where("parent_id=?", departmentID).Find(&departmentList).Error
	return err, departmentList
}

func DeleteDepartment(departmentID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	if err = dbConn.Table("departments").Where("department_id=?", departmentID).Delete(db.Department{}).Error; err != nil {
		return err
	}
	if err = dbConn.Table("department_members").Where("department_id=?", departmentID).Delete(db.DepartmentMember{}).Error; err != nil {
		return err
	}
	return nil
}

func CreateOrganizationUser(organizationUser *db.OrganizationUser) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	organizationUser.CreateTime = time.Now()
	return dbConn.Table("organization_users").Create(organizationUser).Error
}

func UpdateOrganizationUser(organizationUser *db.OrganizationUser, args map[string]interface{}) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	if err = dbConn.Table("organization_users").Where("user_id=?", organizationUser.UserID).Updates(organizationUser).Error; err != nil {
		return err
	}
	if args != nil {
		return dbConn.Table("organization_users").Where("user_id=?", organizationUser.UserID).Updates(args).Error
	}
	return nil
}

func CreateDepartmentMember(departmentMember *db.DepartmentMember) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	departmentMember.CreateTime = time.Now()
	return dbConn.Table("department_members").Create(departmentMember).Error
}

func GetUserInDepartment(userID string) (error, []db.DepartmentMember) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err, nil
	}
	var departmentMemberList []db.DepartmentMember
	err = dbConn.Table("department_members").Where("user_id=?", userID).Find(&departmentMemberList).Error
	return err, departmentMemberList
}

func UpdateUserInDepartment(userInDepartmentList []db.DepartmentMember) error {
	if len(userInDepartmentList) == 0 {
		return errors.New("args failed")
	}
	if err := DeleteUserInAllDepartment(userInDepartmentList[0].UserID); err != nil {
		return err
	}
	for _, v := range userInDepartmentList {
		if err := CreateDepartmentMember(&v); err != nil {
			return err
		}
	}
	return nil
}

func DeleteUserInAllDepartment(userID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	return dbConn.Table("department_members").Where("user_id=?", userID).Delete(db.DepartmentMember{}).Error
}

func DeleteOrganizationUser(OrganizationUserID string) error {
	if err := DeleteUserInAllDepartment(OrganizationUserID); err != nil {
		return err
	}
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	return dbConn.Table("organization_users").Where("user_id=?", OrganizationUserID).Delete(db.OrganizationUser{}).Error
}

func GetDepartmentMemberUserIDList(departmentID string) (error, []string) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err, nil
	}
	var departmentMemberList []db.DepartmentMember
	err = dbConn.Table("department_members").Where("department_id=?", departmentID).Find(&departmentMemberList).Error

	var userIDList []string = make([]string, 0)
	for _, v := range departmentMemberList {
		userIDList = append(userIDList, v.UserID)
	}
	return err, userIDList
}

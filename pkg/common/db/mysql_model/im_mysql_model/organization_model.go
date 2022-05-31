package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"github.com/jinzhu/gorm"
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
	if departmentID == "-1" {
		err = dbConn.Table("departments").Find(&departmentList).Error
	} else {
		err = dbConn.Table("departments").Where("parent_id=?", departmentID).Find(&departmentList).Error
	}

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

func GetOrganizationUser(userID string) (error, *db.OrganizationUser) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err, nil
	}
	organizationUser := db.OrganizationUser{}
	err = dbConn.Table("organization_users").Where("user_id=?", userID).Take(&organizationUser).Error
	return err, &organizationUser
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

func UpdateUserInDepartment(departmentMember *db.DepartmentMember, args map[string]interface{}) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	if err = dbConn.Table("department_members").Where("department_id=? AND user_id=?", departmentMember.DepartmentID, departmentMember.UserID).
		Updates(departmentMember).Error; err != nil {
		return err
	}
	if args != nil {
		return dbConn.Table("department_members").Where("department_id=? AND user_id=?", departmentMember.DepartmentID, departmentMember.UserID).
			Updates(args).Error
	}
	return nil
}

func DeleteUserInDepartment(departmentID, userID string) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	return dbConn.Table("department_members").Where("department_id=? AND user_id=?", departmentID, userID).Delete(db.DepartmentMember{}).Error
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
	err = dbConn.Table("department_members").Where("department_id=?", departmentID).Take(&departmentMemberList).Error
	if err != nil {
		return err, nil
	}
	var userIDList []string = make([]string, 0)
	for _, v := range departmentMemberList {
		userIDList = append(userIDList, v.UserID)
	}
	return err, userIDList
}

func GetDepartmentMemberList(departmentID string) (error, []db.DepartmentMember) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err, nil
	}
	var departmentMemberList []db.DepartmentMember
	if departmentID == "-1" {
		err = dbConn.Table("department_members").Find(&departmentMemberList).Error
	} else {
		err = dbConn.Table("department_members").Where("department_id=?", departmentID).Find(&departmentMemberList).Error
	}

	if err != nil {
		return err, nil
	}
	return err, departmentMemberList
}

func GetAllOrganizationUserID() (error, []string) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err, nil
	}
	var OrganizationUser db.OrganizationUser
	var result []string
	return dbConn.Model(&OrganizationUser).Pluck("user_id", &result).Error, result
}

func GetDepartmentMemberNum(departmentID string) (error, uint32) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return utils.Wrap(err, "DefaultGormDB failed"), 0
	}
	var number uint32
	err = dbConn.Table("department_members").Where("department_id=?", departmentID).Count(&number).Error
	if err != nil {
		return utils.Wrap(err, ""), 0
	}
	return nil, number

}

func GetSubDepartmentNum(departmentID string) (error, uint32) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return utils.Wrap(err, "DefaultGormDB failed"), 0
	}
	var number uint32
	err = dbConn.Table("departments").Where("parent_id=?", departmentID).Count(&number).Error
	if err != nil {
		return utils.Wrap(err, ""), 0
	}
	return nil, number
}

func GetDepartmentRelatedGroupIDList(departmentIDList []string) ([]string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, utils.Wrap(err, "DefaultGormDB failed")
	}
	var groupIDList []string
	err = dbConn.Table("departments").Where("department_id IN (?) ", departmentIDList).Pluck("related_group_id", &groupIDList).Error
	return groupIDList, err
}

func getDepartmentParent(departmentID string, dbConn *gorm.DB) (*db.Department, error) {
	var department db.Department
	var parentID string
	dbConn.LogMode(true)
	// select * from departments where department_id = (select parent_id from departments where department_id= zx234fd);
	err := dbConn.Table("departments").Where("department_id=?", dbConn.Table("departments").Where("department_id=?", departmentID).Pluck("parent_id", parentID)).Find(&department).Error
	return &department, err
}

func GetDepartmentParent(departmentID string, dbConn *gorm.DB, parentIDList *[]string) error {
	department, err := getDepartmentParent(departmentID, dbConn)
	if err != nil {
		return err
	}
	if department.ParentID != "" {
		*parentIDList = append(*parentIDList, department.ParentID)
		err = GetDepartmentParent(departmentID, dbConn, parentIDList)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDepartmentParentIDList(departmentID string) ([]string, error) {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return nil, err
	}
	var parentIDList []string
	err = GetDepartmentParent(departmentID, dbConn, &parentIDList)
	return parentIDList, err
}

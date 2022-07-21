package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"gorm.io/gorm"
	"time"
)

func CreateDepartment(department *db.Department) error {
	department.CreateTime = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Table("departments").Create(department).Error
}

func GetDepartment(departmentID string) (*db.Department, error) {
	var department db.Department
	err := db.DB.MysqlDB.DefaultGormDB().Table("departments").Where("department_id=?", departmentID).Find(&department).Error
	return &department, err
}

func UpdateDepartment(department *db.Department, args map[string]interface{}) error {
	if err := db.DB.MysqlDB.DefaultGormDB().Table("departments").Where("department_id=?", department.DepartmentID).Updates(department).Error; err != nil {
		return err
	}
	if args != nil {
		return db.DB.MysqlDB.DefaultGormDB().Table("departments").Where("department_id=?", department.DepartmentID).Updates(args).Error
	}
	return nil
}

func GetSubDepartmentList(departmentID string) ([]db.Department, error) {
	var departmentList []db.Department
	var err error
	if departmentID == "-1" {
		err = db.DB.MysqlDB.DefaultGormDB().Table("departments").Find(&departmentList).Error
	} else {
		err = db.DB.MysqlDB.DefaultGormDB().Table("departments").Where("parent_id=?", departmentID).Find(&departmentList).Error
	}

	return departmentList, err
}

func DeleteDepartment(departmentID string) error {
	var err error
	if err = db.DB.MysqlDB.DefaultGormDB().Table("departments").Where("department_id=?", departmentID).Delete(db.Department{}).Error; err != nil {
		return err
	}
	if err = db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("department_id=?", departmentID).Delete(db.DepartmentMember{}).Error; err != nil {
		return err
	}
	return nil
}

func CreateOrganizationUser(organizationUser *db.OrganizationUser) error {
	organizationUser.CreateTime = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Table("organization_users").Create(organizationUser).Error
}

func GetOrganizationUser(userID string) (error, *db.OrganizationUser) {
	organizationUser := db.OrganizationUser{}
	err := db.DB.MysqlDB.DefaultGormDB().Table("organization_users").Where("user_id=?", userID).Take(&organizationUser).Error
	return err, &organizationUser
}

func UpdateOrganizationUser(organizationUser *db.OrganizationUser, args map[string]interface{}) error {
	if err := db.DB.MysqlDB.DefaultGormDB().Table("organization_users").Where("user_id=?", organizationUser.UserID).Updates(organizationUser).Error; err != nil {
		return err
	}
	if args != nil {
		return db.DB.MysqlDB.DefaultGormDB().Table("organization_users").Where("user_id=?", organizationUser.UserID).Updates(args).Error
	}
	return nil
}

func CreateDepartmentMember(departmentMember *db.DepartmentMember) error {
	departmentMember.CreateTime = time.Now()
	return db.DB.MysqlDB.DefaultGormDB().Table("department_members").Create(departmentMember).Error
}

func GetUserInDepartment(userID string) (error, []db.DepartmentMember) {
	var departmentMemberList []db.DepartmentMember
	err := db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("user_id=?", userID).Find(&departmentMemberList).Error
	return err, departmentMemberList
}

func UpdateUserInDepartment(departmentMember *db.DepartmentMember, args map[string]interface{}) error {
	if err := db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("department_id=? AND user_id=?", departmentMember.DepartmentID, departmentMember.UserID).
		Updates(departmentMember).Error; err != nil {
		return err
	}
	if args != nil {
		return db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("department_id=? AND user_id=?", departmentMember.DepartmentID, departmentMember.UserID).
			Updates(args).Error
	}
	return nil
}

func DeleteUserInDepartment(departmentID, userID string) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("department_id=? AND user_id=?", departmentID, userID).Delete(db.DepartmentMember{}).Error
}

func DeleteUserInAllDepartment(userID string) error {
	return db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("user_id=?", userID).Delete(db.DepartmentMember{}).Error
}

func DeleteOrganizationUser(OrganizationUserID string) error {
	if err := DeleteUserInAllDepartment(OrganizationUserID); err != nil {
		return err
	}
	return db.DB.MysqlDB.DefaultGormDB().Table("organization_users").Where("user_id=?", OrganizationUserID).Delete(db.OrganizationUser{}).Error
}

func GetDepartmentMemberUserIDList(departmentID string) (error, []string) {
	var departmentMemberList []db.DepartmentMember
	err := db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("department_id=?", departmentID).Take(&departmentMemberList).Error
	if err != nil {
		return err, nil
	}
	var userIDList []string = make([]string, 0)
	for _, v := range departmentMemberList {
		userIDList = append(userIDList, v.UserID)
	}
	return err, userIDList
}

func GetDepartmentMemberList(departmentID string) ([]db.DepartmentMember, error) {
	var departmentMemberList []db.DepartmentMember
	var err error
	if departmentID == "-1" {
		err = db.DB.MysqlDB.DefaultGormDB().Table("department_members").Find(&departmentMemberList).Error
	} else {
		err = db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("department_id=?", departmentID).Find(&departmentMemberList).Error
	}

	if err != nil {
		return nil, err
	}
	return departmentMemberList, err
}

func GetAllOrganizationUserID() (error, []string) {
	var OrganizationUser db.OrganizationUser
	var result []string
	return db.DB.MysqlDB.DefaultGormDB().Model(&OrganizationUser).Pluck("user_id", &result).Error, result
}

func GetDepartmentMemberNum(departmentID string) (error, uint32) {
	var number int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("department_members").Where("department_id=?", departmentID).Count(&number).Error
	if err != nil {
		return utils.Wrap(err, ""), 0
	}
	return nil, uint32(number)

}

func GetSubDepartmentNum(departmentID string) (error, uint32) {
	var number int64
	err := db.DB.MysqlDB.DefaultGormDB().Table("departments").Where("parent_id=?", departmentID).Count(&number).Error
	if err != nil {
		return utils.Wrap(err, ""), 0
	}
	return nil, uint32(number)
}

func SetDepartmentRelatedGroupID(groupID, departmentID string) error {
	department := &db.Department{RelatedGroupID: groupID}
	return db.DB.MysqlDB.DefaultGormDB().Model(&department).Where("department_id=?", departmentID).Updates(department).Error
}

func GetDepartmentRelatedGroupIDList(departmentIDList []string) ([]string, error) {
	var groupIDList []string
	err := db.DB.MysqlDB.DefaultGormDB().Table("departments").Where("department_id IN (?) ", departmentIDList).Pluck("related_group_id", &groupIDList).Error
	return groupIDList, err
}

func getDepartmentParent(departmentID string, dbConn *gorm.DB) (*db.Department, error) {
	var department db.Department
	var parentDepartment db.Department
	//var parentID string
	err := db.DB.MysqlDB.DefaultGormDB().Model(&department).Where("department_id=?", departmentID).Select("parent_id").First(&department).Error
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	if department.ParentID != "" {
		err = dbConn.Model(&parentDepartment).Where("department_id = ?", department.ParentID).Find(&parentDepartment).Error
	}
	return &parentDepartment, utils.Wrap(err, "")
}

func GetDepartmentParent(departmentID string, dbConn *gorm.DB, parentIDList *[]string) error {
	department, err := getDepartmentParent(departmentID, dbConn)
	if err != nil {
		return err
	}
	if department.DepartmentID != "" {
		*parentIDList = append(*parentIDList, department.DepartmentID)
		err = GetDepartmentParent(department.DepartmentID, dbConn, parentIDList)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDepartmentParentIDList(departmentID string) ([]string, error) {
	dbConn := db.DB.MysqlDB.DefaultGormDB()
	var parentIDList []string
	err := GetDepartmentParent(departmentID, dbConn, &parentIDList)
	return parentIDList, err
}

func GetRandomDepartmentID() (string, error) {
	department := &db.Department{}
	err := db.DB.MysqlDB.DefaultGormDB().Model(department).Order("RAND()").Where("related_group_id != ? AND department_id != ? AND department_type = ?", "", "0", 1).First(department).Error
	return department.DepartmentID, err
}

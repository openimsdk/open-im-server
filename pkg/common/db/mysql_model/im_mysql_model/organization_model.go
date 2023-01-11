package im_mysql_model

import (
	"Open_IM/pkg/utils"
	"gorm.io/gorm"
	"time"
)

var OrgDB *gorm.DB

func CreateDepartment(department *Department) error {
	department.CreateTime = time.Now()
	return OrgDB.Table("departments").Create(department).Error
}

func GetDepartment(departmentID string) (*Department, error) {
	var department Department
	err := OrgDB.Table("departments").Where("department_id=?", departmentID).Find(&department).Error
	return &department, err
}

func UpdateDepartment(department *Department, args map[string]interface{}) error {
	if err := OrgDB.Table("departments").Where("department_id=?", department.DepartmentID).Updates(department).Error; err != nil {
		return err
	}
	if args != nil {
		return OrgDB.Table("departments").Where("department_id=?", department.DepartmentID).Updates(args).Error
	}
	return nil
}

func GetSubDepartmentList(departmentID string) ([]Department, error) {
	var departmentList []Department
	var err error
	if departmentID == "-1" {
		err = OrgDB.Table("departments").Find(&departmentList).Error
	} else {
		err = OrgDB.Table("departments").Where("parent_id=?", departmentID).Find(&departmentList).Error
	}

	return departmentList, err
}

func DeleteDepartment(departmentID string) error {
	var err error
	if err = OrgDB.Table("departments").Where("department_id=?", departmentID).Delete(Department{}).Error; err != nil {
		return err
	}
	if err = OrgDB.Table("department_members").Where("department_id=?", departmentID).Delete(DepartmentMember{}).Error; err != nil {
		return err
	}
	return nil
}

func CreateOrganizationUser(organizationUser *OrganizationUser) error {
	organizationUser.CreateTime = time.Now()
	return OrgDB.Table("organization_users").Create(organizationUser).Error
}

func GetOrganizationUser(userID string) (error, *OrganizationUser) {
	organizationUser := OrganizationUser{}
	err := OrgDB.Table("organization_users").Where("user_id=?", userID).Take(&organizationUser).Error
	return err, &organizationUser
}

func GetOrganizationUsers(userIDList []string) ([]*OrganizationUser, error) {
	var organizationUserList []*OrganizationUser
	err := OrgDB.Table("organization_users").Where("user_id in (?)", userIDList).Find(&organizationUserList).Error
	return organizationUserList, err
}

func UpdateOrganizationUser(organizationUser *OrganizationUser, args map[string]interface{}) error {
	if err := OrgDB.Table("organization_users").Where("user_id=?", organizationUser.UserID).Updates(organizationUser).Error; err != nil {
		return err
	}
	if args != nil {
		return OrgDB.Table("organization_users").Where("user_id=?", organizationUser.UserID).Updates(args).Error
	}
	return nil
}

func CreateDepartmentMember(departmentMember *DepartmentMember) error {
	departmentMember.CreateTime = time.Now()
	return OrgDB.Table("department_members").Create(departmentMember).Error
}

func GetUserInDepartment(userID string) (error, []DepartmentMember) {
	var departmentMemberList []DepartmentMember
	err := OrgDB.Where("user_id=?", userID).Find(&departmentMemberList).Error
	return err, departmentMemberList
}

func UpdateUserInDepartment(departmentMember *DepartmentMember, args map[string]interface{}) error {
	if err := OrgDB.Where("department_id=? AND user_id=?", departmentMember.DepartmentID, departmentMember.UserID).
		Updates(departmentMember).Error; err != nil {
		return err
	}
	if args != nil {
		return OrgDB.Where("department_id=? AND user_id=?", departmentMember.DepartmentID, departmentMember.UserID).
			Updates(args).Error
	}
	return nil
}

func DeleteUserInDepartment(departmentID, userID string) error {
	return OrgDB.Table("department_members").Where("department_id=? AND user_id=?", departmentID, userID).Delete(DepartmentMember{}).Error
}

func DeleteUserInAllDepartment(userID string) error {
	return OrgDB.Table("department_members").Where("user_id=?", userID).Delete(DepartmentMember{}).Error
}

func DeleteOrganizationUser(OrganizationUserID string) error {
	if err := DeleteUserInAllDepartment(OrganizationUserID); err != nil {
		return err
	}
	return OrgDB.Table("organization_users").Where("user_id=?", OrganizationUserID).Delete(OrganizationUser{}).Error
}

func GetDepartmentMemberUserIDList(departmentID string) (error, []string) {
	var departmentMemberList []DepartmentMember
	err := OrgDB.Table("department_members").Where("department_id=?", departmentID).Take(&departmentMemberList).Error
	if err != nil {
		return err, nil
	}
	var userIDList []string = make([]string, 0)
	for _, v := range departmentMemberList {
		userIDList = append(userIDList, v.UserID)
	}
	return err, userIDList
}

func GetDepartmentMemberList(departmentID string) ([]DepartmentMember, error) {
	var departmentMemberList []DepartmentMember
	var err error
	if departmentID == "-1" {
		err = OrgDB.Table("department_members").Find(&departmentMemberList).Error
	} else {
		err = OrgDB.Table("department_members").Where("department_id=?", departmentID).Find(&departmentMemberList).Error
	}

	if err != nil {
		return nil, err
	}
	return departmentMemberList, err
}

func GetAllOrganizationUserID() (error, []string) {
	var OrganizationUser OrganizationUser
	var result []string
	return OrgDB.Model(&OrganizationUser).Pluck("user_id", &result).Error, result
}

func GetDepartmentMemberNum(departmentID string) (error, uint32) {
	var number int64
	err := OrgDB.Table("department_members").Where("department_id=?", departmentID).Count(&number).Error
	if err != nil {
		return utils.Wrap(err, ""), 0
	}
	return nil, uint32(number)

}

func GetSubDepartmentNum(departmentID string) (error, uint32) {
	var number int64
	err := OrgDB.Table("departments").Where("parent_id=?", departmentID).Count(&number).Error
	if err != nil {
		return utils.Wrap(err, ""), 0
	}
	return nil, uint32(number)
}

func SetDepartmentRelatedGroupID(groupID, departmentID string) error {
	department := &Department{RelatedGroupID: groupID}
	return OrgDB.Model(&department).Where("department_id=?", departmentID).Updates(department).Error
}

func GetDepartmentRelatedGroupIDList(departmentIDList []string) ([]string, error) {
	var groupIDList []string
	err := OrgDB.Table("departments").Where("department_id IN (?) ", departmentIDList).Pluck("related_group_id", &groupIDList).Error
	return groupIDList, err
}

func getDepartmentParent(departmentID string, dbConn *gorm.DB) (*Department, error) {
	var department Department
	var parentDepartment Department
	//var parentID string
	err := OrgDB.Model(&department).Where("department_id=?", departmentID).Select("parent_id").First(&department).Error
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
	dbConn := OrgDB
	var parentIDList []string
	err := GetDepartmentParent(departmentID, dbConn, &parentIDList)
	return parentIDList, err
}

func GetRandomDepartmentID() (string, error) {
	department := &Department{}
	err := OrgDB.Model(department).Order("RAND()").Where("related_group_id != ? AND department_id != ? AND department_type = ?", "", "0", 1).First(department).Error
	return department.DepartmentID, err
}

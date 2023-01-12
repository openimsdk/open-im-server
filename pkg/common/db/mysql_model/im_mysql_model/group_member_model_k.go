package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

var GroupMemberDB *gorm.DB

type GroupMember struct {
	GroupID        string    `gorm:"column:group_id;primary_key;size:64"`
	UserID         string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:nickname;size:255"`
	FaceURL        string    `gorm:"column:user_group_face_url;size:255"`
	RoleLevel      int32     `gorm:"column:role_level"`
	JoinTime       time.Time `gorm:"column:join_time"`
	JoinSource     int32     `gorm:"column:join_source"`
	InviterUserID  string    `gorm:"column:inviter_user_id;size:64"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	MuteEndTime    time.Time `gorm:"column:mute_end_time"`
	Ex             string    `gorm:"column:ex;size:1024"`
	DB             *gorm.DB
}

func (g *GroupMember) Create(ctx context.Context, groupMemberList []*GroupMember) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupMemberList", groupMemberList)
	}()
	return utils.Wrap(GroupMemberDB.Create(&groupMemberList).Error, "")
}

func (g *GroupMember) Delete(ctx context.Context, groupMembers []*GroupMember) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupMembers", groupMembers)
	}()
	return utils.Wrap(GroupMemberDB.Delete(groupMembers).Error, "")
}

func (g *GroupMember) UpdateByMap(ctx context.Context, groupID string, userID string, args map[string]interface{}) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "args", args)
	}()
	return utils.Wrap(GroupMemberDB.Model(&GroupMember{}).Where("group_id = ? and user_id = ?", groupID, userID).Updates(args).Error, "")
}

func (g *GroupMember) Update(ctx context.Context, groupMembers []*GroupMember) (err error) {
	defer func() { trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupMembers", groupMembers) }()
	return utils.Wrap(GroupMemberDB.Updates(&groupMembers).Error, "")
}

func (g *GroupMember) Find(ctx context.Context, groupMembers []*GroupMember) (groupList []*GroupMember, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupMembers", groupMembers, "groupList", groupList)
	}()
	var where [][]interface{}
	for _, groupMember := range groupMembers {
		where = append(where, []interface{}{groupMember.GroupID, groupMember.UserID})
	}
	err = utils.Wrap(GroupMemberDB.Where("(group_id, user_id) in ?", where).Find(&groupList).Error, "")
	return groupList, err
}

func (g *GroupMember) Take(ctx context.Context, groupID string, userID string) (groupMember *GroupMember, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "groupMember", *groupMember)
	}()
	groupMember = &GroupMember{}
	return groupMember, utils.Wrap(GroupMemberDB.Where("group_id = ? and user_id = ?", groupID, userID).Take(groupMember).Error, "")
}

func (g *GroupMember) TakeOwnerInfo(ctx context.Context, groupID string) (groupMember *GroupMember, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "groupMember", *groupMember)
	}()
	groupMember = &GroupMember{}
	err = GroupMemberDB.Where("group_id = ? and role_level = ?", groupID, constant.GroupOwner).Take(groupMember).Error
	return groupMember, utils.Wrap(err, "")
}

func InsertIntoGroupMember(toInsertInfo GroupMember) error {
	toInsertInfo.JoinTime = time.Now()
	if toInsertInfo.RoleLevel == 0 {
		toInsertInfo.RoleLevel = constant.GroupOrdinaryUsers
	}
	toInsertInfo.MuteEndTime = time.Unix(int64(time.Now().Second()), 0)
	err := GroupMemberDB.Table("group_members").Create(toInsertInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func BatchInsertIntoGroupMember(toInsertInfoList []*GroupMember) error {
	for _, toInsertInfo := range toInsertInfoList {
		toInsertInfo.JoinTime = time.Now()
		if toInsertInfo.RoleLevel == 0 {
			toInsertInfo.RoleLevel = constant.GroupOrdinaryUsers
		}
		toInsertInfo.MuteEndTime = time.Unix(int64(time.Now().Second()), 0)
	}
	return GroupMemberDB.Create(toInsertInfoList).Error

}

func GetGroupMemberListByUserID(userID string) ([]GroupMember, error) {
	var groupMemberList []GroupMember
	err := GroupMemberDB.Table("group_members").Where("user_id=?", userID).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberListByGroupID(groupID string) ([]GroupMember, error) {
	var groupMemberList []GroupMember
	err := GroupMemberDB.Table("group_members").Where("group_id=?", groupID).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberIDListByGroupID(groupID string) ([]string, error) {
	var groupMemberIDList []string
	err := GroupMemberDB.Table("group_members").Where("group_id=?", groupID).Pluck("user_id", &groupMemberIDList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberIDList, nil
}

func GetGroupMemberListByGroupIDAndRoleLevel(groupID string, roleLevel int32) ([]GroupMember, error) {
	var groupMemberList []GroupMember
	err := GroupMemberDB.Table("group_members").Where("group_id=? and role_level=?", groupID, roleLevel).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberInfoByGroupIDAndUserID(groupID, userID string) (*GroupMember, error) {
	var groupMember GroupMember
	err := GroupMemberDB.Table("group_members").Where("group_id=? and user_id=? ", groupID, userID).Limit(1).Take(&groupMember).Error
	if err != nil {
		return nil, err
	}
	return &groupMember, nil
}

func DeleteGroupMemberByGroupIDAndUserID(groupID, userID string) error {
	return GroupMemberDB.Table("group_members").Where("group_id=? and user_id=? ", groupID, userID).Delete(GroupMember{}).Error
}

func DeleteGroupMemberByGroupID(groupID string) error {
	return GroupMemberDB.Table("group_members").Where("group_id=?  ", groupID).Delete(GroupMember{}).Error
}

func UpdateGroupMemberInfo(groupMemberInfo GroupMember) error {
	return GroupMemberDB.Table("group_members").Where("group_id=? and user_id=?", groupMemberInfo.GroupID, groupMemberInfo.UserID).Updates(&groupMemberInfo).Error
}

func UpdateGroupMemberInfoByMap(groupMemberInfo GroupMember, m map[string]interface{}) error {
	return GroupMemberDB.Table("group_members").Where("group_id=? and user_id=?", groupMemberInfo.GroupID, groupMemberInfo.UserID).Updates(m).Error
}

func GetOwnerManagerByGroupID(groupID string) ([]GroupMember, error) {
	var groupMemberList []GroupMember
	err := GroupMemberDB.Table("group_members").Where("group_id=? and role_level>?", groupID, constant.GroupOrdinaryUsers).Find(&groupMemberList).Error
	if err != nil {
		return nil, err
	}
	return groupMemberList, nil
}

func GetGroupMemberNumByGroupID(groupID string) (int64, error) {
	var number int64
	err := GroupMemberDB.Table("group_members").Where("group_id=?", groupID).Count(&number).Error
	if err != nil {
		return 0, utils.Wrap(err, "")
	}
	return number, nil
}

func GetGroupOwnerInfoByGroupID(groupID string) (*GroupMember, error) {
	omList, err := GetOwnerManagerByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	for _, v := range omList {
		if v.RoleLevel == constant.GroupOwner {
			return &v, nil
		}
	}
	return nil, utils.Wrap(constant.ErrGroupNoOwner, "")
}

func IsExistGroupMember(groupID, userID string) bool {
	var number int64
	err := GroupMemberDB.Table("group_members").Where("group_id = ? and user_id = ?", groupID, userID).Count(&number).Error
	if err != nil {
		return false
	}
	if number != 1 {
		return false
	}
	return true
}

func CheckIsExistGroupMember(ctx context.Context, groupID, userID string) error {
	var number int64
	err := GroupMemberDB.Table("group_members").Where("group_id = ? and user_id = ?", groupID, userID).Count(&number).Error
	if err != nil {
		return constant.ErrDB.Wrap()
	}
	if number != 1 {
		return constant.ErrData.Wrap()
	}
	return nil
}

func GetGroupMemberByGroupID(groupID string, filter int32, begin int32, maxNumber int32) ([]GroupMember, error) {
	var memberList []GroupMember
	var err error
	if filter >= 0 {
		memberList, err = GetGroupMemberListByGroupIDAndRoleLevel(groupID, filter) //sorted by join time
	} else {
		memberList, err = GetGroupMemberListByGroupID(groupID)
	}

	if err != nil {
		return nil, err
	}
	if begin >= int32(len(memberList)) {
		return nil, nil
	}

	var end int32
	if begin+int32(maxNumber) < int32(len(memberList)) {
		end = begin + maxNumber
	} else {
		end = int32(len(memberList))
	}
	return memberList[begin:end], nil
}

func GetJoinedGroupIDListByUserID(userID string) ([]string, error) {
	memberList, err := GetGroupMemberListByUserID(userID)
	if err != nil {
		return nil, err
	}
	var groupIDList []string
	for _, v := range memberList {
		groupIDList = append(groupIDList, v.GroupID)
	}
	return groupIDList, nil
}

func IsGroupOwnerAdmin(groupID, UserID string) bool {
	groupMemberList, err := GetOwnerManagerByGroupID(groupID)
	if err != nil {
		return false
	}
	for _, v := range groupMemberList {
		if v.UserID == UserID && v.RoleLevel > constant.GroupOrdinaryUsers {
			return true
		}
	}
	return false
}

func GetGroupMembersByGroupIdCMS(groupId string, userName string, showNumber, pageNumber int32) ([]GroupMember, error) {
	var groupMembers []GroupMember
	err := GroupMemberDB.Table("group_members").Where("group_id=?", groupId).Where(fmt.Sprintf(" nickname like '%%%s%%' ", userName)).Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1))).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}
	return groupMembers, nil
}

func GetGroupMembersCount(groupID, userName string) (int64, error) {
	var count int64
	if err := GroupMemberDB.Table("group_members").Where("group_id=?", groupID).Where(fmt.Sprintf(" nickname like '%%%s%%' ", userName)).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

func UpdateGroupMemberInfoDefaultZero(groupMemberInfo GroupMember, args map[string]interface{}) error {
	return GroupMemberDB.Model(groupMemberInfo).Updates(args).Error
}

package rocksCache

import (
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	userInfoCache            = "USER_INFO_CACHE:"
	friendRelationCache      = "FRIEND_RELATION_CACHE:"
	blackListCache           = "BLACK_LIST_CACHE:"
	groupCache               = "GROUP_CACHE:"
	groupInfoCache           = "GROUP_INFO_CACHE:"
	groupOwnerIDCache        = "GROUP_OWNER_ID:"
	joinedGroupListCache     = "JOINED_GROUP_LIST_CACHE:"
	groupMemberInfoCache     = "GROUP_MEMBER_INFO_CACHE:"
	groupAllMemberInfoCache  = "GROUP_ALL_MEMBER_INFO_CACHE:"
	allFriendInfoCache       = "ALL_FRIEND_INFO_CACHE:"
	allDepartmentCache       = "ALL_DEPARTMENT_CACHE:"
	allDepartmentMemberCache = "ALL_DEPARTMENT_MEMBER_CACHE:"
)

func init() {
	fmt.Println("init to del old keys")
	for _, key := range []string{groupCache, friendRelationCache, blackListCache, userInfoCache, groupInfoCache, groupOwnerIDCache, joinedGroupListCache,
		groupMemberInfoCache, groupAllMemberInfoCache, allFriendInfoCache} {
		var cursor uint64
		var n int
		for {
			var keys []string
			var err error
			keys, cursor, err = db.DB.RDB.Scan(context.Background(), cursor, key+"*", 3000).Result()
			if err != nil {
				panic(err.Error())
			}
			n += len(keys)
			//fmt.Printf("\n %s key found %d keys: %v, current cursor %d\n", key, n, keys, cursor)
			if len(keys) > 0 {
				err = db.DB.RDB.Del(context.Background(), keys...).Err()
				if err != nil {
					panic(err.Error())
				}
			}
			if cursor == 0 {
				break
			}
		}
	}
}

func GetFriendIDListFromCache(userID string) ([]string, error) {
	getFriendIDList := func() (string, error) {
		friendIDList, err := imdb.GetFriendIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(friendIDList)
		return string(bytes), utils.Wrap(err, "")
	}
	friendIDListStr, err := db.DB.Rc.Fetch(friendRelationCache+userID, time.Second*30*60, getFriendIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var friendIDList []string
	err = json.Unmarshal([]byte(friendIDListStr), &friendIDList)
	return friendIDList, utils.Wrap(err, "")
}

func DelFriendIDListFromCache(userID string) error {
	err := db.DB.Rc.TagAsDeleted(friendRelationCache + userID)
	return err
}

func GetBlackListFromCache(userID string) ([]string, error) {
	getBlackIDList := func() (string, error) {
		blackIDList, err := imdb.GetBlackIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(blackIDList)
		return string(bytes), utils.Wrap(err, "")
	}
	blackIDListStr, err := db.DB.Rc.Fetch(blackListCache+userID, time.Second*30*60, getBlackIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var blackIDList []string
	err = json.Unmarshal([]byte(blackIDListStr), &blackIDList)
	return blackIDList, utils.Wrap(err, "")
}

func DelBlackIDListFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(blackListCache + userID)
}

func GetJoinedGroupIDListFromCache(userID string) ([]string, error) {
	getJoinedGroupIDList := func() (string, error) {
		joinedGroupList, err := imdb.GetJoinedGroupIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(joinedGroupList)
		return string(bytes), utils.Wrap(err, "")
	}
	joinedGroupIDListStr, err := db.DB.Rc.Fetch(joinedGroupListCache+userID, time.Second*30*60, getJoinedGroupIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var joinedGroupList []string
	err = json.Unmarshal([]byte(joinedGroupIDListStr), &joinedGroupList)
	return joinedGroupList, utils.Wrap(err, "")
}

func DelJoinedGroupIDListFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(joinedGroupListCache + userID)
}

func GetGroupMemberIDListFromCache(groupID string) ([]string, error) {
	getGroupMemberIDList := func() (string, error) {
		groupMemberIDList, err := imdb.GetGroupMemberIDListByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupMemberIDList)
		return string(bytes), utils.Wrap(err, "")
	}
	groupIDListStr, err := db.DB.WeakRc.Fetch(groupCache+groupID, time.Second*30*60, getGroupMemberIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupMemberIDList []string
	err = json.Unmarshal([]byte(groupIDListStr), &groupMemberIDList)
	return groupMemberIDList, utils.Wrap(err, "")
}

func DelGroupMemberIDListFromCache(userID string) error {
	err := db.DB.WeakRc.TagAsDeleted(groupCache + userID)
	return err
}

func GetUserInfoFromCache(userID string) (*db.User, error) {
	getUserInfo := func() (string, error) {
		userInfo, err := imdb.GetUserByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(userInfo)
		return string(bytes), utils.Wrap(err, "")
	}
	userInfoStr, err := db.DB.Rc.Fetch(userInfoCache+userID, time.Second*30*60, getUserInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	userInfo := &db.User{}
	err = json.Unmarshal([]byte(userInfoStr), userInfo)
	return userInfo, utils.Wrap(err, "")
}

func DelUserInfoFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(userInfoCache + userID)
}

func GetGroupMemberInfoFromCache(groupID, userID string) (*db.GroupMember, error) {
	getGroupMemberInfo := func() (string, error) {
		groupMemberInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupMemberInfo)
		return string(bytes), utils.Wrap(err, "")
	}
	groupMemberInfoStr, err := db.DB.Rc.Fetch(groupMemberInfoCache+groupID+"-"+userID, time.Second*30*60, getGroupMemberInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	groupMember := &db.GroupMember{}
	err = json.Unmarshal([]byte(groupMemberInfoStr), groupMember)
	return groupMember, utils.Wrap(err, "")
}

func DelGroupMemberInfoFromCache(groupID, userID string) error {
	return db.DB.Rc.TagAsDeleted(groupMemberInfoCache + groupID + "-" + userID)
}

func GetAllGroupMembersInfoFromCache(groupID string) ([]*db.GroupMember, error) {
	getGroupMemberInfo := func() (string, error) {
		groupMembers, err := imdb.GetGroupMemberListByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupMembers)
		return string(bytes), utils.Wrap(err, "")
	}
	groupMembersStr, err := db.DB.Rc.Fetch(groupAllMemberInfoCache+groupID, time.Second*30*60, getGroupMemberInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupMembers []*db.GroupMember
	err = json.Unmarshal([]byte(groupMembersStr), &groupMembers)
	return groupMembers, utils.Wrap(err, "")
}

func DelAllGroupMembersInfoFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupAllMemberInfoCache + groupID)
}

func GetGroupInfoFromCache(groupID string) (*db.Group, error) {
	getGroupInfo := func() (string, error) {
		groupInfo, err := imdb.GetGroupInfoByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupInfo)
		return string(bytes), utils.Wrap(err, "")
	}
	groupInfoStr, err := db.DB.Rc.Fetch(groupInfoCache+groupID, time.Second*30*60, getGroupInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	groupInfo := &db.Group{}
	err = json.Unmarshal([]byte(groupInfoStr), groupInfo)
	return groupInfo, utils.Wrap(err, "")
}

func DelGroupInfoFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupInfoCache + groupID)
}

func GetAllFriendsInfoFromCache(userID string) ([]*db.Friend, error) {
	getAllFriendInfo := func() (string, error) {
		friendInfoList, err := imdb.GetFriendListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(friendInfoList)
		return string(bytes), utils.Wrap(err, "")
	}
	allFriendInfoStr, err := db.DB.Rc.Fetch(allFriendInfoCache+userID, time.Second*30*60, getAllFriendInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var friendInfoList []*db.Friend
	err = json.Unmarshal([]byte(allFriendInfoStr), &friendInfoList)
	return friendInfoList, utils.Wrap(err, "")
}

func DelAllFriendsInfoFromCache(userID string) error {
	return db.DB.Rc.TagAsDeleted(allFriendInfoCache + userID)
}

func GetAllDepartmentsFromCache() ([]db.Department, error) {
	getAllDepartments := func() (string, error) {
		departmentList, err := imdb.GetSubDepartmentList("-1")
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(departmentList)
		return string(bytes), utils.Wrap(err, "")
	}
	allDepartmentsStr, err := db.DB.Rc.Fetch(allDepartmentCache, time.Second*30*60, getAllDepartments)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var allDepartments []db.Department
	err = json.Unmarshal([]byte(allDepartmentsStr), &allDepartments)
	return allDepartments, utils.Wrap(err, "")
}

func DelAllDepartmentsFromCache() error {
	return db.DB.Rc.TagAsDeleted(allDepartmentCache)
}

func GetAllDepartmentMembersFromCache() ([]db.DepartmentMember, error) {
	getAllDepartmentMembers := func() (string, error) {
		departmentMembers, err := imdb.GetDepartmentMemberList("-1")
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(departmentMembers)
		return string(bytes), utils.Wrap(err, "")
	}
	allDepartmentMembersStr, err := db.DB.Rc.Fetch(allDepartmentMemberCache, time.Second*30*60, getAllDepartmentMembers)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var allDepartmentMembers []db.DepartmentMember
	err = json.Unmarshal([]byte(allDepartmentMembersStr), &allDepartmentMembers)
	return allDepartmentMembers, utils.Wrap(err, "")
}

func DelAllDepartmentMembersFromCache() error {
	return db.DB.Rc.TagAsDeleted(allDepartmentMemberCache)
}

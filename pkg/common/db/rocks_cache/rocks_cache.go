package rocksCache

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"time"
)

const (
	userInfoCache             = "USER_INFO_CACHE:"
	friendRelationCache       = "FRIEND_RELATION_CACHE:"
	blackListCache            = "BLACK_LIST_CACHE:"
	groupCache                = "GROUP_CACHE:"
	groupInfoCache            = "GROUP_INFO_CACHE:"
	groupOwnerIDCache         = "GROUP_OWNER_ID:"
	joinedGroupListCache      = "JOINED_GROUP_LIST_CACHE:"
	groupMemberInfoCache      = "GROUP_MEMBER_INFO_CACHE:"
	groupAllMemberInfoCache   = "GROUP_ALL_MEMBER_INFO_CACHE:"
	allFriendInfoCache        = "ALL_FRIEND_INFO_CACHE:"
	allDepartmentCache        = "ALL_DEPARTMENT_CACHE:"
	allDepartmentMemberCache  = "ALL_DEPARTMENT_MEMBER_CACHE:"
	joinedSuperGroupListCache = "JOINED_SUPER_GROUP_LIST_CACHE:"
	groupMemberListHashCache  = "GROUP_MEMBER_LIST_HASH_CACHE:"
	groupMemberNumCache       = "GROUP_MEMBER_NUM_CACHE:"
	conversationCache         = "CONVERSATION_CACHE:"
	conversationIDListCache   = "CONVERSATION_ID_LIST_CACHE:"
)

func init() {
	fmt.Println("init to del old keys")
	for _, key := range []string{groupCache, friendRelationCache, blackListCache, userInfoCache, groupInfoCache, groupOwnerIDCache, joinedGroupListCache,
		groupMemberInfoCache, groupAllMemberInfoCache, allFriendInfoCache} {
		fName := utils.GetSelfFuncName()
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
			// for each for redis cluster
			for _, key := range keys {
				if err = db.DB.RDB.Del(context.Background(), key).Err(); err != nil {
					log.NewError("", fName, key, err.Error())
					err = db.DB.RDB.Del(context.Background(), key).Err()
					if err != nil {
						panic(err.Error())
					}
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
	f := func() (string, error) {
		groupInfo, err := GetGroupInfoFromCache(groupID)
		if err != nil {
			return "", utils.Wrap(err, "GetGroupInfoFromCache failed")
		}
		var groupMemberIDList []string
		if groupInfo.GroupType == constant.SuperGroup {
			superGroup, err := db.DB.GetSuperGroup(groupID)
			if err != nil {
				return "", utils.Wrap(err, "")
			}
			groupMemberIDList = superGroup.MemberIDList
		} else {
			groupMemberIDList, err = imdb.GetGroupMemberIDListByGroupID(groupID)
			if err != nil {
				return "", utils.Wrap(err, "")
			}
		}
		bytes, err := json.Marshal(groupMemberIDList)
		return string(bytes), utils.Wrap(err, "")
	}
	groupIDListStr, err := db.DB.Rc.Fetch(groupCache+groupID, time.Second*30*60, f)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	var groupMemberIDList []string
	err = json.Unmarshal([]byte(groupIDListStr), &groupMemberIDList)
	return groupMemberIDList, utils.Wrap(err, "")
}

func DelGroupMemberIDListFromCache(groupID string) error {
	err := db.DB.Rc.TagAsDeleted(groupCache + groupID)
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

func GetGroupMembersInfoFromCache(count, offset int32, groupID string) ([]*db.GroupMember, error) {
	groupMemberIDList, err := GetGroupMemberIDListFromCache(groupID)
	if err != nil {
		return nil, err
	}
	if count < 0 || offset < 0 {
		return nil, nil
	}
	var groupMemberList []*db.GroupMember
	var start, stop int32
	start = offset
	stop = offset + count
	l := int32(len(groupMemberIDList))
	if start > stop {
		return nil, nil
	}
	if start >= l {
		return nil, nil
	}
	if count != 0 {
		if stop >= l {
			stop = l
		}
		groupMemberIDList = groupMemberIDList[start:stop]
	} else {
		if l < 1000 {
			stop = l
		} else {
			stop = 1000
		}
		groupMemberIDList = groupMemberIDList[start:stop]
	}
	//log.NewDebug("", utils.GetSelfFuncName(), "ID list: ", groupMemberIDList)
	for _, userID := range groupMemberIDList {
		groupMembers, err := GetGroupMemberInfoFromCache(groupID, userID)
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), err.Error(), groupID, userID)
			continue
		}
		groupMemberList = append(groupMemberList, groupMembers)
	}
	return groupMemberList, nil
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

func GetJoinedSuperGroupListFromCache(userID string) ([]string, error) {
	getJoinedSuperGroupIDList := func() (string, error) {
		userToSuperGroup, err := db.DB.GetSuperGroupByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		if len(userToSuperGroup.GroupIDList) == 0 {
			return "", errors.New("GroupIDList == 0")
		}
		bytes, err := json.Marshal(userToSuperGroup.GroupIDList)
		return string(bytes), utils.Wrap(err, "")
	}
	joinedSuperGroupListStr, err := db.DB.Rc.Fetch(joinedSuperGroupListCache+userID, time.Second*30*60, getJoinedSuperGroupIDList)
	var joinedSuperGroupList []string
	err = json.Unmarshal([]byte(joinedSuperGroupListStr), &joinedSuperGroupList)
	return joinedSuperGroupList, err
}

func DelJoinedSuperGroupIDListFromCache(userID string) error {
	err := db.DB.Rc.TagAsDeleted(joinedSuperGroupListCache + userID)
	return err
}

func GetGroupMemberListHashFromCache(groupID string) (uint64, error) {
	generateHash := func() (string, error) {
		groupMemberIDList, err := GetGroupMemberIDListFromCache(groupID)
		if err != nil {
			return "", utils.Wrap(err, "GetGroupMemberIDListFromCache failed")
		}
		sort.Strings(groupMemberIDList)
		var all string
		for _, v := range groupMemberIDList {
			all += v
		}
		bi := big.NewInt(0)
		bi.SetString(utils.Md5(all)[0:8], 16)
		return strconv.Itoa(int(bi.Uint64())), nil
	}
	hashCode, err := db.DB.Rc.Fetch(groupMemberListHashCache+groupID, time.Second*30*60, generateHash)
	hashCodeUint64, err := strconv.Atoi(hashCode)
	return uint64(hashCodeUint64), err
}

func DelGroupMemberListHashFromCache(groupID string) error {
	err := db.DB.Rc.TagAsDeleted(groupMemberListHashCache + groupID)
	return err
}

func GetGroupMemberNumFromCache(groupID string) (int64, error) {
	getGroupMemberNum := func() (string, error) {
		num, err := imdb.GetGroupMemberNumByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return strconv.Itoa(int(num)), nil
	}
	groupMember, err := db.DB.Rc.Fetch(groupMemberNumCache+groupID, time.Second*30*60, getGroupMemberNum)
	num, err := strconv.Atoi(groupMember)
	return int64(num), err
}

func DelGroupMemberNumFromCache(groupID string) error {
	return db.DB.Rc.TagAsDeleted(groupMemberNumCache + groupID)
}

func GetUserConversationIDListFromCache(userID string) ([]string, error) {
	getConversationIDList := func() (string, error) {
		conversationIDList, err := imdb.GetConversationIDListByUserID(userID)
		if err != nil {
			return "", utils.Wrap(err, "getConversationIDList failed")
		}
		log.NewDebug("", utils.GetSelfFuncName(), conversationIDList)
		bytes, err := json.Marshal(conversationIDList)
		return string(bytes), utils.Wrap(err, "")
	}
	conversationIDListStr, err := db.DB.Rc.Fetch(conversationIDListCache+userID, time.Second*30*60, getConversationIDList)
	var conversationIDList []string
	err = json.Unmarshal([]byte(conversationIDListStr), &conversationIDList)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	return conversationIDList, nil
}

func DelUserConversationIDListFromCache(userID string) error {
	return utils.Wrap(db.DB.Rc.TagAsDeleted(conversationIDListCache+userID), "DelUserConversationIDListFromCache err")
}

func GetConversationFromCache(ownerUserID, conversationID string) (*db.Conversation, error) {
	getConversation := func() (string, error) {
		conversation, err := imdb.GetConversation(ownerUserID, conversationID)
		if err != nil {
			return "", utils.Wrap(err, "get failed")
		}
		bytes, err := json.Marshal(conversation)
		if err != nil {
			return "", utils.Wrap(err, "Marshal failed")
		}
		return string(bytes), nil
	}
	conversationStr, err := db.DB.Rc.Fetch(conversationCache+ownerUserID+":"+conversationID, time.Second*30*60, getConversation)
	if err != nil {
		return nil, utils.Wrap(err, "Fetch failed")
	}
	conversation := db.Conversation{}
	err = json.Unmarshal([]byte(conversationStr), &conversation)
	if err != nil {
		return nil, utils.Wrap(err, "Unmarshal failed")
	}
	return &conversation, nil
}

func GetConversationsFromCache(ownerUserID string, conversationIDList []string) ([]db.Conversation, error) {
	var conversationList []db.Conversation
	for _, conversationID := range conversationIDList {
		conversation, err := GetConversationFromCache(ownerUserID, conversationID)
		if err != nil {
			return nil, utils.Wrap(err, "GetConversationFromCache failed")
		}
		conversationList = append(conversationList, *conversation)
	}
	return conversationList, nil
}

func GetUserAllConversationList(ownerUserID string) ([]db.Conversation, error) {
	IDList, err := GetUserConversationIDListFromCache(ownerUserID)
	if err != nil {
		return nil, err
	}
	var conversationList []db.Conversation
	log.NewDebug("", utils.GetSelfFuncName(), IDList)
	for _, conversationID := range IDList {
		conversation, err := GetConversationFromCache(ownerUserID, conversationID)
		if err != nil {
			return nil, utils.Wrap(err, "GetConversationFromCache failed")
		}
		conversationList = append(conversationList, *conversation)
	}
	return conversationList, nil
}

func DelConversationFromCache(ownerUserID, conversationID string) error {
	return utils.Wrap(db.DB.Rc.TagAsDeleted(conversationCache+ownerUserID+":"+conversationID), "DelConversationFromCache err")
}

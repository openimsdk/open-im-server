package table

import (
	"strconv"
	"strings"
)

const (
	CSuperGroup       = "super_group"
	CUserToSuperGroup = "user_to_super_group"
	CTag              = "tag"
	CSendLog          = "send_log"
	CWorkMoment       = "work_moment"
	CExtendMsgSet     = "extend_msgs"

	ExtendMsgMaxNum = 100
)

type SuperGroupModel struct {
	GroupID      string   `bson:"group_id" json:"groupID"`
	MemberIDList []string `bson:"member_id_list" json:"memberIDList"`
}

func (SuperGroupModel) TableName() string {
	return CSuperGroup
}

type UserToSuperGroupModel struct {
	UserID      string   `bson:"user_id" json:"userID"`
	GroupIDList []string `bson:"group_id_list" json:"groupIDList"`
}

func (UserToSuperGroupModel) TableName() string {
	return CUserToSuperGroup
}

type TagModel struct {
	UserID   string   `bson:"user_id"`
	TagID    string   `bson:"tag_id"`
	TagName  string   `bson:"tag_name"`
	UserList []string `bson:"user_list"`
}

func (TagModel) TableName() string {
	return CTag
}

type CommonUserModel struct {
	UserID   string `bson:"user_id"`
	UserName string `bson:"user_name"`
}

type TagSendLogModel struct {
	UserList         []CommonUserModel `bson:"tag_list"`
	SendID           string            `bson:"send_id"`
	SenderPlatformID int32             `bson:"sender_platform_id"`
	Content          string            `bson:"content"`
	SendTime         int64             `bson:"send_time"`
}

func (TagSendLogModel) TableName() string {
	return CSendLog
}

type WorkMoment struct {
	WorkMomentID         string             `bson:"work_moment_id"`
	UserID               string             `bson:"user_id"`
	UserName             string             `bson:"user_name"`
	FaceURL              string             `bson:"face_url"`
	Content              string             `bson:"content"`
	LikeUserList         []*CommonUserModel `bson:"like_user_list"`
	AtUserList           []*CommonUserModel `bson:"at_user_list"`
	PermissionUserList   []*CommonUserModel `bson:"permission_user_list"`
	Comments             []*CommonUserModel `bson:"comments"`
	PermissionUserIDList []string           `bson:"permission_user_id_list"`
	Permission           int32              `bson:"permission"`
	CreateTime           int32              `bson:"create_time"`
}

func (WorkMoment) TableName() string {
	return CWorkMoment
}

type Comment struct {
	UserID        string `bson:"user_id" json:"user_id"`
	UserName      string `bson:"user_name" json:"user_name"`
	ReplyUserID   string `bson:"reply_user_id" json:"reply_user_id"`
	ReplyUserName string `bson:"reply_user_name" json:"reply_user_name"`
	ContentID     string `bson:"content_id" json:"content_id"`
	Content       string `bson:"content" json:"content"`
	CreateTime    int32  `bson:"create_time" json:"create_time"`
}

type ExtendMsgSet struct {
	SourceID         string               `bson:"source_id" json:"sourceID"`
	SessionType      int32                `bson:"session_type" json:"sessionType"`
	ExtendMsgs       map[string]ExtendMsg `bson:"extend_msgs" json:"extendMsgs"`
	ExtendMsgNum     int32                `bson:"extend_msg_num" json:"extendMsgNum"`
	CreateTime       int64                `bson:"create_time" json:"createTime"`               // this block's create time
	MaxMsgUpdateTime int64                `bson:"max_msg_update_time" json:"maxMsgUpdateTime"` // index find msg
}

type KeyValue struct {
	TypeKey          string `bson:"type_key" json:"typeKey"`
	Value            string `bson:"value" json:"value"`
	LatestUpdateTime int64  `bson:"latest_update_time" json:"latestUpdateTime"`
}

type ExtendMsg struct {
	ReactionExtensionList map[string]KeyValue `bson:"reaction_extension_list" json:"reactionExtensionList"`
	ClientMsgID           string              `bson:"client_msg_id" json:"clientMsgID"`
	MsgFirstModifyTime    int64               `bson:"msg_first_modify_time" json:"msgFirstModifyTime"` // this extendMsg create time
	AttachedInfo          string              `bson:"attached_info" json:"attachedInfo"`
	Ex                    string              `bson:"ex" json:"ex"`
}

func (ExtendMsgSet) TableName() string {
	return CExtendMsgSet
}

func (ExtendMsgSet) GetExtendMsgMaxNum() int32 {
	return ExtendMsgMaxNum
}

func (ExtendMsgSet) GetSourceID(ID string, index int32) string {
	return ID + ":" + strconv.Itoa(int(index))
}

func (e *ExtendMsgSet) SplitSourceIDAndGetIndex() int32 {
	l := strings.Split(e.SourceID, ":")
	index, _ := strconv.Atoi(l[len(l)-1])
	return int32(index)
}

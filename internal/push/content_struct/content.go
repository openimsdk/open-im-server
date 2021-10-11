/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/27 11:24).
 */
package content_struct

import (
	"encoding/json"
)

type Content struct {
	IsDisplay int32  `json:"isDisplay"`
	ID        string `json:"id"`
	Text      string `json:"text"`
}

func NewContentStructString(isDisplay int32, ID string, text string) string {
	c := Content{IsDisplay: isDisplay, ID: ID, Text: text}
	return c.contentToString()
}

func (c *Content) contentToString() string {
	data, _ := json.Marshal(c)
	dataString := string(data)
	return dataString
}

type groupMemberFullInfo struct {
	GroupId  string `json:"groupID"`
	UserId   string `json:"userId"`
	Role     int    `json:"role"`
	JoinTime uint64 `json:"joinTime"`
	NickName string `json:"nickName"`
	FaceUrl  string `json:"faceUrl"`
}

type AgreeOrRejectGroupMember struct {
	GroupId  string `json:"groupID"`
	UserId   string `json:"userId"`
	Role     int    `json:"role"`
	JoinTime uint64 `json:"joinTime"`
	NickName string `json:"nickName"`
	FaceUrl  string `json:"faceUrl"`
	Reason   string `json:"reason"`
}
type AtTextContent struct {
	Text       string   `json:"text"`
	AtUserList []string `json:"atUserList"`
	IsAtSelf   bool     `json:"isAtSelf"`
}

type CreateGroupSysMsg struct {
	uIdCreator     string                `creatorUid`
	initMemberList []groupMemberFullInfo `json: initMemberList`
	CreateTime     uint64                `json:"CreateTime"`
	Text           string                `json:"text"`
}

type NotificationContent struct {
	IsDisplay   int32  `json:"isDisplay"`
	DefaultTips string `json:"defaultTips"`
	Detail      string `json:"detail"`
}

func (c *NotificationContent) ContentToString() string {
	data, _ := json.Marshal(c)
	dataString := string(data)
	return dataString
}

type KickGroupMemberApiReq struct {
	GroupID     string   `json:"groupID"`
	UidList     []string `json:"uidList"`
	Reason      string   `json:"reason"`
	OperationID string   `json:"operationID"`
}

func NewCreateGroupSysMsgString(create *CreateGroupSysMsg, text string) string {
	create.Text = text
	jstring, _ := json.Marshal(create)

	return string(jstring)
}

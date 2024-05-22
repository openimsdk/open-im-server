package listdemo

import (
	"time"
)

var (
	_ Elem    = (*FriendElem)(nil)
	_ ListDoc = (*Friend)(nil)
)

type FriendElem struct {
	FriendUserID   string     `bson:"friend_user_id"`
	Nickname       string     `bson:"nickname"`
	FaceURL        string     `bson:"face_url"`
	Remark         string     `bson:"remark"`
	CreateTime     time.Time  `bson:"create_time"`
	AddSource      int32      `bson:"add_source"`
	OperatorUserID string     `bson:"operator_user_id"`
	Ex             string     `bson:"ex"`
	IsPinned       bool       `bson:"is_pinned"`
	Version        uint       `bson:"version"`
	DeleteTime     *time.Time `bson:"delete_time"`
}

func (f *FriendElem) IDName() string {
	return "friend_user_id"
}

func (f *FriendElem) IDValue() any {
	return f.FriendUserID
}

func (f *FriendElem) VersionName() string {
	return "version"
}

func (f *FriendElem) DeletedName() string {
	return "delete_time"
}

func (f *FriendElem) ToMap() map[string]any {
	return map[string]any{
		"friend_user_id":   f.FriendUserID,
		"nickname":         f.Nickname,
		"face_url":         f.FaceURL,
		"remark":           f.Remark,
		"create_time":      f.CreateTime,
		"add_source":       f.AddSource,
		"operator_user_id": f.OperatorUserID,
		"ex":               f.Ex,
		"is_pinned":        f.IsPinned,
		"version":          f.Version,
		"delete_time":      f.DeleteTime,
	}
}

type Friend struct {
	UserID        string        `bson:"user_id"`
	Friends       []*FriendElem `bson:"friends"`
	Version       uint          `bson:"version"`
	DeleteVersion uint          `bson:"delete_version"`
}

func (f *Friend) BuildDoc(lid any, e Elem) any {
	return &Friend{
		UserID:  lid.(string),
		Friends: []*FriendElem{e.(*FriendElem)},
	}
}

func (f *Friend) ElemsID() string {
	return "user_id"
}

func (f *Friend) IDName() string {
	return "user_id"
}

func (f *Friend) ElemsName() string {
	return "friends"
}

func (f *Friend) VersionName() string {
	return "version"
}

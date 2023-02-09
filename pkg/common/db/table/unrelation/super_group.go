package unrelation

import (
	"context"
	"strconv"
)

const (
	CSuperGroup       = "super_group"
	CUserToSuperGroup = "user_to_super_group"
)

type SuperGroupModel struct {
	GroupID   string   `bson:"group_id" json:"groupID"`
	MemberIDs []string `bson:"member_id_list" json:"memberIDList"`
}

func (SuperGroupModel) TableName() string {
	return CSuperGroup
}

type UserToSuperGroupModel struct {
	UserID   string   `bson:"user_id" json:"userID"`
	GroupIDs []string `bson:"group_id_list" json:"groupIDList"`
}

func (UserToSuperGroupModel) TableName() string {
	return CUserToSuperGroup
}

type SuperGroupModelInterface interface {
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDs []string, tx ...any) error
	FindSuperGroup(ctx context.Context, groupIDs []string, tx ...any) (groups []*SuperGroupModel, err error)
	AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string, tx ...any) error
	RemoverUserFromSuperGroup(ctx context.Context, groupID string, userIDs []string, tx ...any) error
	GetSuperGroupByUserID(ctx context.Context, userID string, tx ...any) (*UserToSuperGroupModel, error)
	DeleteSuperGroup(ctx context.Context, groupID string, tx ...any) error
	RemoveGroupFromUser(ctx context.Context, groupID string, userIDs []string, tx ...any) error
}

func superGroupIndexGen(groupID string, seqSuffix uint32) string {
	return "super_group_" + groupID + ":" + strconv.FormatInt(int64(seqSuffix), 10)
}

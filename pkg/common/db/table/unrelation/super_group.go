package unrelation

import (
	"context"
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
	// tx is your transaction object
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDs []string, tx ...interface{}) error
	GetSuperGroup(ctx context.Context, groupID string) (SuperGroupModel, error)
	AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string, tx ...interface{}) error
	RemoverUserFromSuperGroup(ctx context.Context, groupID string, userIDs []string, tx ...interface{}) error
	GetSuperGroupByUserID(ctx context.Context, userID string) (*UserToSuperGroupModel, error)
	DeleteSuperGroup(ctx context.Context, groupID string, tx ...interface{}) error
}

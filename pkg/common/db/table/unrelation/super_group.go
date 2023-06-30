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
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDs []string) error
	TakeSuperGroup(ctx context.Context, groupID string) (group *SuperGroupModel, err error)
	FindSuperGroup(ctx context.Context, groupIDs []string) (groups []*SuperGroupModel, err error)
	AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string) error
	RemoverUserFromSuperGroup(ctx context.Context, groupID string, userIDs []string) error
	GetSuperGroupByUserID(ctx context.Context, userID string) (*UserToSuperGroupModel, error)
	DeleteSuperGroup(ctx context.Context, groupID string) error
	RemoveGroupFromUser(ctx context.Context, groupID string, userIDs []string) error
}

package relation

import (
	"context"
	"time"
)

const (
	GroupMemberModelTableName = "group_members"
)

type GroupMemberModel struct {
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
}

func (GroupMemberModel) TableName() string {
	return GroupMemberModelTableName
}

type GroupMemberModelInterface interface {
	Create(ctx context.Context, groupMemberList []*GroupMemberModel, tx ...any) (err error)
	Delete(ctx context.Context, groupID string, userIDs []string, tx ...any) (err error)
	DeleteGroup(ctx context.Context, groupIDs []string, tx ...any) (err error)
	//UpdateByMap(ctx context.Context, groupID string, userID string, args map[string]interface{}, tx ...any) (err error)
	Update(ctx context.Context, groupID string, userID string, data map[string]any, tx ...any) (err error)
	UpdateRoleLevel(ctx context.Context, groupID string, userID string, roleLevel int32, tx ...any) (rowsAffected int64, err error)
	Find(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, tx ...any) (groupList []*GroupMemberModel, err error)
	//FindGroupUser(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, tx ...any) (groupList []*GroupMemberModel, err error)
	Take(ctx context.Context, groupID string, userID string, tx ...any) (groupMember *GroupMemberModel, err error)
	TakeOwner(ctx context.Context, groupID string, tx ...any) (groupMember *GroupMemberModel, err error)
	SearchMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32, tx ...any) (total int32, groupList []*GroupMemberModel, err error)
	MapGroupMemberNum(ctx context.Context, groupIDs []string, tx ...any) (count map[string]uint32, err error)
	FindJoinUserID(ctx context.Context, groupIDs []string, tx ...any) (groupUsers map[string][]string, err error)
}

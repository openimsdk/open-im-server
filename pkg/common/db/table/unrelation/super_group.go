package unrelation

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
}

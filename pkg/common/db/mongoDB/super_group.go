package mongoDB

type SuperGroup struct {
	GroupID      string   `bson:"group_id" json:"groupID"`
	MemberIDList []string `bson:"member_id_list" json:"memberIDList"`
}

type UserToSuperGroup struct {
	UserID      string   `bson:"user_id" json:"userID"`
	GroupIDList []string `bson:"group_id_list" json:"groupIDList"`
}

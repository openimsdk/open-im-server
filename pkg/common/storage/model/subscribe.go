package model

// SubscribeUserTableName collection constant.
const (
	SubscribeUserTableName = "subscribe_user"
)

// SubscribeUser collection structure.
type SubscribeUser struct {
	UserID     string   `bson:"user_id"      json:"userID"`
	UserIDList []string `bson:"user_id_list" json:"userIDList"`
}

func (SubscribeUser) TableName() string {
	return SubscribeUserTableName
}

package db

type ExtendMsgSet struct {
	ID               string       `bson:"id" json:"ID"`
	ExtendMsg        []*ExtendMsg `bson:"extend_msg" json:"extendMsg"`
	LatestUpdateTime int32        `bson:"latest_update_time" json:"latestUpdateTime"`
	AttachedInfo     string       `bson:"attached_info" json:"attachedInfo"`
	Ex               string       `bson:"ex" json:"ex"`
	ExtendMsgNum     int32        `bson:"extend_msg_num" json:"extendMsgNum"`
	CreateTime       int32        `bson:"create_time" json:"createTime"`
}

type ExtendMsg struct {
	SendID            string              `bson:"send_id" json:"sendID"`
	ServerMsgID       string              `bson:"server_msg_id" json:"serverMsgID"`
	Ex                string              `bson:"ex" json:"ex"`
	AttachedInfo      string              `bson:"attached_info" json:"attachedInfo"`
	LikeUserIDList    []string            `bson:"like_user_id_list" json:"likeUserIDList"`
	Content           string              `bson:"content" json:"content"`
	ExtendMsgComments []*ExtendMsgComment `bson:"extend_msg_comments" json:"extendMsgComment"`
	Vote              *Vote               `bson:"vote" json:"vote"`
	Urls              []string            `bson:"urls" json:"urls"`
	CreateTime        int32               `bson:"create_time" json:"createTime"`
}

type Vote struct {
	Content      string     `bson:"content" json:"content"`
	AttachedInfo string     `bson:"attached_info" json:"attachedInfo"`
	Ex           string     `bson:"ex" json:"ex"`
	Options      []*Options `bson:"options" json:"options"`
}

type Options struct {
	Content        string   `bson:"content" json:"content"`
	AttachedInfo   string   `bson:"attached_info" json:"attachedInfo"`
	Ex             string   `bson:"ex" json:"ex"`
	VoteUserIDList []string `bson:"vote_user_id_list" json:"voteUserIDList"`
}

type ExtendMsgComment struct {
	UserID         string `bson:"user_id" json:"userID"`
	ReplyUserID    string `bson:"reply_user_id" json:"replyUserID"`
	ReplyContentID string `bson:"reply_content_id" json:"replyContentID"`
	ContentID      string `bson:"content_id" json:"contentID"`
	Content        string `bson:"content" json:"content"`
	CreateTime     int32  `bson:"create_time" json:"createTime"`
	AttachedInfo   string `bson:"attached_info" json:"attachedInfo"`
	Ex             string `bson:"ex" json:"ex"`
}

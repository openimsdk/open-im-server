package model

type SeqUser struct {
	UserID         string `bson:"user_id"`
	ConversationID string `bson:"conversation_id"`
	MinSeq         int64  `bson:"min_seq"`
	MaxSeq         int64  `bson:"max_seq"`
	ReadSeq        int64  `bson:"read_seq"`
}

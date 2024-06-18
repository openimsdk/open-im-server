package model

type SeqConversation struct {
	ConversationID string `bson:"conversation_id"`
	MaxSeq         int64  `bson:"max_seq"`
	MinSeq         int64  `bson:"min_seq"`
}

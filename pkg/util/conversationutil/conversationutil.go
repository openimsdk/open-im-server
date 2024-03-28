package conversationutil

import (
	"sort"
	"strings"
)

func GenConversationIDForSingle(sendID, recvID string) string {
	l := []string{sendID, recvID}
	sort.Strings(l)
	return "si_" + strings.Join(l, "_")
}

func GenConversationUniqueKeyForGroup(groupID string) string {
	return groupID
}

func GenGroupConversationID(groupID string) string {
	return "sg_" + groupID
}

func GenConversationUniqueKeyForSingle(sendID, recvID string) string {
	l := []string{sendID, recvID}
	sort.Strings(l)
	return strings.Join(l, "_")
}

func GetNotificationConversationIDByConversationID(conversationID string) string {
	l := strings.Split(conversationID, "_")
	if len(l) > 1 {
		l[0] = "n"
		return strings.Join(l, "_")
	}
	return ""
}

func GetSelfNotificationConversationID(userID string) string {
	return "n_" + userID + "_" + userID
}

func GetSeqsBeginEnd(seqs []int64) (int64, int64) {
	if len(seqs) == 0 {
		return 0, 0
	}
	return seqs[0], seqs[len(seqs)-1]
}

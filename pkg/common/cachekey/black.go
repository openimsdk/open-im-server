package cachekey

const (
	BlackIDsKey = "BLACK_IDS:"
	IsBlackKey  = "IS_BLACK:" // local cache
)

func GetBlackIDsKey(ownerUserID string) string {
	return BlackIDsKey + ownerUserID

}

func GetIsBlackIDsKey(possibleBlackUserID, userID string) string {
	return IsBlackKey + userID + "-" + possibleBlackUserID
}

package cachekey

const (
	blackIDsKey = "BLACK_IDS:"
	isBlackKey  = "IS_BLACK:"
)

func GetBlackIDsKey(ownerUserID string) string {
	return blackIDsKey + ownerUserID

}

func GetIsBlackIDsKey(possibleBlackUserID, userID string) string {
	return isBlackKey + possibleBlackUserID + "-" + userID
}

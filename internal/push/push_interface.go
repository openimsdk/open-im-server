package push

type offlinePusher interface {
	auth(apiKey, secretKey string, timeStamp int64) (token string, err error)
	push(userIDList []string, alert, detailContent, platform string) (resp string, err error)
}

package getui

type Getui struct {
}

func (g *Getui) Push(userIDList []string, alert, detailContent, platform string) (resp string, err error) {
	return "", nil
}

func (g *Getui) Auth(apiKey, secretKey string, timeStamp int64) (token string, err error) {
	return "", nil
}

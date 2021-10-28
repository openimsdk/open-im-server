package common

import (
	"encoding/base64"
	"fmt"
)

const (
	PushUrl = "https://api.jpush.cn/v3/push"
)

func GetAuthorization(Appkey string, MasterSecret string) string {
	str := fmt.Sprintf("%s:%s", Appkey, MasterSecret)
	buf := []byte(str)
	Authorization := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(buf))
	return Authorization
}

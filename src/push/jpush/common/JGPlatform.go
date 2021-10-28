package common

import (
	"encoding/base64"
	"fmt"
)

func GetAuthorization(Appkey string, MasterSecret string) string {
	str := fmt.Sprintf("%s:%s", Appkey, MasterSecret)
	buf := []byte(str)
	Authorization := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(buf))
	return Authorization
}

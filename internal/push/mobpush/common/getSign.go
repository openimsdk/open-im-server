package common

import (
	"Open_IM/pkg/common/config"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

func GetSign(paramsStr string) string {
	h := md5.New()
	io.WriteString(h, paramsStr)
	io.WriteString(h, config.Config.Push.Mob.AppSecret)
	fmt.Printf("%x", h.Sum(nil))

	return hex.EncodeToString(h.Sum(nil))
}

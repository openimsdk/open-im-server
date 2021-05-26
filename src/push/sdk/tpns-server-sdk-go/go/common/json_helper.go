package common

import "encoding/json"

func ToJson(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}

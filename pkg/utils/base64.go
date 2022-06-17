package utils

import "encoding/base64"

func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func Base64Decode(data string) string {
	decodedByte, _ := base64.StdEncoding.DecodeString(data)
	return string(decodedByte)
}

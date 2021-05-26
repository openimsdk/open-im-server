package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Auther struct {
	AccessID  string
	SecretKey string
}

var UseSignAuthored = true

func (a *Auther) Auth(req *http.Request, useSignAuthored bool, auth Auther, reqBody string) {

	if useSignAuthored {
		now := time.Now()
		timeStamp := now.Unix()
		req.Header.Add("AccessId", auth.AccessID)
		req.Header.Add("TimeStamp", strconv.Itoa(int(timeStamp)))
		sign := GenSign(uint64(timeStamp), auth.AccessID, auth.SecretKey, reqBody)
		req.Header.Add("Sign", sign)
	} else {
		author := makeAuthHeader(a.AccessID, a.SecretKey)
		//log.Printf("author string:%v", author)
		req.Header.Add("Authorization", author)
	}
	//req.Header.Add("Content-Type", "application/json")
}

func makeAuthHeader(appID, secretKey string) string {
	base64Str := base64.StdEncoding.EncodeToString(
		[]byte(
			fmt.Sprintf("%s:%s", appID, secretKey),
		),
	)
	return fmt.Sprintf("Basic %s", base64Str)
}

func GenSign(timeStamp uint64, accessId string, secretKey, requestBody string) string {
	signBody := strconv.Itoa(int(timeStamp)) + accessId + requestBody
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secretKey))
	// Write Data to it
	h.Write([]byte(signBody))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))
	//fmt.Println()
	//fmt.Println("timeStamp: " + strconv.Itoa(int(timeStamp)) + " accessID:" + accessId + " body:" + requestBody)
	sEnc := b64.StdEncoding.EncodeToString([]byte(sha))
	//fmt.Println("final Result " + sEnc)
	return sEnc
}

package mw

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)

var (
	once  sync.Once
	block cipher.Block
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func initAesKey() {
	once.Do(func() {
		key := md5.Sum([]byte("openim:" + config.Config.TokenPolicy.AccessSecret))
		var err error
		block, err = aes.NewCipher(key[:])
		if err != nil {
			panic(err)
		}
	})
}

func genReqKey(args []string) string {
	initAesKey()
	plaintext := md5.Sum([]byte(strings.Join(args, ":")))
	iv := make([]byte, aes.BlockSize, aes.BlockSize+md5.Size)
	if _, err := rand.Read(iv); err != nil {
		panic(err)
	}
	ciphertext := make([]byte, md5.Size)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, plaintext[:])
	return base64.StdEncoding.EncodeToString(append(iv, ciphertext...))
}

func verifyReqKey(args []string, key string) error {
	initAesKey()
	k, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return fmt.Errorf("invalid key %v", err)
	}
	if len(k) != aes.BlockSize+md5.Size {
		return errors.New("invalid key")
	}
	plaintext := make([]byte, md5.Size)
	cipher.NewCBCDecrypter(block, k[:aes.BlockSize]).CryptBlocks(plaintext, k[aes.BlockSize:])
	sum := md5.Sum([]byte(strings.Join(args, ":")))
	if string(plaintext) != string(sum[:]) {
		return errors.New("mismatch key")
	}
	return nil
}

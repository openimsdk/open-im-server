package mw

import (
	"fmt"
	"testing"
)

func TestCheck(t *testing.T) {
	// config.Config.TokenPolicy.Secret = "123456"

	args := []string{"1", "2", "3"}

	key := genReqKey(args)
	fmt.Println("key:", key)
	err := verifyReqKey(args, key)

	fmt.Println(err)

	args = []string{"4", "5", "6"}

	key = genReqKey(args)
	fmt.Println("key:", key)
	err = verifyReqKey(args, key)

	fmt.Println(err)

}

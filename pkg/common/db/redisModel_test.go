package db

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SetTokenMapByUidPid(t *testing.T) {
	m := make(map[string]int, 0)
	m["哈哈"] = 1
	m["heihei"] = 2
	m["2332"] = 4
	_ = DB.SetTokenMapByUidPid("1234", 2, m)

}
func Test_GetTokenMapByUidPid(t *testing.T) {
	m, err := DB.GetTokenMapByUidPid("1234", "Android")
	assert.Nil(t, err)
	fmt.Println(m)
}

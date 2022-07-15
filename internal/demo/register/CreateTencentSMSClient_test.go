package register

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateTencentSMSClient(t *testing.T) {
	result, err := CreateTencentSMSClient()
	assert.Nil(t, err)
	fmt.Println("return result is ", result)

}

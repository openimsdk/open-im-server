package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_UserRegister(t *testing.T) {
	res := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(res)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"uidList": []}`))

	GetUserInfo(c)
	assert.Equal(t, 400, res.Code)

	res = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(res)
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"operationID": "1", "uidList": []}`))

	GetUserInfo(c)
	assert.Equal(t, 200, res.Code)
}

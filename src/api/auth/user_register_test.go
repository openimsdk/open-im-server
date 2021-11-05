package apiAuth

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
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(`{"secret": "111", "platform": 1, "uid": "1", "name": "1"}`))

	UserRegister(c)

	assert.Equal(t, res.Code, 200)
}

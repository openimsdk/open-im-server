package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func GinParseOperationID(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		operationID := c.Request.Header.Get("operationID")
		if operationID == "" {
			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.String(400, "read request body error: "+err.Error())
				c.Abort()
				return
			}
			req := struct {
				OperationID string `json:"operationID"`
			}{}
			if err := json.Unmarshal(body, &req); err != nil {
				c.String(400, "get operationID error: "+err.Error())
				c.Abort()
				return
			}
			if req.OperationID == "" {
				c.String(400, "operationID empty")
				c.Abort()
				return
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
			operationID = req.OperationID
			c.Request.Header.Set("operationID", operationID)
		}
		c.Set("operationID", operationID)
		c.Next()
		return
	}
	c.Next()
}

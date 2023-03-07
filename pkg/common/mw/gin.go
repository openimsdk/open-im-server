package mw

import (
	"OpenIM/pkg/common/constant"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func CorsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Methods", "*")
		context.Header("Access-Control-Allow-Headers", "*")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
		context.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
		context.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
		context.Header("content-type", "application/json")                                                                                                                                                           // 设置返回格式是json
		//Release all option pre-requests
		if context.Request.Method == http.MethodOptions {
			context.JSON(http.StatusOK, "Options Request!")
		}
		context.Next()
	}
}

func GinParseOperationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			operationID := c.Request.Header.Get(constant.OperationID)
			if operationID == "" {
				body, err := io.ReadAll(c.Request.Body)
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
				c.Request.Body = io.NopCloser(bytes.NewReader(body))
				operationID = req.OperationID
				c.Request.Header.Set(constant.OperationID, operationID)
			}
			c.Set(constant.OperationID, operationID)
			c.Next()
			return
		}
		c.Next()
	}
}

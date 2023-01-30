package middleware

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, userID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), "")
		// log.NewInfo("0", utils.GetSelfFuncName(), "userID: ", userID)
		c.Set("userID", userID)
		if !ok {
			log.NewError("", "GetUserIDFromToken false ", c.Request.Header.Get("token"))
			c.Abort()
			c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": errInfo})
			return
		} else {
			if !utils.IsContain(userID, config.Config.Manager.AppManagerUid) {
				c.Abort()
				c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": "user is not admin"})
				return
			}
			log.NewInfo("0", utils.GetSelfFuncName(), "failed: ", errInfo)
		}
	}
}

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

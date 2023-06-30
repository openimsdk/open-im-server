package mw

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type GinMwOptions func(*gin.RouterGroup)

func WithRecovery() GinMwOptions {
	return func(group *gin.RouterGroup) {
		group.Use(gin.Recovery())
	}
}

func WithCorsHandler() GinMwOptions {
	return func(group *gin.RouterGroup) {
		group.Use(CorsHandler())
	}
}

func WithGinParseOperationID() GinMwOptions {
	return func(group *gin.RouterGroup) {
		group.Use(GinParseOperationID())
	}
}

func WithGinParseToken(rdb redis.UniversalClient) GinMwOptions {
	return func(group *gin.RouterGroup) {
		group.Use(GinParseToken(rdb))
	}
}

func NewRouterGroup(routerGroup *gin.RouterGroup, route string, options ...GinMwOptions) *gin.RouterGroup {
	routerGroup = routerGroup.Group(route)
	for _, option := range options {
		option(routerGroup)
	}
	return routerGroup
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

func GinParseOperationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			operationID := c.Request.Header.Get(constant.OperationID)
			if operationID == "" {
				body, err := io.ReadAll(c.Request.Body)
				if err != nil {
					log.ZWarn(c, "read request body error", errs.ErrArgs.Wrap("read request body error: "+err.Error()))
					apiresp.GinError(c, errs.ErrArgs.Wrap("read request body error: "+err.Error()))
					c.Abort()
					return
				}
				req := struct {
					OperationID string `json:"operationID"`
				}{}
				if err := json.Unmarshal(body, &req); err != nil {
					log.ZWarn(c, "json unmarshal error", errs.ErrArgs.Wrap(err.Error()))
					apiresp.GinError(c, errs.ErrArgs.Wrap("json unmarshal error"+err.Error()))
					c.Abort()
					return
				}
				if req.OperationID == "" {
					err := errors.New("header must have operationID")
					log.ZWarn(c, "header must have operationID", err)
					apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error()))
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
func GinParseToken(rdb redis.UniversalClient) gin.HandlerFunc {
	dataBase := controller.NewAuthDatabase(cache.NewMsgCacheModel(rdb), config.Config.TokenPolicy.AccessSecret, config.Config.TokenPolicy.AccessExpire)
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost:
			token := c.Request.Header.Get(constant.Token)
			if token == "" {
				log.ZWarn(c, "header get token error", errs.ErrArgs.Wrap("header must have token"))
				apiresp.GinError(c, errs.ErrArgs.Wrap("header must have token"))
				c.Abort()
				return
			}
			claims, err := tokenverify.GetClaimFromToken(token)
			if err != nil {
				log.ZWarn(c, "jwt get token error", errs.ErrTokenUnknown.Wrap())
				apiresp.GinError(c, errs.ErrTokenUnknown.Wrap())
				c.Abort()
				return
			}
			m, err := dataBase.GetTokensWithoutError(c, claims.UserID, claims.PlatformID)
			if err != nil {
				log.ZWarn(c, "cache get token error", errs.ErrTokenNotExist.Wrap())
				apiresp.GinError(c, errs.ErrTokenNotExist.Wrap())
				c.Abort()
				return
			}
			if len(m) == 0 {
				log.ZWarn(c, "cache do not exist token error", errs.ErrTokenNotExist.Wrap())
				apiresp.GinError(c, errs.ErrTokenNotExist.Wrap())
				c.Abort()
				return
			}
			if v, ok := m[token]; ok {
				switch v {
				case constant.NormalToken:
				case constant.KickedToken:
					log.ZWarn(c, "cache kicked token error", errs.ErrTokenKicked.Wrap())
					apiresp.GinError(c, errs.ErrTokenKicked.Wrap())
					c.Abort()
					return
				default:
					log.ZWarn(c, "cache unknown token error", errs.ErrTokenUnknown.Wrap())
					apiresp.GinError(c, errs.ErrTokenUnknown.Wrap())
					c.Abort()
					return
				}
			} else {
				apiresp.GinError(c, errs.ErrTokenNotExist.Wrap())
				return
			}
			c.Set(constant.OpUserPlatform, constant.PlatformIDToName(claims.PlatformID))
			c.Set(constant.OpUserID, claims.UserID)
			c.Next()
		}
	}
}

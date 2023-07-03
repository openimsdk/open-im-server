package mw

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header(
			"Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar",
		) // 跨域关键设置 让浏览器可以解析
		c.Header(
			"Access-Control-Max-Age",
			"172800",
		) // 缓存请求信息 单位为秒
		c.Header(
			"Access-Control-Allow-Credentials",
			"false",
		) //  跨域请求是否需要带cookie信息 默认设置为true
		c.Header(
			"content-type",
			"application/json",
		) // 设置返回格式是json
		//Release all option pre-requests
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, "Options Request!")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GinParseOperationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			operationID := c.Request.Header.Get(constant.OperationID)
			if operationID == "" {
				err := errors.New("header must have operationID")
				apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error()))
				c.Abort()
				return
			}
			c.Set(constant.OperationID, operationID)
		}
		c.Next()
	}
}

func GinParseToken(rdb redis.UniversalClient) gin.HandlerFunc {
	dataBase := controller.NewAuthDatabase(
		cache.NewMsgCacheModel(rdb),
		config.Config.Secret,
		config.Config.TokenPolicy.Expire,
	)
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
				c.Abort()
				return
			}
			c.Set(constant.OpUserPlatform, constant.PlatformIDToName(claims.PlatformID))
			c.Set(constant.OpUserID, claims.UserID)
			c.Next()
		}
	}
}

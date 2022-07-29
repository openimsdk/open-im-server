package apiThird

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * FCM第三方上报Token
 */
type FcmUpdateTokenReq struct {
	OperationID string `json:"operationID"`
	Platform    int    `json:"platform" binding:"required,min=1,max=2"` //only for ios + android
	FcmToken    string `json:"fcmToken"`
}

func FcmUpdateToken(c *gin.Context) {
	var (
		req FcmUpdateTokenReq
	)
	if err := c.Bind(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)

	ok, UserId, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req, UserId)
	//逻辑处理开始
	err := db.DB.SetFcmToken(UserId, req.Platform, req.FcmToken, 0)
	if err != nil {
		errMsg := req.OperationID + " " + "SetFcmToken failed " + err.Error() + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	//逻辑处理完毕
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
	return
}

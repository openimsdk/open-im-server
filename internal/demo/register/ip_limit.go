package register

import (
	//api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"

	"github.com/gin-gonic/gin"

	"net/http"
	"time"
)

type QueryIPRegisterReq struct {
	OperationID string `json:"operationID"`
	IP          string `json:"ip"`
}

type QueryIPRegisterResp struct {
	IP          string   `json:"ip"`
	RegisterNum int      `json:"num"`
	Status      int      `json:"status"`
	UserIDList  []string `json:"userIDList"`
}

func QueryIPRegister(c *gin.Context) {
	req := QueryIPRegisterReq{}
	resp := QueryIPRegisterResp{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	userIDList, err := imdb.GetRegisterUserNum(req.IP)
	if err != nil {
		log.NewError(req.OperationID, "GetInvitationCode failed", req.IP)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB.ErrCode, "errMsg": "GetRegisterUserNum error!"})
		return
	}
	resp.IP = req.IP
	resp.RegisterNum = len(userIDList)
	resp.UserIDList = userIDList
	ipLimit, err := imdb.QueryIPLimits(req.IP)
	if err != nil {
		log.NewError(req.OperationID, "QueryIPLimits failed", req.IP, err.Error())
	} else {
		if ipLimit != nil {
			if ipLimit.Ip != "" {
				resp.Status = 1
			}
		}

	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp:", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
}

type AddIPLimitReq struct {
	OperationID string `json:"operationID"`
	IP          string `json:"ip"`
	LimitTime   int32  `json:"limitTime"`
}

type AddIPLimitResp struct {
}

func AddIPLimit(c *gin.Context) {
	req := AddIPLimitReq{}
	//resp := AddIPLimitResp{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	if err := imdb.InsertOneIntoIpLimits(db.IpLimit{
		Ip:            req.IP,
		LimitRegister: 1,
		LimitLogin:    1,
		CreateTime:    time.Now(),
		LimitTime:     utils.UnixSecondToTime(int64(req.LimitTime)),
	}); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.IP, req.LimitTime)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB.ErrCode, "errMsg": "InsertOneIntoIpLimits error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
}

type RemoveIPLimitReq struct {
	OperationID string `json:"operationID"`
	IP          string `json:"ip"`
}

type RemoveIPLimitResp struct {
}

func RemoveIPLimit(c *gin.Context) {
	req := RemoveIPLimitReq{}
	//resp := AddIPLimitResp{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.ErrArgs, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	if err := imdb.DeleteOneFromIpLimits(req.IP); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.IP)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB.ErrCode, "errMsg": "InsertOneIntoIpLimits error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
}

// ===========================================sk ==========================

type QueryUserIDIPLimitLoginReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

//type QueryUserIDIPLimitLoginResp struct {
//	UserIpLimit []db.UserIpLimit `json:"userIpLimit"`
//}

func QueryUserIDLimitLogin(c *gin.Context) {
	req := QueryUserIDIPLimitLoginReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	resp, err := imdb.GetIpLimitsLoginByUserID(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB.ErrCode, "errMsg": "GetIpLimitsByUserID error!"})
		return
	}
	if len(resp) > 0 {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp:", resp)
		c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": gin.H{"limit": resp}})
}

type AddUserIPLimitLoginReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
	IP          string `json:"ip"`
}

type AddUserIPLimitLoginResp struct {
}

// 添加ip 特定用户才能登录 user_ip_limits 表
func AddUserIPLimitLogin(c *gin.Context) {
	req := AddUserIPLimitLoginReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	userIp := db.UserIpLimit{UserID: req.UserID, Ip: req.IP}
	err := imdb.UpdateUserInfo(db.User{
		UserID: req.UserID,
		// LoginLimit: 1,
	})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "InsertUserIpLimitsLogin error!"})
		return
	}
	err = imdb.InsertUserIpLimitsLogin(&userIp)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "InsertUserIpLimitsLogin error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
}

type RemoveUserIPLimitReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
	IP          string `json:"ip"`
}

type RemoveUserIPLimitResp struct {
}

// 删除ip 特定用户才能登录 user_ip_limits 表
func RemoveUserIPLimitLogin(c *gin.Context) {
	req := RemoveUserIPLimitReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	err := imdb.DeleteUserIpLimitsLogin(req.UserID, req.IP)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "DeleteUserIpLimitsLogin error!"})
		return
	}
	ips, err := imdb.GetIpLimitsLoginByUserID(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserID)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "GetIpLimitsLoginByUserID error!"})
		return
	}
	if len(ips) == 0 {
		err := imdb.UpdateUserInfoByMap(db.User{
			UserID: req.UserID,
		}, map[string]interface{}{"login_limit": 0})
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserID)
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "UpdateUserInfo error!"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
}

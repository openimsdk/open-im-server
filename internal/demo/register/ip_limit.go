package register

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type QueryIPReq struct {
	OperationID string `json:"operationID"`
	IP          string `json:"ip"`
}

type QueryIPResp struct {
	IP          string   `json:"ip"`
	RegisterNum int      `json:"num"`
	UserIDList  []string `json:"userIDList"`
	Status      int
}

func QueryIP(c *gin.Context) {
	req := QueryIPReq{}
	resp := QueryIPResp{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	ips, err := imdb.QueryUserIPLimits(req.IP)
	if err != nil {
		log.NewError(req.OperationID, "GetInvitationCode failed", req.IP)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB.ErrCode, "errMsg": "QueryUserIPLimits error!"})
		return
	}
	resp.IP = req.IP
	resp.RegisterNum = len(ips)
	for _, ip := range ips {
		resp.UserIDList = append(resp.UserIDList, ip.UserID)
	}
	b, _ := imdb.IsLimitLoginIp(req.IP)
	if b == true {
		resp.Status = 1
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp:", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
}

type GetIPListReq struct {
}

type GetIPListResp struct {
}

func GetIPList(c *gin.Context) {

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
		LimitTime:     time.Time{},
	}); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.IP, req.LimitTime)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB.ErrCode, "errMsg": "InsertOneIntoIpLimits error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
}

type RemoveIPLimitReq struct {
}

type RemoveIPLimitResp struct {
}

func RemoveIPLimit(c *gin.Context) {

}

// ===========================================sk ==========================

type QueryUserIDIPLimitLoginReq struct {
	UserID      string `json:"userID" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

//type QueryUserIDIPLimitLoginResp struct {
//	UserIpLimit []db.UserIpLimit `json:"userIpLimit"`
//}

func QueryUserIPLimitLogin(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": "[]"})
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
	err := imdb.InsertUserIpLimitsLogin(&userIp)
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
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
}

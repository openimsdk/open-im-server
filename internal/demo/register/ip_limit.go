package register

import (
	api "Open_IM/pkg/base_info"
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
	UserIDList  []string `json:"userIDList"`
	Status      int      `json:"status"`
}

func QueryIPRegister(c *gin.Context) {
	req := QueryIPRegisterReq{}
	resp := QueryIPRegisterResp{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	ips, err := imdb.QueryUserIPLimits(req.IP)
	if err != nil {
		log.NewError(req.OperationID, "GetInvitationCode failed", req.IP)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "QueryUserIPLimits error!"})
		return
	}
	resp.IP = req.IP
	resp.RegisterNum = len(ips)
	for _, ip := range ips {
		resp.UserIDList = append(resp.UserIDList, ip.UserID)
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
		LimitTime:     time.Time{},
	}); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.IP, req.LimitTime)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "InsertOneIntoIpLimits error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": ""})
}

type RemoveIPLimitReq struct {
}

type RemoveIPLimitResp struct {
}

func RemoveIPLimit(c *gin.Context) {
	//DeleteOneFromIpLimits
}

// ===========================================sk 写

type QueryUserIDIPLimitReq struct {
	UserID string `json:"userID" binding:"required"`
}

type QueryUserIDIPLimitResp struct {
}

func QueryUserIDIPLimit(c *gin.Context) {

}

type AddUserIPLimitReq struct {
}

type AddUserIPLimitResp struct {
}

// 添加ip 特定用户才能登录 user_ip_limits 表
func AddUserIPLimit(c *gin.Context) {

}

type RemoveUserIPLimitReq struct {
}

type RemoveUserIPLimitResp struct {
}

// 删除ip 特定用户才能登录 user_ip_limits 表
func RemoveUserIPLimit(c *gin.Context) {

}

package register

import (
	apiStruct "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type InvitationCode struct {
	InvitationCode string    `json:"invitationCode"`
	CreateTime     time.Time `json:"createTime"`
	UserID         string    `json:"userID"`
	LastTime       time.Time `json:"lastTime"`
	Status         int32     `json:"status"`
}

type GenerateInvitationCodeReq struct {
	CodesNum    int    `json:"codesNum" binding:"required"`
	CodeLen     int    `json:"codeLen" binding:"required"`
	OperationID string `json:"operationID" binding:"required"`
}

type GenerateInvitationCodeResp struct {
	Codes []string `json:"codes"`
}

func GenerateInvitationCode(c *gin.Context) {
	req := GenerateInvitationCodeReq{}
	resp := GenerateInvitationCodeResp{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	var err error
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	resp.Codes, err = imdb.BatchCreateInvitationCodes(req.CodesNum, req.CodeLen)
	if err != nil {
		log.NewError(req.OperationID, "BatchCreateInvitationCodes failed", req.CodesNum, req.CodeLen)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "Verification code error!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
}

type QueryInvitationCodeReq struct {
	Code        string `json:"code"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type QueryInvitationCodeResp struct {
	InvitationCode
}

func QueryInvitationCode(c *gin.Context) {
	req := QueryInvitationCodeReq{}
	resp := QueryInvitationCodeResp{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	invitation, err := imdb.GetInvitationCode(req.Code)
	if err != nil {
		log.NewError(req.OperationID, "GetInvitationCode failed", req.Code)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "Verification code error!"})
		return
	}
	resp.UserID = invitation.UserID
	resp.CreateTime = invitation.CreateTime
	resp.Status = invitation.Status
	resp.LastTime = invitation.LastTime
	resp.InvitationCode.InvitationCode = invitation.InvitationCode
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp:", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
}

type GetInvitationCodesReq struct {
	Status      int32  `json:"status"`
	OperationID string `json:"operationID"  binding:"required"`
	apiStruct.Pagination
}

type GetInvitationCodesResp struct {
	apiStruct.Pagination
	Codes    []InvitationCode `json:"codes"`
	CodeNums int64            `json:"codeNums"`
}

func GetInvitationCodes(c *gin.Context) {
	req := GetInvitationCodesReq{}
	resp := GetInvitationCodesResp{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	codes, count, err := imdb.GetInvitationCodes(req.ShowNumber, req.PageNumber, req.Status)
	if err != nil {
		log.NewError(req.OperationID, "GetInvitationCode failed", req.ShowNumber, req.PageNumber, req.Status)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrDB, "errMsg": "Verification code error!"})
		return
	}
	resp.Pagination.PageNumber = req.PageNumber
	resp.Pagination.ShowNumber = req.ShowNumber
	for _, v := range codes {
		resp.Codes = append(resp.Codes, InvitationCode{
			InvitationCode: v.InvitationCode,
			CreateTime:     v.CreateTime,
			UserID:         v.UserID,
			LastTime:       v.LastTime,
			Status:         v.Status,
		})
	}
	resp.CodeNums = count
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp:", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
}

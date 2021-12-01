package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type ParamsSetPassword struct {
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	Password         string `json:"password"`
	VerificationCode string `json:"verificationCode"`
}

type Data struct {
	ExpiredTime int64  `json:"expiredTime"`
	Token       string `json:"token"`
	Uid         string `json:"uid"`
}

type IMRegisterResp struct {
	Data    Data   `json:"data"`
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

func SetPassword(c *gin.Context) {
	log.InfoByKv("setPassword api is statrting...", "")
	params := ParamsSetPassword{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	var account string
	if params.Email != "" {
		account = params.Email
	} else {
		account = params.PhoneNumber
	}

	log.InfoByKv("begin store redis", account)
	v, err := redis.String(db.DB.Exec("GET", account))

	if params.VerificationCode == config.Config.Demo.SuperCode {
		goto openIMRegisterTab
	}

	fmt.Println("Get Redis:", v, err)
	if err != nil {
		log.ErrorByKv("password Verification code expired", account, "err", err.Error())
		data := make(map[string]interface{})
		data["phoneNumber"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "Verification expired!", "data": data})
		return
	}
	if v != params.VerificationCode {
		log.InfoByKv("password Verification code error", account, params.VerificationCode)
		data := make(map[string]interface{})
		data["PhoneNumber"] = account
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "Verification code error!", "data": data})
		return
	}

openIMRegisterTab:
	log.InfoByKv("openIM register begin", account)
	resp, err := OpenIMRegister(account)

	log.InfoByKv("openIM register resp", account, resp, err)
	if err != nil {
		log.ErrorByKv("request openIM register error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.HttpError, "errMsg": err.Error()})
		return
	}
	response, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.IoErrot, "errMsg": err.Error()})
		return
	}
	imrep := IMRegisterResp{}
	err = json.Unmarshal(response, &imrep)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	if imrep.ErrCode != 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.HttpError, "errMsg": imrep.ErrMsg})
		return
	}

	queryParams := im_mysql_model.SetPasswordParams{
		Account:  account,
		Password: params.Password,
	}

	log.InfoByKv("begin store mysql", account, params.Password)
	_, err = im_mysql_model.SetPassword(&queryParams)
	if err != nil {
		log.ErrorByKv("set phone number password error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.DatabaseError, "errMsg": err.Error()})
		return
	}

	log.InfoByKv("end setPassword", account)
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": imrep.Data})
	return
}

func OpenIMRegister(account string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:10000/auth/user_register", utils.ServerIP)
	fmt.Println("1:", config.Config.Secret)

	client := &http.Client{}

	params := make(map[string]interface{})

	params["secret"] = config.Config.Secret
	params["platform"] = 7
	params["uid"] = account
	params["name"] = account
	params["icon"] = ""
	params["gender"] = 0

	params["mobile"] = ""

	params["email"] = ""
	params["birth"] = ""
	params["ex"] = ""
	con, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	log.InfoByKv("openIM register params", account, params)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(con)))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	return resp, err
}

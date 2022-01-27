package register

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type ParamsLogin struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Platform    int32  `json:"platform"`
}

func Login(c *gin.Context) {
	log.NewDebug("Login api is statrting...")
	params := ParamsLogin{}
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

	log.InfoByKv("api Login get params", account)

	queryParams := im_mysql_model.Register{
		Account:  account,
		Password: params.Password,
	}

	canLogin := im_mysql_model.Login(&queryParams)
	if canLogin == 1 {
		log.ErrorByKv("Incorrect phone number password", account, "err", "Mobile phone number is not registered")
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "Mobile phone number is not registered"})
		return
	}
	if canLogin == 2 {
		log.ErrorByKv("Incorrect phone number password", account, "err", "Incorrect password")
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "Incorrect password"})
		return
	}

	resp, err := OpenIMToken(account, params.Platform)
	if err != nil {
		log.ErrorByKv("get token by phone number err", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.HttpError, "errMsg": err.Error()})
		return
	}
	response, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.ErrorByKv("Failed to read file", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.IoError, "errMsg": err.Error()})
		return
	}
	imRep := IMRegisterResp{}
	err = json.Unmarshal(response, &imRep)
	if err != nil {
		log.ErrorByKv("json parsing failed", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	if imRep.ErrCode != 0 {
		log.ErrorByKv("openIM Login request failed", account, "err")
		c.JSON(http.StatusOK, gin.H{"errCode": constant.HttpError, "errMsg": imRep.ErrMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": imRep.Data})
	return

}

func OpenIMToken(Account string, platform int32) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:10000/auth/user_token", utils.ServerIP)

	client := &http.Client{}
	params := make(map[string]interface{})

	params["secret"] = config.Config.Secret
	params["platform"] = platform
	params["uid"] = Account
	con, err := json.Marshal(params)
	if err != nil {
		log.ErrorByKv("json parsing failed", Account, "err", err.Error())
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(con)))
	if err != nil {
		log.ErrorByKv("request error", "/auth/user_token", "err", err.Error())
		return nil, err
	}

	resp, err := client.Do(req)
	return resp, err
}

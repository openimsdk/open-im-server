package admin

import (
	apiStruct "Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdmin "Open_IM/pkg/proto/admin_cms"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"net/http"
	"path"
	"strings"

	"Open_IM/internal/api/third"
	"github.com/gin-gonic/gin"
)

// register
func AdminLogin(c *gin.Context) {
	var (
		req   apiStruct.AdminLoginRequest
		resp  apiStruct.AdminLoginResponse
		reqPb pbAdmin.AdminLoginReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.Secret = req.Secret
	reqPb.AdminID = req.AdminName
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName)
	client := pbAdmin.NewAdminCMSClient(etcdConn)
	respPb, err := client.AdminLogin(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.Token = respPb.Token
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func UploadUpdateApp(c *gin.Context) {
	var (
		req  apiStruct.UploadUpdateAppReq
		resp apiStruct.UploadUpdateAppResp
	)
	if err := c.Bind(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	newFileName, newYamlName, err := utils.GetUploadAppNewName(req.Type, req.Version, req.File.Filename, req.Yaml.Filename)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUploadAppNewName failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "invalid file type" + err.Error()})
		return
	}
	//fileObj, err := req.File.Open()
	//if err != nil {
	//	log.NewError(req.OperationID, utils.GetSelfFuncName(), "Open file error", err.Error())
	//	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "Open file error" + err.Error()})
	//	return
	//}
	//yamlObj, err := req.Yaml.Open()
	//if err != nil {
	//	log.NewError(req.OperationID, utils.GetSelfFuncName(), "Open Yaml error", err.Error())
	//	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "Open Yaml error" + err.Error()})
	//	return
	//}

	// v2.0.9_app_linux v2.0.9_yaml_linux
	file, err := c.FormFile("file")
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "FormFile failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "missing file arg: " + err.Error()})
		return
	}
	fileObj, err := file.Open()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "Open file error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "invalid file path" + err.Error()})
		return
	}

	yaml, err := c.FormFile("yaml")
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "FormFile failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "missing file arg: " + err.Error()})
		return
	}
	yamlObj, err := yaml.Open()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "Open file error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "invalid file path" + err.Error()})
		return
	}

	fmt.Println(req.OperationID, utils.GetSelfFuncName(), "name: ", config.Config.Credential.Minio.AppBucket, newFileName, fileObj, file.Size)
	fmt.Println(req.OperationID, utils.GetSelfFuncName(), "name: ", config.Config.Credential.Minio.AppBucket, newYamlName, yamlObj, yaml.Size)

	_, err = apiThird.MinioClient.PutObject(context.Background(), config.Config.Credential.Minio.AppBucket, newFileName, fileObj, file.Size, minio.PutObjectOptions{ContentType: path.Ext(newFileName)})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "PutObject file error")
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "PutObject file error" + err.Error()})
		return
	}
	_, err = apiThird.MinioClient.PutObject(context.Background(), config.Config.Credential.Minio.AppBucket, newYamlName, yamlObj, file.Size, minio.PutObjectOptions{ContentType: path.Ext(newYamlName)})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "PutObject yaml error")
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "PutObject yaml error" + err.Error()})
		return
	}
	if err := imdb.UpdateAppVersion(req.Type, req.Version, req.ForceUpdate, newFileName, newYamlName); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateAppVersion error", err.Error())
		resp.ErrCode = http.StatusInternalServerError
		resp.ErrMsg = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName())
	c.JSON(http.StatusOK, resp)
}

func GetDownloadURL(c *gin.Context) {
	var (
		req  apiStruct.GetDownloadURLReq
		resp apiStruct.GetDownloadURLResp
	)
	defer func() {
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	}()
	if err := c.Bind(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	//fileName, yamlName, err := utils.GetUploadAppNewName(req.Type, req.Version, req.)
	//if err != nil {
	//	log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUploadAppNewName failed", err.Error())
	//	c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "invalid file type" + err.Error()})
	//	return
	//}
	app, err := imdb.GetNewestVersion(req.Type)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "getNewestVersion failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "getNewestVersion failed" + err.Error()})
		return
	}
	if app.Version != req.Version {
		resp.Data.HasNewVersion = true
		if app.ForceUpdate == true {
			resp.Data.ForceUpdate = true
		}
		resp.Data.YamlURL = config.Config.Credential.Minio.Endpoint + "/" + config.Config.Credential.Minio.AppBucket + "/" + app.YamlName
		resp.Data.FileURL = config.Config.Credential.Minio.Endpoint + "/" + config.Config.Credential.Minio.AppBucket + "/" + app.FileName
		c.JSON(http.StatusOK, resp)
	} else {
		resp.Data.HasNewVersion = false
		c.JSON(http.StatusOK, resp)
	}
}

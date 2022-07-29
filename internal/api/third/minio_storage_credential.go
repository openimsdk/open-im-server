package apiThird

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	_ "Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	_ "github.com/minio/minio-go/v7"
	cr "github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
)

// @Summary minio上传文件(web api)
// @Description minio上传文件(web api), 请注意本api请求为form并非json
// @Tags 第三方服务相关
// @ID MinioUploadFile
// @Accept json
// @Param token header string true "im token"
// @Param file formData file true "要上传的文件文件"
// @Param fileType formData int true "文件类型"
// @Param operationID formData string true "操作唯一ID"
// @Produce json
// @Success 0 {object} api.MinioUploadFileResp ""
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /third/minio_upload [post]
func MinioUploadFile(c *gin.Context) {
	var (
		req  api.MinioUploadFileReq
		resp api.MinioUploadFile
	)
	defer func() {
		if r := recover(); r != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), r)
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "missing file or snapShot args"})
			return
		}
	}()
	if err := c.Bind(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	switch req.FileType {
	// videoType upload snapShot
	case constant.VideoType:
		snapShotFile, err := c.FormFile("snapShot")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "missing snapshot arg: " + err.Error()})
			return
		}
		snapShotFileObj, err := snapShotFile.Open()
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "Open file error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
			return
		}
		snapShotNewName, snapShotNewType := utils.GetNewFileNameAndContentType(snapShotFile.Filename, constant.ImageType)
		log.Debug(req.OperationID, utils.GetSelfFuncName(), snapShotNewName, snapShotNewType)
		_, err = MinioClient.PutObject(context.Background(), config.Config.Credential.Minio.Bucket, snapShotNewName, snapShotFileObj, snapShotFile.Size, minio.PutObjectOptions{ContentType: snapShotNewType})
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "PutObject snapShotFile error", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
			return
		}
		resp.SnapshotURL = config.Config.Credential.Minio.Endpoint + "/" + config.Config.Credential.Minio.Bucket + "/" + snapShotNewName
		resp.SnapshotNewName = snapShotNewName
	}
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
	newName, newType := utils.GetNewFileNameAndContentType(file.Filename, req.FileType)
	log.Debug(req.OperationID, utils.GetSelfFuncName(), config.Config.Credential.Minio.Bucket, newName, fileObj, file.Size, newType, MinioClient.EndpointURL())
	_, err = MinioClient.PutObject(context.Background(), config.Config.Credential.Minio.Bucket, newName, fileObj, file.Size, minio.PutObjectOptions{ContentType: newType})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "upload file error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "upload file error" + err.Error()})
		return
	}
	resp.NewName = newName
	resp.URL = config.Config.Credential.Minio.Endpoint + "/" + config.Config.Credential.Minio.Bucket + "/" + newName
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
	return
}

func MinioStorageCredential(c *gin.Context) {
	var (
		req  api.MinioStorageCredentialReq
		resp api.MiniostorageCredentialResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	var stsOpts cr.STSAssumeRoleOptions
	stsOpts.AccessKey = config.Config.Credential.Minio.AccessKeyID
	stsOpts.SecretKey = config.Config.Credential.Minio.SecretAccessKey
	stsOpts.DurationSeconds = constant.MinioDurationTimes
	var endpoint string
	if config.Config.Credential.Minio.EndpointInnerEnable {
		endpoint = config.Config.Credential.Minio.EndpointInner
	} else {
		endpoint = config.Config.Credential.Minio.Endpoint
	}
	li, err := cr.NewSTSAssumeRole(endpoint, stsOpts)
	if err != nil {
		log.NewError("", utils.GetSelfFuncName(), "NewSTSAssumeRole failed", err.Error(), stsOpts, config.Config.Credential.Minio.Endpoint)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	v, err := li.Get()
	if err != nil {
		log.NewError("0", utils.GetSelfFuncName(), "li.Get error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	resp.SessionToken = v.SessionToken
	resp.SecretAccessKey = v.SecretAccessKey
	resp.AccessKeyID = v.AccessKeyID
	resp.BucketName = config.Config.Credential.Minio.Bucket
	resp.StsEndpointURL = config.Config.Credential.Minio.Endpoint
	resp.StorageTime = config.Config.Credential.Minio.StorageTime
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
}

func UploadUpdateApp(c *gin.Context) {
	var (
		req  api.UploadUpdateAppReq
		resp api.UploadUpdateAppResp
	)
	if err := c.Bind(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	var yamlName string
	if req.Yaml == nil {
		yamlName = ""
	} else {
		yamlName = req.Yaml.Filename
	}
	fileObj, err := req.File.Open()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "Open file error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "Open file error" + err.Error()})
		return
	}
	_, err = MinioClient.PutObject(context.Background(), config.Config.Credential.Minio.AppBucket, req.File.Filename, fileObj, req.File.Size, minio.PutObjectOptions{})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "PutObject file error")
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "PutObject file error" + err.Error()})
		return
	}
	if yamlName != "" {
		yamlObj, err := req.Yaml.Open()
		if err == nil {
			_, err = MinioClient.PutObject(context.Background(), config.Config.Credential.Minio.AppBucket, yamlName, yamlObj, req.Yaml.Size, minio.PutObjectOptions{})
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "PutObject yaml error")
				c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "PutObject yaml error" + err.Error()})
				return
			}
		} else {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		}
	}
	if err := imdb.UpdateAppVersion(req.Type, req.Version, req.ForceUpdate, req.File.Filename, yamlName, req.UpdateLog); err != nil {
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
		req  api.GetDownloadURLReq
		resp api.GetDownloadURLResp
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
	app, err := imdb.GetNewestVersion(req.Type)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "getNewestVersion failed", err.Error())
	}
	log.Debug(req.OperationID, utils.GetSelfFuncName(), "app: ", app)
	if app != nil {
		if app.Version != req.Version && app.Version != "" {
			resp.Data.HasNewVersion = true
			if app.ForceUpdate == true {
				resp.Data.ForceUpdate = true
			}
			if app.YamlName != "" {
				resp.Data.YamlURL = config.Config.Credential.Minio.Endpoint + "/" + config.Config.Credential.Minio.AppBucket + "/" + app.YamlName
			}
			resp.Data.FileURL = config.Config.Credential.Minio.Endpoint + "/" + config.Config.Credential.Minio.AppBucket + "/" + app.FileName
			resp.Data.Version = app.Version
			resp.Data.UpdateLog = app.UpdateLog
			c.JSON(http.StatusOK, resp)
			return
		} else {
			resp.Data.HasNewVersion = false
			c.JSON(http.StatusOK, resp)
			return
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"errCode": 0, "errMsg": "not found app version"})
}

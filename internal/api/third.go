package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apistruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/minio/minio-go/v7"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type ThirdApi rpcclient.Third

func NewThirdApi(discov discoveryregistry.SvcDiscoveryRegistry) ThirdApi {
	return ThirdApi(*rpcclient.NewThird(discov))
}

func (o *ThirdApi) FcmUpdateToken(c *gin.Context) {
	a2r.Call(third.ThirdClient.FcmUpdateToken, o.Client, c)
}

func (o *ThirdApi) SetAppBadge(c *gin.Context) {
	a2r.Call(third.ThirdClient.SetAppBadge, o.Client, c)
}

// #################### s3 ####################

func (o *ThirdApi) PartLimit(c *gin.Context) {
	a2r.Call(third.ThirdClient.PartLimit, o.Client, c)
}

func (o *ThirdApi) PartSize(c *gin.Context) {
	a2r.Call(third.ThirdClient.PartSize, o.Client, c)
}

func (o *ThirdApi) InitiateMultipartUpload(c *gin.Context) {
	a2r.Call(third.ThirdClient.InitiateMultipartUpload, o.Client, c)
}

func (o *ThirdApi) AuthSign(c *gin.Context) {
	a2r.Call(third.ThirdClient.AuthSign, o.Client, c)
}

func (o *ThirdApi) CompleteMultipartUpload(c *gin.Context) {
	a2r.Call(third.ThirdClient.CompleteMultipartUpload, o.Client, c)
}

func (o *ThirdApi) AccessURL(c *gin.Context) {
	a2r.Call(third.ThirdClient.AccessURL, o.Client, c)
}

func (o *ThirdApi) ObjectRedirect(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.String(http.StatusBadRequest, "name is empty")
		return
	}
	if name[0] == '/' {
		name = name[1:]
	}
	operationID := c.Query("operationID")
	if operationID == "" {
		operationID = strconv.Itoa(rand.Int())
	}
	ctx := mcontext.SetOperationID(c, operationID)
	resp, err := o.Client.AccessURL(ctx, &third.AccessURLReq{Name: name})
	if err != nil {
		if errs.ErrArgs.Is(err) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if errs.ErrRecordNotFound.Is(err) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, resp.Url)
}

func (o *ThirdApi) MinioUploadFile(c *gin.Context) {
	var (
		req  apistruct.MinioUploadFileReq
		resp apistruct.MinioUploadFile
	)

	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
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
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
			return
		}
		snapShotNewName, snapShotNewType := utils.GetNewFileNameAndContentType(snapShotFile.Filename, constant.ImageType)
		_, err = o.MinioClient.PutObject(c, config.Config.Object.Minio.Bucket, snapShotNewName, snapShotFileObj, snapShotFile.Size, minio.PutObjectOptions{ContentType: snapShotNewType})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
			return
		}
		resp.SnapshotURL = config.Config.Object.Minio.Endpoint + "/" + config.Config.Object.Minio.Bucket + "/" + snapShotNewName
		resp.SnapshotNewName = snapShotNewName
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "missing file arg: " + err.Error()})
		return
	}
	fileObj, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "invalid file path" + err.Error()})
		return
	}
	newName, newType := utils.GetNewFileNameAndContentType(file.Filename, req.FileType)
	_, err = o.MinioClient.PutObject(c, config.Config.Object.Minio.Bucket, newName, fileObj, file.Size, minio.PutObjectOptions{ContentType: newType})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "upload file error" + err.Error()})
		return
	}
	resp.NewName = newName
	resp.URL = config.Config.Object.Minio.Endpoint + "/" + config.Config.Object.Minio.Bucket + "/" + newName
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
	return
}

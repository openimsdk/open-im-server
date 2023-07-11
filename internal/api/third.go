package api

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apistruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/minio/minio-go/v7"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

type ThirdApi rpcclient.Third

func NewThirdApi(discov discoveryregistry.SvcDiscoveryRegistry) ThirdApi {
	return ThirdApi(*rpcclient.NewThird(discov))
}

func (o *ThirdApi) ApplyPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ApplyPut, o.Client, c)
}

func (o *ThirdApi) GetPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetPut, o.Client, c)
}

func (o *ThirdApi) ConfirmPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ConfirmPut, o.Client, c)
}

func (o *ThirdApi) GetHash(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetHashInfo, o.Client, c)
}

func (o *ThirdApi) FcmUpdateToken(c *gin.Context) {
	a2r.Call(third.ThirdClient.FcmUpdateToken, o.Client, c)
}

func (o *ThirdApi) SetAppBadge(c *gin.Context) {
	a2r.Call(third.ThirdClient.SetAppBadge, o.Client, c)
}

func (o *ThirdApi) GetURL(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		a2r.Call(third.ThirdClient.GetUrl, o.Client, c)
		return
	}
	name := c.Query("name")
	if name == "" {
		c.String(http.StatusBadRequest, "name is empty")
		return
	}
	operationID := c.Query("operationID")
	if operationID == "" {
		operationID = "auto_" + strconv.Itoa(rand.Int())
	}
	expires, _ := strconv.ParseInt(c.Query("expires"), 10, 64)
	if expires <= 0 {
		expires = 3600 * 1000
	}
	attachment, _ := strconv.ParseBool(c.Query("attachment"))
	c.Set(constant.OperationID, operationID)
	resp, err := o.Client.GetUrl(mcontext.SetOperationID(c, operationID), &third.GetUrlReq{Name: name, Expires: expires, Attachment: attachment})
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
		_, err = o.MinioClient.PutObject(context.Background(), config.Config.Object.Minio.DataBucket, snapShotNewName, snapShotFileObj, snapShotFile.Size, minio.PutObjectOptions{ContentType: snapShotNewType})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
			return
		}
		resp.SnapshotURL = config.Config.Object.Minio.Endpoint + "/" + config.Config.Object.Minio.DataBucket + "/" + snapShotNewName
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
	_, err = o.MinioClient.PutObject(context.Background(), config.Config.Object.Minio.DataBucket, newName, fileObj, file.Size, minio.PutObjectOptions{ContentType: newType})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "upload file error" + err.Error()})
		return
	}
	resp.NewName = newName
	resp.URL = config.Config.Object.Minio.Endpoint + "/" + config.Config.Object.Minio.DataBucket + "/" + newName
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": resp})
	return
}

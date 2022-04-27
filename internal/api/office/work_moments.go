package office

import (
	apiStruct "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbOffice "Open_IM/pkg/proto/office"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CreateOneWorkMoment(c *gin.Context) {
	var (
		req    apiStruct.CreateOneWorkMomentReq
		resp   apiStruct.CreateOneWorkMomentResp
		reqPb  pbOffice.CreateOneWorkMomentReq
		respPb *pbOffice.CreateOneWorkMomentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.WorkMoment.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.CreateOneWorkMoment(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateOneWorkMoment rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CreateOneWorkMoment rpc server failed" + err.Error()})
		return
	}
	resp.CommResp = apiStruct.CommResp{
		ErrCode: respPb.CommonResp.ErrCode,
		ErrMsg:  respPb.CommonResp.ErrMsg,
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func DeleteOneWorkMoment(c *gin.Context) {
	var (
		req    apiStruct.DeleteOneWorkMomentReq
		resp   apiStruct.DeleteOneWorkMomentResp
		reqPb  pbOffice.DeleteOneWorkMomentReq
		respPb *pbOffice.DeleteOneWorkMomentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.DeleteOneWorkMoment(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DeleteOneWorkMoment rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "DeleteOneWorkMoment rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func LikeOneWorkMoment(c *gin.Context) {
	var (
		req    apiStruct.LikeOneWorkMomentReq
		resp   apiStruct.LikeOneWorkMomentResp
		reqPb  pbOffice.LikeOneWorkMomentReq
		respPb *pbOffice.LikeOneWorkMomentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.LikeOneWorkMoment(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "LikeOneWorkMoment rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "LikeOneWorkMoment rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func CommentOneWorkMoment(c *gin.Context) {
	var (
		req    apiStruct.CommentOneWorkMomentReq
		resp   apiStruct.CommentOneWorkMomentResp
		reqPb  pbOffice.CommentOneWorkMomentReq
		respPb *pbOffice.CommentOneWorkMomentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.CommentOneWorkMoment(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CommentOneWorkMoment rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CommentOneWorkMoment rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetWorkMomentByID(c *gin.Context) {
	var (
		req    apiStruct.GetWorkMomentByIDReq
		resp   apiStruct.GetWorkMomentByIDResp
		reqPb  pbOffice.GetWorkMomentByIDReq
		respPb *pbOffice.GetWorkMomentByIDResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.OpUserID = userID
	reqPb.WorkMomentID = req.WorkMomentID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.GetWorkMomentByID(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserWorkMoments rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserWorkMoments rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	resp.Data.WorkMoment = &apiStruct.WorkMoment{LikeUserList: []*apiStruct.WorkMomentUser{}, Comments: []*apiStruct.Comment{}, AtUserList: []*apiStruct.WorkMomentUser{}}
	if err := utils.CopyStructFields(&resp.Data.WorkMoment, respPb.WorkMoment); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetUserWorkMoments(c *gin.Context) {
	var (
		req    apiStruct.GetUserWorkMomentsReq
		resp   apiStruct.GetUserWorkMomentsResp
		reqPb  pbOffice.GetUserWorkMomentsReq
		respPb *pbOffice.GetUserWorkMomentsResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, opUserID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.Pagination = &pbCommon.RequestPagination{
		PageNumber: req.PageNumber,
		ShowNumber: req.ShowNumber,
	}
	reqPb.OpUserID = opUserID
	reqPb.UserID = req.UserID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.GetUserWorkMoments(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserWorkMoments rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserWorkMoments rpc server failed" + err.Error()})
		return
	}
	resp.Data.WorkMoments = []*apiStruct.WorkMoment{}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	//if err := utils.CopyStructFields(&resp.Data.WorkMoments, respPb.WorkMoments); err != nil {
	//	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	//}
	for _, v := range respPb.WorkMoments {
		workMoment := apiStruct.WorkMoment{
			WorkMomentID: v.WorkMomentID,
			UserID:       v.UserID,
			Content:      v.Content,
			FaceURL:      v.FaceURL,
			UserName:     v.UserName,
			CreateTime:   v.CreateTime,
			Comments:     make([]*apiStruct.Comment, len(v.Comments)),
			LikeUserList: make([]*apiStruct.WorkMomentUser, len(v.LikeUserList)),
			AtUserList:   make([]*apiStruct.WorkMomentUser, len(v.AtUserList)),
		}
		for i, comment := range v.Comments {
			workMoment.Comments[i] = &apiStruct.Comment{
				UserID:        comment.UserID,
				UserName:      comment.UserName,
				ReplyUserID:   comment.ReplyUserID,
				ReplyUserName: comment.ReplyUserName,
				ContentID:     comment.ContentID,
				Content:       comment.Content,
				CreateTime:    comment.CreateTime,
			}
		}
		for i, likeUser := range v.LikeUserList {
			workMoment.LikeUserList[i] = &apiStruct.WorkMomentUser{
				UserID:   likeUser.UserID,
				UserName: likeUser.UserName,
			}
		}
		for i, atUser := range v.AtUserList {
			workMoment.AtUserList[i] = &apiStruct.WorkMomentUser{
				UserID:   atUser.UserID,
				UserName: atUser.UserName,
			}
		}
		resp.Data.WorkMoments = append(resp.Data.WorkMoments, &workMoment)
	}
	resp.Data.ShowNumber = respPb.Pagination.ShowNumber
	resp.Data.CurrentPage = respPb.Pagination.CurrentPage
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetUserFriendWorkMoments(c *gin.Context) {
	var (
		req    apiStruct.GetUserFriendWorkMomentsReq
		resp   apiStruct.GetUserFriendWorkMomentsResp
		reqPb  pbOffice.GetUserFriendWorkMomentsReq
		respPb *pbOffice.GetUserFriendWorkMomentsResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.Pagination = &pbCommon.RequestPagination{
		PageNumber: req.PageNumber,
		ShowNumber: req.ShowNumber,
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.GetUserFriendWorkMoments(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserFriendWorkMoments rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserFriendWorkMoments rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	//if err := utils.CopyStructFields(&resp.Data.WorkMoments, respPb.WorkMoments); err != nil {
	//	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	//}
	resp.Data.WorkMoments = []*apiStruct.WorkMoment{}
	for _, v := range respPb.WorkMoments {
		workMoment := apiStruct.WorkMoment{
			WorkMomentID: v.WorkMomentID,
			UserID:       v.UserID,
			Content:      v.Content,
			FaceURL:      v.FaceURL,
			UserName:     v.UserName,
			CreateTime:   v.CreateTime,
			Comments:     make([]*apiStruct.Comment, len(v.Comments)),
			LikeUserList: make([]*apiStruct.WorkMomentUser, len(v.LikeUserList)),
			AtUserList:   make([]*apiStruct.WorkMomentUser, len(v.AtUserList)),
		}
		for i, comment := range v.Comments {
			workMoment.Comments[i] = &apiStruct.Comment{
				UserID:        comment.UserID,
				UserName:      comment.UserName,
				ReplyUserID:   comment.ReplyUserID,
				ReplyUserName: comment.ReplyUserName,
				ContentID:     comment.ContentID,
				Content:       comment.Content,
				CreateTime:    comment.CreateTime,
			}
		}
		for i, likeUser := range v.LikeUserList {
			workMoment.LikeUserList[i] = &apiStruct.WorkMomentUser{
				UserID:   likeUser.UserID,
				UserName: likeUser.UserName,
			}
		}
		for i, atUser := range v.AtUserList {
			workMoment.AtUserList[i] = &apiStruct.WorkMomentUser{
				UserID:   atUser.UserID,
				UserName: atUser.UserName,
			}
		}
		resp.Data.WorkMoments = append(resp.Data.WorkMoments, &workMoment)
	}
	resp.Data.ShowNumber = respPb.Pagination.ShowNumber
	resp.Data.CurrentPage = respPb.Pagination.CurrentPage
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func SetUserWorkMomentsLevel(c *gin.Context) {
	var (
		req    apiStruct.SetUserWorkMomentsLevelReq
		resp   apiStruct.SetUserWorkMomentsLevelResp
		reqPb  pbOffice.SetUserWorkMomentsLevelReq
		respPb *pbOffice.SetUserWorkMomentsLevelResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	ok, userID := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.SetUserWorkMomentsLevel(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetUserWorkMomentsLevel rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "SetUserWorkMomentsLevel rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

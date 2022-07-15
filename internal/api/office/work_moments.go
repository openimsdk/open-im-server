package office

import (
	api "Open_IM/pkg/base_info"
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

// @Summary 创建一条工作圈
// @Description 用户创建一条工作圈
// @Tags 工作圈
// @ID CreateOneWorkMoment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CreateOneWorkMomentReq true "请求 atUserList likeUserList permissionGroupList permissionUserList 字段中userName可以不填"
// @Produce json
// @Success 0 {object} api.CreateOneWorkMomentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/create_one_work_moment [post]
func CreateOneWorkMoment(c *gin.Context) {
	var (
		req    api.CreateOneWorkMomentReq
		resp   api.CreateOneWorkMomentResp
		reqPb  pbOffice.CreateOneWorkMomentReq
		respPb *pbOffice.CreateOneWorkMomentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.WorkMoment.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.CreateOneWorkMoment(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateOneWorkMoment rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CreateOneWorkMoment rpc server failed" + err.Error()})
		return
	}
	resp.CommResp = api.CommResp{
		ErrCode: respPb.CommonResp.ErrCode,
		ErrMsg:  respPb.CommonResp.ErrMsg,
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 删除一条工作圈
// @Description 根据用户工作圈ID删除一条工作圈
// @Tags 工作圈
// @ID DeleteOneWorkMoment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.DeleteOneWorkMomentReq true "请求"
// @Produce json
// @Success 0 {object} api.DeleteOneWorkMomentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/delete_one_work_moment [post]
func DeleteOneWorkMoment(c *gin.Context) {
	var (
		req    api.DeleteOneWorkMomentReq
		resp   api.DeleteOneWorkMomentResp
		reqPb  pbOffice.DeleteOneWorkMomentReq
		respPb *pbOffice.DeleteOneWorkMomentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

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

// @Summary 点赞一条工作圈
// @Description 工作圈ID点赞一条工作圈
// @Tags 工作圈
// @ID LikeOneWorkMoment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.LikeOneWorkMomentReq true "请求"
// @Produce json
// @Success 0 {object} api.LikeOneWorkMomentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/like_one_work_moment [post]
func LikeOneWorkMoment(c *gin.Context) {
	var (
		req    api.LikeOneWorkMomentReq
		resp   api.LikeOneWorkMomentResp
		reqPb  pbOffice.LikeOneWorkMomentReq
		respPb *pbOffice.LikeOneWorkMomentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 评论一条工作圈
// @Description 评论一条工作圈
// @Tags 工作圈
// @ID CommentOneWorkMoment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CommentOneWorkMomentReq true "请求"
// @Produce json
// @Success 0 {object} api.CommentOneWorkMomentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/comment_one_work_moment [post]
func CommentOneWorkMoment(c *gin.Context) {
	var (
		req    api.CommentOneWorkMomentReq
		resp   api.CommentOneWorkMomentResp
		reqPb  pbOffice.CommentOneWorkMomentReq
		respPb *pbOffice.CommentOneWorkMomentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

// @Summary 删除一条评论
// @Description 删除一条评论
// @Tags 工作圈
// @ID DeleteComment
// @Accept json
// @Param token header string true "im token"
// @Param req body api.DeleteCommentReq true "请求"
// @Produce json
// @Success 0 {object} api.DeleteCommentResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/delete_comment [post]
func DeleteComment(c *gin.Context) {
	var (
		req    api.DeleteCommentReq
		resp   api.DeleteCommentResp
		reqPb  pbOffice.DeleteCommentReq
		respPb *pbOffice.DeleteCommentResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), err.Error())
	}

	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.DeleteComment(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DeleteComment rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "DeleteComment rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 通过ID获取工作圈
// @Description 通过ID获取工作圈
// @Tags 工作圈
// @ID GetWorkMomentByID
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetWorkMomentByIDReq true "请求"
// @Produce json
// @Success 0 {object} api.GetWorkMomentByIDResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/get_work_moment_by_id [post]
func GetWorkMomentByID(c *gin.Context) {
	var (
		req    api.GetWorkMomentByIDReq
		resp   api.GetWorkMomentByIDResp
		reqPb  pbOffice.GetWorkMomentByIDReq
		respPb *pbOffice.GetWorkMomentByIDResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.OpUserID = userID
	reqPb.WorkMomentID = req.WorkMomentID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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
	resp.Data.WorkMoment = &api.WorkMoment{LikeUserList: []*api.WorkMomentUser{}, Comments: []*api.Comment{},
		AtUserList: []*api.WorkMomentUser{}, PermissionUserList: []*api.WorkMomentUser{}}
	if err := utils.CopyStructFields(&resp.Data.WorkMoment, respPb.WorkMoment); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 查询用户工作圈
// @Description 查询用户工作圈
// @Tags 工作圈
// @ID GetUserWorkMoments
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUserWorkMomentsReq true "请求"
// @Produce json
// @Success 0 {object} api.GetUserWorkMomentsResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/get_user_work_moments [post]
func GetUserWorkMoments(c *gin.Context) {
	var (
		req    api.GetUserWorkMomentsReq
		resp   api.GetUserWorkMomentsResp
		reqPb  pbOffice.GetUserWorkMomentsReq
		respPb *pbOffice.GetUserWorkMomentsResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)

	var ok bool
	var errInfo string
	var opUserID string
	ok, opUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.Pagination = &pbCommon.RequestPagination{
		PageNumber: req.PageNumber,
		ShowNumber: req.ShowNumber,
	}
	reqPb.OpUserID = opUserID
	reqPb.UserID = req.UserID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.GetUserWorkMoments(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserWorkMoments rpc failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserWorkMoments rpc server failed" + err.Error()})
		return
	}
	resp.Data.WorkMoments = []*api.WorkMoment{}
	if err := utils.CopyStructFields(&resp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	for _, v := range respPb.WorkMoments {
		workMoment := api.WorkMoment{
			WorkMomentID:       v.WorkMomentID,
			UserID:             v.UserID,
			Content:            v.Content,
			FaceURL:            v.FaceURL,
			UserName:           v.UserName,
			CreateTime:         v.CreateTime,
			Comments:           make([]*api.Comment, len(v.Comments)),
			LikeUserList:       make([]*api.WorkMomentUser, len(v.LikeUserList)),
			AtUserList:         make([]*api.WorkMomentUser, len(v.AtUserList)),
			PermissionUserList: make([]*api.WorkMomentUser, len(v.PermissionUserList)),
			Permission:         v.Permission,
		}
		for i, comment := range v.Comments {
			workMoment.Comments[i] = &api.Comment{
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
			workMoment.LikeUserList[i] = &api.WorkMomentUser{
				UserID:   likeUser.UserID,
				UserName: likeUser.UserName,
			}
		}
		for i, atUser := range v.AtUserList {
			workMoment.AtUserList[i] = &api.WorkMomentUser{
				UserID:   atUser.UserID,
				UserName: atUser.UserName,
			}
		}
		for i, permissionUser := range v.PermissionUserList {
			workMoment.PermissionUserList[i] = &api.WorkMomentUser{
				UserID:   permissionUser.UserID,
				UserName: permissionUser.UserName,
			}
		}
		resp.Data.WorkMoments = append(resp.Data.WorkMoments, &workMoment)
	}
	resp.Data.ShowNumber = respPb.Pagination.ShowNumber
	resp.Data.CurrentPage = respPb.Pagination.CurrentPage
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 查询自己大工作圈页面
// @Description 查询用户工作圈页面
// @Tags 工作圈
// @ID GetUserFriendWorkMoments
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUserFriendWorkMomentsReq true "请求"
// @Produce json
// @Success 0 {object} api.GetUserFriendWorkMomentsResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/get_user_friend_work_moments [post]
func GetUserFriendWorkMoments(c *gin.Context) {
	var (
		req    api.GetUserFriendWorkMomentsReq
		resp   api.GetUserFriendWorkMomentsResp
		reqPb  pbOffice.GetUserFriendWorkMomentsReq
		respPb *pbOffice.GetUserFriendWorkMomentsResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.OperationID = req.OperationID
	reqPb.Pagination = &pbCommon.RequestPagination{
		PageNumber: req.PageNumber,
		ShowNumber: req.ShowNumber,
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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
	resp.Data.WorkMoments = []*api.WorkMoment{}
	for _, v := range respPb.WorkMoments {
		workMoment := api.WorkMoment{
			WorkMomentID:       v.WorkMomentID,
			UserID:             v.UserID,
			Content:            v.Content,
			FaceURL:            v.FaceURL,
			UserName:           v.UserName,
			CreateTime:         v.CreateTime,
			Comments:           make([]*api.Comment, len(v.Comments)),
			LikeUserList:       make([]*api.WorkMomentUser, len(v.LikeUserList)),
			AtUserList:         make([]*api.WorkMomentUser, len(v.AtUserList)),
			PermissionUserList: make([]*api.WorkMomentUser, len(v.PermissionUserList)),
			Permission:         v.Permission,
		}
		for i, comment := range v.Comments {
			workMoment.Comments[i] = &api.Comment{
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
			workMoment.LikeUserList[i] = &api.WorkMomentUser{
				UserID:   likeUser.UserID,
				UserName: likeUser.UserName,
			}
		}
		for i, atUser := range v.AtUserList {
			workMoment.AtUserList[i] = &api.WorkMomentUser{
				UserID:   atUser.UserID,
				UserName: atUser.UserName,
			}
		}
		for i, permissionUser := range v.PermissionUserList {
			workMoment.PermissionUserList[i] = &api.WorkMomentUser{
				UserID:   permissionUser.UserID,
				UserName: permissionUser.UserName,
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
		req    api.SetUserWorkMomentsLevelReq
		resp   api.SetUserWorkMomentsLevelResp
		reqPb  pbOffice.SetUserWorkMomentsLevelReq
		respPb *pbOffice.SetUserWorkMomentsLevelResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
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

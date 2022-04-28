package office

import (
	"Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbOffice "Open_IM/pkg/proto/office"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"
)

type officeServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewOfficeServer(port int) *officeServer {
	log.NewPrivateLog(constant.LogFileName)
	return &officeServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImOfficeName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *officeServer) Run() {
	log.NewInfo("0", "officeServer rpc start ")
	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("0", "Listen failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("0", "listen network success, ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbOffice.RegisterOfficeServiceServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *officeServer) GetUserTags(_ context.Context, req *pbOffice.GetUserTagsReq) (resp *pbOffice.GetUserTagsResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req ", req.String())
	resp = &pbOffice.GetUserTagsResp{
		CommonResp: &pbOffice.CommonResp{},
		Tags:       []*pbOffice.Tag{},
	}
	tags, err := db.DB.GetUserTags(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "tags: ", tags)
	for _, v := range tags {
		tag := &pbOffice.Tag{
			TagID:   v.TagID,
			TagName: v.TagName,
		}
		for _, userID := range v.UserList {
			UserName, err := im_mysql_model.GetUserNameByUserID(userID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserNameByUserID failed", err.Error())
				continue
			}
			tag.UserList = append(tag.UserList, &pbOffice.TagUser{
				UserID:   userID,
				UserName: UserName,
			})
		}
		resp.Tags = append(resp.Tags, tag)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp ", resp.String())
	return resp, nil
}

func (s *officeServer) CreateTag(_ context.Context, req *pbOffice.CreateTagReq) (resp *pbOffice.CreateTagResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "CreateTag req", req.String())
	userIDList := utils.RemoveRepeatedStringInList(req.UserIDList)
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "userIDList: ", userIDList)
	resp = &pbOffice.CreateTagResp{CommonResp: &pbOffice.CommonResp{}}
	if err := db.DB.CreateTag(req.UserID, req.TagName, userIDList); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp", resp.String())
	return resp, nil
}

func (s *officeServer) DeleteTag(_ context.Context, req *pbOffice.DeleteTagReq) (resp *pbOffice.DeleteTagResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.DeleteTagResp{CommonResp: &pbOffice.CommonResp{}}
	if err := db.DB.DeleteTag(req.UserID, req.TagID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DeleteTag failed", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) SetTag(_ context.Context, req *pbOffice.SetTagReq) (resp *pbOffice.SetTagResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.SetTagResp{CommonResp: &pbOffice.CommonResp{}}
	IncreaseUserIDList := utils.RemoveRepeatedStringInList(req.IncreaseUserIDList)
	reduceUserIDList := utils.RemoveRepeatedStringInList(req.ReduceUserIDList)
	if err := db.DB.SetTag(req.UserID, req.TagID, req.NewName, IncreaseUserIDList, reduceUserIDList); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetTag failed", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) SendMsg2Tag(_ context.Context, req *pbOffice.SendMsg2TagReq) (resp *pbOffice.SendMsg2TagResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.SendMsg2TagResp{CommonResp: &pbOffice.CommonResp{}}
	var tagUserIDList []string
	for _, tagID := range req.TagList {
		userIDList, err := db.DB.GetUserIDListByTagID(req.SendID, tagID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserIDListByTagID failed", err.Error())
			continue
		}
		tagUserIDList = append(tagUserIDList, userIDList...)
	}
	var groupUserIDList []string
	for _, groupID := range req.GroupList {
		userIDList, err := im_mysql_model.GetGroupMemberIDListByGroupID(groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMemberIDListByGroupID failed", err.Error())
			continue
		}
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), userIDList)
		groupUserIDList = append(groupUserIDList, userIDList...)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), groupUserIDList, req.GroupList)
	var userIDList []string
	userIDList = append(userIDList, tagUserIDList...)
	userIDList = append(userIDList, groupUserIDList...)
	userIDList = append(userIDList, req.UserList...)
	userIDList = utils.RemoveRepeatedStringInList(userIDList)
	for i, userID := range userIDList {
		if userID == req.SendID || userID == "" {
			userIDList = append(userIDList[:i], userIDList[i+1:]...)
		}
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "total userIDList result: ", userIDList)
	us, err := imdb.GetUserByUserID(req.SendID)
	if err != nil {
		log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), req.SendID)
	}
	for _, userID := range userIDList {
		msg.TagSendMessage(req.OperationID, us, userID, req.Content, req.SenderPlatformID)
	}
	var tagSendLogs db.TagSendLog
	for _, userID := range userIDList {
		userName, err := im_mysql_model.GetUserNameByUserID(userID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserNameByUserID failed", err.Error())
			continue
		}
		tagSendLogs.UserList = append(tagSendLogs.UserList, db.TagUser{
			UserID:   userID,
			UserName: userName,
		})
	}
	tagSendLogs.SendID = req.SendID
	tagSendLogs.Content = req.Content
	tagSendLogs.SenderPlatformID = req.SenderPlatformID
	tagSendLogs.SendTime = time.Now().Unix()
	if err := db.DB.SaveTagSendLog(&tagSendLogs); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SaveTagSendLog failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) GetTagSendLogs(_ context.Context, req *pbOffice.GetTagSendLogsReq) (resp *pbOffice.GetTagSendLogsResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.GetTagSendLogsResp{
		CommonResp: &pbOffice.CommonResp{},
		Pagination: &pbCommon.ResponsePagination{
			CurrentPage: req.Pagination.PageNumber,
			ShowNumber:  req.Pagination.ShowNumber,
		},
		TagSendLogs: []*pbOffice.TagSendLog{},
	}
	tagSendLogs, err := db.DB.GetTagSendLogs(req.UserID, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTagSendLogs", err.Error())
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		return resp, nil
	}
	if err := utils.CopyStructFields(&resp.TagSendLogs, tagSendLogs); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) GetUserTagByID(_ context.Context, req *pbOffice.GetUserTagByIDReq) (resp *pbOffice.GetUserTagByIDResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.GetUserTagByIDResp{
		CommonResp: &pbOffice.CommonResp{},
		Tag:        &pbOffice.Tag{},
	}
	tag, err := db.DB.GetTagByID(req.UserID, req.TagID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTagByID failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	for _, userID := range tag.UserList {
		userName, err := im_mysql_model.GetUserNameByUserID(userID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserNameByUserID failed", err.Error())
			continue
		}
		resp.Tag.UserList = append(resp.Tag.UserList, &pbOffice.TagUser{
			UserID:   userID,
			UserName: userName,
		})
	}
	resp.Tag.TagID = tag.TagID
	resp.Tag.TagName = tag.TagName
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) CreateOneWorkMoment(_ context.Context, req *pbOffice.CreateOneWorkMomentReq) (resp *pbOffice.CreateOneWorkMomentResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.CreateOneWorkMomentResp{CommonResp: &pbOffice.CommonResp{}}
	workMoment := db.WorkMoment{
		Comments:     []*db.Comment{},
		LikeUserList: []*db.LikeUser{},
	}
	createUser, err := imdb.GetUserByUserID(req.WorkMoment.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserByUserID", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err := utils.CopyStructFields(&workMoment, req.WorkMoment); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	workMoment.UserName = createUser.Nickname
	workMoment.FaceURL = createUser.FaceURL
	workMoment.PermissionUserIDList = s.getPermissionUserIDList(req.OperationID, req.WorkMoment.PermissionGroupList, req.WorkMoment.PermissionUserList)
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "workMoment to create", workMoment)
	err = db.DB.CreateOneWorkMoment(&workMoment)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CreateOneWorkMoment", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}

	// send notification to at users
	for _, atUser := range req.WorkMoment.AtUserList {
		workMomentNotificationMsg := &pbOffice.WorkMomentNotificationMsg{
			NotificationMsgType: constant.WorkMomentAtUserNotification,
			WorkMomentID:        workMoment.WorkMomentID,
			WorkMomentContent:   workMoment.Content,
			UserID:              workMoment.UserID,
			FaceURL:             createUser.FaceURL,
			UserName:            createUser.Nickname,
			CreateTime:          workMoment.CreateTime,
		}
		msg.WorkMomentSendNotification(req.OperationID, atUser.UserID, workMomentNotificationMsg)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

// count and distinct permission users
func (s *officeServer) getPermissionUserIDList(operationID string, groupList []*pbOffice.PermissionGroup, userList []*pbOffice.WorkMomentUser) []string {
	var permissionUserIDList []string
	for _, group := range groupList {
		GroupMemberIDList, err := imdb.GetGroupMemberIDListByGroupID(group.GroupID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "GetGroupMemberIDListByGroupID failed", group, err.Error())
			continue
		}
		permissionUserIDList = append(permissionUserIDList, GroupMemberIDList...)
	}
	var userIDList []string
	for _, user := range userList {
		userIDList = append(userIDList, user.UserID)
	}
	permissionUserIDList = append(permissionUserIDList, userIDList...)
	permissionUserIDList = utils.RemoveRepeatedStringInList(permissionUserIDList)
	return permissionUserIDList
}

func (s *officeServer) DeleteOneWorkMoment(_ context.Context, req *pbOffice.DeleteOneWorkMomentReq) (resp *pbOffice.DeleteOneWorkMomentResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.DeleteOneWorkMomentResp{CommonResp: &pbOffice.CommonResp{}}
	workMoment, err := db.DB.GetWorkMomentByID(req.WorkMomentID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetWorkMomentByID failed", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "workMoment", workMoment)
	if workMoment.UserID != req.UserID {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "workMoment.UserID != req.WorkMomentID, delete failed", workMoment, req.WorkMomentID)
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}
		return resp, nil
	}
	err = db.DB.DeleteOneWorkMoment(req.WorkMomentID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DeleteOneWorkMoment", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func isUserCanSeeWorkMoment(userID string, workMoment db.WorkMoment) bool {
	if userID != workMoment.UserID {
		switch workMoment.Permission {
		case constant.WorkMomentPublic:
			return true
		case constant.WorkMomentPrivate:
			return false
		case constant.WorkMomentPermissionCanSee:
			return utils.IsContain(userID, workMoment.PermissionUserIDList)
		case constant.WorkMomentPermissionCantSee:
			return !utils.IsContain(userID, workMoment.PermissionUserIDList)
		}
		return false
	}
	return true
}

func (s *officeServer) LikeOneWorkMoment(_ context.Context, req *pbOffice.LikeOneWorkMomentReq) (resp *pbOffice.LikeOneWorkMomentResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.LikeOneWorkMomentResp{CommonResp: &pbOffice.CommonResp{}}
	user, err := imdb.GetUserByUserID(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserNameByUserID failed", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	workMoment, like, err := db.DB.LikeOneWorkMoment(req.UserID, user.Nickname, req.WorkMomentID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "LikeOneWorkMoment failed ", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	workMomentNotificationMsg := &pbOffice.WorkMomentNotificationMsg{
		NotificationMsgType: constant.WorkMomentLikeNotification,
		WorkMomentID:        workMoment.WorkMomentID,
		WorkMomentContent:   workMoment.Content,
		UserID:              user.UserID,
		FaceURL:             user.FaceURL,
		UserName:            user.Nickname,
		CreateTime:          int32(time.Now().Unix()),
	}
	// send notification
	if like && workMoment.UserID != req.UserID {
		msg.WorkMomentSendNotification(req.OperationID, workMoment.UserID, workMomentNotificationMsg)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) CommentOneWorkMoment(_ context.Context, req *pbOffice.CommentOneWorkMomentReq) (resp *pbOffice.CommentOneWorkMomentResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.CommentOneWorkMomentResp{CommonResp: &pbOffice.CommonResp{}}
	commentUser, err := imdb.GetUserByUserID(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserNameByUserID commentUserName failed", req.UserID, err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	var replyUserName string
	if req.ReplyUserID != "" {
		replyUserName, err = imdb.GetUserNameByUserID(req.ReplyUserID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserNameByUserID get replyUserName failed", req.ReplyUserID, err.Error())
			resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
	}
	comment := &db.Comment{
		UserID:        req.UserID,
		UserName:      commentUser.Nickname,
		ReplyUserID:   req.ReplyUserID,
		ReplyUserName: replyUserName,
		Content:       req.Content,
		CreateTime:    int32(time.Now().Unix()),
	}
	workMoment, err := db.DB.CommentOneWorkMoment(comment, req.WorkMomentID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CommentOneWorkMoment failed", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	workMomentNotificationMsg := &pbOffice.WorkMomentNotificationMsg{
		NotificationMsgType: constant.WorkMomentCommentNotification,
		WorkMomentID:        workMoment.WorkMomentID,
		WorkMomentContent:   workMoment.Content,
		UserID:              commentUser.UserID,
		FaceURL:             commentUser.FaceURL,
		UserName:            commentUser.Nickname,
		ReplyUserID:         comment.ReplyUserID,
		ReplyUserName:       comment.ReplyUserName,
		ContentID:           comment.ContentID,
		Content:             comment.Content,
		CreateTime:          comment.CreateTime,
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "msg: ", *workMomentNotificationMsg)
	if req.UserID != workMoment.UserID {
		msg.WorkMomentSendNotification(req.OperationID, workMoment.UserID, workMomentNotificationMsg)
	}
	if req.ReplyUserID != "" && req.ReplyUserID != workMoment.UserID && req.ReplyUserID != req.UserID {
		msg.WorkMomentSendNotification(req.OperationID, req.ReplyUserID, workMomentNotificationMsg)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) GetWorkMomentByID(_ context.Context, req *pbOffice.GetWorkMomentByIDReq) (resp *pbOffice.GetWorkMomentByIDResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.GetWorkMomentByIDResp{
		WorkMoment: &pbOffice.WorkMoment{},
		CommonResp: &pbOffice.CommonResp{},
	}
	workMoment, err := db.DB.GetWorkMomentByID(req.WorkMomentID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetWorkMomentByID failed", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	canSee := isUserCanSeeWorkMoment(req.OpUserID, *workMoment)
	log.Debug(req.OperationID, utils.GetSelfFuncName(), canSee, req.OpUserID, *workMoment)
	if !canSee {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "workMoments not access to user", canSee, workMoment, req.OpUserID)
	}
	if err := utils.CopyStructFields(resp.WorkMoment, workMoment); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) GetUserWorkMoments(_ context.Context, req *pbOffice.GetUserWorkMomentsReq) (resp *pbOffice.GetUserWorkMomentsResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.GetUserWorkMomentsResp{CommonResp: &pbOffice.CommonResp{}, WorkMoments: []*pbOffice.WorkMoment{}}
	resp.Pagination = &pbCommon.ResponsePagination{CurrentPage: req.Pagination.PageNumber, ShowNumber: req.Pagination.ShowNumber}
	var workMoments []db.WorkMoment
	if req.UserID == req.OpUserID {
		workMoments, err = db.DB.GetUserSelfWorkMoments(req.UserID, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	} else {
		workMoments, err = db.DB.GetUserWorkMoments(req.OpUserID, req.UserID, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	}
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserWorkMoments failed", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err := utils.CopyStructFields(&resp.WorkMoments, workMoments); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) GetUserFriendWorkMoments(_ context.Context, req *pbOffice.GetUserFriendWorkMomentsReq) (resp *pbOffice.GetUserFriendWorkMomentsResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.GetUserFriendWorkMomentsResp{CommonResp: &pbOffice.CommonResp{}, WorkMoments: []*pbOffice.WorkMoment{}}
	resp.Pagination = &pbCommon.ResponsePagination{CurrentPage: req.Pagination.PageNumber, ShowNumber: req.Pagination.ShowNumber}
	workMoments, err := db.DB.GetUserFriendWorkMoments(req.Pagination.ShowNumber, req.Pagination.PageNumber, req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserFriendWorkMoments", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err := utils.CopyStructFields(&resp.WorkMoments, workMoments); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) SetUserWorkMomentsLevel(_ context.Context, req *pbOffice.SetUserWorkMomentsLevelReq) (resp *pbOffice.SetUserWorkMomentsLevelResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.SetUserWorkMomentsLevelResp{CommonResp: &pbOffice.CommonResp{}}
	if err := db.DB.SetUserWorkMomentsLevel(req.UserID, req.Level); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetUserWorkMomentsLevel failed", err.Error())
		resp.CommonResp = &pbOffice.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *officeServer) ChangeWorkMomentPermission(_ context.Context, req *pbOffice.ChangeWorkMomentPermissionReq) (resp *pbOffice.ChangeWorkMomentPermissionResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbOffice.ChangeWorkMomentPermissionResp{CommonResp: &pbOffice.CommonResp{}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

package group

import (
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/tools"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"time"
)

func getDBGroupRequest(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (dbGroupRequest *imdb.GroupRequest) {
	dbGroupRequest = &imdb.GroupRequest{}
	utils.CopyStructFields(&dbGroupRequest, req)
	dbGroupRequest.UserID = req.FromUserID
	dbGroupRequest.HandleUserID = tools.OpUserID(ctx)
	dbGroupRequest.HandledTime = time.Now()
	return dbGroupRequest
}

func getDBGroupMember(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (dbGroupMember *imdb.GroupMember) {
	dbGroupMember = &imdb.GroupMember{}
	utils.CopyStructFields(&dbGroupRequest, req)
	dbGroupRequest.UserID = req.FromUserID
	dbGroupRequest.HandleUserID = tools.OpUserID(ctx)
	dbGroupRequest.HandledTime = time.Now()

	member := imdb.GroupMember{}
	member.GroupID = req.GroupID
	member.UserID = req.FromUserID
	member.RoleLevel = constant.GroupOrdinaryUsers
	member.OperatorUserID = tools.OpUserID(ctx)

	member.FaceURL = user.FaceURL
	member.Nickname = user.Nickname
	member.JoinSource = request.JoinSource
	member.InviterUserID = request.InviterUserID
	member.MuteEndTime = time.Unix(int64(time.Now().Second()), 0)

	return dbGroupRequest
}

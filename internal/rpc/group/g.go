package group

import (
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql"
	"Open_IM/pkg/common/tools"
	pbGroup "Open_IM/pkg/proto/group"
	sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"math/big"
	"strconv"
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

func getDBGroupMember(ctx context.Context, groupID, userID string) (dbGroupMember *imdb.GroupMember, err error) {
	dbGroupMember = &imdb.GroupMember{}

	member := imdb.GroupMember{}
	member.GroupID = groupID
	member.UserID = userID
	member.RoleLevel = constant.GroupOrdinaryUsers
	member.OperatorUserID = tools.OpUserID(ctx)

	member.FaceURL = user.FaceURL
	member.Nickname = user.Nickname
	member.JoinSource = request.JoinSource
	member.InviterUserID = request.InviterUserID
	member.MuteEndTime = time.Unix(int64(time.Now().Second()), 0)

	return dbGroupMember, nil
}

func getUsersInfo(ctx context.Context, userIDs []string) ([]*sdk.UserInfo, error) {
	return nil, nil
}

func genGroupID(ctx context.Context, groupID string) string {
	if groupID != "" {
		return groupID
	}
	groupID = utils.Md5(tools.OperationID(ctx) + strconv.FormatInt(time.Now().UnixNano(), 10))
	bi := big.NewInt(0)
	bi.SetString(groupID[0:8], 16)
	groupID = bi.String()
	return groupID
}

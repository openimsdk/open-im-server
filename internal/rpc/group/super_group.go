package group

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbGroup "Open_IM/pkg/proto/group"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetJoinedSuperGroupListResp{CommonResp: &pbGroup.CommonResp{}}
	userToSuperGroup, err := db.DB.GetSuperGroupByUserID(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroupByUserID failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	for _, groupID := range userToSuperGroup.GroupIDList {
		groupInfoDB, err := imdb.GetGroupInfoByGroupID(groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupInfoByGroupID failed", groupID, err.Error())
			continue
		}
		groupInfo := &commonPb.GroupInfo{}
		if err := utils.CopyStructFields(groupInfo, groupInfoDB); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		}
		resp.GroupList = append(resp.GroupList, groupInfo)
	}
	log.NewError(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

package group

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	cp "Open_IM/pkg/common/utils"
	pbGroup "Open_IM/pkg/proto/group"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetJoinedSuperGroupListResp{CommonResp: &pbGroup.CommonResp{}}
	userToSuperGroup, err := db.DB.GetSuperGroupByUserID(req.UserID)
	if err == mongo.ErrNoDocuments {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroupByUserID failed ", err.Error(), req.UserID)
		return resp, nil
	}
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroupByUserID failed ", err.Error(), req.UserID)
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
		group, err := db.DB.GetSuperGroup(groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroup failed", groupID, err.Error())
			continue
		}
		groupInfo.MemberCount = uint32(len(group.MemberIDList))
		resp.GroupList = append(resp.GroupList, groupInfo)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *groupServer) GetSuperGroupsInfo(_ context.Context, req *pbGroup.GetSuperGroupsInfoReq) (resp *pbGroup.GetSuperGroupsInfoResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbGroup.GetSuperGroupsInfoResp{CommonResp: &pbGroup.CommonResp{}}
	groupsInfoList := make([]*commonPb.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromMysql, err := imdb.GetGroupInfoByGroupID(groupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), groupID)
			continue
		}
		var groupInfo commonPb.GroupInfo
		cp.GroupDBCopyOpenIM(&groupInfo, groupInfoFromMysql)
		groupsInfoList = append(groupsInfoList, &groupInfo)
	}
	resp.GroupInfoList = groupsInfoList
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

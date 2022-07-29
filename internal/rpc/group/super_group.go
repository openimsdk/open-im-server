package group

import (
	"Open_IM/pkg/common/constant"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	cp "Open_IM/pkg/common/utils"
	pbGroup "Open_IM/pkg/proto/group"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/go-redis/redis/v8"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetJoinedSuperGroupListResp{CommonResp: &pbGroup.CommonResp{}}
	//userToSuperGroup, err := db.DB.GetSuperGroupByUserID(req.UserID)
	groupIDList, err := rocksCache.GetJoinedSuperGroupListFromCache(req.UserID)
	if err != nil {
		if err == redis.Nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroupByUserID nil ", err.Error(), req.UserID)
			return resp, nil
		}
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroupByUserID failed ", err.Error(), req.UserID)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	for _, groupID := range groupIDList {
		groupInfoFromCache, err := rocksCache.GetGroupInfoFromCache(groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupInfoByGroupID failed", groupID, err.Error())
			continue
		}
		groupInfo := &commonPb.GroupInfo{}
		if err := utils.CopyStructFields(groupInfo, groupInfoFromCache); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		}
		groupMemberIDList, err := rocksCache.GetGroupMemberIDListFromCache(groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSuperGroup failed", groupID, err.Error())
			continue
		}
		groupInfo.MemberCount = uint32(len(groupMemberIDList))
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
		groupInfoFromRedis, err := rocksCache.GetGroupInfoFromCache(groupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), groupID)
			continue
		}
		var groupInfo commonPb.GroupInfo
		cp.GroupDBCopyOpenIM(&groupInfo, groupInfoFromRedis)
		groupsInfoList = append(groupsInfoList, &groupInfo)
	}
	resp.GroupInfoList = groupsInfoList
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

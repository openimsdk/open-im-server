package utils

import (
	"Open_IM/pkg/common/config"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"strings"
	"sync"
)

type GroupMemberUserIDListHash struct {
	MemberListHash uint64
	UserIDList     []string
}

var CacheGroupMemberUserIDList = make(map[string]*GroupMemberUserIDListHash, 0)
var CacheGroupMtx sync.RWMutex

func GetGroupMemberUserIDList(groupID string, operationID string) ([]string, error) {
	groupHashRemote, err := GetGroupMemberUserIDListHashFromRemote(groupID)
	if err != nil {
		CacheGroupMtx.Lock()
		defer CacheGroupMtx.Unlock()
		delete(CacheGroupMemberUserIDList, groupID)
		log.Error(operationID, "GetGroupMemberUserIDListHashFromRemote failed ", err.Error(), groupID)
		return nil, utils.Wrap(err, groupID)
	}

	CacheGroupMtx.Lock()
	defer CacheGroupMtx.Unlock()

	if groupHashRemote == 0 {
		log.Info(operationID, "groupHashRemote == 0  ", groupID)
		delete(CacheGroupMemberUserIDList, groupID)
		return []string{}, nil
	}

	groupInLocalCache, ok := CacheGroupMemberUserIDList[groupID]
	if ok && groupInLocalCache.MemberListHash == groupHashRemote {
		log.Debug(operationID, "in local cache ", groupID)
		return groupInLocalCache.UserIDList, nil
	}
	log.Debug(operationID, "not in local cache or hash changed", groupID, " remote hash ", groupHashRemote, " in cache ", ok)
	memberUserIDListRemote, err := GetGroupMemberUserIDListFromRemote(groupID, operationID)
	if err != nil {
		log.Error(operationID, "GetGroupMemberUserIDListFromRemote failed ", err.Error(), groupID)
		return nil, utils.Wrap(err, groupID)
	}
	CacheGroupMemberUserIDList[groupID] = &GroupMemberUserIDListHash{MemberListHash: groupHashRemote, UserIDList: memberUserIDListRemote}
	return memberUserIDListRemote, nil
}

func GetGroupMemberUserIDListHashFromRemote(groupID string) (uint64, error) {
	return rocksCache.GetGroupMemberListHashFromCache(groupID)
}

func GetGroupMemberUserIDListFromRemote(groupID string, operationID string) ([]string, error) {
	getGroupMemberIDListFromCacheReq := &pbCache.GetGroupMemberIDListFromCacheReq{OperationID: operationID, GroupID: groupID}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, operationID)
	if etcdConn == nil {
		errMsg := operationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(operationID, errMsg)
		return nil, errors.New("errMsg")
	}
	client := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := client.GetGroupMemberIDListFromCache(context.Background(), getGroupMemberIDListFromCacheReq)
	if err != nil {
		log.NewError(operationID, "GetGroupMemberIDListFromCache rpc call failed ", err.Error())
		return nil, utils.Wrap(err, "GetGroupMemberIDListFromCache rpc call failed")
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		errMsg := operationID + "GetGroupMemberIDListFromCache rpc logic call failed " + cacheResp.CommonResp.ErrMsg
		log.NewError(operationID, errMsg)
		return nil, errors.New("errMsg")
	}
	return cacheResp.UserIDList, nil
}

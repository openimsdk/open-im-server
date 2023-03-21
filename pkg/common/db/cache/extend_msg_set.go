package cache

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/dtm-labs/rockscache"
	"time"
)

const (
	extendMsgSetCache = "EXTEND_MSG_SET_CACHE:"
	extendMsgCache    = "EXTEND_MSG_CACHE:"
)

type ExtendMsgSetCache struct {
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func (e *ExtendMsgSetCache) getKey(clientMsgID string) string {
	return extendMsgCache + clientMsgID
}

func (e *ExtendMsgSetCache) GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, firstModifyTime int64) (extendMsg *unrelation.ExtendMsgModel, err error) {
	//getExtendMsg := func() (string, error) {
	//	extendMsg, err := db.DB.GetExtendMsg(sourceID, sessionType, clientMsgID, firstModifyTime)
	//	if err != nil {
	//		return "", utils.Wrap(err, "GetExtendMsgList failed")
	//	}
	//	bytes, err := json.Marshal(extendMsg)
	//	if err != nil {
	//		return "", utils.Wrap(err, "Marshal failed")
	//	}
	//	return string(bytes), nil
	//}
	//defer func() {
	//	mcontext.SetCtxDebug(ctx, utils.GetFuncName(1), err, "sourceID", sourceID, "sessionType",
	//		sessionType, "clientMsgID", clientMsgID, "firstModifyTime", firstModifyTime, "extendMsg", extendMsg)
	//}()
	//extendMsgStr, err := db.DB.Rc.Fetch(extendMsgCache+clientMsgID, time.Second*30*60, getExtendMsg)
	//if err != nil {
	//	return nil, utils.Wrap(err, "Fetch failed")
	//}
	//extendMsg = &mongoDB.ExtendMsg{}
	//err = json.Unmarshal([]byte(extendMsgStr), extendMsg)
	//return extendMsg, utils.Wrap(err, "Unmarshal failed")
	return GetCache(ctx, e.rcClient, e.getKey(clientMsgID), e.expireTime, func(ctx context.Context) (*unrelation.ExtendMsgModel, error) {
		panic("")
	})

}

func (e *ExtendMsgSetCache) DelExtendMsg(ctx context.Context, clientMsgID string) (err error) {
	return utils.Wrap(e.rcClient.TagAsDeleted(e.getKey(clientMsgID)), "DelExtendMsg err")
}

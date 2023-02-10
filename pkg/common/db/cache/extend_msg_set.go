package cache

import (
	"Open_IM/pkg/common/db/table/unrelation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
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

func (e *ExtendMsgSetCache) GetExtendMsg(ctx context.Context, sourceID string, sessionType int32, clientMsgID string, firstModifyTime int64) (extendMsg *unrelation.ExtendMsg, err error) {
	getExtendMsg := func() (string, error) {
		extendMsg, err := db.DB.GetExtendMsg(sourceID, sessionType, clientMsgID, firstModifyTime)
		if err != nil {
			return "", utils.Wrap(err, "GetExtendMsgList failed")
		}
		bytes, err := json.Marshal(extendMsg)
		if err != nil {
			return "", utils.Wrap(err, "Marshal failed")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "sourceID", sourceID, "sessionType",
			sessionType, "clientMsgID", clientMsgID, "firstModifyTime", firstModifyTime, "extendMsg", extendMsg)
	}()
	extendMsgStr, err := db.DB.Rc.Fetch(extendMsgCache+clientMsgID, time.Second*30*60, getExtendMsg)
	if err != nil {
		return nil, utils.Wrap(err, "Fetch failed")
	}
	extendMsg = &mongoDB.ExtendMsg{}
	err = json.Unmarshal([]byte(extendMsgStr), extendMsg)
	return extendMsg, utils.Wrap(err, "Unmarshal failed")
}

func (e *ExtendMsgSetCache) DelExtendMsg(ctx context.Context, clientMsgID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "clientMsgID", clientMsgID)
	}()
	return utils.Wrap(db.DB.Rc.TagAsDeleted(extendMsgCache+clientMsgID), "DelExtendMsg err")
}

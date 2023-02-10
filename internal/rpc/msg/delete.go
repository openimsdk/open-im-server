package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/proto/msg"
	common "Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"time"
)

func (rpc *rpcChat) DelMsgList(_ context.Context, req *common.DelMsgListReq) (*common.DelMsgListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &common.DelMsgListResp{}
	select {
	case rpc.delMsgCh <- deleteMsg{
		UserID:      req.UserID,
		OpUserID:    req.OpUserID,
		SeqList:     req.SeqList,
		OperationID: req.OperationID,
	}:
	case <-time.After(1 * time.Second):
		resp.ErrCode = constant.ErrSendLimit.ErrCode
		resp.ErrMsg = constant.ErrSendLimit.ErrMsg
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}
func (rpc *rpcChat) DelSuperGroupMsg(ctx context.Context, req *msg.DelSuperGroupMsgReq) (*msg.DelSuperGroupMsgResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	if !tokenverify.CheckAccess(ctx, req.OpUserID, req.UserID) {
		log.NewError(req.OperationID, "CheckAccess false ", req.OpUserID, req.UserID)
		return &msg.DelSuperGroupMsgResp{ErrCode: constant.ErrNoPermission.ErrCode, ErrMsg: constant.ErrNoPermission.ErrMsg}, nil
	}
	resp := &msg.DelSuperGroupMsgResp{}
	groupMaxSeq, err := db.DB.GetGroupMaxSeq(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMaxSeq false ", req.OpUserID, req.UserID, req.GroupID)
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
		return resp, nil
	}
	err = db.DB.SetGroupUserMinSeq(req.GroupID, req.UserID, groupMaxSeq)
	if err != nil {
		log.NewError(req.OperationID, "SetGroupUserMinSeq false ", req.OpUserID, req.UserID, req.GroupID)
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (rpc *rpcChat) ClearMsg(_ context.Context, req *pbChat.ClearMsgReq) (*pbChat.ClearMsgResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc req: ", req.String())
	if req.OpUserID != req.UserID && !tokenverify.IsManagerUserID(req.UserID) {
		errMsg := "No permission" + req.OpUserID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrNoPermission.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CleanUpOneUserAllMsgFromRedis args", req.UserID)
	err := db.DB.CleanUpOneUserAllMsgFromRedis(req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpOneUserAllMsgFromRedis failed " + err.Error() + req.OperationID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDatabase.ErrCode, ErrMsg: errMsg}, nil
	}
	log.Debug(req.OperationID, "CleanUpUserMsgFromMongo args", req.UserID)
	err = db.DB.CleanUpUserMsgFromMongo(req.UserID, req.OperationID)
	if err != nil {
		errMsg := "CleanUpUserMsgFromMongo failed " + err.Error() + req.OperationID + req.UserID
		log.Error(req.OperationID, errMsg)
		return &pbChat.ClearMsgResp{ErrCode: constant.ErrDatabase.ErrCode, ErrMsg: errMsg}, nil
	}

	resp := pbChat.ClearMsgResp{ErrCode: 0}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return &resp, nil
}

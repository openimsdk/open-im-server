package user

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbUser "Open_IM/pkg/proto/user"
	"context"
)

func (s *userServer) SetReceiveMessageOpt(ctx context.Context, req *pbUser.SetReceiveMessageOptReq) (*pbUser.SetReceiveMessageOptResp, error) {
	m := make(map[string]int, len(req.ConversationId))
	for _, v := range req.ConversationId {
		m[v] = int(req.Opt)
	}
	err := db.DB.SetMultiConversationMsgOpt(req.UId, m)
	if err != nil {
		log.NewError(req.OperationID, "SetMultiConversationMsgOpt failed ", err.Error(), req)
		return &pbUser.SetReceiveMessageOptResp{ErrCode: constant.DatabaseError, ErrMsg: err.Error()}, nil
	}
	var resp pbUser.SetReceiveMessageOptResp
	resp.ErrCode = 0

	for _, v := range req.ConversationId {
		resp.OptResult = append(resp.OptResult, &pbUser.OptResult{ConversationId: v, Result: 0})
	}
	log.NewInfo(req.OperationID, "SetReceiveMessageOpt req, resp ", req, resp)
	return &resp, nil
}

func (s *userServer) GetReceiveMessageOpt(ctx context.Context, req *pbUser.GetReceiveMessageOptReq) (*pbUser.GetReceiveMessageOptResp, error) {
	m, err := db.DB.GetMultiConversationMsgOpt(req.UId, req.ConversationId)
	if err != nil {
		log.NewError(req.OperationID, "GetMultiConversationMsgOpt failed ", err.Error(), req)
		return &pbUser.GetReceiveMessageOptResp{ErrCode: constant.DatabaseError, ErrMsg: err.Error()}, nil
	}
	var resp pbUser.GetReceiveMessageOptResp
	resp.ErrCode = 0
	for k, v := range m {
		resp.ConversationOptResult = append(resp.ConversationOptResult, &pbUser.OptResult{ConversationId: k, Result: int32(v)})
	}
	log.NewInfo(req.OperationID, "GetReceiveMessageOpt, req, resp", req, resp)
	return &resp, nil
}

func (s *userServer) GetAllConversationMsgOpt(ctx context.Context, req *pbUser.GetAllConversationMsgOptReq) (*pbUser.GetAllConversationMsgOptResp, error) {
	m, err := db.DB.GetAllConversationMsgOpt(req.UId)
	if err != nil {
		log.NewError(req.OperationID, "GetAllConversationMsgOpt failed ", err.Error(), req)
		return &pbUser.GetAllConversationMsgOptResp{ErrCode: constant.DatabaseError, ErrMsg: err.Error()}, nil
	}
	var resp pbUser.GetAllConversationMsgOptResp
	resp.ErrCode = 0
	for k, v := range m {
		resp.ConversationOptResult = append(resp.ConversationOptResult, &pbUser.OptResult{ConversationId: k, Result: int32(v)})
	}
	log.NewInfo(req.OperationID, "GetAllConversationMsgOpt, req, resp", req, resp)
	return &resp, nil
}

package group

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
)

func (s *groupServer) GetJoinedSuperGroupList(ctx context.Context, req *pbGroup.GetJoinedSuperGroupListReq) (*pbGroup.GetJoinedSuperGroupListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetJoinedSuperGroupListResp{}
	_, err := db.DB.GetSuperGroupByUserID(req.UserID)
	if err != nil {
		return resp, nil
	}
	log.NewError(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

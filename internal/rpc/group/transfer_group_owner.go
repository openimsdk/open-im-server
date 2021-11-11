package group

import (
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/group"
	"context"
)

func (s *groupServer) TransferGroupOwner(_ context.Context, pb *group.TransferGroupOwnerReq) (*group.TransferGroupOwnerResp, error) {
	log.Info("", "", "rpc TransferGroupOwner call start..., [pb: %s]", pb.String())

	reply, err := im_mysql_model.TransferGroupOwner(pb)
	if err != nil {
		log.Error("", "", "rpc TransferGroupOwner call..., im_mysql_model.TransferGroupOwner fail [pb: %s] [err: %s]", pb.String(), err.Error())
		return nil, err
	}
	log.Info("", "", "rpc TransferGroupOwner call..., im_mysql_model.TransferGroupOwner")

	return reply, nil
}

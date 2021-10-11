package group

import (
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/group"
	"context"
)

func (s *groupServer) GetGroupApplicationList(_ context.Context, pb *group.GetGroupApplicationListReq) (*group.GetGroupApplicationListResp, error) {
	log.Info("", "", "rpc GetGroupApplicationList call start..., [pb: %s]", pb.String())

	reply, err := im_mysql_model.GetGroupApplicationList(pb.UID)
	if err != nil {
		log.Error("", "", "rpc GetGroupApplicationList call..., im_mysql_model.GetGroupApplicationList fail [uid: %s] [err: %s]", pb.UID, err.Error())
		return &group.GetGroupApplicationListResp{ErrCode: 701, ErrMsg: "GetGroupApplicationList failed"}, nil
	}
	log.Info("", "", "rpc GetGroupApplicationList call..., im_mysql_model.GetGroupApplicationList")

	return reply, nil
}

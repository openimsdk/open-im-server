package group

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/group"
	"context"
)

func (s *groupServer) GroupApplicationResponse(_ context.Context, pb *group.GroupApplicationResponseReq) (*group.GroupApplicationResponseResp, error) {
	log.Info("", "", "rpc GroupApplicationResponse call start..., [pb: %s]", pb.String())

	reply, err := im_mysql_model.GroupApplicationResponse(pb)
	if err != nil {
		log.Error("", "", "rpc GroupApplicationResponse call..., im_mysql_model.GroupApplicationResponse fail [pb: %s] [err: %s]", pb.String(), err.Error())
		return &group.GroupApplicationResponseResp{ErrCode: 702, ErrMsg: "rpc GroupApplicationResponse failed"}, nil
	}

	if pb.HandleResult == 1 {
		if pb.ToUserID == "0" {
			err = db.DB.AddGroupMember(pb.GroupID, pb.FromUserID)
			if err != nil {
				log.Error("", "", "rpc GroupApplicationResponse call..., db.DB.AddGroupMember fail [pb: %s] [err: %s]", pb.String(), err.Error())
				return nil, err
			}
		} else {
			err = db.DB.AddGroupMember(pb.GroupID, pb.ToUserID)
			if err != nil {
				log.Error("", "", "rpc GroupApplicationResponse call..., db.DB.AddGroupMember fail [pb: %s] [err: %s]", pb.String(), err.Error())
				return nil, err
			}
		}
	}

	log.Info("", "", "rpc GroupApplicationResponse call..., im_mysql_model.GroupApplicationResponse")

	return reply, nil
}

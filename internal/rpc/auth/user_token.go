package auth

import (
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbAuth "Open_IM/src/proto/auth"
	"Open_IM/src/utils"
	"context"
)

func (rpc *rpcAuth) UserToken(_ context.Context, pb *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	log.Info("", "", "rpc user_token call start..., [pbTokenReq: %s]", pb.String())

	_, err := im_mysql_model.FindUserByUID(pb.UID)
	if err != nil {
		log.Error("", "", "rpc user_token call..., im_mysql_model.AppServerFindFromUserByUserID fail [uid: %s] [err: %s]", pb.UID, err.Error())
		return &pbAuth.UserTokenResp{ErrCode: 500, ErrMsg: err.Error()}, err
	}
	log.Info("", "", "rpc user_token call..., im_mysql_model.AppServerFindFromUserByUserID")

	tokens, expTime, err := utils.CreateToken(pb.UID, "", pb.Platform)
	if err != nil {
		log.Error("", "", "rpc user_token call..., utils.CreateToken fail [uid: %s] [err: %s]", pb.UID, err.Error())
		return &pbAuth.UserTokenResp{ErrCode: 500, ErrMsg: err.Error()}, err
	}
	log.Info("", "", "rpc user_token success return, [uid: %s] [tokens: %s]", pb.UID, tokens)

	return &pbAuth.UserTokenResp{Token: tokens, ExpiredTime: expTime}, nil
}

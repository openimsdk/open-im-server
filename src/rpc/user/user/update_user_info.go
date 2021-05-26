package user

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	log2 "Open_IM/src/common/log"
	pbUser "Open_IM/src/proto/user"
	"Open_IM/src/utils"
	"context"
)

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbUser.UpdateUserInfoReq) (*pbUser.CommonResp, error) {
	log2.Info(req.Token, req.OperationID, "rpc modify user is server,args=%s", req.String())
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log2.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbUser.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: err.Error()}, nil
	}
	err = im_mysql_model.UpDateUserInfo(claims.UID, req.Name, req.Icon, req.Mobile, req.Birth, req.Email, req.Ex, req.Gender)
	if err != nil {
		log2.Error(req.Token, req.OperationID, "update user some attribute failed,err=%s", err.Error())
		return &pbUser.CommonResp{ErrorCode: config.ErrModifyUserInfo.ErrCode, ErrorMsg: config.ErrModifyUserInfo.ErrMsg}, nil
	}
	return &pbUser.CommonResp{}, nil
}

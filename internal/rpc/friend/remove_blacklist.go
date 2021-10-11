package friend

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/src/utils"
	"context"
)

func (s *friendServer) RemoveBlacklist(ctx context.Context, req *pbFriend.RemoveBlacklistReq) (*pbFriend.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc remove blacklist is server,userid=%s", req.Uid)
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	err = im_mysql_model.RemoveBlackList(claims.UID, req.Uid)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,remove blacklist failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrMysql.ErrCode, ErrorMsg: config.ErrMysql.ErrMsg}, nil
	}
	log.Info(req.Token, req.OperationID, "rpc remove blacklist success return,userid=%s", req.Uid)
	return &pbFriend.CommonResp{}, nil
}

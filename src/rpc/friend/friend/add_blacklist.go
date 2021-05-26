package friend

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbFriend "Open_IM/src/proto/friend"
	"Open_IM/src/utils"
	"context"
)

func (s *friendServer) AddBlacklist(ctx context.Context, req *pbFriend.AddBlacklistReq) (*pbFriend.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc add blacklist is server,args=%s", req.String())
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	} else {
		err := im_mysql_model.InsertInToUserBlackList(claims.UID, req.Uid)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s,Failed to add blacklist", err.Error())
			return &pbFriend.CommonResp{ErrorCode: config.ErrMysql.ErrCode, ErrorMsg: config.ErrMysql.ErrMsg}, nil
		}
	}
	log.Info(req.Token, req.OperationID, "rpc add blacklist success return,uid=%s", req.Uid)
	return &pbFriend.CommonResp{}, nil
}

package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
)

func (s *friendServer) AddBlacklist(ctx context.Context, req *pbFriend.AddBlacklistReq) (*pbFriend.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc add blacklist is server,args=%s", req.String())
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}

	isMagagerFlag := 0
	tokenUid := claims.UID

	if utils.IsContain(tokenUid, config.Config.Manager.AppManagerUid) {
		isMagagerFlag = 1
	}

	if isMagagerFlag == 0 {
		err = im_mysql_model.InsertInToUserBlackList(claims.UID, req.Uid)
		if err != nil {
			log.Error(req.Token, req.OperationID, "err=%s,Failed to add blacklist", err.Error())
			return &pbFriend.CommonResp{ErrorCode: config.ErrMysql.ErrCode, ErrorMsg: config.ErrMysql.ErrMsg}, nil
		}
		log.Info(req.Token, req.OperationID, "rpc add blacklist success return,uid=%s", req.Uid)
		return &pbFriend.CommonResp{}, nil
	}

	err = im_mysql_model.InsertInToUserBlackList(req.OwnerUid, req.Uid)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,Failed to add blacklist", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrMysql.ErrCode, ErrorMsg: config.ErrMysql.ErrMsg}, nil
	}
	log.Info(req.Token, req.OperationID, "rpc add blacklist success return,uid=%s", req.Uid)
	return &pbFriend.CommonResp{}, nil
}

package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
)

func (s *friendServer) SetFriendComment(ctx context.Context, req *pbFriend.SetFriendCommentReq) (*pbFriend.CommonResp, error) {
	log.Info(req.Token, req.OperationID, "rpc set friend comment is server,params=%s", req.String())
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	err = im_mysql_model.UpdateFriendComment(claims.UID, req.Uid, req.Comment)
	if err != nil {
		log.Error(req.Token, req.OperationID, "set friend comment failed,err=%s", err.Error())
		return &pbFriend.CommonResp{ErrorCode: config.ErrSetFriendComment.ErrCode, ErrorMsg: config.ErrSetFriendComment.ErrMsg}, nil
	}
	log.Info(req.Token, req.OperationID, "rpc set friend comment is success return")
	return &pbFriend.CommonResp{}, nil
}

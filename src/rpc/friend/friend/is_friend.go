package friend

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	"Open_IM/src/common/db/mysql_model/im_mysql_model"
	"Open_IM/src/common/log"
	pbFriend "Open_IM/src/proto/friend"
	"Open_IM/src/utils"
	"context"
	"fmt"
)

func (s *friendServer) IsFriend(ctx context.Context, req *pbFriend.IsFriendReq) (*pbFriend.IsFriendResp, error) {
	log.InfoByArgs("rpc is friend is server,args=%s", req.String())
	var isFriend int32
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.IsFriendResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	_, err = im_mysql_model.FindFriendRelationshipFromFriend(claims.UID, req.ReceiveUid)
	if err == nil {
		isFriend = constant.FriendFlag
	} else {
		isFriend = constant.ApplicationFriendFlag
	}
	log.InfoByArgs(fmt.Sprintf("rpc is friend success return"))
	return &pbFriend.IsFriendResp{ShipType: isFriend}, nil
}

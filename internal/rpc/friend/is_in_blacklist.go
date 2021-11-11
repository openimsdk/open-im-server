package friend

import (
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"
	"fmt"
)

func (s *friendServer) IsInBlackList(ctx context.Context, req *pbFriend.IsInBlackListReq) (*pbFriend.IsInBlackListResp, error) {
	log.InfoByArgs("rpc is in blacklist is server,args=%s", req.String())
	var isInBlacklist = false
	err := im_mysql_model.FindRelationshipFromBlackList(req.ReceiveUid, req.SendUid)
	if err == nil {
		isInBlacklist = true
	}
	log.InfoByArgs(fmt.Sprintf("rpc is in blackList success return"))
	return &pbFriend.IsInBlackListResp{Response: isInBlacklist}, nil
}

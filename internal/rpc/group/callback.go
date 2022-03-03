package group

import (
	pbGroup "Open_IM/pkg/proto/group"
)

func callbackBeforeCreateGroup(req *pbGroup.CreateGroupReq) (bool, error) {
	return true, nil
}

/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/15 10:28).
 */
package user

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
)

func (s *userServer) DeleteUsers(_ context.Context, req *pbUser.DeleteUsersReq) (*pbUser.DeleteUsersResp, error) {
	log.InfoByKv("rpc DeleteUsers arrived server", req.OperationID, "args", req.String())
	var resp pbUser.DeleteUsersResp
	var common pbUser.CommonResp
	c, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.ErrorByKv("parse token failed", req.OperationID, "err", err.Error())
		return &pbUser.DeleteUsersResp{CommonResp: &pbUser.CommonResp{ErrorCode: constant.ErrParseToken.ErrCode, ErrorMsg: err.Error()}, FailedUidList: req.DeleteUidList}, nil
	}
	if !utils.IsContain(c.UID, config.Config.Manager.AppManagerUid) {
		log.ErrorByKv(" Authentication failed", req.OperationID, "args", c)
		return &pbUser.DeleteUsersResp{CommonResp: &pbUser.CommonResp{ErrorCode: 401, ErrorMsg: "not authorized"}, FailedUidList: req.DeleteUidList}, nil
	}
	for _, uid := range req.DeleteUidList {
		err = im_mysql_model.UserDelete(uid)
		if err != nil {
			common.ErrorCode = 201
			common.ErrorMsg = "some uid deleted failed"
			resp.FailedUidList = append(resp.FailedUidList, uid)
		}
	}
	resp.CommonResp = &common
	return &resp, nil

}

func (s *userServer) GetAllUsersUid(_ context.Context, req *pbUser.GetAllUsersUidReq) (*pbUser.GetAllUsersUidResp, error) {
	log.InfoByKv("rpc GetAllUsersUid arrived server", req.OperationID, "args", req.String())
	c, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.InfoByKv("parse token failed", req.OperationID, "err", err.Error())
		return &pbUser.GetAllUsersUidResp{CommonResp: &pbUser.CommonResp{ErrorCode: constant.ErrParseToken.ErrCode, ErrorMsg: err.Error()}}, nil
	}
	if !utils.IsContain(c.UID, config.Config.Manager.AppManagerUid) {
		log.ErrorByKv(" Authentication failed", req.OperationID, "args", c)
		return &pbUser.GetAllUsersUidResp{CommonResp: &pbUser.CommonResp{ErrorCode: 401, ErrorMsg: "not authorized"}}, nil
	}
	uidList, err := im_mysql_model.SelectAllUID()
	if err != nil {
		log.ErrorByKv("db get failed", req.OperationID, "err", err.Error())
		return &pbUser.GetAllUsersUidResp{CommonResp: &pbUser.CommonResp{ErrorCode: constant.ErrMysql.ErrCode, ErrorMsg: err.Error()}}, nil
	} else {
		return &pbUser.GetAllUsersUidResp{CommonResp: &pbUser.CommonResp{ErrorCode: 0, ErrorMsg: ""}, UidList: uidList}, nil
	}

}
func (s *userServer) AccountCheck(_ context.Context, req *pbUser.AccountCheckReq) (*pbUser.AccountCheckResp, error) {
	log.InfoByKv("rpc AccountCheck arrived server", req.OperationID, "args", req.String())
	c, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.InfoByKv("parse token failed", req.OperationID, "err", err.Error())
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrorCode: constant.ErrParseToken.ErrCode, ErrorMsg: err.Error()}}, nil
	}
	if !utils.IsContain(c.UID, config.Config.Manager.AppManagerUid) {
		log.ErrorByKv(" Authentication failed", req.OperationID, "args", c)
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrorCode: 401, ErrorMsg: "not authorized"}}, nil
	}
	uidList, err := im_mysql_model.SelectSomeUID(req.UidList)
	if err != nil {
		log.ErrorByKv("db get SelectSomeUID failed", req.OperationID, "err", err.Error())
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrorCode: constant.ErrMysql.ErrCode, ErrorMsg: err.Error()}}, nil
	} else {
		var r []*pbUser.AccountCheckResp_SingleUserStatus
		for _, v := range req.UidList {
			temp := new(pbUser.AccountCheckResp_SingleUserStatus)
			temp.UserID = v
			if utils.IsContain(v, uidList) {
				temp.AccountStatus = constant.Registered
			} else {
				temp.AccountStatus = constant.UnRegistered
			}
			r = append(r, temp)
		}
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrorCode: 0, ErrorMsg: ""}, Result: r}, nil
	}

}

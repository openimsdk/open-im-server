package admin_cms

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdminCMS "Open_IM/pkg/proto/admin_cms"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type adminCMSServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewAdminCMSServer(port int) *adminCMSServer {
	log.NewPrivateLog(constant.LogFileName)
	return &adminCMSServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAdminCMSName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *adminCMSServer) Run() {
	log.NewInfo("0", "AdminCMS rpc start ")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

	//listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("0", "listen network success, ", address, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbAdminCMS.RegisterAdminCMSServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		panic(utils.Wrap(err, "register admin module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *adminCMSServer) AdminLogin(_ context.Context, req *pbAdminCMS.AdminLoginReq) (*pbAdminCMS.AdminLoginResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.AdminLoginResp{CommonResp: &pbAdminCMS.CommonResp{}}
	for i, adminID := range config.Config.Manager.AppManagerUid {
		if adminID == req.AdminID && config.Config.Manager.Secrets[i] == req.Secret {
			token, expTime, err := token_verify.CreateToken(adminID, constant.SingleChatType)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "generate token failed", "adminID: ", adminID, err.Error())
				resp.CommonResp.ErrCode = constant.ErrTokenUnknown.ErrCode
				resp.CommonResp.ErrMsg = err.Error()
				return resp, nil
			}
			log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "generate token success", "token: ", token, "expTime:", expTime)
			resp.Token = token
			break
		}
	}
	if resp.Token == "" {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "failed")
		resp.CommonResp.ErrCode = constant.ErrTokenUnknown.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrTokenMalformed.ErrMsg
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *adminCMSServer) AddUserRegisterAddFriendIDList(_ context.Context, req *pbAdminCMS.AddUserRegisterAddFriendIDListReq) (*pbAdminCMS.AddUserRegisterAddFriendIDListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.AddUserRegisterAddFriendIDListResp{CommonResp: &pbAdminCMS.CommonResp{}}
	if err := imdb.AddUserRegisterAddFriendIDList(req.UserIDList...); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserIDList)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", req.String())
	return resp, nil
}

func (s *adminCMSServer) ReduceUserRegisterAddFriendIDList(_ context.Context, req *pbAdminCMS.ReduceUserRegisterAddFriendIDListReq) (*pbAdminCMS.ReduceUserRegisterAddFriendIDListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.ReduceUserRegisterAddFriendIDListResp{CommonResp: &pbAdminCMS.CommonResp{}}
	if req.Operation == 0 {
		if err := imdb.ReduceUserRegisterAddFriendIDList(req.UserIDList...); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserIDList)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
	} else {
		if err := imdb.DeleteAllRegisterAddFriendIDList(); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserIDList)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", req.String())
	return resp, nil
}

func (s *adminCMSServer) GetUserRegisterAddFriendIDList(_ context.Context, req *pbAdminCMS.GetUserRegisterAddFriendIDListReq) (*pbAdminCMS.GetUserRegisterAddFriendIDListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.GetUserRegisterAddFriendIDListResp{CommonResp: &pbAdminCMS.CommonResp{}}
	userIDList, err := imdb.GetRegisterAddFriendList(req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	userList, err := imdb.GetUsersByUserIDList(userIDList)
	if err != nil {
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), userList, userIDList)
	resp.Pagination = &server_api_params.ResponsePagination{
		CurrentPage: req.Pagination.PageNumber,
		ShowNumber:  req.Pagination.ShowNumber,
	}
	utils.CopyStructFields(&resp.UserInfoList, userList)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", req.String())
	return resp, nil
}

func (s *adminCMSServer) GetChatLogs(_ context.Context, req *pbAdminCMS.GetChatLogsReq) (*pbAdminCMS.GetChatLogsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetChatLogs", req.String())
	resp := &pbAdminCMS.GetChatLogsResp{CommonResp: &pbAdminCMS.CommonResp{}, Pagination: &server_api_params.ResponsePagination{}}
	time, err := utils.TimeStringToTime(req.SendTime)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "time string parse error", err.Error())
		resp.CommonResp.ErrCode = constant.ErrArgs.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	chatLog := db.ChatLog{
		Content:     req.Content,
		SendTime:    time,
		ContentType: req.ContentType,
		SessionType: req.SessionType,
		RecvID:      req.RecvID,
		SendID:      req.SendID,
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "chat_log: ", chatLog)
	nums, err := imdb.GetChatLogCount(chatLog)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetChatLogCount", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.ChatLogsNum = int32(nums)
	chatLogs, err := imdb.GetChatLog(chatLog, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetChatLog", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	for _, chatLog := range chatLogs {
		pbChatLog := &pbAdminCMS.ChatLog{}
		utils.CopyStructFields(pbChatLog, chatLog)
		pbChatLog.SendTime = chatLog.SendTime.Unix()
		pbChatLog.CreateTime = chatLog.CreateTime.Unix()
		if chatLog.SenderNickname == "" {
			sendUser, err := imdb.GetUserByUserID(chatLog.SendID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserByUserID failed", err.Error())
				continue
			}
			pbChatLog.SenderNickname = sendUser.Nickname
		}
		switch chatLog.SessionType {
		case constant.SingleChatType:
			recvUser, err := imdb.GetUserByUserID(chatLog.RecvID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserByUserID failed", err.Error())
				continue
			}
			pbChatLog.SenderNickname = recvUser.Nickname

		case constant.GroupChatType, constant.SuperGroupChatType:
			group, err := imdb.GetGroupInfoByGroupID(chatLog.RecvID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupById failed")
				continue
			}
			pbChatLog.RecvID = group.GroupID
			pbChatLog.GroupName = group.GroupName
		}
		resp.ChatLogs = append(resp.ChatLogs, pbChatLog)
	}
	resp.Pagination = &server_api_params.ResponsePagination{
		CurrentPage: req.Pagination.PageNumber,
		ShowNumber:  req.Pagination.ShowNumber,
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp output: ", resp.String())
	return resp, nil
}

func (s *adminCMSServer) GetActiveGroup(_ context.Context, req *pbAdminCMS.GetActiveGroupReq) (*pbAdminCMS.GetActiveGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req", req.String())
	resp := &pbAdminCMS.GetActiveGroupResp{CommonResp: &pbAdminCMS.CommonResp{}}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ParseTimeFromTo failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "time: ", fromTime, toTime)
	activeGroups, err := imdb.GetActiveGroups(fromTime, toTime, 12)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetActiveGroups failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	for _, activeGroup := range activeGroups {
		resp.Groups = append(resp.Groups,
			&pbAdminCMS.GroupResp{
				GroupName:  activeGroup.Name,
				GroupId:    activeGroup.Id,
				MessageNum: int32(activeGroup.MessageNum),
			})
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}

func (s *adminCMSServer) GetActiveUser(_ context.Context, req *pbAdminCMS.GetActiveUserReq) (*pbAdminCMS.GetActiveUserResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.String())
	resp := &pbAdminCMS.GetActiveUserResp{CommonResp: &pbAdminCMS.CommonResp{}}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ParseTimeFromTo failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "time: ", fromTime, toTime)
	activeUsers, err := imdb.GetActiveUsers(fromTime, toTime, 12)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetActiveUsers failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	for _, activeUser := range activeUsers {
		resp.Users = append(resp.Users,
			&pbAdminCMS.UserResp{
				UserID:     activeUser.ID,
				NickName:   activeUser.Name,
				MessageNum: int32(activeUser.MessageNum),
			},
		)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}

func ParseTimeFromTo(from, to string) (time.Time, time.Time, error) {
	var fromTime time.Time
	var toTime time.Time
	fromTime, err := utils.TimeStringToTime(from)
	if err != nil {
		return fromTime, toTime, err
	}
	toTime, err = utils.TimeStringToTime(to)
	if err != nil {
		return fromTime, toTime, err
	}
	return fromTime, toTime, nil
}

func isInOneMonth(from, to time.Time) bool {
	return from.Month() == to.Month() && from.Year() == to.Year()
}

func GetRangeDate(from, to time.Time) [][2]time.Time {
	interval := to.Sub(from)
	var times [][2]time.Time
	switch {
	// today
	case interval == 0:
		times = append(times, [2]time.Time{
			from, from.Add(time.Hour * 24),
		})
	// days
	case isInOneMonth(from, to):
		for i := 0; ; i++ {
			fromTime := from.Add(time.Hour * 24 * time.Duration(i))
			toTime := from.Add(time.Hour * 24 * time.Duration(i+1))
			if toTime.After(to.Add(time.Hour * 24)) {
				break
			}
			times = append(times, [2]time.Time{
				fromTime, toTime,
			})
		}
	// month
	case !isInOneMonth(from, to):
		if to.Sub(from) < time.Hour*24*30 {
			for i := 0; ; i++ {
				fromTime := from.Add(time.Hour * 24 * time.Duration(i))
				toTime := from.Add(time.Hour * 24 * time.Duration(i+1))
				if toTime.After(to.Add(time.Hour * 24)) {
					break
				}
				times = append(times, [2]time.Time{
					fromTime, toTime,
				})
			}
		} else {
			for i := 0; ; i++ {
				if i == 0 {
					fromTime := from
					toTime := getFirstDateOfNextNMonth(fromTime, 1)
					times = append(times, [2]time.Time{
						fromTime, toTime,
					})
				} else {
					fromTime := getFirstDateOfNextNMonth(from, i)
					toTime := getFirstDateOfNextNMonth(fromTime, 1)
					if toTime.After(to) {
						toTime = to
						times = append(times, [2]time.Time{
							fromTime, toTime,
						})
						break
					}
					times = append(times, [2]time.Time{
						fromTime, toTime,
					})
				}

			}
		}
	}
	return times
}

func getFirstDateOfNextNMonth(currentTime time.Time, n int) time.Time {
	lastOfMonth := time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, currentTime.Location()).AddDate(0, n, 0)
	return lastOfMonth
}

func (s *adminCMSServer) GetGroupStatistics(_ context.Context, req *pbAdminCMS.GetGroupStatisticsReq) (*pbAdminCMS.GetGroupStatisticsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.String())
	resp := &pbAdminCMS.GetGroupStatisticsResp{CommonResp: &pbAdminCMS.CommonResp{}}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupStatistics failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	increaseGroupNum, err := imdb.GetIncreaseGroupNum(fromTime, toTime.Add(time.Hour*24))

	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseGroupNum failed", err.Error(), fromTime, toTime)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	totalGroupNum, err := imdb.GetTotalGroupNum()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.IncreaseGroupNum = increaseGroupNum
	resp.TotalGroupNum = totalGroupNum
	times := GetRangeDate(fromTime, toTime)
	log.NewDebug(req.OperationID, "times:", times)
	wg := &sync.WaitGroup{}
	resp.IncreaseGroupNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	resp.TotalGroupNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	wg.Add(len(times))
	for i, v := range times {
		go func(wg *sync.WaitGroup, index int, v [2]time.Time) {
			defer wg.Done()
			num, err := imdb.GetIncreaseGroupNum(v[0], v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.IncreaseGroupNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  num,
			}
			num, err = imdb.GetGroupNum(v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.TotalGroupNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  num,
			}
		}(wg, i, v)
	}
	wg.Wait()
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	return resp, nil
}

func (s *adminCMSServer) GetMessageStatistics(_ context.Context, req *pbAdminCMS.GetMessageStatisticsReq) (*pbAdminCMS.GetMessageStatisticsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req.String())
	resp := &pbAdminCMS.GetMessageStatisticsResp{CommonResp: &pbAdminCMS.CommonResp{}}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "times: ", fromTime, toTime)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ParseTimeFromTo failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	privateMessageNum, err := imdb.GetPrivateMessageNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetPrivateMessageNum failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	groupMessageNum, err := imdb.GetGroupMessageNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMessageNum failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), privateMessageNum, groupMessageNum)
	resp.PrivateMessageNum = privateMessageNum
	resp.GroupMessageNum = groupMessageNum
	times := GetRangeDate(fromTime, toTime)
	resp.GroupMessageNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	resp.PrivateMessageNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	wg := &sync.WaitGroup{}
	wg.Add(len(times))
	for i, v := range times {
		go func(wg *sync.WaitGroup, index int, v [2]time.Time) {
			defer wg.Done()

			num, err := imdb.GetPrivateMessageNum(v[0], v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.PrivateMessageNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  num,
			}
			num, err = imdb.GetGroupMessageNum(v[0], v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.GroupMessageNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  num,
			}
		}(wg, i, v)
	}
	wg.Wait()
	return resp, nil
}

func (s *adminCMSServer) GetUserStatistics(_ context.Context, req *pbAdminCMS.GetUserStatisticsReq) (*pbAdminCMS.GetUserStatisticsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.GetUserStatisticsResp{CommonResp: &pbAdminCMS.CommonResp{}}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ParseTimeFromTo failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	activeUserNum, err := imdb.GetActiveUserNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetActiveUserNum failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	increaseUserNum, err := imdb.GetIncreaseUserNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseUserNum failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	totalUserNum, err := imdb.GetTotalUserNum()
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTotalUserNum failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.ActiveUserNum = activeUserNum
	resp.TotalUserNum = totalUserNum
	resp.IncreaseUserNum = increaseUserNum
	times := GetRangeDate(fromTime, toTime)
	resp.TotalUserNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	resp.ActiveUserNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	resp.IncreaseUserNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	wg := &sync.WaitGroup{}
	wg.Add(len(times))
	for i, v := range times {
		go func(wg *sync.WaitGroup, index int, v [2]time.Time) {
			defer wg.Done()
			num, err := imdb.GetActiveUserNum(v[0], v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.ActiveUserNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  num,
			}

			num, err = imdb.GetTotalUserNumByDate(v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTotalUserNumByDate", v, err.Error())
			}
			resp.TotalUserNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  num,
			}
			num, err = imdb.GetIncreaseUserNum(v[0], v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseUserNum", v, err.Error())
			}
			resp.IncreaseUserNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  num,
			}
		}(wg, i, v)
	}
	wg.Wait()
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	return resp, nil
}

func (s *adminCMSServer) GetUserFriends(_ context.Context, req *pbAdminCMS.GetUserFriendsReq) (*pbAdminCMS.GetUserFriendsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.GetUserFriendsResp{CommonResp: &pbAdminCMS.CommonResp{}, Pagination: &server_api_params.ResponsePagination{CurrentPage: req.Pagination.PageNumber, ShowNumber: req.Pagination.ShowNumber}}
	var friendList []*imdb.FriendUser
	var err error
	if req.FriendUserID != "" {
		friend, err := imdb.GetFriendByIDCMS(req.UserID, req.FriendUserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserID, req.FriendUserID)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		friendList = append(friendList, friend)
	} else {
		friendList, err = imdb.GetUserFriendsCMS(req.UserID, req.FriendUserName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserID, req.FriendUserName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
	}
	for _, v := range friendList {
		friendInfo := &server_api_params.FriendInfo{}
		userInfo := &server_api_params.UserInfo{UserID: v.FriendUserID, Nickname: v.Nickname}
		utils.CopyStructFields(friendInfo, v)
		friendInfo.FriendUser = userInfo
		resp.FriendInfoList = append(resp.FriendInfoList, friendInfo)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *adminCMSServer) GenerateInvitationCode(_ context.Context, req *pbAdminCMS.GenerateInvitationCodeReq) (*pbAdminCMS.GenerateInvitationCodeResp, error) {
	return nil, nil
}

func (s *adminCMSServer) GetInvitationCodes(_ context.Context, req *pbAdminCMS.GetInvitationCodesReq) (*pbAdminCMS.GetInvitationCodesResp, error) {
	return nil, nil
}

func (s *adminCMSServer) QueryIPRegister(_ context.Context, req *pbAdminCMS.QueryIPRegisterReq) (*pbAdminCMS.QueryIPRegisterResp, error) {
	return nil, nil
}

func (s *adminCMSServer) AddIPLimit(_ context.Context, req *pbAdminCMS.AddIPLimitReq) (*pbAdminCMS.AddIPLimitResp, error) {
	return nil, nil
}

func (s *adminCMSServer) RemoveIPLimit(_ context.Context, req *pbAdminCMS.RemoveIPLimitReq) (*pbAdminCMS.RemoveIPLimitResp, error) {
	return nil, nil
}

func (s *adminCMSServer) QueryUserIDIPLimitLogin(_ context.Context, req *pbAdminCMS.QueryUserIDIPLimitLoginReq) (*pbAdminCMS.QueryUserIDIPLimitLoginResp, error) {
	return nil, nil
}

func (s *adminCMSServer) AddUserIPLimitLogin(_ context.Context, req *pbAdminCMS.AddUserIPLimitLoginReq) (*pbAdminCMS.AddUserIPLimitLoginResp, error) {
	return nil, nil
}

func (s *adminCMSServer) RemoveUserIPLimit(_ context.Context, req *pbAdminCMS.RemoveUserIPLimitReq) (*pbAdminCMS.RemoveUserIPLimitResp, error) {
	return nil, nil
}

func (s *adminCMSServer) GetClientInitConfig(_ context.Context, req *pbAdminCMS.GetClientInitConfigReq) (*pbAdminCMS.GetClientInitConfigResp, error) {
	return nil, nil
}

func (s *adminCMSServer) SetClientInitConfig(_ context.Context, req *pbAdminCMS.SetClientInitConfigReq) (*pbAdminCMS.SetClientInitConfigResp, error) {
	return nil, nil
}

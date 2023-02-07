package admin_cms

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/getcdv3"
	pbAdminCMS "Open_IM/pkg/proto/admin_cms"
	server_api_params "Open_IM/pkg/proto/sdk_ws"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

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
	rpcPort           int
	rpcRegisterName   string
	etcdSchema        string
	etcdAddr          []string
	adminCMSInterface controller.AdminCMSInterface
	groupInterface    controller.GroupInterface
	userInterface     controller.UserInterface
	chatLogInterface  controller.ChatLogInterface
}

func NewAdminCMSServer(port int) *adminCMSServer {
	log.NewPrivateLog(constant.LogFileName)
	admin := &adminCMSServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAdminCMSName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	var mysql relation.Mysql
	var redis cache.RedisClient
	mysql.InitConn()
	redis.InitRedis()
	admin.userInterface = controller.NewUserController(mysql.GormConn())
	admin.groupInterface = controller.NewGroupController(mysql.GormConn(), redis.GetClient(), nil)
	admin.adminCMSInterface = controller.NewAdminCMSController(mysql.GormConn())
	admin.chatLogInterface = controller.NewChatLogController(mysql.GormConn())
	return admin
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
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
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
	log.NewInfo("", "rpcRegisterIP ", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10, "")
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

func (s *adminCMSServer) AdminLogin(ctx context.Context, req *pbAdminCMS.AdminLoginReq) (*pbAdminCMS.AdminLoginResp, error) {
	resp := &pbAdminCMS.AdminLoginResp{}
	for i, adminID := range config.Config.Manager.AppManagerUid {
		if adminID == req.AdminID && config.Config.Manager.Secrets[i] == req.Secret {
			token, expTime, err := token_verify.CreateToken(adminID, constant.LinuxPlatformID)
			if err != nil {
				log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "generate token failed", "adminID: ", adminID, err.Error())
				return nil, err
			}
			log.NewInfo(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "generate token success", "token: ", token, "expTime:", expTime)
			resp.Token = token
			break
		}
	}
	if resp.Token == "" {
		log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "failed")
		return nil, constant.ErrInternalServer
	}
	admin, err := s.userInterface.Take(ctx, req.AdminID)
	if err != nil {
		return nil, err
	}
	resp.UserName = admin.Nickname
	resp.FaceURL = admin.FaceURL
	return resp, nil
}

func (s *adminCMSServer) GetUserToken(ctx context.Context, req *pbAdminCMS.GetUserTokenReq) (*pbAdminCMS.GetUserTokenResp, error) {
	token, expTime, err := token_verify.CreateToken(req.UserID, int(req.PlatformID))
	if err != nil {
		return nil, err
	}
	resp := &pbAdminCMS.GetUserTokenResp{Token: token, ExpTime: expTime}
	return resp, nil
}

func (s *adminCMSServer) GetChatLogs(ctx context.Context, req *pbAdminCMS.GetChatLogsReq) (*pbAdminCMS.GetChatLogsResp, error) {
	chatLog := relation.ChatLog{
		Content:     req.Content,
		ContentType: req.ContentType,
		SessionType: req.SessionType,
		RecvID:      req.RecvID,
		SendID:      req.SendID,
	}
	if req.SendTime != "" {
		sendTime, err := utils.TimeStringToTime(req.SendTime)
		if err != nil {
			log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "time string parse error", err.Error())
			return nil, err
		}
		chatLog.SendTime = sendTime
	}
	num, chatLogs, err := s.chatLogInterface.GetChatLog(&chatLog, req.Pagination.PageNumber, req.Pagination.ShowNumber, []int32{
		constant.Text,
		constant.Picture,
		constant.Voice,
		constant.Video,
		constant.File,
		constant.AtText,
		constant.Merger,
		constant.Card,
		constant.Location,
		constant.Custom,
		constant.Revoke,
		constant.Quote,
		constant.AdvancedText,
		constant.AdvancedRevoke,
		constant.CustomNotTriggerConversation,
	})
	if err != nil {
		return nil, err
	}
	resp := &pbAdminCMS.GetChatLogsResp{}
	resp.ChatLogsNum = int32(num)
	for _, chatLog := range chatLogs {
		pbChatLog := &pbAdminCMS.ChatLog{}
		utils.CopyStructFields(pbChatLog, chatLog)
		pbChatLog.SendTime = chatLog.SendTime.Unix()
		pbChatLog.CreateTime = chatLog.CreateTime.Unix()
		if chatLog.SenderNickname == "" {
			sendUser, err := s.userInterface.Take(ctx, chatLog.SendID)
			if err != nil {
				return nil, err
			}
			pbChatLog.SenderNickname = sendUser.Nickname
		}
		switch chatLog.SessionType {
		case constant.SingleChatType:
			recvUser, err := s.userInterface.Take(ctx, chatLog.RecvID)
			if err != nil {
				return nil, err
			}
			pbChatLog.SenderNickname = recvUser.Nickname

		case constant.GroupChatType, constant.SuperGroupChatType:
			group, err := s.groupInterface.TakeGroup(ctx, chatLog.RecvID)
			if err != nil {
				return nil, err
			}
			pbChatLog.RecvID = group.GroupID
			pbChatLog.GroupName = group.GroupName
		}
		resp.ChatLogs = append(resp.ChatLogs, pbChatLog)
	}
	return resp, nil
}

func (s *adminCMSServer) GetActiveGroup(_ context.Context, req *pbAdminCMS.GetActiveGroupReq) (*pbAdminCMS.GetActiveGroupResp, error) {
	resp := &pbAdminCMS.GetActiveGroupResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		return nil, err
	}
	activeGroups, err := s.adminCMSInterface.GetActiveGroups(fromTime, toTime, 12)
	if err != nil {
		return nil, err
	}
	for _, activeGroup := range activeGroups {
		resp.Groups = append(resp.Groups,
			&pbAdminCMS.GroupResp{
				GroupName:  activeGroup.Name,
				GroupID:    activeGroup.ID,
				MessageNum: int32(activeGroup.MessageNum),
			})
	}
	return resp, nil
}

func (s *adminCMSServer) GetActiveUser(ctx context.Context, req *pbAdminCMS.GetActiveUserReq) (*pbAdminCMS.GetActiveUserResp, error) {
	resp := &pbAdminCMS.GetActiveUserResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "ParseTimeFromTo failed", err.Error())
		return nil, err
	}
	activeUsers, err := s.adminCMSInterface.GetActiveUsers(fromTime, toTime, 12)
	if err != nil {
		return nil, err
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

func (s *adminCMSServer) GetGroupStatistics(ctx context.Context, req *pbAdminCMS.GetGroupStatisticsReq) (*pbAdminCMS.GetGroupStatisticsResp, error) {
	resp := &pbAdminCMS.GetGroupStatisticsResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		return nil, err
	}
	increaseGroupNum, err := s.adminCMSInterface.GetIncreaseGroupNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		return nil, err
	}
	totalGroupNum, err := s.adminCMSInterface.GetTotalGroupNum()
	if err != nil {
		return nil, err
	}
	resp.IncreaseGroupNum = int32(increaseGroupNum)
	resp.TotalGroupNum = int32(totalGroupNum)
	times := GetRangeDate(fromTime, toTime)
	wg := &sync.WaitGroup{}
	resp.IncreaseGroupNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	resp.TotalGroupNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	wg.Add(len(times))
	for i, v := range times {
		go func(wg *sync.WaitGroup, index int, v [2]time.Time) {
			defer wg.Done()
			num, err := s.adminCMSInterface.GetIncreaseGroupNum(v[0], v[1])
			if err != nil {
				log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.IncreaseGroupNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  int32(num),
			}
			num, err = s.adminCMSInterface.GetGroupNum(v[1])
			if err != nil {
				log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.TotalGroupNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  int32(num),
			}
		}(wg, i, v)
	}
	wg.Wait()
	return resp, nil
}

func (s *adminCMSServer) GetMessageStatistics(ctx context.Context, req *pbAdminCMS.GetMessageStatisticsReq) (*pbAdminCMS.GetMessageStatisticsResp, error) {
	resp := &pbAdminCMS.GetMessageStatisticsResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	log.NewDebug(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "times: ", fromTime, toTime)
	if err != nil {
		log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "ParseTimeFromTo failed", err.Error())
		return nil, err
	}
	privateMessageNum, err := s.adminCMSInterface.GetSingleChatMessageNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		return nil, err
	}
	groupMessageNum, err := s.adminCMSInterface.GetGroupMessageNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		return nil, err
	}
	log.NewDebug(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), privateMessageNum, groupMessageNum)
	resp.PrivateMessageNum = int32(privateMessageNum)
	resp.GroupMessageNum = int32(groupMessageNum)
	times := GetRangeDate(fromTime, toTime)
	resp.GroupMessageNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	resp.PrivateMessageNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	wg := &sync.WaitGroup{}
	wg.Add(len(times))
	for i, v := range times {
		go func(wg *sync.WaitGroup, index int, v [2]time.Time) {
			defer wg.Done()
			num, err := s.adminCMSInterface.GetSingleChatMessageNum(v[0], v[1])
			if err != nil {
				log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.PrivateMessageNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  int32(num),
			}
			num, err = s.adminCMSInterface.GetGroupMessageNum(v[0], v[1])
			if err != nil {
				log.NewError(tracelog.GetOperationID(ctx), utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.GroupMessageNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  int32(num),
			}
		}(wg, i, v)
	}
	wg.Wait()
	return resp, nil
}

func (s *adminCMSServer) GetUserStatistics(_ context.Context, req *pbAdminCMS.GetUserStatisticsReq) (*pbAdminCMS.GetUserStatisticsResp, error) {
	resp := &pbAdminCMS.GetUserStatisticsResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ParseTimeFromTo failed", err.Error())
		return nil, err
	}
	activeUserNum, err := s.adminCMSInterface.GetActiveUserNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		return nil, err
	}
	increaseUserNum, err := s.adminCMSInterface.GetIncreaseUserNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		return nil, err
	}
	totalUserNum, err := s.adminCMSInterface.GetTotalUserNum()
	if err != nil {
		return nil, err
	}
	resp.ActiveUserNum = int32(activeUserNum)
	resp.TotalUserNum = int32(totalUserNum)
	resp.IncreaseUserNum = int32(increaseUserNum)
	times := GetRangeDate(fromTime, toTime)
	resp.TotalUserNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	resp.ActiveUserNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	resp.IncreaseUserNumList = make([]*pbAdminCMS.DateNumList, len(times), len(times))
	wg := &sync.WaitGroup{}
	wg.Add(len(times))
	for i, v := range times {
		go func(wg *sync.WaitGroup, index int, v [2]time.Time) {
			defer wg.Done()
			num, err := s.adminCMSInterface.GetActiveUserNum(v[0], v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
			}
			resp.ActiveUserNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  int32(num),
			}

			num, err = s.adminCMSInterface.GetTotalUserNumByDate(v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTotalUserNumByDate", v, err.Error())
			}
			resp.TotalUserNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  int32(num),
			}
			num, err = s.adminCMSInterface.GetIncreaseUserNum(v[0], v[1])
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetIncreaseUserNum", v, err.Error())
			}
			resp.IncreaseUserNumList[index] = &pbAdminCMS.DateNumList{
				Date: v[0].String(),
				Num:  int32(num),
			}
		}(wg, i, v)
	}
	wg.Wait()
	return resp, nil
}

func (s *adminCMSServer) GetUserFriends(ctx context.Context, req *pbAdminCMS.GetUserFriendsReq) (*pbAdminCMS.GetUserFriendsResp, error) {
	resp := &pbAdminCMS.GetUserFriendsResp{}
	var friendList []*relation.FriendUser
	var err error
	if req.FriendUserID != "" {
		friend, err := s.adminCMSInterface.GetFriendByIDCMS(req.UserID, req.FriendUserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			return nil, err
		}
		friendList = append(friendList, friend)
		resp.FriendNums = 1
	} else {
		var count int64
		friendList, count, err = s.adminCMSInterface.GetUserFriendsCMS(req.UserID, req.FriendUserName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
		if err != nil {
			return nil, err
		}
		resp.FriendNums = int32(count)
	}
	for _, v := range friendList {
		friendInfo := &server_api_params.FriendInfo{}
		userInfo := &server_api_params.UserInfo{UserID: v.FriendUserID, Nickname: v.Nickname}
		utils.CopyStructFields(friendInfo, v)
		friendInfo.FriendUser = userInfo
		resp.FriendInfoList = append(resp.FriendInfoList, friendInfo)
	}
	return resp, nil
}

func (s *adminCMSServer) GetUserIDByEmailAndPhoneNumber(ctx context.Context, req *pbAdminCMS.GetUserIDByEmailAndPhoneNumberReq) (*pbAdminCMS.GetUserIDByEmailAndPhoneNumberResp, error) {
	resp := &pbAdminCMS.GetUserIDByEmailAndPhoneNumberResp{}
	userIDList, err := s.userInterface.GetUserIDsByEmailAndID(req.PhoneNumber, req.Email)
	if err != nil {
		return resp, nil
	}
	resp.UserIDList = userIDList
	return resp, nil
}

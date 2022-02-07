package statistics

import (
	"Open_IM/pkg/common/config"
	"context"
	"fmt"
	"time"

	//"Open_IM/pkg/common/constant"
	//"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"

	//cp "Open_IM/pkg/common/utils"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbStatistics "Open_IM/pkg/proto/statistics"

	//open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	//"context"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type statisticsServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewStatisticsServer(port int) *statisticsServer {
	log.NewPrivateLog("Statistics")
	return &statisticsServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImStatisticsName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *statisticsServer) Run() {
	log.NewInfo("0", "Statistics rpc start ")
	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("0", "Listen failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("0", "listen network success, ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbStatistics.RegisterUserServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "statistics rpc success")
}

func (s *statisticsServer) GetActiveGroup(_ context.Context, req *pbStatistics.GetActiveGroupReq) (*pbStatistics.GetActiveGroupResp, error) {
	log.NewInfo("GetActiveGroup", utils.GetSelfFuncName(), req.String())
	resp := &pbStatistics.GetActiveGroupResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError("GetActiveGroup", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	activeGroups, err := imdb.GetActiveGroups(fromTime, toTime, 12)
	if err != nil {
		log.NewError("GetActiveGroup", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	for _, activeGroup := range activeGroups {
		resp.Groups = append(resp.Groups,
			&pbStatistics.GroupResp{
				GroupName:  activeGroup.Name,
				GroupId:    activeGroup.Id,
				MessageNum: int32(activeGroup.MessageNum),
			})
	}
	return resp, nil
}

func (s *statisticsServer) GetActiveUser(_ context.Context, req *pbStatistics.GetActiveUserReq) (*pbStatistics.GetActiveUserResp, error) {
	log.NewInfo("GetActiveUser", utils.GetSelfFuncName(), req.String())
	resp := &pbStatistics.GetActiveUserResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError("GetActiveUser", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	activeUsers, err := imdb.GetActiveUsers(fromTime, toTime, 12)
	if err != nil {
		log.NewError("GetActiveUser", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	for _, activeUser := range activeUsers {
		resp.Users = append(resp.Users,
			&pbStatistics.UserResp{
				UserId:     activeUser.Id,
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
	return from.Month() == to.Month()
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
	return times
}

func getFirstDateOfNextNMonth(currentTime time.Time, n int) time.Time {
	lastOfMonth := time.Date(currentTime.Year(), currentTime.Month(), 1, 0, 0, 0, 0, currentTime.Location()).AddDate(0, n, 0)
	return lastOfMonth
}

func (s *statisticsServer) GetGroupStatistics(_ context.Context, req *pbStatistics.GetGroupStatisticsReq) (*pbStatistics.GetGroupStatisticsResp, error) {
	log.NewInfo("GetGroupStatistics", utils.GetSelfFuncName(), req.String())
	resp := &pbStatistics.GetGroupStatisticsResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError("GetMessageStatistics", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	increaseGroupNum, err := imdb.GetIncreaseGroupNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError("GetMessageStatistics", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	totalGroupNum, err := imdb.GetTotalGroupNum()
	if err != nil {
		log.NewError("GetMessageStatistics", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	resp.IncreaseGroupNum = increaseGroupNum
	resp.TotalGroupNum = totalGroupNum
	times := GetRangeDate(fromTime, toTime)
	log.NewInfo("", "times:", times)
	for _, v := range times {
		num, err := imdb.GetIncreaseGroupNum(v[0], v[1])
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
		}
		resp.IncreaseGroupNumList = append(resp.IncreaseGroupNumList, &pbStatistics.DateNumList{
			Date: v[0].String(),
			Num:  num,
		})
		num, err = imdb.GetGroupNum(v[0])
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
		}
		resp.TotalGroupNumList = append(resp.TotalGroupNumList, &pbStatistics.DateNumList{
			Date: v[0].String(),
			Num:  num,
		})
	}
	return resp, nil
}

func (s *statisticsServer) GetMessageStatistics(_ context.Context, req *pbStatistics.GetMessageStatisticsReq) (*pbStatistics.GetMessageStatisticsResp, error) {
	log.NewInfo("GetMessageStatistics", utils.GetSelfFuncName(), req.String())
	resp := &pbStatistics.GetMessageStatisticsResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		log.NewError("GetMessageStatistics", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	privateMessageNum, err := imdb.GetPrivateMessageNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError("GetMessageStatistics", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	groupMessageNum, err := imdb.GetGroupMessageNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError("GetMessageStatistics", utils.GetSelfFuncName(), err.Error())
	}
	resp.PrivateMessageNum = privateMessageNum
	resp.GroupMessageNum = groupMessageNum
	times := GetRangeDate(fromTime, toTime)
	fmt.Println(times)
	for _, v := range times {
		num, err := imdb.GetPrivateMessageNum(v[0], v[1])
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
		}
		resp.PrivateMessageNumList = append(resp.PrivateMessageNumList, &pbStatistics.DateNumList{
			Date: v[0].String(),
			Num:  num,
		})

		num, err = imdb.GetGroupMessageNum(v[0], v[1])
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
		}
		resp.GroupMessageNumList = append(resp.GroupMessageNumList, &pbStatistics.DateNumList{
			Date: v[0].String(),
			Num:  num,
		})
	}
	fmt.Println(resp)
	return resp, nil
}

func (s *statisticsServer) GetUserStatistics(_ context.Context, req *pbStatistics.GetUserStatisticsReq) (*pbStatistics.GetUserStatisticsResp, error) {
	log.NewInfo("GetUserStatistics", utils.GetSelfFuncName(), req.String())
	resp := &pbStatistics.GetUserStatisticsResp{}
	fromTime, toTime, err := ParseTimeFromTo(req.StatisticsReq.From, req.StatisticsReq.To)
	if err != nil {
		return resp, err
	}
	activeUserNum, err := imdb.GetActiveUserNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError("GetUserStatistics", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	increaseUserNum, err := imdb.GetIncreaseUserNum(fromTime, toTime.Add(time.Hour*24))
	if err != nil {
		log.NewError("GetUserStatistics", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	totalUserNum, err := imdb.GetTotalUserNum()
	if err != nil {
		log.NewError("GetUserStatistics", utils.GetSelfFuncName(), err.Error())
		return resp, err
	}
	resp.ActiveUserNum = activeUserNum
	resp.TotalUserNum = totalUserNum
	resp.IncreaseUserNum = increaseUserNum

	times := GetRangeDate(fromTime, toTime)
	for _, v := range times {
		num, err := imdb.GetActiveUserNum(v[0], v[1])
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), "GetIncreaseGroupNum", v, err.Error())
		}
		resp.ActiveUserNumList = append(resp.ActiveUserNumList, &pbStatistics.DateNumList{
			Date: v[0].String(),
			Num:  num,
		})
		num, err = imdb.GetTotalUserNumByDate(v[0])
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), "GetTotalUserNumByDate", v, err.Error())
		}
		resp.TotalUserNumList = append(resp.TotalUserNumList, &pbStatistics.DateNumList{
			Date: v[0].String(),
			Num:  num,
		})
		num, err = imdb.GetIncreaseUserNum(v[0], v[1])
		if err != nil {
			log.NewError("", utils.GetSelfFuncName(), "GetIncreaseUserNum", v, err.Error())
		}
		resp.IncreaseUserNumList = append(resp.IncreaseUserNumList, &pbStatistics.DateNumList{
			Date: v[0].String(),
			Num:  num,
		})
	}
	return resp, nil
}

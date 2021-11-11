package chat

import (
	"context"
	"github.com/garyburd/redigo/redis"

	commonDB "Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"

	"sort"
	"strings"

	pbMsg "Open_IM/pkg/proto/chat"
)

func (rpc *rpcChat) GetMaxAndMinSeq(_ context.Context, in *pbMsg.GetMaxAndMinSeqReq) (*pbMsg.GetMaxAndMinSeqResp, error) {
	log.InfoByKv("rpc getMaxAndMinSeq is arriving", in.OperationID, in.String())
	//seq, err := model.GetBiggestSeqFromReceive(in.UserID)
	maxSeq, err1 := commonDB.DB.GetUserMaxSeq(in.UserID)
	minSeq, err2 := commonDB.DB.GetUserMinSeq(in.UserID)
	resp := new(pbMsg.GetMaxAndMinSeqResp)
	if err1 == nil {
		resp.MaxSeq = maxSeq
	} else if err1 == redis.ErrNil {
		resp.MaxSeq = 0
	} else {
		log.NewError(in.OperationID, "getMaxSeq from redis error", in.String(), err1.Error())
		resp.MaxSeq = -1
		resp.ErrCode = 200
		resp.ErrMsg = "redis get err"
	}
	if err2 == nil {
		resp.MinSeq = minSeq
	} else if err2 == redis.ErrNil {
		resp.MinSeq = 0
	} else {
		log.NewError(in.OperationID, "getMaxSeq from redis error", in.String(), err2.Error())
		resp.MinSeq = -1
		resp.ErrCode = 201
		resp.ErrMsg = "redis get err"
	}
	return resp, nil
}
func (rpc *rpcChat) PullMessage(_ context.Context, in *pbMsg.PullMessageReq) (*pbMsg.PullMessageResp, error) {
	log.InfoByKv("rpc pullMessage is arriving", in.OperationID, "args", in.String())
	resp := new(pbMsg.PullMessageResp)
	var respSingleMsgFormat []*pbMsg.GatherFormat
	var respGroupMsgFormat []*pbMsg.GatherFormat
	SingleMsgFormat, GroupMsgFormat, MaxSeq, MinSeq, err := commonDB.DB.GetMsgBySeqRange(in.UserID, in.SeqBegin, in.SeqEnd)
	if err != nil {
		log.ErrorByKv("pullMsg data error", in.OperationID, in.String())
		resp.ErrCode = 1
		resp.ErrMsg = err.Error()
		return resp, nil
	}
	respSingleMsgFormat = singleMsgHandleByUser(SingleMsgFormat, in.UserID)
	respGroupMsgFormat = groupMsgHandleByUser(GroupMsgFormat)
	return &pbMsg.PullMessageResp{
		ErrCode:       0,
		ErrMsg:        "",
		MaxSeq:        MaxSeq,
		MinSeq:        MinSeq,
		SingleUserMsg: respSingleMsgFormat,
		GroupUserMsg:  respGroupMsgFormat,
	}, nil
}
func (rpc *rpcChat) PullMessageBySeqList(_ context.Context, in *pbMsg.PullMessageBySeqListReq) (*pbMsg.PullMessageResp, error) {
	log.NewInfo(in.OperationID, "rpc PullMessageBySeqList is arriving", in.String())
	resp := new(pbMsg.PullMessageResp)
	var respSingleMsgFormat []*pbMsg.GatherFormat
	var respGroupMsgFormat []*pbMsg.GatherFormat
	SingleMsgFormat, GroupMsgFormat, MaxSeq, MinSeq, err := commonDB.DB.GetMsgBySeqList(in.UserID, in.SeqList)
	if err != nil {
		log.ErrorByKv("PullMessageBySeqList data error", in.OperationID, in.String())
		resp.ErrCode = 1
		resp.ErrMsg = err.Error()
		return resp, nil
	}
	respSingleMsgFormat = singleMsgHandleByUser(SingleMsgFormat, in.UserID)
	respGroupMsgFormat = groupMsgHandleByUser(GroupMsgFormat)
	return &pbMsg.PullMessageResp{
		ErrCode:       0,
		ErrMsg:        "",
		MaxSeq:        MaxSeq,
		MinSeq:        MinSeq,
		SingleUserMsg: respSingleMsgFormat,
		GroupUserMsg:  respGroupMsgFormat,
	}, nil
}
func singleMsgHandleByUser(allMsg []*pbMsg.MsgFormat, ownerId string) []*pbMsg.GatherFormat {
	var userid string
	var respMsgFormat []*pbMsg.GatherFormat
	m := make(map[string]MsgFormats)
	//Gather messages in the dimension of users
	for _, v := range allMsg {
		if v.RecvID != ownerId {
			userid = v.RecvID
		} else {
			userid = v.SendID
		}
		if value, ok := m[userid]; !ok {
			var t MsgFormats
			m[userid] = append(t, v)
		} else {
			m[userid] = append(value, v)
		}
	}
	//Return in pb format
	for user, msg := range m {
		tempUserMsg := new(pbMsg.GatherFormat)
		tempUserMsg.ID = user
		tempUserMsg.List = msg
		sort.Sort(msg)
		respMsgFormat = append(respMsgFormat, tempUserMsg)
	}
	return respMsgFormat
}
func groupMsgHandleByUser(allMsg []*pbMsg.MsgFormat) []*pbMsg.GatherFormat {
	var respMsgFormat []*pbMsg.GatherFormat
	m := make(map[string]MsgFormats)
	//Gather messages in the dimension of users
	for _, v := range allMsg {
		//Get group ID
		groupID := strings.Split(v.RecvID, " ")[1]
		if value, ok := m[groupID]; !ok {
			var t MsgFormats
			m[groupID] = append(t, v)
		} else {
			m[groupID] = append(value, v)
		}

	}
	//Return in pb format
	for groupID, msg := range m {
		tempUserMsg := new(pbMsg.GatherFormat)
		tempUserMsg.ID = groupID
		tempUserMsg.List = msg
		sort.Sort(msg)
		respMsgFormat = append(respMsgFormat, tempUserMsg)
	}
	return respMsgFormat
}

type MsgFormats []*pbMsg.MsgFormat

// Implement the sort.Interface interface to get the number of elements method
func (s MsgFormats) Len() int {
	return len(s)
}

//Implement the sort.Interface interface comparison element method
func (s MsgFormats) Less(i, j int) bool {
	return s[i].SendTime < s[j].SendTime
}

//Implement the sort.Interface interface exchange element method
func (s MsgFormats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

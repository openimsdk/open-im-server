package rpcChat

import (
	"context"

	commonDB "Open_IM/src/common/db"
	"Open_IM/src/common/log"

	"sort"
	"strings"

	pbMsg "Open_IM/src/proto/chat"
)

func (rpc *rpcChat) GetNewSeq(_ context.Context, in *pbMsg.GetNewSeqReq) (*pbMsg.GetNewSeqResp, error) {
	log.InfoByKv("rpc getNewSeq is arriving", in.OperationID, in.String())
	//seq, err := model.GetBiggestSeqFromReceive(in.UserID)
	seq, err := commonDB.DB.GetUserSeq(in.UserID)
	resp := new(pbMsg.GetNewSeqResp)
	if err == nil {
		resp.Seq = seq
		resp.ErrCode = 0
		resp.ErrMsg = ""
		return resp, err
	} else {
		log.ErrorByKv("getSeq from redis error", in.OperationID, "args", in.String(), "err", err.Error())
		resp.Seq = 0
		resp.ErrCode = 0
		resp.ErrMsg = ""
		return resp, nil
	}

}

func (rpc *rpcChat) PullMessage(_ context.Context, in *pbMsg.PullMessageReq) (*pbMsg.PullMessageResp, error) {
	log.InfoByKv("rpc pullMessage is arriving", in.OperationID, "args", in.String())
	resp := new(pbMsg.PullMessageResp)
	var respSingleMsgFormat []*pbMsg.GatherFormat
	var respGroupMsgFormat []*pbMsg.GatherFormat
	SingleMsgFormat, GroupMsgFormat, MaxSeq, MinSeq, err := commonDB.DB.GetUserChat(in.UserID, in.SeqBegin, in.SeqEnd)
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

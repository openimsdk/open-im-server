//实现pb定义的rpc服务
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

//func (s *MsgServer) PullMessage(_ context.Context, in *pbMsg.PullMessageReq) (*pbMsg.PullMessageResp, error) {
//	log.InfoByArgs(fmt.Sprintf("rpc pullMessage is arriving,args=%s", in.String()))
//	resp := new(pbMsg.PullMessageResp)
//	var respMsgFormat []*pbMsg.MsgFormat
//	var respUserMsgFormat []*pbMsg.UserMsgFormat
//	conn := db.NewDbConnection()
//	rows, err := conn.Table("receive r").Select("c.sender_id,c.receiver_id,"+
//		"c.msg_type,c.push_msg_type,c.chat_type,c.msg_id,c.send_content,r.seq,c.send_time,c.sender_nickname,c.receiver_nickname,c.sender_head_url,c.receiver_head_url").
//		Joins("INNER JOIN chat_log c ON r.msg_id = c.msg_id AND r.user_id = ? AND seq BETWEEN ? AND ?",
//			in.UserID, in.SeqBegin, in.SeqEnd).Rows()
//	if err != nil {
//		fmt.Printf("pullMsg data error: %v\n", err)
//		resp.ErrCode = 1
//		resp.ErrMsg = err.Error()
//		return resp, nil
//	}
//	defer rows.Close()
//	for rows.Next() {
//		tempResp := new(pbMsg.MsgFormat)
//		rows.Scan(&tempResp.SendID, &tempResp.RecvID, &tempResp.MsgType, &tempResp.PushMsgType, &tempResp.ChatType,
//			&tempResp.MsgID, &tempResp.Msg, &tempResp.Seq, &tempResp.Time, &tempResp.SendNickName, &tempResp.RecvNickName,
//			&tempResp.SendHeadUrl, &tempResp.RecvHeadUrl)
//		respMsgFormat = append(respMsgFormat, tempResp)
//	}
//	respUserMsgFormat = msgHandleByUser(respMsgFormat, in.UserID)
//	return &pbMsg.PullMessageResp{
//		ErrCode: 0,
//		ErrMsg:  "",
//		UserMsg: respUserMsgFormat,
//	}, nil
//}
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
	//将消息以用户为维度聚集
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
	//形成pb格式返回
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
	//将消息以用户为维度聚集
	for _, v := range allMsg {
		//获得群ID
		groupID := strings.Split(v.RecvID, " ")[1]
		if value, ok := m[groupID]; !ok {
			var t MsgFormats
			m[groupID] = append(t, v)
		} else {
			m[groupID] = append(value, v)
		}

	}
	//形成pb格式返回
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

// 实现sort.Interface接口取元素数量方法
func (s MsgFormats) Len() int {
	return len(s)
}

// 实现sort.Interface接口比较元素方法
func (s MsgFormats) Less(i, j int) bool {
	return s[i].SendTime < s[j].SendTime
}

// 实现sort.Interface接口交换元素方法
func (s MsgFormats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

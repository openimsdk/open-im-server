package msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/internal/model"
	"github.com/OpenIMSDK/Open-IM-Server/internal/logic/service"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/protocol"
)

type msgServer struct {
	pb.UnimplementedMsgServer
}

func (s *msgServer) SearchMsg(ctx context.Context, req *pb.SearchMsgReq) (*pb.SearchMsgResp, error) {
	// Assume that chatLogs is the result of the database query that retrieves the chatLogs
 	chatLogs, err := service.NewMsgService().SearchMsgFromDB(ctx, req)
	if err != nil {
		return nil, err
	}

	// Handle the case where the query returns no results
	if len(chatLogs) == 0 {
		return &pb.SearchMsgResp{}, nil
	}

	// Initialize a slice of pb.ChatLog objects
	var pbChatLogs []*pb.ChatLog

	// Iterate over the chatLogs slice
	for _, chatLog := range chatLogs {
		// Create a pb.ChatLog object for each chatLog and append it to the pbChatLogs slice
		pbChatLog := &pb.ChatLog{
			SenderID:    chatLog.SenderID,
			ReceiverID:  chatLog.ReceiverID,
			Content:     chatLog.Content,
			ContentType: chatLog.ContentType,
			Timestamp:   chatLog.Timestamp,
		}
		pbChatLogs = append(pbChatLogs, pbChatLog)
	}

	// Return the pbChatLogs slice in the pb.SearchMsgResp object
	return &pb.SearchMsgResp{ChatLogs: pbChatLogs}, nil
}


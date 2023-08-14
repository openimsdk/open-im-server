package msg

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/internal/logic/model"
	"github.com/OpenIMSDK/Open-IM-Server/internal/logic/service"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/pb"
)

type msgServer struct {
	pb.UnimplementedMsgServer
}

func (s *msgServer) SearchMsg(ctx context.Context, req *pb.SearchMsgReq) (*pb.SearchMsgResp, error) {
	// Assume that chatLogs is the result of the database query that retrieves the chatLogs
	chatLogs, err := service.NewMsgService().SearchMsg(ctx, req)
	if err != nil {
		return nil, err
	}

	// Handle the case where the query returns no results
	if len(chatLogs) == 0 {
		return &pb.SearchMsgResp{}, nil
	}

	// Rest of the function implementation
	// ...
}


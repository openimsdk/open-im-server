package msg

import (
	"testing"
	"github.com/OpenIMSDK/Open-IM-Server/internal/logic/service"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/pb"
)

func TestSearchMsg(t *testing.T) {
	// Ensure that the database is empty
	// clearDatabase()

	// Call the SearchMsg function
	req := &pb.SearchMsgReq{}
	resp, err := service.NewMsgService().SearchMsg(context.Background(), req)

	// Check that the error is nil and the response is empty
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if resp != &pb.SearchMsgResp{} {
		t.Errorf("Expected empty response, got %v", resp)
	}
}


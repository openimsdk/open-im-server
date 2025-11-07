package cmd

import (
	"math"
	"testing"

	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/utils/jsonutil"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockRootCmd is a mock type for the RootCmd type
type MockRootCmd struct {
	mock.Mock
}

func (m *MockRootCmd) Execute() error {
	args := m.Called()
	return args.Error(0)
}

func TestName(t *testing.T) {
	resp := &apiresp.ApiResponse{
		ErrCode: 1234,
		ErrMsg:  "test",
		ErrDlt:  "4567",
		Data: &auth.GetUserTokenResp{
			Token:             "1234567",
			ExpireTimeSeconds: math.MaxInt64,
		},
	}
	data, err := resp.MarshalJSON()
	if err != nil {
		panic(err)
	}
	t.Log(string(data))

	var rReso apiresp.ApiResponse
	rReso.Data = &auth.GetUserTokenResp{}

	if err := jsonutil.JsonUnmarshal(data, &rReso); err != nil {
		panic(err)
	}

	t.Logf("%+v\n", rReso)

}

func TestName1(t *testing.T) {
	t.Log(primitive.NewObjectID().String())
	t.Log(primitive.NewObjectID().Hex())

}

// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

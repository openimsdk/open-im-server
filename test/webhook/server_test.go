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

package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	sdkws "github.com/openimsdk/protocol/sdkws"
	"github.com/stretchr/testify/require"
)

func TestWebhookServer_DefaultEndpoints(t *testing.T) {
	t.Parallel()
	srv := NewServer()

	tests := []struct {
		name      string
		url       string
		payload   any
		expectRes int
	}{
		{
			name: "CallbackBeforeUserRegister with prefix",
			url:  "/callbackExample/" + cbapi.CallbackBeforeUserRegisterCommand,
			payload: &cbapi.CallbackBeforeUserRegisterReq{
				CallbackCommand: cbapi.CallbackCommand(cbapi.CallbackBeforeUserRegisterCommand),
				Users:           []*sdkws.UserInfo{{UserID: "u100", Nickname: "Alice"}},
			},
			expectRes: http.StatusOK,
		},
		{
			name: "CallbackAfterMsgSaveDB without prefix",
			url:  "/" + cbapi.CallbackAfterMsgSaveDBCommand,
			payload: &cbapi.CallbackAfterMsgSaveDBReq{
				CommonCallbackReq: cbapi.CommonCallbackReq{
					SendID:          "u100",
					CallbackCommand: cbapi.CallbackAfterMsgSaveDBCommand,
				},
				RecvID: "u200",
			},
			expectRes: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			body, err := json.Marshal(tt.payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			srv.ServeHTTP(w, req)
			require.Equal(t, tt.expectRes, w.Code)
		})
	}
}

func TestWebhookServer_CustomHandlerOverride(t *testing.T) {
	t.Parallel()
	srv := NewServer()

	// Override BeforeMsgModify to mutate content
	filteredContent := `{"text":"***"}`
	RegisterCallback(
		srv,
		cbapi.CallbackBeforeMsgModifyCommand,
		func(c *gin.Context, req *cbapi.CallbackMsgModifyCommandReq) (*cbapi.CallbackMsgModifyCommandResp, error) {
			return &cbapi.CallbackMsgModifyCommandResp{
				CommonCallbackResp: cbapi.CommonCallbackResp{
					ErrCode: 0,
					ErrMsg:  "filtered",
				},
				Content: &filteredContent,
			}, nil
		},
	)

	reqPayload := cbapi.CallbackMsgModifyCommandReq{
		CommonCallbackReq: cbapi.CommonCallbackReq{
			SendID:          "u1",
			CallbackCommand: cbapi.CallbackBeforeMsgModifyCommand,
		},
	}
	body, err := json.Marshal(reqPayload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/"+cbapi.CallbackBeforeMsgModifyCommand, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var resp cbapi.CallbackMsgModifyCommandResp
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.NotNil(t, resp.Content)
	require.Equal(t, `{"text":"***"}`, *resp.Content)
	require.Equal(t, "filtered", resp.ErrMsg)
}

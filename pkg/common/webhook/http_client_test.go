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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/tools/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SyncPost_FailedContinue(t *testing.T) {
	respond := func(body string) string {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(body))
		}))
		t.Cleanup(srv.Close)
		return srv.URL
	}
	dropped := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if conn, _, err := w.(http.Hijacker).Hijack(); err == nil {
			_ = conn.Close()
		}
	}))
	t.Cleanup(dropped.Close)
	bob := "bob"

	tests := []struct {
		name           string
		url            string
		failedContinue bool
		wantErr        errs.CodeError
		wantNickname   *string
	}{
		{"connection dropped", dropped.URL, false, servererrs.ErrNetwork, nil},
		{"connection dropped continue", dropped.URL, true, nil, nil},
		{"malformed response", respond("<html>bad gateway</html>"), false, servererrs.ErrData, nil},
		{"malformed response continue", respond("<html>bad gateway</html>"), true, nil, nil},
		{"partially decoded response continue", respond(`{"nickName":"bob","actionCode":"x"}`), true, nil, nil},
		{"accepted continue", respond(`{"actionCode":0,"nextCode":0,"nickName":"bob"}`), true, nil, &bob},
		{"rejected continue", respond(`{"actionCode":0,"nextCode":1,"errCode":5001,"errMsg":"denied"}`), true, errs.NewCodeError(5001, "denied"), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &callbackstruct.CallbackBeforeUpdateUserInfoResp{}
			before := &config.BeforeConfig{Enable: true, Timeout: 5, FailedContinue: tt.failedContinue}
			err := NewWebhookClient(tt.url).SyncPost(context.Background(), "callbackBeforeUpdateUserInfoCommand", &callbackstruct.CallbackBeforeUpdateUserInfoReq{}, resp, before)
			if tt.wantErr == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.wantNickname, resp.Nickname)
				return
			}
			require.Error(t, err)
			assert.True(t, tt.wantErr.Is(err), "got %v", err)
		})
	}
}

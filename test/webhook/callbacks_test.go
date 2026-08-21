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

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	sdkws "github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/stretchr/testify/require"
)

func runWebhookTest(t *testing.T, srv *Server, command string, reqObj any) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(reqObj)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/callbackExample/"+command, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	return w
}

func TestAllWebhookCallbacks_UserAndStatus(t *testing.T) {
	t.Parallel()
	srv := NewServer()

	t.Run("UserRegister", func(t *testing.T) {
		w := runWebhookTest(t, srv, cbapi.CallbackBeforeUserRegisterCommand, &cbapi.CallbackBeforeUserRegisterReq{
			CallbackCommand: cbapi.CallbackCommand(cbapi.CallbackBeforeUserRegisterCommand),
			Users:           []*sdkws.UserInfo{{UserID: "u1", Nickname: "User1"}},
		})
		var resp cbapi.CallbackBeforeUserRegisterResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.Len(t, resp.Users, 1)
		require.Equal(t, "User1", resp.Users[0].Nickname)

		runWebhookTest(t, srv, cbapi.CallbackAfterUserRegisterCommand, &cbapi.CallbackAfterUserRegisterReq{
			CallbackCommand: cbapi.CallbackCommand(cbapi.CallbackAfterUserRegisterCommand),
			Users:           []*sdkws.UserInfo{{UserID: "u1"}},
		})
	})

	t.Run("UpdateUserInfo", func(t *testing.T) {
		newNick := "NewNick"
		w := runWebhookTest(t, srv, cbapi.CallbackBeforeUpdateUserInfoCommand, &cbapi.CallbackBeforeUpdateUserInfoReq{
			CallbackCommand: cbapi.CallbackCommand(cbapi.CallbackBeforeUpdateUserInfoCommand),
			UserID:          "u1",
			Nickname:        &newNick,
		})
		var resp cbapi.CallbackBeforeUpdateUserInfoResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.NotNil(t, resp.Nickname)
		require.Equal(t, "NewNick", *resp.Nickname)

		runWebhookTest(t, srv, cbapi.CallbackAfterUpdateUserInfoCommand, &cbapi.CallbackAfterUpdateUserInfoReq{
			CallbackCommand: cbapi.CallbackCommand(cbapi.CallbackAfterUpdateUserInfoCommand),
			UserID:          "u1",
			Nickname:        "NewNick",
		})
	})

	t.Run("UpdateUserInfoEx", func(t *testing.T) {
		w := runWebhookTest(t, srv, cbapi.CallbackBeforeUpdateUserInfoExCommand, &cbapi.CallbackBeforeUpdateUserInfoExReq{
			CallbackCommand: cbapi.CallbackCommand(cbapi.CallbackBeforeUpdateUserInfoExCommand),
			UserID:          "u1",
			Ex:              wrapperspb.String("new-ex"),
		})
		var resp cbapi.CallbackBeforeUpdateUserInfoExResp
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.NotNil(t, resp.Ex)
		require.Equal(t, "new-ex", resp.Ex.Value)

		runWebhookTest(t, srv, cbapi.CallbackAfterUpdateUserInfoExCommand, &cbapi.CallbackAfterUpdateUserInfoExReq{
			CallbackCommand: cbapi.CallbackCommand(cbapi.CallbackAfterUpdateUserInfoExCommand),
			UserID:          "u1",
			Ex:              wrapperspb.String("new-ex"),
		})
	})

	t.Run("UserStatusGateway", func(t *testing.T) {
		for _, cmd := range []string{
			cbapi.CallbackAfterUserOnlineCommand,
			cbapi.CallbackAfterUserOfflineCommand,
			cbapi.CallbackAfterUserKickOffCommand,
		} {
			runWebhookTest(t, srv, cmd, &cbapi.CallbackUserOnlineReq{
				UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
					UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
						CallbackCommand: cmd,
						Platform:        "iOS",
					},
					UserID: "u1",
				},
			})
		}
	})
}

func TestAllWebhookCallbacks_FriendAndGroup(t *testing.T) {
	t.Parallel()
	srv := NewServer()

	t.Run("FriendCallbacks", func(t *testing.T) {
		friendCmds := []string{
			cbapi.CallbackBeforeAddFriendCommand,
			cbapi.CallbackAfterAddFriendCommand,
			cbapi.CallbackBeforeAddFriendAgreeCommand,
			cbapi.CallbackAfterAddFriendAgreeCommand,
			cbapi.CallbackAfterDeleteFriendCommand,
			cbapi.CallbackBeforeImportFriendsCommand,
			cbapi.CallbackAfterImportFriendsCommand,
			cbapi.CallbackBeforeSetFriendRemarkCommand,
			cbapi.CallbackAfterSetFriendRemarkCommand,
			cbapi.CallbackBeforeAddBlackCommand,
			cbapi.CallbackAfterRemoveBlackCommand,
		}
		for _, cmd := range friendCmds {
			runWebhookTest(t, srv, cmd, &cbapi.CallbackBeforeAddFriendReq{
				CallbackCommand: cbapi.CallbackCommand(cmd),
				FromUserID:      "u1",
				ToUserID:        "u2",
			})
		}
	})

	t.Run("GroupCallbacks", func(t *testing.T) {
		groupCmds := []string{
			cbapi.CallbackBeforeCreateGroupCommand,
			cbapi.CallbackAfterCreateGroupCommand,
			cbapi.CallbackBeforeMembersJoinGroupCommand,
			cbapi.CallbackBeforeJoinGroupCommand,
			cbapi.CallbackAfterJoinGroupCommand,
			cbapi.CallbackAfterQuitGroupCommand,
			cbapi.CallbackAfterKickGroupCommand,
			cbapi.CallbackAfterDisMissGroupCommand,
			cbapi.CallbackAfterTransferGroupOwnerCommand,
			cbapi.CallbackBeforeInviteJoinGroupCommand,
			cbapi.CallbackBeforeSetGroupInfoCommand,
			cbapi.CallbackAfterSetGroupInfoCommand,
			cbapi.CallbackBeforeSetGroupInfoExCommand,
			cbapi.CallbackAfterSetGroupInfoExCommand,
			cbapi.CallbackBeforeSetGroupMemberInfoCommand,
			cbapi.CallbackAfterSetGroupMemberInfoCommand,
		}
		for _, cmd := range groupCmds {
			runWebhookTest(t, srv, cmd, &cbapi.CallbackQuitGroupReq{
				CallbackCommand: cbapi.CallbackCommand(cmd),
				GroupID:         "g1",
				UserID:          "u1",
			})
		}
	})
}

func TestAllWebhookCallbacks_MsgAndPush(t *testing.T) {
	t.Parallel()
	srv := NewServer()

	msgCmds := []string{
		cbapi.CallbackBeforeSendSingleMsgCommand,
		cbapi.CallbackAfterSendSingleMsgCommand,
		cbapi.CallbackBeforeSendGroupMsgCommand,
		cbapi.CallbackAfterSendGroupMsgCommand,
		cbapi.CallbackBeforeMsgModifyCommand,
		cbapi.CallbackAfterSingleMsgReadCommand,
		cbapi.CallbackAfterGroupMsgReadCommand,
		cbapi.CallbackAfterMsgSaveDBCommand,
		cbapi.CallbackAfterRevokeMsgCommand,
		cbapi.CallbackBeforeOfflinePushCommand,
		cbapi.CallbackBeforeOnlinePushCommand,
		cbapi.CallbackBeforeGroupOnlinePushCommand,
	}

	for _, cmd := range msgCmds {
		cmd := cmd
		t.Run(cmd, func(t *testing.T) {
			t.Parallel()
			runWebhookTest(t, srv, cmd, &cbapi.CallbackAfterMsgSaveDBReq{
				CommonCallbackReq: cbapi.CommonCallbackReq{
					SendID:          "u1",
					CallbackCommand: cmd,
				},
				RecvID:  "u2",
				GroupID: "g1",
			})
		})
	}
}

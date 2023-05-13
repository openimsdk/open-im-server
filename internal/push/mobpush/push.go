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

package mobpush

import (
	"Open_IM/internal/push"
	"Open_IM/internal/push/mobpush/common"
	"Open_IM/internal/push/mobpush/requestParams"
	"Open_IM/pkg/common/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	MobPushClient *MobPush
)

func init() {
	MobPushClient = newGetuiClient()
}

type MobPush struct{}

func newGetuiClient() *MobPush {
	return &MobPush{}
}

func (j *MobPush) Push(accounts []string, alert, detailContent, operationID string, opts push.PushOpts) (string, error) {

	var target requestParams.PushTarget

	target.SetAlias(accounts)
	target.SetTarget(2)

	var no requestParams.PushNotify
	no.SetType(1)
	no.SetIosProduction(1)
	no.SetPlats([]int{1, 2})
	no.SetContent(alert)

	var forward requestParams.PushForward
	forward.SetNextType(2)
	forward.SetScheme(config.Config.Push.Mob.Scheme)

	var po requestParams.PushObj
	po.SetSource("webapi")
	po.SetAppkey(config.Config.Push.Mob.AppKey)
	po.SetPushTarget(&target)
	po.SetPushNotify(&no)
	po.SetPushForward(&forward)

	con, err := json.Marshal(po)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", config.Config.Push.Mob.PushUrl, strings.NewReader(string(con)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("key", config.Config.Push.Mob.AppKey)
	req.Header.Set("sign", common.GetSign(string(con)))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

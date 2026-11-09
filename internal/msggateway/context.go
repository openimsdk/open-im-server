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

package msggateway

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/utils/encrypt"
	"github.com/openimsdk/tools/utils/timeutil"
)

type UserConnContextInfo struct {
	Token        string `json:"token"`
	UserID       string `json:"userID"`
	PlatformID   int    `json:"platformID"`
	OperationID  string `json:"operationID"`
	Compression  string `json:"compression"`
	SDKType      string `json:"sdkType"`
	SendResponse bool   `json:"sendResponse"`
	Background   bool   `json:"background"`
	SDKVersion   string `json:"sdkVersion"`
}

type UserConnContext struct {
	RespWriter http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	RemoteAddr string
	ConnID     string
	info       *UserConnContextInfo
}

func (c *UserConnContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *UserConnContext) Done() <-chan struct{} {
	return nil
}

func (c *UserConnContext) Err() error {
	return nil
}

func (c *UserConnContext) Value(key any) any {
	switch key {
	case constant.OpUserID:
		return c.GetUserID()
	case constant.OperationID:
		return c.GetOperationID()
	case constant.ConnID:
		return c.GetConnID()
	case constant.OpUserPlatform:
		return c.GetPlatformID()
	case constant.RemoteAddr:
		return c.RemoteAddr
	case SDKVersion:
		return c.info.SDKVersion
	default:
		return ""
	}
}

func newContext(respWriter http.ResponseWriter, req *http.Request) *UserConnContext {
	remoteAddr := req.RemoteAddr
	if forwarded := req.Header.Get("X-Forwarded-For"); forwarded != "" {
		remoteAddr += "_" + forwarded
	}
	return &UserConnContext{
		RespWriter: respWriter,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		RemoteAddr: remoteAddr,
		ConnID:     encrypt.Md5(req.RemoteAddr + "_" + strconv.Itoa(int(timeutil.GetCurrentTimestampByMill()))),
	}
}

func newTempContext() *UserConnContext {
	return &UserConnContext{
		Req:  &http.Request{URL: &url.URL{}},
		info: &UserConnContextInfo{},
	}
}

func (c *UserConnContext) ParseEssentialArgs() error {
	query := c.Req.URL.Query()
	if data := query.Get("v"); data != "" {
		return c.parseByJson(data)
	} else {
		return c.parseByQuery(query, c.Req.Header)
	}
}

func (c *UserConnContext) parseByQuery(query url.Values, header http.Header) error {
	info := UserConnContextInfo{
		Token:       query.Get(Token),
		UserID:      query.Get(WsUserID),
		OperationID: query.Get(OperationID),
		Compression: query.Get(Compression),
		SDKType:     query.Get(SDKType),
		SDKVersion:  query.Get(SDKVersion),
	}
	platformID, err := strconv.Atoi(query.Get(PlatformID))
	if err != nil {
		return servererrs.ErrConnArgsErr.WrapMsg("platformID is not int")
	}
	info.PlatformID = platformID
	if val := query.Get(SendResponse); val != "" {
		ok, err := strconv.ParseBool(val)
		if err != nil {
			return servererrs.ErrConnArgsErr.WrapMsg("isMsgResp is not bool")
		}
		info.SendResponse = ok
	}
	if info.Compression == "" {
		info.Compression = header.Get(Compression)
	}
	background, err := strconv.ParseBool(query.Get(BackgroundStatus))
	if err != nil {
		return err
	}
	info.Background = background
	return c.checkInfo(&info)
}

func (c *UserConnContext) parseByJson(data string) error {
	reqInfo, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return servererrs.ErrConnArgsErr.WrapMsg("data is not base64")
	}
	var info UserConnContextInfo
	if err := json.Unmarshal(reqInfo, &info); err != nil {
		return servererrs.ErrConnArgsErr.WrapMsg("data is not json", "info", err.Error())
	}
	return c.checkInfo(&info)
}

func (c *UserConnContext) checkInfo(info *UserConnContextInfo) error {
	if info.OperationID == "" {
		return servererrs.ErrConnArgsErr.WrapMsg("operationID is empty")
	}
	if info.Token == "" {
		return servererrs.ErrConnArgsErr.WrapMsg("token is empty")
	}
	if info.UserID == "" {
		return servererrs.ErrConnArgsErr.WrapMsg("sendID is empty")
	}
	if _, ok := constant.PlatformID2Name[info.PlatformID]; !ok {
		return servererrs.ErrConnArgsErr.WrapMsg("platformID is invalid")
	}
	switch info.SDKType {
	case "":
		info.SDKType = GoSDK
	case GoSDK, JsSDK:
	default:
		return servererrs.ErrConnArgsErr.WrapMsg("sdkType is invalid")
	}
	c.info = info
	return nil
}

func (c *UserConnContext) GetRemoteAddr() string {
	return c.RemoteAddr
}

func (c *UserConnContext) SetHeader(key, value string) {
	c.RespWriter.Header().Set(key, value)
}

func (c *UserConnContext) ErrReturn(error string, code int) {
	http.Error(c.RespWriter, error, code)
}

func (c *UserConnContext) GetConnID() string {
	return c.ConnID
}

func (c *UserConnContext) GetUserID() string {
	if c == nil || c.info == nil {
		return ""
	}
	return c.info.UserID
}

func (c *UserConnContext) GetPlatformID() int {
	if c == nil || c.info == nil {
		return 0
	}
	return c.info.PlatformID
}

func (c *UserConnContext) GetOperationID() string {
	if c == nil || c.info == nil {
		return ""
	}
	return c.info.OperationID
}

func (c *UserConnContext) SetOperationID(operationID string) {
	if c.info == nil {
		c.info = &UserConnContextInfo{}
	}
	c.info.OperationID = operationID
}

func (c *UserConnContext) GetToken() string {
	if c == nil || c.info == nil {
		return ""
	}
	return c.info.Token
}

func (c *UserConnContext) GetCompression() bool {
	return c != nil && c.info != nil && c.info.Compression == GzipCompressionProtocol
}

func (c *UserConnContext) GetSDKType() string {
	if c == nil || c.info == nil {
		return GoSDK
	}
	switch c.info.SDKType {
	case "", GoSDK:
		return GoSDK
	case JsSDK:
		return JsSDK
	default:
		return ""
	}
}

func (c *UserConnContext) GetSDKVersion() string {
	if c == nil || c.info == nil {
		return ""
	}
	return c.info.SDKVersion
}

func (c *UserConnContext) ShouldSendResp() bool {
	return c != nil && c.info != nil && c.info.SendResponse
}

func (c *UserConnContext) SetToken(token string) {
	if c.info == nil {
		c.info = &UserConnContextInfo{}
	}
	c.info.Token = token
}

func (c *UserConnContext) GetBackground() bool {
	return c != nil && c.info != nil && c.info.Background
}

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

package oss

import (
	"net/http"
	"net/url"
	_ "unsafe"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//go:linkname signHeader github.com/aliyun/aliyun-oss-go-sdk/oss.Conn.signHeader
func signHeader(c oss.Conn, req *http.Request, canonicalizedResource string)

//go:linkname getURLParams github.com/aliyun/aliyun-oss-go-sdk/oss.Conn.getURLParams
func getURLParams(c oss.Conn, params map[string]any) string

//go:linkname getURL github.com/aliyun/aliyun-oss-go-sdk/oss.urlMaker.getURL
func getURL(um urlMaker, bucket, object, params string) *url.URL

type urlMaker struct {
	Scheme  string
	NetLoc  string
	Type    int
	IsProxy bool
}

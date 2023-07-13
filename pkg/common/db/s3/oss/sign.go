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
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"hash"
	"io"
	"net/http"
	"sort"
	"strings"
)

func (o *OSS) getAdditionalHeaderKeys(req *http.Request) ([]string, map[string]string) {
	var keysList []string
	keysMap := make(map[string]string)
	srcKeys := make(map[string]string)

	for k := range req.Header {
		srcKeys[strings.ToLower(k)] = ""
	}

	for _, v := range o.bucket.Client.Config.AdditionalHeaders {
		if _, ok := srcKeys[strings.ToLower(v)]; ok {
			keysMap[strings.ToLower(v)] = ""
		}
	}

	for k := range keysMap {
		keysList = append(keysList, k)
	}
	sort.Strings(keysList)
	return keysList, keysMap
}

func (o *OSS) getSignedStr(req *http.Request, canonicalizedResource string, keySecret string) string {
	// Find out the "x-oss-"'s address in header of the request
	ossHeadersMap := make(map[string]string)
	additionalList, additionalMap := o.getAdditionalHeaderKeys(req)
	for k, v := range req.Header {
		if strings.HasPrefix(strings.ToLower(k), "x-oss-") {
			ossHeadersMap[strings.ToLower(k)] = v[0]
		} else if o.bucket.Client.Config.AuthVersion == oss.AuthV2 {
			if _, ok := additionalMap[strings.ToLower(k)]; ok {
				ossHeadersMap[strings.ToLower(k)] = v[0]
			}
		}
	}
	hs := newHeaderSorter(ossHeadersMap)

	// Sort the ossHeadersMap by the ascending order
	hs.Sort()

	// Get the canonicalizedOSSHeaders
	canonicalizedOSSHeaders := ""
	for i := range hs.Keys {
		canonicalizedOSSHeaders += hs.Keys[i] + ":" + hs.Vals[i] + "\n"
	}

	// Give other parameters values
	// when sign URL, date is expires
	date := req.Header.Get(oss.HTTPHeaderDate)
	contentType := req.Header.Get(oss.HTTPHeaderContentType)
	contentMd5 := req.Header.Get(oss.HTTPHeaderContentMD5)

	// default is v1 signature
	signStr := req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + canonicalizedResource
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(keySecret))

	// v2 signature
	if o.bucket.Client.Config.AuthVersion == oss.AuthV2 {
		signStr = req.Method + "\n" + contentMd5 + "\n" + contentType + "\n" + date + "\n" + canonicalizedOSSHeaders + strings.Join(additionalList, ";") + "\n" + canonicalizedResource
		h = hmac.New(func() hash.Hash { return sha256.New() }, []byte(keySecret))
	}
	_, _ = io.WriteString(h, signStr)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signedStr
}

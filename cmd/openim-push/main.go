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

package main

import (
	"github.com/openimsdk/open-im-server/v3/internal/push"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func main() {
	pushCmd := cmd.NewRpcCmd(cmd.RpcPushServer)
	pushCmd.AddPortFlag()
	pushCmd.AddPrometheusPortFlag()
	if err := pushCmd.Exec(); err != nil {# How do I contribute code to OpenIM

		<p align="center">
		  <a href="./CONTRIBUTING.md">Englist</a> · 
		  <a href="./CONTRIBUTING-zh_CN.md">中文</a> · 
		  <a href="docs/contributing/CONTRIBUTING-UA.md">Українська</a> · 
		  <a href="docs/contributing/CONTRIBUTING-CS.md">Česky</a> · 
		  <a href="docs/contributing/CONTRIBUTING-HU.md">Magyar</a> · 
		  <a href="docs/contributing/CONTRIBUTING-ES.md">Español</a> · 
		  <a href="docs/contributing/CONTRIBUTING-FA.md">فارسی</a> · 
		  <a href="docs/contributing/CONTRIBUTING-FR.md">Français</a> · 
		  <a href="docs/contributing/CONTRIBUTING-DE.md">Deutsch</a> · 
		  <a href="docs/contributing/CONTRIBUTING-PL.md">Polski</a> · 
		  <a href="docs/contributing/CONTRIBUTING-ID.md">Indonesian</a> · 
		  <a href="docs/contributing/CONTRIBUTING-FI.md">Suomi</a> · 
		  <a href="docs/contributing/CONTRIBUTING-ML.md">മലയാളം</a> · 
		  <a href="docs/contributing/CONTRIBUTING-JP.md">日本語</a> · 
		  <a href="docs/contributing/CONTRIBUTING-NL.md">Nederlands</a> · 
		  <a href="docs/contributing/CONTRIBUTING-IT.md">Italiano</a> · 
		  <a href="docs/contributing/CONTRIBUTING-RU.md">Русский</a> · 
		  <a href="docs/contributing/CONTRIBUTING-PTBR.md">Português (Brasil)</a> · 
		  <a href="docs/contributing/CONTRIBUTING-EO.md">Esperanto</a> · 
		  <a href="docs/contributing/CONTRIBUTING-KR.md">한국어</a> · 
		  <a href="docs/contributing/CONTRIBUTING-AR.md">العربي</a> · 
		  <a href="docs/contributing/CONTRIBUTING-VN.md">Tiếng Việt</a> · 
		  <a href="docs/contributing/CONTRIBUTING-DA.md">Dansk</a> · 
		  <a href="docs/contributing/CONTRIBUTING-GR.md">Ελληνικά</a> · 
		  <a href="docs/contributing/CONTRIBUTING-TR.md">Türkçe</a>
		</p>
		
		</div>
		
		</p>
		
		panic(err.Error())
	}
	if err := pushCmd.StartSvr(config.Config.RpcRegisterName.OpenImPushName, push.Start); err != nil {
		panic(err.Error())
	}
}

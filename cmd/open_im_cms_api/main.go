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
	"Open_IM/internal/cms_api"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/utils"
	"flag"
	"fmt"
	"strconv"

	"Open_IM/pkg/common/config"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := cms_api.NewGinRouter()
	router.Use(utils.CorsHandler())
	defaultPorts := config.Config.CmsApi.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10006 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.Config.CmsApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start cms api server, address: ", address, ", OpenIM version: ", constant.CurrentVersion, "\n")
	router.Run(address)
}

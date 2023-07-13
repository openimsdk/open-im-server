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

package rpcclient

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/url"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
)

type Third struct {
	conn        grpc.ClientConnInterface
	Client      third.ThirdClient
	discov      discoveryregistry.SvcDiscoveryRegistry
	MinioClient *minio.Client
}

func NewThird(discov discoveryregistry.SvcDiscoveryRegistry) *Third {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImThirdName)
	if err != nil {
		panic(err)
	}
	client := third.NewThirdClient(conn)
	minioClient, err := minioInit()
	return &Third{discov: discov, Client: client, conn: conn, MinioClient: minioClient}
}

func minioInit() (*minio.Client, error) {
	minioClient := &minio.Client{}
	var initUrl string
	initUrl = config.Config.Object.Minio.Endpoint
	minioUrl, err := url.Parse(initUrl)
	if err != nil {
		return nil, err
	}
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(config.Config.Object.Minio.AccessKeyID, config.Config.Object.Minio.SecretAccessKey, ""),
		//Region: config.Config.Credential.Minio.Location,
	}
	if minioUrl.Scheme == "http" {
		opts.Secure = false
	} else if minioUrl.Scheme == "https" {
		opts.Secure = true
	}
	minioClient, err = minio.New(minioUrl.Host, opts)
	if err != nil {
		return nil, err
	}
	return minioClient, nil
}

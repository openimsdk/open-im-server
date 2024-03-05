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
	"github.com/OpenIMSDK/tools/errs"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/OpenIMSDK/protocol/third"
	"github.com/OpenIMSDK/tools/discoveryregistry"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	util "github.com/openimsdk/open-im-server/v3/pkg/util/genutil"
	"google.golang.org/grpc"
)

type Third struct {
	conn        grpc.ClientConnInterface
	Client      third.ThirdClient
	discov      discoveryregistry.SvcDiscoveryRegistry
	MinioClient *minio.Client
	Config      *config.GlobalConfig
}

func NewThird(discov discoveryregistry.SvcDiscoveryRegistry, config *config.GlobalConfig) *Third {
	conn, err := discov.GetConn(context.Background(), config.RpcRegisterName.OpenImThirdName)
	if err != nil {
		util.ExitWithError(err)
	}
	client := third.NewThirdClient(conn)
	minioClient, err := minioInit(config)
	if err != nil {
		util.ExitWithError(err)
	}
	return &Third{discov: discov, Client: client, conn: conn, MinioClient: minioClient, Config: config}
}

func minioInit(config *config.GlobalConfig) (*minio.Client, error) {
	minioClient := &minio.Client{}
	initUrl := config.Object.Minio.Endpoint
	minioUrl, err := url.Parse(initUrl)
	if err != nil {
		return nil, errs.Wrap(err, "minioInit: failed to parse MinIO endpoint URL")
	}
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(config.Object.Minio.AccessKeyID, config.Object.Minio.SecretAccessKey, ""),
		// Region: config.Credential.Minio.Location,
	}
	if minioUrl.Scheme == "http" {
		opts.Secure = false
	} else if minioUrl.Scheme == "https" {
		opts.Secure = true
	}
	minioClient, err = minio.New(minioUrl.Host, opts)
	if err != nil {
		return nil, errs.Wrap(err, "minioInit: failed to create MinIO client")
	}
	return minioClient, nil
}

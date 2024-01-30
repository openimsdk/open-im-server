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
	"os"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/protocol/third"
	"github.com/OpenIMSDK/tools/discoveryregistry"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
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
	initUrl = getMinioAddr("MINIO_ENDPOINT", "MINIO_ADDRESS", "MINIO_PORT", config.Config.Object.Minio.Endpoint)
	minioUrl, err := url.Parse(initUrl)
	if err != nil {
		return nil, err
	}
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(config.Config.Object.Minio.AccessKeyID, config.Config.Object.Minio.SecretAccessKey, ""),
		// Region: config.Config.Credential.Minio.Location,
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

func getMinioAddr(key1, key2, key3, fallback string) string {
	// Prioritize environment variables
	endpoint, endpointExist := os.LookupEnv(key1)
	if !endpointExist {
		endpoint = fallback
	}
	address, addressExist := os.LookupEnv(key2)
	port, portExist := os.LookupEnv(key3)
	if portExist && addressExist {
		endpoint = "http://" + address + ":" + port
	}
	return endpoint
}

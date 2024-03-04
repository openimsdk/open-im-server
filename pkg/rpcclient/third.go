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
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/grpc"

	"github.com/OpenIMSDK/protocol/third"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	util "github.com/openimsdk/open-im-server/v3/pkg/util/genutil"
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
		util.ExitWithError(err)
	}
	client := third.NewThirdClient(conn)
	minioClient, err := minioInit()
	if err != nil {
		util.ExitWithError(err)
	}
	return &Third{discov: discov, Client: client, conn: conn, MinioClient: minioClient}
}

func minioInit() (*minio.Client, error) {
	// Retrieve MinIO configuration details
	endpoint := config.Config.Object.Minio.Endpoint
	accessKeyID := config.Config.Object.Minio.AccessKeyID
	secretAccessKey := config.Config.Object.Minio.SecretAccessKey

	// Parse the MinIO URL to determine if the connection should be secure
	minioURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, errs.Wrap(err, "minioInit: failed to parse MinIO endpoint URL")
	}

	// Determine the security of the connection based on the scheme
	secure := minioURL.Scheme == "https"

	// Setup MinIO client options
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	}

	// Initialize MinIO client
	minioClient, err := minio.New(minioURL.Host, opts)
	if err != nil {
		return nil, errs.Wrap(err, "minioInit: failed to create MinIO client")
	}

	return minioClient, nil
}

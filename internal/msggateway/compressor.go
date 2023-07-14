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
	"bytes"
	"compress/gzip"
	"io"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type Compressor interface {
	Compress(rawData []byte) ([]byte, error)
	DeCompress(compressedData []byte) ([]byte, error)
}
type GzipCompressor struct {
	compressProtocol string
}

func NewGzipCompressor() *GzipCompressor {
	return &GzipCompressor{compressProtocol: "gzip"}
}

func (g *GzipCompressor) Compress(rawData []byte) ([]byte, error) {
	gzipBuffer := bytes.Buffer{}
	gz := gzip.NewWriter(&gzipBuffer)
	if _, err := gz.Write(rawData); err != nil {
		return nil, utils.Wrap(err, "")
	}
	if err := gz.Close(); err != nil {
		return nil, utils.Wrap(err, "")
	}
	return gzipBuffer.Bytes(), nil
}

func (g *GzipCompressor) DeCompress(compressedData []byte) ([]byte, error) {
	buff := bytes.NewBuffer(compressedData)
	reader, err := gzip.NewReader(buff)
	if err != nil {
		return nil, utils.Wrap(err, "NewReader failed")
	}
	compressedData, err = io.ReadAll(reader)
	if err != nil {
		return nil, utils.Wrap(err, "ReadAll failed")
	}
	_ = reader.Close()
	return compressedData, nil
}

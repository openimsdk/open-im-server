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
	"errors"
	"io"
	"sync"

	"github.com/OpenIMSDK/tools/utils"
)

var (
	gzipWriterPool = sync.Pool{New: func() any { return gzip.NewWriter(nil) }}
	gzipReaderPool = sync.Pool{New: func() any { return new(gzip.Reader) }}
)

type Compressor interface {
	Compress(rawData []byte) ([]byte, error)
	CompressWithPool(rawData []byte) ([]byte, error)
	DeCompress(compressedData []byte) ([]byte, error)
	DecompressWithPool(compressedData []byte) ([]byte, error)
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

func (g *GzipCompressor) CompressWithPool(rawData []byte) ([]byte, error) {
	gz := gzipWriterPool.Get().(*gzip.Writer)
	defer gzipWriterPool.Put(gz)

	gzipBuffer := bytes.Buffer{}
	gz.Reset(&gzipBuffer)

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

func (g *GzipCompressor) DecompressWithPool(compressedData []byte) ([]byte, error) {
	reader := gzipReaderPool.Get().(*gzip.Reader)
	if reader == nil {
		return nil, errors.New("NewReader failed")
	}
	defer gzipReaderPool.Put(reader)

	err := reader.Reset(bytes.NewReader(compressedData))
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

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

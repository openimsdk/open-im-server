package minio

import (
	"bytes"
	"context"
	"errors"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	config.Config.Object.Minio.Bucket = "openim"
	config.Config.Object.Minio.AccessKeyID = "root"
	config.Config.Object.Minio.SecretAccessKey = "openIM123"
	config.Config.Object.Minio.Endpoint = "http://172.16.8.38:10005"
	tmp, err := NewMinio(nil)
	if err != nil {
		panic(err)
	}
	min := tmp.(*Minio)

	text := []byte("hello world!")
	name := "posttest.txt"

	u, err := min.FormData(context.Background(), "posttest.txt", int64(len(text)), "image/png", time.Second*1000)
	if err != nil {
		panic(err)
	}
	t.Log(u.URL)
	for k, v := range u.FormData {
		t.Log(k, v)
	}
	if err := PostFile(u, name, text); err != nil {
		t.Error(err)
	}
}

func PostFile(fd *s3.FormData, name string, data []byte) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for k, v := range fd.FormData {
		if err := writer.WriteField(k, v); err != nil {
			return err
		}
	}
	fileWriter, err := writer.CreateFormFile(fd.File, path.Base(name))
	if err != nil {
		return err
	}
	if _, err := fileWriter.Write(data); err != nil {
		return nil
	}
	defer writer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	reqBody := body.Bytes()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fd.URL, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = int64(len(reqBody))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(string(respBody))
	}
	return nil
}

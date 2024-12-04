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

package pkg

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/errs"
)

type Upload struct {
	URL         string `json:"url"`
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
}

type Task struct {
	Index  int
	Upload Upload
}

type PartInfo struct {
	ContentType string
	PartSize    int64
	PartNum     int
	FileMd5     string
	PartMd5     string
	PartSizes   []int64
	PartMd5s    []string
}

func Run(conf Config) error {
	m := &Manage{
		prefix: time.Now().Format("20060102150405"),
		conf:   &conf,
		ctx:    context.Background(),
	}
	return m.Run()
}

type Manage struct {
	conf      *Config
	ctx       context.Context
	api       *Api
	partLimit *third.PartLimitResp
	prefix    string
	tasks     chan Task
	id        uint64
	success   int64
	failed    int64
}

func (m *Manage) tempFilePath() string {
	return filepath.Join(m.conf.TempDir, fmt.Sprintf("%s_%d", m.prefix, atomic.AddUint64(&m.id, 1)))
}

func (m *Manage) Run() error {
	defer func(start time.Time) {
		log.Printf("run time %s\n", time.Since(start))
	}(time.Now())
	m.api = &Api{
		Api:    m.conf.Api,
		UserID: m.conf.UserID,
		Secret: m.conf.Secret,
		Client: &http.Client{Timeout: m.conf.Timeout},
	}
	var err error
	ctx := context.WithValue(m.ctx, "operationID", fmt.Sprintf("%s_init", m.prefix))
	m.api.Token, err = m.api.GetAdminToken(ctx)
	if err != nil {
		return err
	}
	m.partLimit, err = m.api.GetPartLimit(ctx)
	if err != nil {
		return err
	}
	progress, err := ReadProgress(m.conf.ProgressPath)
	if err != nil {
		return err
	}
	progressFile, err := os.OpenFile(m.conf.ProgressPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	var mutex sync.Mutex
	writeSuccessIndex := func(index int) {
		mutex.Lock()
		defer mutex.Unlock()
		if _, err := progressFile.Write([]byte(strconv.Itoa(index) + "\n")); err != nil {
			log.Printf("write progress err: %v\n", err)
		}
	}
	file, err := os.Open(m.conf.TaskPath)
	if err != nil {
		return err
	}
	m.tasks = make(chan Task, m.conf.Concurrency*2)
	go func() {
		defer file.Close()
		defer close(m.tasks)
		scanner := bufio.NewScanner(file)
		var (
			index int
			num   int
		)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			index++
			if progress.IsUploaded(index) {
				log.Printf("index: %d already uploaded %s\n", index, line)
				continue
			}
			var upload Upload
			if err := json.Unmarshal([]byte(line), &upload); err != nil {
				log.Printf("index: %d json.Unmarshal(%s) err: %v", index, line, err)
				continue
			}
			num++
			m.tasks <- Task{
				Index:  index,
				Upload: upload,
			}
		}
		if num == 0 {
			log.Println("mark all completed")
		}
	}()
	var wg sync.WaitGroup
	wg.Add(m.conf.Concurrency)
	for i := 0; i < m.conf.Concurrency; i++ {
		go func(tid int) {
			defer wg.Done()
			for task := range m.tasks {
				var success bool
				for n := 0; n < m.conf.Retry; n++ {
					ctx := context.WithValue(m.ctx, "operationID", fmt.Sprintf("%s_%d_%d_%d", m.prefix, tid, task.Index, n+1))
					if urlRaw, err := m.RunTask(ctx, task); err == nil {
						writeSuccessIndex(task.Index)
						log.Println("index:", task.Index, "upload success", "urlRaw", urlRaw)
						success = true
						break
					} else {
						log.Printf("index: %d upload: %+v err: %v", task.Index, task.Upload, err)
					}
				}
				if success {
					atomic.AddInt64(&m.success, 1)
				} else {
					atomic.AddInt64(&m.failed, 1)
					log.Printf("index: %d upload: %+v failed", task.Index, task.Upload)
				}
			}
		}(i + 1)
	}
	wg.Wait()
	log.Printf("execution completed success %d failed %d\n", m.success, m.failed)
	return nil
}

func (m *Manage) RunTask(ctx context.Context, task Task) (string, error) {
	resp, err := m.HttpGet(ctx, task.Upload.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	reader, err := NewReader(resp.Body, m.conf.CacheSize, m.tempFilePath())
	if err != nil {
		return "", err
	}
	defer reader.Close()
	part, err := m.getPartInfo(ctx, reader, reader.Size())
	if err != nil {
		return "", err
	}
	var contentType string
	if task.Upload.ContentType == "" {
		contentType = part.ContentType
	} else {
		contentType = task.Upload.ContentType
	}
	initiateMultipartUploadResp, err := m.api.InitiateMultipartUpload(ctx, &third.InitiateMultipartUploadReq{
		Hash:        part.PartMd5,
		Size:        reader.Size(),
		PartSize:    part.PartSize,
		MaxParts:    -1,
		Cause:       "batch-import",
		Name:        task.Upload.Name,
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}
	if initiateMultipartUploadResp.Upload == nil {
		return initiateMultipartUploadResp.Url, nil
	}
	if _, err := reader.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	uploadParts := make([]*third.SignPart, part.PartNum)
	for _, part := range initiateMultipartUploadResp.Upload.Sign.Parts {
		uploadParts[part.PartNumber-1] = part
	}
	for i, currentPartSize := range part.PartSizes {
		md5Reader := NewMd5Reader(io.LimitReader(reader, currentPartSize))
		if err := m.doPut(ctx, m.api.Client, initiateMultipartUploadResp.Upload.Sign, uploadParts[i], md5Reader, currentPartSize); err != nil {
			return "", err
		}
		if md5val := md5Reader.Md5(); md5val != part.PartMd5s[i] {
			return "", fmt.Errorf("upload part %d failed, md5 not match, expect %s, got %s", i, part.PartMd5s[i], md5val)
		}
	}
	urlRaw, err := m.api.CompleteMultipartUpload(ctx, &third.CompleteMultipartUploadReq{
		UploadID:    initiateMultipartUploadResp.Upload.UploadID,
		Parts:       part.PartMd5s,
		Name:        task.Upload.Name,
		ContentType: contentType,
		Cause:       "batch-import",
	})
	if err != nil {
		return "", err
	}
	return urlRaw, nil
}

func (m *Manage) partSize(size int64) (int64, error) {
	if size <= 0 {
		return 0, errs.New("size must be greater than 0")
	}
	if size > m.partLimit.MaxPartSize*int64(m.partLimit.MaxNumSize) {
		return 0, errs.New("size must be less than", "size", m.partLimit.MaxPartSize*int64(m.partLimit.MaxNumSize))
	}
	if size <= m.partLimit.MinPartSize*int64(m.partLimit.MaxNumSize) {
		return m.partLimit.MinPartSize, nil
	}
	partSize := size / int64(m.partLimit.MaxNumSize)
	if size%int64(m.partLimit.MaxNumSize) != 0 {
		partSize++
	}
	return partSize, nil
}

func (m *Manage) partMD5(parts []string) string {
	s := strings.Join(parts, ",")
	md5Sum := md5.Sum([]byte(s))
	return hex.EncodeToString(md5Sum[:])
}

func (m *Manage) getPartInfo(ctx context.Context, r io.Reader, fileSize int64) (*PartInfo, error) {
	partSize, err := m.partSize(fileSize)
	if err != nil {
		return nil, err
	}
	partNum := int(fileSize / partSize)
	if fileSize%partSize != 0 {
		partNum++
	}
	partSizes := make([]int64, partNum)
	for i := 0; i < partNum; i++ {
		partSizes[i] = partSize
	}
	partSizes[partNum-1] = fileSize - partSize*(int64(partNum)-1)
	partMd5s := make([]string, partNum)
	buf := make([]byte, 1024*8)
	fileMd5 := md5.New()
	var contentType string
	for i := 0; i < partNum; i++ {
		h := md5.New()
		r := io.LimitReader(r, partSize)
		for {
			if n, err := r.Read(buf); err == nil {
				if contentType == "" {
					contentType = http.DetectContentType(buf[:n])
				}
				h.Write(buf[:n])
				fileMd5.Write(buf[:n])
			} else if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		partMd5s[i] = hex.EncodeToString(h.Sum(nil))
	}
	partMd5Val := m.partMD5(partMd5s)
	fileMd5val := hex.EncodeToString(fileMd5.Sum(nil))
	return &PartInfo{
		ContentType: contentType,
		PartSize:    partSize,
		PartNum:     partNum,
		FileMd5:     fileMd5val,
		PartMd5:     partMd5Val,
		PartSizes:   partSizes,
		PartMd5s:    partMd5s,
	}, nil
}

func (m *Manage) doPut(ctx context.Context, client *http.Client, sign *third.AuthSignParts, part *third.SignPart, reader io.Reader, size int64) error {
	rawURL := part.Url
	if rawURL == "" {
		rawURL = sign.Url
	}
	if len(sign.Query)+len(part.Query) > 0 {
		u, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		query := u.Query()
		for i := range sign.Query {
			v := sign.Query[i]
			query[v.Key] = v.Values
		}
		for i := range part.Query {
			v := part.Query[i]
			query[v.Key] = v.Values
		}
		u.RawQuery = query.Encode()
		rawURL = u.String()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, rawURL, reader)
	if err != nil {
		return err
	}
	for i := range sign.Header {
		v := sign.Header[i]
		req.Header[v.Key] = v.Values
	}
	for i := range part.Header {
		v := part.Header[i]
		req.Header[v.Key] = v.Values
	}
	req.ContentLength = size
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/200 != 1 {
		return fmt.Errorf("PUT %s part %d failed, status code %d, body %s", rawURL, part.PartNumber, resp.StatusCode, string(body))
	}
	return nil
}

func (m *Manage) HttpGet(ctx context.Context, url string) (*http.Response, error) {
	reqUrl := url
	for {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
		if err != nil {
			return nil, err
		}
		DefaultRequestHeader(request.Header)
		response, err := m.api.Client.Do(request)
		if err != nil {
			return nil, err
		}
		if response.StatusCode != http.StatusOK {
			_ = response.Body.Close()
			return nil, fmt.Errorf("webhook get %s status %s", url, response.Status)
		}
		return response, nil
	}
}

package objstorage

import (
	"net/http"
	"time"
)

type PutRes struct {
	URL           string
	Bucket        string
	Name          string
	EffectiveTime time.Time
}

type FragmentPutArgs struct {
	PutArgs
	FragmentSize int64 // 分片大小
}

type PutArgs struct {
	Name          string        // 文件名
	Size          int64         // 大小
	Hash          string        // md5
	Prefix        string        // 前缀
	ClearTime     time.Duration // 自动清理时间
	EffectiveTime time.Duration // 申请有效时间
	Header        http.Header   // header
}

type BucketFile struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

type ObjectInfo struct {
	URL  string
	Size int64
	Hash string
}

//type PutSpace struct {
//	URL           string
//	EffectiveTime time.Time
//}

type PutAddr struct {
	ResourceURL   string
	PutID         string
	FragmentSize  int64
	EffectiveTime time.Time
	PutURLs       []string
}

type KVData struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

type PutResp struct {
	URL  string
	Time *time.Time
}

type ApplyPutArgs struct {
	Bucket    string
	Name      string
	Effective time.Duration // 申请有效时间
	Header    http.Header   // header
}

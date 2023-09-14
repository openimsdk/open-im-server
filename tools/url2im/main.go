package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/openimsdk/open-im-server/v3/tools/url2im/pkg"
)

/*take.txt
{"url":"http://xxx/xxxx","name":"xxxx","contentType":"image/jpeg"}
{"url":"http://xxx/xxxx","name":"xxxx","contentType":"image/jpeg"}
{"url":"http://xxx/xxxx","name":"xxxx","contentType":"image/jpeg"}
*/

func main() {
	var conf pkg.Config                                                  // 后面带*的为必填项
	flag.StringVar(&conf.TaskPath, "task", "take.txt", "task path")      // 任务日志文件*
	flag.StringVar(&conf.ProgressPath, "progress", "", "progress path")  // 进度日志文件
	flag.IntVar(&conf.Concurrency, "concurrency", 1, "concurrency num")  // 并发数
	flag.IntVar(&conf.Retry, "retry", 1, "retry num")                    // 重试次数
	flag.StringVar(&conf.TempDir, "temp", "", "temp dir")                // 临时文件夹
	flag.Int64Var(&conf.CacheSize, "cache", 1024*1024*100, "cache size") // 缓存大小(超过时,下载到磁盘)
	flag.Int64Var((*int64)(&conf.Timeout), "timeout", 5000, "timeout")   // 请求超时时间(毫秒)
	flag.StringVar(&conf.Api, "api", "http://127.0.0.1:10002", "api")    // im地址*
	flag.StringVar(&conf.UserID, "userID", "openIM123456", "userID")     // im管理员
	flag.StringVar(&conf.Secret, "secret", "openIM123", "secret")        // im config secret
	flag.Parse()
	if !filepath.IsAbs(conf.TaskPath) {
		var err error
		conf.TaskPath, err = filepath.Abs(conf.TaskPath)
		if err != nil {
			log.Println("get abs path err:", err)
			return
		}
	}
	if conf.ProgressPath == "" {
		conf.ProgressPath = conf.TaskPath + ".progress.txt"
	} else if !filepath.IsAbs(conf.ProgressPath) {
		var err error
		conf.ProgressPath, err = filepath.Abs(conf.ProgressPath)
		if err != nil {
			log.Println("get abs path err:", err)
			return
		}
	}
	if conf.TempDir == "" {
		conf.TempDir = conf.TaskPath + ".temp"
	}
	if info, err := os.Stat(conf.TempDir); err == nil {
		if !info.IsDir() {
			log.Printf("temp dir %s is not dir\n", err)
			return
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(conf.TempDir, os.ModePerm); err != nil {
			log.Printf("mkdir temp dir %s err %+v\n", conf.TempDir, err)
			return
		}
		defer os.RemoveAll(conf.TempDir)
	} else {
		log.Println("get temp dir err:", err)
		return
	}
	if conf.Concurrency <= 0 {
		conf.Concurrency = 1
	}
	if conf.Retry <= 0 {
		conf.Retry = 1
	}
	if conf.CacheSize <= 0 {
		conf.CacheSize = 1024 * 1024 * 100 // 100M
	}
	if conf.Timeout <= 0 {
		conf.Timeout = 5000
	}
	conf.Timeout = conf.Timeout * time.Millisecond
	if err := pkg.Run(conf); err != nil {
		log.Println("main err:", err)
	}
}

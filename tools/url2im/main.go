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
	var conf pkg.Config // Configuration object, '*' denotes required fields

	// *Required*: Path for the task log file
	flag.StringVar(&conf.TaskPath, "task", "take.txt", "Path for the task log file")

	// Optional: Path for the progress log file
	flag.StringVar(&conf.ProgressPath, "progress", "", "Path for the progress log file")

	// Number of concurrent operations
	flag.IntVar(&conf.Concurrency, "concurrency", 1, "Number of concurrent operations")

	// Number of retry attempts
	flag.IntVar(&conf.Retry, "retry", 1, "Number of retry attempts")

	// Optional: Path for the temporary directory
	flag.StringVar(&conf.TempDir, "temp", "", "Path for the temporary directory")

	// Cache size in bytes (downloads move to disk when exceeded)
	flag.Int64Var(&conf.CacheSize, "cache", 1024*1024*100, "Cache size in bytes")

	// Request timeout in milliseconds
	flag.Int64Var((*int64)(&conf.Timeout), "timeout", 5000, "Request timeout in milliseconds")

	// *Required*: API endpoint for the IM service
	flag.StringVar(&conf.Api, "api", "http://127.0.0.1:10002", "API endpoint for the IM service")

	// IM administrator's user ID
	flag.StringVar(&conf.UserID, "userID", "openIM123456", "IM administrator's user ID")

	// Secret for the IM configuration
	flag.StringVar(&conf.Secret, "secret", "openIM123", "Secret for the IM configuration")

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

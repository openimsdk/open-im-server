package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/s3"
	"github.com/openimsdk/tools/s3/aws"
	"github.com/openimsdk/tools/s3/cos"
	"github.com/openimsdk/tools/s3/kodo"
	"github.com/openimsdk/tools/s3/minio"
	"github.com/openimsdk/tools/s3/oss"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

const defaultTimeout = time.Second * 10

func readConf(path string, val any) error {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	fn := func(config *mapstructure.DecoderConfig) {
		config.TagName = "mapstructure"
	}
	return v.Unmarshal(val, fn)
}

func getS3(path string, name string, thirdConf *config.Third) (s3.Interface, error) {
	switch name {
	case "minio":
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		var minioConf config.Minio
		if err := readConf(filepath.Join(path, minioConf.GetConfigFileName()), &minioConf); err != nil {
			return nil, err
		}
		var redisConf config.Redis
		if err := readConf(filepath.Join(path, redisConf.GetConfigFileName()), &redisConf); err != nil {
			return nil, err
		}
		rdb, err := redisutil.NewRedisClient(ctx, redisConf.Build())
		if err != nil {
			return nil, err
		}
		return minio.NewMinio(ctx, redis.NewMinioCache(rdb), *minioConf.Build())
	case "cos":
		return cos.NewCos(*thirdConf.Object.Cos.Build())
	case "oss":
		return oss.NewOSS(*thirdConf.Object.Oss.Build())
	case "kodo":
		return kodo.NewKodo(*thirdConf.Object.Kodo.Build())
	case "aws":
		return aws.NewAws(*thirdConf.Object.Aws.Build())
	default:
		return nil, fmt.Errorf("invalid object enable: %s", name)
	}
}

func getMongo(path string) (database.ObjectInfo, error) {
	var mongoConf config.Mongo
	if err := readConf(filepath.Join(path, mongoConf.GetConfigFileName()), &mongoConf); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	mgocli, err := mongoutil.NewMongoDB(ctx, mongoConf.Build())
	if err != nil {
		return nil, err
	}
	return mgo.NewS3Mongo(mgocli.GetDB())
}

func Main(path string, engine string) error {
	var thirdConf config.Third
	if err := readConf(filepath.Join(path, thirdConf.GetConfigFileName()), &thirdConf); err != nil {
		return err
	}
	if thirdConf.Object.Enable == engine {
		return errors.New("same s3 storage")
	}
	s3db, err := getMongo(path)
	if err != nil {
		return err
	}
	oldS3, err := getS3(path, engine, &thirdConf)
	if err != nil {
		return err
	}
	newS3, err := getS3(path, thirdConf.Object.Enable, &thirdConf)
	if err != nil {
		return err
	}
	count, err := getEngineCount(s3db, oldS3.Engine())
	if err != nil {
		return err
	}
	log.Printf("engine %s count: %d", oldS3.Engine(), count)
	var skip int
	for i := 1; i <= count+1; i++ {
		log.Printf("start %d/%d", i, count)
		start := time.Now()
		res, err := doObject(s3db, newS3, oldS3, skip)
		if err != nil {
			log.Printf("end [%s] %d/%d error %s", time.Since(start), i, count, err)
			return err
		}
		log.Printf("end [%s] %d/%d result %+v", time.Since(start), i, count, *res)
		if res.Skip {
			skip++
		}
		if res.End {
			break
		}
	}
	return nil
}

func getEngineCount(db database.ObjectInfo, name string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	count, err := db.GetEngineCount(ctx, name)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func doObject(db database.ObjectInfo, newS3, oldS3 s3.Interface, skip int) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	infos, err := db.GetEngineInfo(ctx, oldS3.Engine(), 1, skip)
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return &Result{End: true}, nil
	}
	obj := infos[0]
	if _, err := db.Take(ctx, newS3.Engine(), obj.Name); err == nil {
		return &Result{Skip: true}, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	downloadURL, err := oldS3.AccessURL(ctx, obj.Key, time.Hour, &s3.AccessURLOption{})
	if err != nil {
		return nil, err
	}
	putURL, err := newS3.PresignedPutObject(ctx, obj.Key, time.Hour, &s3.PutOption{ContentType: obj.ContentType})
	if err != nil {
		return nil, err
	}
	downloadResp, err := http.Get(downloadURL)
	if err != nil {
		return nil, err
	}
	defer downloadResp.Body.Close()
	switch downloadResp.StatusCode {
	case http.StatusNotFound:
		return &Result{Skip: true}, nil
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("download object failed %s", downloadResp.Status)
	}
	log.Printf("file size %d", obj.Size)
	request, err := http.NewRequest(http.MethodPut, putURL.URL, downloadResp.Body)
	if err != nil {
		return nil, err
	}
	putResp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer putResp.Body.Close()
	if putResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("put object failed %s", putResp.Status)
	}
	ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	if err := db.UpdateEngine(ctx, obj.Engine, obj.Name, newS3.Engine()); err != nil {
		return nil, err
	}
	return &Result{}, nil
}

type Result struct {
	Skip bool
	End  bool
}

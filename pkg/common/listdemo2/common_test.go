package listdemo

import (
	"context"
	"errors"
	"github.com/openimsdk/tools/db/mongoutil"
	"testing"
)

func Result[V any](val V, err error) V {
	if err != nil {
		panic(err)
	}
	return val
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func TestName(t *testing.T) {
	cli := Result(mongoutil.NewMongoDB(context.Background(), &mongoutil.Config{Uri: "mongodb://openIM:openIM123@172.16.8.48:37017/openim_v3?maxPoolSize=100", Database: "openim_v3"}))

	db := cli.GetDB()
	tx := cli.GetTx()

	const num = 1
	lm := &LogModel{coll: db.Collection("friend_version")}

	err := tx.Transaction(context.Background(), func(ctx context.Context) error {
		err := tx.Transaction(ctx, func(ctx context.Context) error {
			return lm.WriteLogBatch(ctx, "100", []string{"1000", "2000"}, true)
		})
		if err != nil {
			t.Log("--------->")
			return err
		}
		return errors.New("1234")
	})
	t.Log(err)

	//start := time.Now()
	//eIds := make([]string, 0, num)
	//for i := 0; i < num; i++ {
	//	eIds = append(eIds, strconv.Itoa(1000+(i)))
	//}
	//lm.WriteLogBatch1(context.Background(), "100", eIds, false)
	//end := time.Now()
	//t.Log(end.Sub(start))       // 509.962208ms
	//t.Log(end.Sub(start) / num) // 511.496µs

	//start := time.Now()
	//wll, err := lm.FindChangeLog(context.Background(), "100", 3, 100)
	//if err != nil {
	//	panic(err)
	//}
	//t.Log(time.Since(start))
	//t.Log(wll)
}

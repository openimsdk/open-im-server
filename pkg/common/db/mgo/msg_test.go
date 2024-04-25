package mgo

import (
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	//conf := config.Mongo{
	//	Address:  []string{"localhost:37017"},
	//	Username: "openIM",
	//	Password: "openIM123",
	//	Database: "demo",
	//}
	//conf.URI = `mongodb://openIM:openIM123@localhost:37017/demo?maxPoolSize=100&authSource=admin`
	//cli, err := mongoutil.NewMongoDB(context.Background(), conf.Build())
	//if err != nil {
	//	panic(err)
	//}
	//msg, _ := NewMsgMongo(cli.GetDB())
	//count, err := msg.ClearMsg(context.Background(), time.Now().Add(-time.Hour*24*61))
	//if err != nil {
	//	t.Log("error", err)
	//	return
	//}
	//t.Log("count", count)

	s := `si_5300327160_9129042887:0123`

	t.Log(s[:strings.LastIndex(s, ":")])

}

func TestName2(t *testing.T) {
	m := map[string]string{
		"1": "1",
		"2": "2",
	}
	t.Log(m)
	clear(m)
	t.Log(m)
	a := []string{"1", "2"}
	t.Log(a)
	clear(a)
	t.Log(a)
}

package mgo

import (
	"context"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	ctx = mcontext.SetOpUserID(ctx, "imAdmin")
	ctx = mcontext.SetOperationID(ctx, "test123456")

	conn, err := grpc.DialContext(ctx, "172.16.8.48:10130", grpc.WithTransportCredentials(insecure.NewCredentials()), mw.GrpcClient())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cli := msg.NewMsgClient(conn)
	var ts int64

	ts = time.Now().UnixMilli()

	if _, err := cli.ClearMsg(ctx, &msg.ClearMsgReq{Timestamp: ts}); err != nil {
		panic(err)
	}
	t.Log("success!")
}

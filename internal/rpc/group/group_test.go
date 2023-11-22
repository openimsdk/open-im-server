package group

import (
	"context"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/wrapperspb"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo"
	tx2 "github.com/openimsdk/open-im-server/v3/pkg/common/db/tx"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/grouphash"
	"github.com/redis/go-redis/v9"
	"net"
	"testing"
	"time"
)

var (
	rdb redis.UniversalClient
	mgo *unrelation.Mongo
	ctx context.Context

	groupDB  *newmgo.GroupMgo
	memberDB *newmgo.GroupMemberMgo
	gs       *groupServer
)

func InitDB() error {
	addr := "172.16.8.142"
	pwd := "openIM123"

	config.Config.Redis.Address = []string{net.JoinHostPort(addr, "16379")}
	config.Config.Redis.Password = pwd
	config.Config.Mongo.Address = []string{net.JoinHostPort(addr, "37017")}
	config.Config.Mongo.Database = "openIM_v3"
	config.Config.Mongo.Username = "root"
	config.Config.Mongo.Password = pwd
	config.Config.Mongo.MaxPoolSize = 100
	var err error
	rdb, err = cache.NewRedis()
	if err != nil {
		return err
	}
	mgo, err = unrelation.NewMongo()
	if err != nil {
		return err
	}
	tx, err := tx2.NewAuto(context.Background(), mgo.GetClient())
	if err != nil {
		return err
	}

	config.Config.Object.Enable = "minio"
	config.Config.Object.ApiURL = "http://" + net.JoinHostPort(addr, "10002")
	config.Config.Object.Minio.Bucket = "openim"
	config.Config.Object.Minio.Endpoint = "http://" + net.JoinHostPort(addr, "10005")
	config.Config.Object.Minio.AccessKeyID = "root"
	config.Config.Object.Minio.SecretAccessKey = pwd
	config.Config.Object.Minio.SignEndpoint = config.Config.Object.Minio.Endpoint

	config.Config.Manager.UserID = []string{"openIM123456"}
	config.Config.Manager.Nickname = []string{"openIM123456"}

	ctx = context.WithValue(context.Background(), constant.OperationID, "debugOperationID")
	ctx = context.WithValue(context.Background(), constant.OpUserID, config.Config.Manager.UserID[0])

	if err := log.InitFromConfig("", "", 6, true, false, "", 2, 1); err != nil {
		panic(err)
	}

	g, err := newmgo.NewGroupMongo(mgo.GetDatabase())
	if err != nil {
		return err
	}
	gm, err := newmgo.NewGroupMember(mgo.GetDatabase())
	if err != nil {
		return err
	}
	gr, err := newmgo.NewGroupRequestMgo(mgo.GetDatabase())
	if err != nil {
		return err
	}
	groupDB = g.(*newmgo.GroupMgo)
	memberDB = gm.(*newmgo.GroupMemberMgo)

	gs = &groupServer{}
	gs.db = controller.NewGroupDatabase(rdb, groupDB, memberDB, gr, tx, grouphash.NewGroupHashFromGroupServer(gs))
	return nil
}

func init() {
	if err := InitDB(); err != nil {
		panic(err)
	}
}

func TestName(t *testing.T) {
	ctx = context.WithValue(context.Background(), constant.OpUserID, "4454892084")
	resp, err := gs.SetGroupMemberInfo(ctx, &group.SetGroupMemberInfoReq{
		Members: []*group.SetGroupMemberInfo{
			{
				GroupID:   "1731877545",
				UserID:    "3332974673",
				Nickname:  wrapperspb.String("nickname"),
				RoleLevel: wrapperspb.Int32(constant.GroupAdmin),
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func Test_GroupCreateCount(t *testing.T) {
	ctx = context.WithValue(context.Background(), constant.OpUserID, "4454892084")
	resp, err := gs.GroupCreateCount(ctx, &group.GroupCreateCountReq{
		Start: 0,
		End:   time.Now().Add(time.Hour).UnixMilli(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

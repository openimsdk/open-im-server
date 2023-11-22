package user

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"
	tablerelation "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	tx2 "github.com/openimsdk/open-im-server/v3/pkg/common/db/tx"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net"
	"strconv"
	"testing"
)

var (
	rdb redis.UniversalClient
	mgo *unrelation.Mongo
	ctx context.Context
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

	users := make([]*tablerelation.UserModel, 0)
	if len(config.Config.Manager.UserID) != len(config.Config.Manager.Nickname) {
		return errors.New("len(config.Config.Manager.AppManagerUid) != len(config.Config.Manager.Nickname)")
	}
	for k, v := range config.Config.Manager.UserID {
		users = append(users, &tablerelation.UserModel{UserID: v, Nickname: config.Config.Manager.Nickname[k], AppMangerLevel: constant.AppAdmin})
	}
	userDB, err := mgo.NewUserMongo(mgo.GetDatabase())
	if err != nil {
		return err
	}

	//var client registry.SvcDiscoveryRegistry
	//_= client
	cache := cache.NewUserCacheRedis(rdb, userDB, cache.GetDefaultOpt())
	userMongoDB := unrelation.NewUserMongoDriver(mgo.GetDatabase())
	database := controller.NewUserDatabase(userDB, cache, tx, userMongoDB)
	//friendRpcClient := rpcclient.NewFriendRpcClient(client)
	//groupRpcClient := rpcclient.NewGroupRpcClient(client)
	//msgRpcClient := rpcclient.NewMessageRpcClient(client)

	userSrv = &userServer{
		UserDatabase: database,
		//RegisterCenter:           client,
		//friendRpcClient:          &friendRpcClient,
		//groupRpcClient:           &groupRpcClient,
		//friendNotificationSender: notification.NewFriendNotificationSender(&msgRpcClient, notification.WithDBFunc(database.FindWithError)),
		//userNotificationSender:   notification.NewUserNotificationSender(&msgRpcClient, notification.WithUserFunc(database.FindWithError)),
	}

	return nil
}

func init() {
	if err := InitDB(); err != nil {
		panic(err)
	}
}

var userSrv *userServer

func TestName(t *testing.T) {
	userID := strconv.Itoa(int(rand.Uint32()))
	res, err := userSrv.UserRegister(ctx, &user.UserRegisterReq{
		Secret: config.Config.Secret,
		Users: []*sdkws.UserInfo{
			{
				UserID:   userID,
				Nickname: userID,
				FaceURL:  "",
				Ex:       "",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	t.Log(res)

}

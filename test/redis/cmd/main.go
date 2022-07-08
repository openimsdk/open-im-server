package main

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/proto/user"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/rockscache"
	"github.com/gin-gonic/gin"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/lithammer/shortuuid"
	"time"
)

func main() {
	rc := rockscache.NewClient(go_redis.NewClient(&go_redis.Options{
		Addr:     "43.128.5.63:16379",
		Password: "openIM", // no password set
		DB:       3,                              // use default DB
		PoolSize: 100,                            // 连接池大小
	}), rockscache.NewDefaultOptions())
	rc.Options.StrongConsistency = true

	//gid := dtmcli.MustGenGid(config.Config.Dtm.ServerURL)
	//
	//getGroupInfo := func()(string, error) {
	//	fmt.Println("start to fetch")
	//	time.Sleep(time.Second*5)
	//	fmt.Println("stop to fetch")
	//	return "value2312", nil
	//}
	//
	//v, err := rc.Fetch("keyasd347", time.Second*120, getGroupInfo)
	//fmt.Println("set success")


	//gid := dtmcli.MustGenGid(config.Config.Dtm.ServerURL)
	//msg := dtmcli.NewMsg(DtmServer, shortuuid.New()).
	//	Add(busi.Busi+"/SagaBTransIn", &TransReq{ Amount: 30 })
	//
	//dtmcli.DB()
	//time.Sleep(time.Second*12)
	//v2, _ := rc.Fetch("keyasd347", time.Second*10, func()(string, error) {
	//	// fetch data from database or other sources
	//	return "value1", nil
	//})



	//fmt.Println(v, err)
	//fmt.Println(v2)
	//err = rc.TagAsDeleted("keyasd347")
	//fmt.Println(err)
}

func UpdateUserInfo(c *gin.Context) {
	gid := dtmcli.MustGenGid(config.Config.Dtm.ServerURL)
	var pb user.GetAllUserIDReq
	msg := dtmgrpc.NewMsgGrpc(config.Config.Dtm.ServerURL, gid).Add("url1", &pb)
	msg.DoAndSubmit("/QueryPreparedB", func(bb *dtmcli.BranchBarrier) error {
		return bb.CallWithDB(db, func(tx *sql.Tx) error {
			return UpdateInTx(tx, &DBRow{
				K:        body["key"].(string),
				V:        body["value"].(string),
				TimeCost: body["time_cost"].(string),
			})
		})
	})
}))

}

func GetUserInfo(c *gin.Context) {
	// rpc logic
	v2, _ := db.DB.Rc.Fetch("keyasd347", time.Second*10, func()(string, error) {
		// fetch data from database or other sources
		return "value1", nil
	})
	c.JSON(200, map[string]interface{}{"s":v2})
}

// QueryPreparedB
func QueryPreparedB() {
	app.GET(BusiAPI+"/QueryPreparedB", dtmutil.WrapHandler2(func(c *gin.Context) interface{} {
		bb := MustBarrierFromGin(c)
		return bb.QueryPrepared(dbGet())
	}))
}
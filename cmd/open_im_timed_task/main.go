package main

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("start delete mongodb expired record")
		timeUnixBegin := time.Now().Unix()
		count, _ := db.DB.MgoUserCount()
		fmt.Println("mongodb record count: ", count)
		for i := 0; i < count; i++ {
			time.Sleep(1 * time.Millisecond)
			uid, _ := db.DB.MgoSkipUID(i)
			fmt.Println("operate uid: ", uid)
			err := db.DB.DelUserChat(uid)
			if err != nil {
				fmt.Println("operate uid failed: ", uid, err.Error())
			}
		}

		timeUnixEnd := time.Now().Unix()
		costTime := timeUnixEnd - timeUnixBegin
		if costTime > int64(config.Config.Mongo.DBRetainChatRecords*24*3600) {
			continue
		} else {
			sleepTime := 0
			if int64(config.Config.Mongo.DBRetainChatRecords*24*3600)-costTime > 24*3600 {
				sleepTime = 24 * 3600
			} else {
				sleepTime = config.Config.Mongo.DBRetainChatRecords*24*3600 - int(costTime)
			}
			fmt.Println("sleep: ", sleepTime)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}

}

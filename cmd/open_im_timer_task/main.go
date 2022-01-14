package main

import (
	commonDB "Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"time"
)

func main() {
	log.NewPrivateLog("timer")
	//for {
	//	fmt.Println("start delete mongodb expired record")
	//	timeUnixBegin := time.Now().Unix()
	//	count, _ := db.DB.MgoUserCount()
	//	fmt.Println("mongodb record count: ", count)
	//	for i := 0; i < count; i++ {
	//		time.Sleep(1 * time.Millisecond)
	//		uid, _ := db.DB.MgoSkipUID(i)
	//		fmt.Println("operate uid: ", uid)
	//		err := db.DB.DelUserChat(uid)
	//		if err != nil {
	//			fmt.Println("operate uid failed: ", uid, err.Error())
	//		}
	//	}
	//
	//	timeUnixEnd := time.Now().Unix()
	//	costTime := timeUnixEnd - timeUnixBegin
	//	if costTime > int64(config.Config.Mongo.DBRetainChatRecords*24*3600) {
	//		continue
	//	} else {
	//		sleepTime := 0
	//		if int64(config.Config.Mongo.DBRetainChatRecords*24*3600)-costTime > 24*3600 {
	//			sleepTime = 24 * 3600
	//		} else {
	//			sleepTime = config.Config.Mongo.DBRetainChatRecords*24*3600 - int(costTime)
	//		}
	//		fmt.Println("sleep: ", sleepTime)
	//		time.Sleep(time.Duration(sleepTime) * time.Second)
	//	}
	//}
	for {
		uidList, err := im_mysql_model.SelectAllUserID()
		if err != nil {
			//log.NewError("999999", err.Error())
		} else {
			for _, v := range uidList {
				minSeq, err := commonDB.DB.GetMinSeqFromMongo(v)
				if err != nil {
					//log.NewError("999999", "get user minSeq err", err.Error(), v)
					continue
				} else {
					err := commonDB.DB.SetUserMinSeq(v, minSeq)
					if err != nil {
						//log.NewError("999999", "set user minSeq err", err.Error(), v)
					}
				}
				time.Sleep(time.Duration(100) * time.Millisecond)
			}

		}

	}

}

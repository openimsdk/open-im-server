package im_mysql_msg_model

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"hash/crc32"
	"strconv"
)

func getHashMsgDBAddr(userID string) string {
	hCode := crc32.ChecksumIEEE([]byte(userID))
	return config.Config.Mysql.DBAddress[hCode%uint32(len(config.Config.Mysql.DBAddress))]
}

func getHashMsgTableIndex(userID string) int {
	hCode := crc32.ChecksumIEEE([]byte(userID))
	return int(hCode % uint32(config.Config.Mysql.DBMsgTableNum))
}

func QueryUserMsgID(userID string) ([]string, error) {
	dbAddress, dbTableIndex := getHashMsgDBAddr(userID), getHashMsgTableIndex(userID)
	dbTableName := "receive" + strconv.Itoa(dbTableIndex)

	dbConn, _ := db.DB.MysqlDB.GormDB(dbAddress, config.Config.Mysql.DBTableName)

	var msgID string
	var msgIDList []string
	rows, _ := dbConn.Raw("select msg_id from ? where user_id = ?", dbTableName, userID).Rows()
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&msgID)
		msgIDList = append(msgIDList, msgID)
	}

	return msgIDList, nil
}

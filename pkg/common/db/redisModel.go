package db

import (
	log2 "Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"github.com/garyburd/redigo/redis"
)

const (
	userIncrSeq      = "REDIS_USER_INCR_SEQ:" // user incr seq
	appleDeviceToken = "DEVICE_TOKEN"
	lastGetSeq       = "LAST_GET_SEQ"
	userMinSeq       = "REDIS_USER_MIN_SEQ:"
	uidPidToken      = "UID_PID_TOKEN:"
)

func (d *DataBases) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := d.redisPool.Get()
	if err := con.Err(); err != nil {
		log2.Error("", "", "redis cmd = %v, err = %v", cmd, err)
		return nil, err
	}
	defer con.Close()

	params := make([]interface{}, 0)
	params = append(params, key)

	if len(args) > 0 {
		for _, v := range args {
			params = append(params, v)
		}
	}

	return con.Do(cmd, params...)
}

//Perform seq auto-increment operation of user messages
func (d *DataBases) IncrUserSeq(uid string) (int64, error) {
	key := userIncrSeq + uid
	return redis.Int64(d.Exec("INCR", key))
}

//Get the largest Seq
func (d *DataBases) GetUserMaxSeq(uid string) (int64, error) {
	key := userIncrSeq + uid
	return redis.Int64(d.Exec("GET", key))
}

//Set the user's minimum seq
func (d *DataBases) SetUserMinSeq(uid string, minSeq int64) (err error) {
	key := userMinSeq + uid
	_, err = d.Exec("SET", key, minSeq)
	return err
}

//Get the smallest Seq
func (d *DataBases) GetUserMinSeq(uid string) (int64, error) {
	key := userMinSeq + uid
	return redis.Int64(d.Exec("GET", key))
}

//Store Apple's device token to redis
func (d *DataBases) SetAppleDeviceToken(accountAddress, value string) (err error) {
	key := appleDeviceToken + accountAddress
	_, err = d.Exec("SET", key, value)
	return err
}

//Delete Apple device token
func (d *DataBases) DelAppleDeviceToken(accountAddress string) (err error) {
	key := appleDeviceToken + accountAddress
	_, err = d.Exec("DEL", key)
	return err
}

//Store userid and platform class to redis
func (d *DataBases) AddTokenFlag(userID string, platformID int32, token string, flag int) error {
	key := uidPidToken + userID + ":" + utils.Int32ToString(platformID)
	m, err := redis.IntMap(d.Exec("GET", key))
	if err != nil && err != redis.ErrNil {
		return err
	}
	if err == redis.ErrNil {
		m = make(map[string]int, 5)
	}
	m[token] = flag
	_, err1 := d.Exec("SET", key, m)
	return err1
}

func (d *DataBases) GetTokenMapByUidPid(userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	return redis.IntMap(d.Exec("GET", key))
}
func (d *DataBases) SetTokenMapByUidPid(userID, platformID string, m map[string]int) error {
	key := uidPidToken + userID + ":" + platformID
	_, err := d.Exec("SET", key, m)
	return err
}

//Check exists userid and platform class from redis
func (d *DataBases) ExistsUserIDAndPlatform(userID, platformClass string) (int64, error) {
	key := userID + platformClass
	return redis.Int64(d.Exec("EXISTS", key))
}

//Get platform class Token
func (d *DataBases) GetPlatformToken(userID, platformClass string) (string, error) {
	key := userID + platformClass
	return redis.String(d.Exec("GET", key))
}

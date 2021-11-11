package db

import (
	log2 "Open_IM/pkg/common/log"
	"github.com/garyburd/redigo/redis"
)

const (
	userIncrSeq      = "REDIS_USER_INCR_SEQ:" // user incr seq
	appleDeviceToken = "DEVICE_TOKEN"
	lastGetSeq       = "LAST_GET_SEQ"
	userMinSeq       = "REDIS_USER_MIN_SEQ:"
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

//Record the last time the user actively pulled the value of Seq
func (d *DataBases) SetLastGetSeq(uid string) (err error) {
	key := lastGetSeq + uid
	_, err = d.Exec("SET", key)
	return err
}

//Get the value of the user's last active pull Seq
func (d *DataBases) GetLastGetSeq(uid string) (int64, error) {
	key := lastGetSeq + uid
	return redis.Int64(d.Exec("GET", key))
}

//Store userid and platform class to redis
func (d *DataBases) SetUserIDAndPlatform(userID, platformClass, value string, ttl int64) error {
	key := userID + platformClass
	_, err := d.Exec("SET", key, value, "EX", ttl)
	return err
}

//Check exists userid and platform class from redis
func (d *DataBases) ExistsUserIDAndPlatform(userID, platformClass string) (interface{}, error) {
	key := userID + platformClass
	exists, err := d.Exec("EXISTS", key)
	return exists, err
}

//Get platform class Token
func (d *DataBases) GetPlatformToken(userID, platformClass string) (interface{}, error) {
	key := userID + platformClass
	token, err := d.Exec("GET", key)
	return token, err
}

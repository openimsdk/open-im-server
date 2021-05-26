package db

import (
	"github.com/garyburd/redigo/redis"
)

const (
	userIncrSeq      = "REDIS_USER_INCR_SEQ:" // user incr seq
	appleDeviceToken = "DEVICE_TOKEN"
	lastGetSeq       = "LAST_GET_SEQ"
)

//执行用户消息的seq自增操作
func (d *DataBases) IncrUserSeq(uid string) (int64, error) {
	key := userIncrSeq + uid
	return redis.Int64(d.RedisDB.Exec("INCR", key))
}

//获取最新的seq
func (d *DataBases) GetUserSeq(uid string) (int64, error) {
	key := userIncrSeq + uid
	return redis.Int64(d.RedisDB.Exec("GET", key))
}

//存储苹果的设备token到redis
func (d *DataBases) SetAppleDeviceToken(accountAddress, value string) (err error) {
	key := appleDeviceToken + accountAddress
	_, err = d.RedisDB.Exec("SET", key, value)
	return err
}

//删除苹果设备token
func (d *DataBases) DelAppleDeviceToken(accountAddress string) (err error) {
	key := appleDeviceToken + accountAddress
	_, err = d.RedisDB.Exec("DEL", key)
	return err
}

//记录用户上一次主动拉取Seq的值
func (d *DataBases) SetLastGetSeq(uid string) (err error) {
	key := lastGetSeq + uid
	_, err = d.RedisDB.Exec("SET", key)
	return err
}

//获取用户上一次主动拉取Seq的值
func (d *DataBases) GetLastGetSeq(uid string) (int64, error) {
	key := userIncrSeq + uid
	return redis.Int64(d.RedisDB.Exec("GET", key))
}

//Store userid and platform class to redis
func (d *DataBases) SetUserIDAndPlatform(userID, platformClass, value string, ttl int64) error {
	key := userID + platformClass
	_, err := d.RedisDB.Exec("SET", key, value, "EX", ttl)
	return err
}

//Check exists userid and platform class from redis
func (d *DataBases) ExistsUserIDAndPlatform(userID, platformClass string) (interface{}, error) {
	key := userID + platformClass
	exists, err := d.RedisDB.Exec("EXISTS", key)
	return exists, err
}

//Get platform class Token
func (d *DataBases) GetPlatformToken(userID, platformClass string) (interface{}, error) {
	key := userID + platformClass
	token, err := d.RedisDB.Exec("GET", key)
	return token, err
}

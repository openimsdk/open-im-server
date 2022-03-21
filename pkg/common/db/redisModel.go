package db

import (
	"Open_IM/pkg/common/constant"
	log2 "Open_IM/pkg/common/log"
	"github.com/garyburd/redigo/redis"
)

const (
	AccountTempCode               = "ACCOUNT_TEMP_CODE"
	resetPwdTempCode              = "RESET_PWD_TEMP_CODE"
	userIncrSeq                   = "REDIS_USER_INCR_SEQ:" // user incr seq
	appleDeviceToken              = "DEVICE_TOKEN"
	userMinSeq                    = "REDIS_USER_MIN_SEQ:"
	uidPidToken                   = "UID_PID_TOKEN_STATUS:"
	conversationReceiveMessageOpt = "CON_RECV_MSG_OPT:"
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
func (d *DataBases) JudgeAccountEXISTS(account string) (bool, error) {
	key := AccountTempCode + account
	return redis.Bool(d.Exec("EXISTS", key))
}
func (d *DataBases) SetAccountCode(account string, code, ttl int) (err error) {
	key := AccountTempCode + account
	_, err = d.Exec("SET", key, code, "ex", ttl)
	return err
}
func (d *DataBases) GetAccountCode(account string) (string, error) {
	key := AccountTempCode + account
	return redis.String(d.Exec("GET", key))
}

//Perform seq auto-increment operation of user messages
func (d *DataBases) IncrUserSeq(uid string) (uint64, error) {
	key := userIncrSeq + uid
	return redis.Uint64(d.Exec("INCR", key))
}

//Get the largest Seq
func (d *DataBases) GetUserMaxSeq(uid string) (uint64, error) {
	key := userIncrSeq + uid
	return redis.Uint64(d.Exec("GET", key))
}

//Set the user's minimum seq
func (d *DataBases) SetUserMinSeq(uid string, minSeq uint32) (err error) {
	key := userMinSeq + uid
	_, err = d.Exec("SET", key, minSeq)
	return err
}

//Get the smallest Seq
func (d *DataBases) GetUserMinSeq(uid string) (uint64, error) {
	key := userMinSeq + uid
	return redis.Uint64(d.Exec("GET", key))
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
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	log2.NewDebug("", "add token key is ", key)
	_, err1 := d.Exec("HSet", key, token, flag)
	return err1
}

func (d *DataBases) GetTokenMapByUidPid(userID, platformID string) (map[string]int, error) {
	key := uidPidToken + userID + ":" + platformID
	log2.NewDebug("", "get token key is ", key)
	return redis.IntMap(d.Exec("HGETALL", key))
}
func (d *DataBases) SetTokenMapByUidPid(userID string, platformID int32, m map[string]int) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	_, err := d.Exec("hmset", key, redis.Args{}.Add().AddFlat(m)...)
	return err
}
func (d *DataBases) DeleteTokenByUidPid(userID string, platformID int32, fields []string) error {
	key := uidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
	_, err := d.Exec("HDEL", key, redis.Args{}.Add().AddFlat(fields)...)
	return err
}
func (d *DataBases) SetSingleConversationMsgOpt(userID, conversationID string, opt int) error {
	key := conversationReceiveMessageOpt + userID
	_, err := d.Exec("HSet", key, conversationID, opt)
	return err
}
func (d *DataBases) GetSingleConversationMsgOpt(userID, conversationID string) (int, error) {
	key := conversationReceiveMessageOpt + userID
	return redis.Int(d.Exec("HGet", key, conversationID))
}
func (d *DataBases) GetAllConversationMsgOpt(userID string) (map[string]int, error) {
	key := conversationReceiveMessageOpt + userID
	return redis.IntMap(d.Exec("HGETALL", key))
}
func (d *DataBases) SetMultiConversationMsgOpt(userID string, m map[string]int) error {
	key := conversationReceiveMessageOpt + userID
	_, err := d.Exec("hmset", key, redis.Args{}.Add().AddFlat(m)...)
	return err
}
func (d *DataBases) GetMultiConversationMsgOpt(userID string, conversationIDs []string) (m map[string]int, err error) {
	m = make(map[string]int)
	key := conversationReceiveMessageOpt + userID
	i, err := redis.Ints(d.Exec("hmget", key, redis.Args{}.Add().AddFlat(conversationIDs)...))
	if err != nil {
		return m, err
	}
	for k, v := range conversationIDs {
		m[v] = i[k]
	}
	return m, nil

}

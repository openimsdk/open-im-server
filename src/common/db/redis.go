package db

import (
	"Open_IM/src/common/config"
	log2 "Open_IM/src/common/log"
	"github.com/garyburd/redigo/redis"
	"time"
)

type redisDB struct {
	pool *redis.Pool
}

func (r *redisDB) newPool() {
	r.pool = &redis.Pool{
		MaxIdle:   config.Config.Redis.DBMaxIdle,
		MaxActive: config.Config.Redis.DBMaxActive,

		IdleTimeout: time.Duration(config.Config.Redis.DBIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				config.Config.Redis.DBAddress[0],
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialDatabase(0),
				redis.DialPassword(config.Config.Redis.DBPassWord),
			)
		},
	}
}

func (r *redisDB) Exec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := r.pool.Get()
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

package db

import (
	"Open_IM/src/common/config"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"sync"
	"time"
)

type mongoDB struct {
	sync.RWMutex
	dbMap map[string]*mgo.Session
}

func (m *mongoDB) mgoSession(dbName string) *mgo.Session {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.dbMap[dbName]; !ok {
		if err := m.newMgoSession(dbName); err != nil {
			panic(err)
			return nil
		}
	}
	return m.dbMap[dbName]
}

func (m *mongoDB) newMgoSession(dbName string) error {
	dailInfo := &mgo.DialInfo{
		Addrs:     config.Config.Mongo.DBAddress,
		Direct:    config.Config.Mongo.DBDirect,
		Timeout:   time.Second * time.Duration(config.Config.Mongo.DBTimeout),
		Database:  dbName,
		Source:    config.Config.Mongo.DBSource,
		Username:  config.Config.Mongo.DBUserName,
		Password:  config.Config.Mongo.DBPassword,
		PoolLimit: config.Config.Mongo.DBMaxPoolSize,
	}
	session, err := mgo.DialWithInfo(dailInfo)
	if err != nil {
		return errors.New(fmt.Sprintf("mongo DialWithInfo fail, err= %s", err.Error()))
	}

	if m.dbMap == nil {
		m.dbMap = make(map[string]*mgo.Session)
	}

	m.dbMap[dbName] = session
	return nil
}

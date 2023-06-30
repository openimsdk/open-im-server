package zookeeper

import (
	"github.com/go-zookeeper/zk"
)

func (s *ZkClient) RegisterConf2Registry(key string, conf []byte) error {
	exists, _, err := s.conn.Exists(s.getPath(key))
	if err != nil {
		return err
	}
	if exists {
		if err := s.conn.Delete(s.getPath(key), 0); err != nil {
			return err
		}
	}
	_, err = s.conn.Create(s.getPath(key), conf, 0, zk.WorldACL(zk.PermAll))
	if err != zk.ErrNodeExists {
		return err
	}
	return nil
}

func (s *ZkClient) GetConfFromRegistry(key string) ([]byte, error) {
	bytes, _, err := s.conn.Get(s.getPath(key))
	return bytes, err
}

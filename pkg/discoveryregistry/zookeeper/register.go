package zookeeper

import (
	"time"

	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func (s *ZkClient) Register(rpcRegisterName, host string, port int, opts ...grpc.DialOption) error {
	if err := s.ensureName(rpcRegisterName); err != nil {
		return err
	}
	addr := s.getAddr(host, port)
	_, err := grpc.Dial(addr, opts...)
	if err != nil {
		return err
	}
	node, err := s.conn.CreateProtectedEphemeralSequential(s.getPath(rpcRegisterName)+"/"+addr+"_", []byte(addr), zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	s.node = node
	return nil
}

func (s *ZkClient) UnRegister() error {
	s.lock.Lock()
	defer s.lock.Unlock()
	err := s.conn.Delete(s.node, -1)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	s.node = ""
	s.localConns = make(map[string][]resolver.Address)
	s.resolvers = make(map[string]*Resolver)
	return nil
}

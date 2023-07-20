// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zookeeper

import (
	"time"

	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func (s *ZkClient) CreateRpcRootNodes(serviceNames []string) error {
	for _, serviceName := range serviceNames {
		if err := s.ensureName(serviceName); err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s *ZkClient) CreateTempNode(rpcRegisterName, addr string) (node string, err error) {
	return s.conn.CreateProtectedEphemeralSequential(
		s.getPath(rpcRegisterName)+"/"+addr+"_",
		[]byte(addr),
		zk.WorldACL(zk.PermAll),
	)
}

func (s *ZkClient) Register(rpcRegisterName, host string, port int, opts ...grpc.DialOption) error {
	if err := s.ensureName(rpcRegisterName); err != nil {
		return err
	}
	addr := s.getAddr(host, port)
	_, err := grpc.Dial(addr, opts...)
	if err != nil {
		return err
	}
	node, err := s.CreateTempNode(rpcRegisterName, addr)
	if err != nil {
		return err
	}
	s.rpcRegisterName = rpcRegisterName
	s.rpcRegisterAddr = addr
	s.node = node
	s.isRegistered = true
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
	s.rpcRegisterName = ""
	s.rpcRegisterAddr = ""
	s.isRegistered = false
	s.localConns = make(map[string][]grpc.ClientConnInterface)
	s.resolvers = make(map[string]*Resolver)
	return nil
}

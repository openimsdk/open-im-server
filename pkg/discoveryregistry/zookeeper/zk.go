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
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const (
	defaultFreq = time.Minute * 30
	timeout     = 5
)

type Logger interface {
	Printf(string, ...interface{})
}

type ZkClient struct {
	zkServers []string
	zkRoot    string
	userName  string
	password  string

	scheme string

	timeout   int
	conn      *zk.Conn
	eventChan <-chan zk.Event
	node      string
	ticker    *time.Ticker

	lock    sync.Locker
	options []grpc.DialOption

	resolvers  map[string]*Resolver
	localConns map[string][]grpc.ClientConnInterface

	balancerName string

	logger Logger
}

type ZkOption func(*ZkClient)

func WithRoundRobin() ZkOption {
	return func(client *ZkClient) {
		client.balancerName = "round_robin"
	}
}

func WithUserNameAndPassword(userName, password string) ZkOption {
	return func(client *ZkClient) {
		client.userName = userName
		client.password = password
	}
}

func WithOptions(opts ...grpc.DialOption) ZkOption {
	return func(client *ZkClient) {
		client.options = opts
	}
}

func WithFreq(freq time.Duration) ZkOption {
	return func(client *ZkClient) {
		client.ticker = time.NewTicker(freq)
	}
}

func WithTimeout(timeout int) ZkOption {
	return func(client *ZkClient) {
		client.timeout = timeout
	}
}

func WithLogger(logger Logger) ZkOption {
	return func(client *ZkClient) {
		client.logger = logger
	}
}

func NewClient(zkServers []string, zkRoot string, options ...ZkOption) (*ZkClient, error) {
	client := &ZkClient{
		zkServers:  zkServers,
		zkRoot:     "/",
		scheme:     zkRoot,
		timeout:    timeout,
		localConns: make(map[string][]grpc.ClientConnInterface),
		resolvers:  make(map[string]*Resolver),
		lock:       &sync.Mutex{},
	}
	client.ticker = time.NewTicker(defaultFreq)
	for _, option := range options {
		option(client)
	}
	conn, eventChan, err := zk.Connect(
		zkServers,
		time.Duration(client.timeout)*time.Second,
		zk.WithLogInfo(true),
		zk.WithLogger(client.logger),
	)
	if err != nil {
		return nil, err
	}
	if client.userName != "" && client.password != "" {
		if err := conn.AddAuth("digest", []byte(client.userName+":"+client.password)); err != nil {
			return nil, err
		}
	}
	client.zkRoot += zkRoot
	client.eventChan = eventChan
	client.conn = conn
	if err := client.ensureRoot(); err != nil {
		client.CloseZK()
		return nil, err
	}
	resolver.Register(client)
	go client.refresh()
	go client.watch()
	return client, nil
}

func (s *ZkClient) CloseZK() {
	s.conn.Close()
}

func (s *ZkClient) ensureAndCreate(node string) error {
	exists, _, err := s.conn.Exists(node)
	if err != nil {
		return err
	}
	if !exists {
		_, err := s.conn.Create(node, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			return err
		}
	}
	return nil
}

func (s *ZkClient) refresh() {
	for range s.ticker.C {
		s.logger.Printf("refresh local conns")
		s.lock.Lock()
		for rpcName := range s.resolvers {
			s.flushResolver(rpcName)
		}
		for rpcName := range s.localConns {
			delete(s.localConns, rpcName)
		}
		s.lock.Unlock()
		s.logger.Printf("refresh local conns success")
	}
}

func (s *ZkClient) flushResolverAndDeleteLocal(serviceName string) {
	s.logger.Printf("start flush %s", serviceName)
	s.flushResolver(serviceName)
	delete(s.localConns, serviceName)
}

func (s *ZkClient) flushResolver(serviceName string) {
	r, ok := s.resolvers[serviceName]
	if ok {
		r.ResolveNowZK(resolver.ResolveNowOptions{})
	}
}

func (s *ZkClient) GetZkConn() *zk.Conn {
	return s.conn
}

func (s *ZkClient) GetRootPath() string {
	return s.zkRoot
}

func (s *ZkClient) GetNode() string {
	return s.node
}

func (s *ZkClient) ensureRoot() error {
	return s.ensureAndCreate(s.zkRoot)
}

func (s *ZkClient) ensureName(rpcRegisterName string) error {
	return s.ensureAndCreate(s.getPath(rpcRegisterName))
}

func (s *ZkClient) getPath(rpcRegisterName string) string {
	return s.zkRoot + "/" + rpcRegisterName
}

func (s *ZkClient) getAddr(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func (s *ZkClient) AddOption(opts ...grpc.DialOption) {
	s.options = append(s.options, opts...)
}

func (s *ZkClient) GetClientLocalConns() map[string][]grpc.ClientConnInterface {
	return s.localConns
}

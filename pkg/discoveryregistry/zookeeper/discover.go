package zookeeper

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var ErrConnIsNil = errors.New("conn is nil")
var ErrConnIsNilButLocalNotNil = errors.New("conn is nil, but local is not nil")

func (s *ZkClient) watch(wg *sync.WaitGroup) {
	for {
		event := <-s.eventChan
		switch event.Type {
		case zk.EventSession:
			s.logger.Printf("zk session event: %+v", event)
		case zk.EventNodeChildrenChanged:
			s.logger.Printf("zk event: %s", event.Path)
			l := strings.Split(event.Path, "/")
			if len(l) > 1 {
				s.lock.Lock()
				rpcName := l[len(l)-1]
				s.flushResolver(rpcName)
				if len(s.localConns[rpcName]) != 0 {
					delete(s.localConns, rpcName)
				}
				s.lock.Unlock()
			}
			s.logger.Printf("zk event handle success: %s", event.Path)
		case zk.EventNodeDataChanged:
		case zk.EventNodeCreated:
		case zk.EventNodeDeleted:
		case zk.EventNotWatching:
		}
	}

}

func (s *ZkClient) GetConnsRemote(serviceName string) (conns []resolver.Address, err error) {
	path := s.getPath(serviceName)
	childNodes, _, err := s.conn.Children(path)
	if err != nil {
		return nil, errors.Wrap(err, "get children error")
	}
	for _, child := range childNodes {
		fullPath := path + "/" + child
		data, _, err := s.conn.Get(fullPath)
		if err != nil {
			if err == zk.ErrNoNode {
				return nil, errors.Wrap(err, "this is zk ErrNoNode")
			}
			return nil, errors.Wrap(err, "get children error")
		}
		conns = append(conns, resolver.Address{Addr: string(data), ServerName: serviceName})
	}
	_, _, _, err = s.conn.ChildrenW(s.getPath(serviceName))
	if err != nil {
		return nil, errors.Wrap(err, "children watch error")
	}
	if len(conns) == 0 {
		return nil, fmt.Errorf("no conn for service %s, grpc server may not exist, local conn is %v, please check zookeeper server %v, path: %s", serviceName, s.localConns, s.zkServers, s.zkRoot)
	}
	return conns, nil
}

func (s *ZkClient) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]grpc.ClientConnInterface, error) {
	s.logger.Printf("get conns from client, serviceName: %s", serviceName)
	s.lock.Lock()
	defer s.lock.Unlock()
	opts = append(s.options, opts...)
	conns := s.localConns[serviceName]
	if len(conns) == 0 {
		var err error
		s.logger.Printf("get conns from zk remote, serviceName: %s", serviceName)
		conns, err = s.GetConnsRemote(serviceName)
		if err != nil {
			return nil, err
		}
		s.localConns[serviceName] = conns
	}
	var ret []grpc.ClientConnInterface
	s.logger.Printf("get conns from zk success, serviceName: %s", serviceName)
	for _, conn := range conns {
		cc, err := grpc.DialContext(ctx, conn.Addr, append(s.options, opts...)...)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("conns dialContext error, conn: %s", conn.Addr))
		}
		ret = append(ret, newClientConnInterface(cc))
	}
	s.logger.Printf("dial ctx success, serviceName: %s", serviceName)
	return ret, nil
}

func (s *ZkClient) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	newOpts := append(s.options, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, s.balancerName)))
	cc, err := grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", s.scheme, serviceName), append(newOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	return newClientConnInterface(cc), nil
}

func (s *ZkClient) CloseConn(conn grpc.ClientConnInterface) {
	if closer, ok := conn.(io.Closer); ok {
		closer.Close()
	}
}

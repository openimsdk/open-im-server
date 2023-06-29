package zookeeper

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/pkg/errors"

	"github.com/go-zookeeper/zk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var ErrConnIsNil = errors.New("conn is nil")
var ErrConnIsNilButLocalNotNil = errors.New("conn is nil, but local is not nil")

func (s *ZkClient) watch() {
	for {
		event := <-s.eventChan
		switch event.Type {
		case zk.EventSession:
			s.logger.Printf("zk session event: %+v", event)
		case zk.EventNodeChildrenChanged:
			s.logger.Printf("zk event: %s", event.Path)
			l := strings.Split(event.Path, "/")
			if len(l) > 1 {
				serviceName := l[len(l)-1]
				s.lock.Lock()
				s.flushResolverAndDeleteLocal(serviceName)
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
	_, _, _, err = s.conn.ChildrenW(path)
	if err != nil {
		return nil, errors.Wrap(err, "children watch error")
	}
	childNodes, _, err := s.conn.Children(path)
	if err != nil {
		return nil, errors.Wrap(err, "get children error")
	} else {
		for _, child := range childNodes {
			fullPath := path + "/" + child
			data, _, err := s.conn.Get(fullPath)
			if err != nil {
				if err == zk.ErrNoNode {
					return nil, errors.Wrap(err, "this is zk ErrNoNode")
				}
				return nil, errors.Wrap(err, "get children error")
			}
			log.ZDebug(context.Background(), "get conns from remote", "conn", string(data))
			conns = append(conns, resolver.Address{Addr: string(data), ServerName: serviceName})
		}
	}
	return conns, nil
}

func (s *ZkClient) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]grpc.ClientConnInterface, error) {
	s.logger.Printf("get conns from client, serviceName: %s", serviceName)
	s.lock.Lock()
	opts = append(s.options, opts...)
	conns := s.localConns[serviceName]
	if len(conns) == 0 {
		var err error
		s.logger.Printf("get conns from zk remote, serviceName: %s", serviceName)
		conns, err = s.GetConnsRemote(serviceName)
		if err != nil {
			s.lock.Unlock()
			return nil, err
		}
		if len(conns) == 0 {
			return nil, fmt.Errorf("no conn for service %s, grpc server may not exist, local conn is %v, please check zookeeper server %v, path: %s", serviceName, s.localConns, s.zkServers, s.zkRoot)
		}
		s.localConns[serviceName] = conns
	}
	s.lock.Unlock()
	var ret []grpc.ClientConnInterface
	s.logger.Printf("get conns from zk success, serviceName: %s", serviceName)
	for _, conn := range conns {
		cc, err := grpc.DialContext(ctx, conn.Addr, append(s.options, opts...)...)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("conns dialContext error, conn: %s", conn.Addr))
		}
		ret = append(ret, cc)
	}
	s.logger.Printf("dial ctx success, serviceName: %s", serviceName)
	return ret, nil
}

func (s *ZkClient) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	newOpts := append(s.options, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, s.balancerName)))
	return grpc.DialContext(ctx, fmt.Sprintf("%s:///%s", s.scheme, serviceName), append(newOpts, opts...)...)
}

func (s *ZkClient) CloseConn(conn grpc.ClientConnInterface) {
	if closer, ok := conn.(io.Closer); ok {
		closer.Close()
	}
}

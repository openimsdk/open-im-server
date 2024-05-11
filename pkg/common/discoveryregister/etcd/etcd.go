package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	gresolver "google.golang.org/grpc/resolver"

	"log"
	"time"
)

// SvcDiscoveryRegistryImpl implementation
type SvcDiscoveryRegistryImpl struct {
	client      *clientv3.Client
	resolver    gresolver.Builder
	dialOptions []grpc.DialOption
	serviceKey  string
	endpointMgr endpoints.Manager
	leaseID     clientv3.LeaseID
	schema      string
}

func NewSvcDiscoveryRegistry(schema string, endpoints []string) (*SvcDiscoveryRegistryImpl, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}
	r, err := resolver.NewBuilder(client)
	if err != nil {
		return nil, err
	}
	return &SvcDiscoveryRegistryImpl{
		client:   client,
		resolver: r,
		schema:   schema,
	}, nil
}

func (r *SvcDiscoveryRegistryImpl) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	return "", nil
}
func (r *SvcDiscoveryRegistryImpl) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	target := fmt.Sprintf("%s:///%s", r.schema, serviceName)
	conn, err := grpc.DialContext(ctx, target, append(r.dialOptions, opts...)...)
	if err != nil {
		return nil, err
	}
	return []*grpc.ClientConn{conn}, nil
}

func (r *SvcDiscoveryRegistryImpl) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("%s:///%s", r.schema, serviceName)
	return grpc.DialContext(ctx, target, append(r.dialOptions, opts...)...)
}

func (r *SvcDiscoveryRegistryImpl) GetSelfConnTarget() string {
	return fmt.Sprintf("%s:///%s", r.schema, r.serviceKey)
}

func (r *SvcDiscoveryRegistryImpl) AddOption(opts ...grpc.DialOption) {
	r.dialOptions = append(r.dialOptions, opts...)
}

func (r *SvcDiscoveryRegistryImpl) CloseConn(conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}

func (r *SvcDiscoveryRegistryImpl) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	r.serviceKey = fmt.Sprintf("%s/%s:%d", serviceName, host, port)
	em, err := endpoints.NewManager(r.client, serviceName)
	if err != nil {
		return err
	}
	r.endpointMgr = em

	leaseResp, err := r.client.Grant(context.Background(), 30)
	if err != nil {
		return err
	}
	r.leaseID = leaseResp.ID

	endpoint := endpoints.Endpoint{Addr: fmt.Sprintf("%s:%d", host, port)}
	err = em.AddEndpoint(context.TODO(), r.serviceKey, endpoint, clientv3.WithLease(leaseResp.ID))
	return err
}

func (r *SvcDiscoveryRegistryImpl) UnRegister() error {
	if r.endpointMgr == nil {
		return fmt.Errorf("endpoint manager is not initialized")
	}
	return r.endpointMgr.DeleteEndpoint(context.TODO(), r.serviceKey)
}

func (r *SvcDiscoveryRegistryImpl) Close() {
	if r.client != nil {
		_ = r.client.Close()
	}
}

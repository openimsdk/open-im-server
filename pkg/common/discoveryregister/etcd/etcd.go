package etcd

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	gresolver "google.golang.org/grpc/resolver"
	"time"
)

// ZkOption defines a function type for modifying clientv3.Config
type ZkOption func(*clientv3.Config)

// SvcDiscoveryRegistryImpl implementation
type SvcDiscoveryRegistryImpl struct {
	client            *clientv3.Client
	resolver          gresolver.Builder
	dialOptions       []grpc.DialOption
	serviceKey        string
	endpointMgr       endpoints.Manager
	leaseID           clientv3.LeaseID
	rpcRegisterTarget string

	rootDirectory string
}

// NewSvcDiscoveryRegistry creates a new service discovery registry implementation
func NewSvcDiscoveryRegistry(rootDirectory string, endpoints []string, options ...ZkOption) (*SvcDiscoveryRegistryImpl, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		// Increase keep-alive queue capacity and message size
		PermitWithoutStream: true,
		MaxCallSendMsgSize:  10 * 1024 * 1024, // 10 MB
	}

	// Apply provided options to the config
	for _, opt := range options {
		opt(&cfg)
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
		client:        client,
		resolver:      r,
		rootDirectory: rootDirectory,
	}, nil
}

// WithDialTimeout sets a custom dial timeout for the etcd client
func WithDialTimeout(timeout time.Duration) ZkOption {
	return func(cfg *clientv3.Config) {
		cfg.DialTimeout = timeout
	}
}

// WithMaxCallSendMsgSize sets a custom max call send message size for the etcd client
func WithMaxCallSendMsgSize(size int) ZkOption {
	return func(cfg *clientv3.Config) {
		cfg.MaxCallSendMsgSize = size
	}
}

// WithUsernameAndPassword sets a username and password for the etcd client
func WithUsernameAndPassword(username, password string) ZkOption {
	return func(cfg *clientv3.Config) {
		cfg.Username = username
		cfg.Password = password
	}
}

// GetUserIdHashGatewayHost returns the gateway host for a given user ID hash
func (r *SvcDiscoveryRegistryImpl) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	return "", nil
}

// GetConns returns gRPC client connections for a given service name
func (r *SvcDiscoveryRegistryImpl) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	var conns []*grpc.ClientConn
	// Construct the full key for the service
	fullServiceKey := fmt.Sprintf("%s/%s", r.rootDirectory, serviceName)

	// List all endpoints for the service
	resp, err := r.client.Get(ctx, fullServiceKey, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("GetConns get ", fullServiceKey, err.Error())
		return nil, err
	}

	for _, kv := range resp.Kvs {
		endpoint := string(kv.Key[len(fullServiceKey)+1:]) // Extract the endpoint address
		//target := fmt.Sprintf("etcd://%s/%s/%s", r.rootDirectory, serviceName, endpoint)
		target := endpoint
		conn, err := grpc.DialContext(ctx, target, append(append(r.dialOptions, opts...), grpc.WithResolvers(r.resolver))...)
		if err != nil {
			fmt.Println("DialContext ", target, err.Error())
			return nil, err
		}
		conns = append(conns, conn)
		fmt.Println("GetConns detail ", *conn)

	}
	return conns, nil
}

// GetConn returns a single gRPC client connection for a given service name
func (r *SvcDiscoveryRegistryImpl) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("etcd:///%s/%s", r.rootDirectory, serviceName)
	return grpc.DialContext(ctx, target, append(append(r.dialOptions, opts...), grpc.WithResolvers(r.resolver))...)
}

// GetSelfConnTarget returns the connection target for the current service
func (r *SvcDiscoveryRegistryImpl) GetSelfConnTarget() string {
	return r.rpcRegisterTarget
	//	return fmt.Sprintf("etcd:///%s", r.serviceKey)

}

// AddOption appends gRPC dial options to the existing options
func (r *SvcDiscoveryRegistryImpl) AddOption(opts ...grpc.DialOption) {
	r.dialOptions = append(r.dialOptions, opts...)
}

// CloseConn closes a given gRPC client connection
func (r *SvcDiscoveryRegistryImpl) CloseConn(conn *grpc.ClientConn) {
	if err := conn.Close(); err != nil {
		fmt.Printf("Failed to close connection: %v\n", err)
	}
}

// Register registers a new service endpoint with etcd
func (r *SvcDiscoveryRegistryImpl) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	r.serviceKey = fmt.Sprintf("%s/%s/%s:%d", r.rootDirectory, serviceName, host, port)
	em, err := endpoints.NewManager(r.client, r.rootDirectory+"/"+serviceName)
	if err != nil {
		return err
	}
	r.endpointMgr = em

	leaseResp, err := r.client.Grant(context.Background(), 60) // Increase TTL time
	if err != nil {
		return err
	}
	r.leaseID = leaseResp.ID

	r.rpcRegisterTarget = fmt.Sprintf("%s:%d", host, port)
	endpoint := endpoints.Endpoint{Addr: r.rpcRegisterTarget}

	err = em.AddEndpoint(context.TODO(), r.serviceKey, endpoint, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	go r.keepAliveLease(r.leaseID)

	return nil
}

// keepAliveLease maintains the lease alive by sending keep-alive requests
func (r *SvcDiscoveryRegistryImpl) keepAliveLease(leaseID clientv3.LeaseID) {
	ch, err := r.client.KeepAlive(context.Background(), leaseID)
	if err != nil {
		fmt.Printf("Failed to keep lease alive: %v\n", err)
		return
	}

	for ka := range ch {
		if ka != nil {
		} else {
			fmt.Printf("Lease keep-alive response channel closed\n")
			return
		}
	}
}

// UnRegister removes the service endpoint from etcd
func (r *SvcDiscoveryRegistryImpl) UnRegister() error {
	if r.endpointMgr == nil {
		return fmt.Errorf("endpoint manager is not initialized")
	}
	return r.endpointMgr.DeleteEndpoint(context.TODO(), r.serviceKey)
}

// Close closes the etcd client connection
func (r *SvcDiscoveryRegistryImpl) Close() {
	if r.client != nil {
		_ = r.client.Close()
	}
}

// Check verifies if etcd is running by checking the existence of the root node and optionally creates it with a lease
func Check(ctx context.Context, etcdServers []string, etcdRoot string, createIfNotExist bool, options ...ZkOption) error {
	// Configure the etcd client with default settings
	cfg := clientv3.Config{
		Endpoints: etcdServers,
	}

	// Apply provided options to the config
	for _, opt := range options {
		opt(&cfg)
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to connect to etcd")
	}
	defer client.Close()

	// Determine timeout for context
	var opCtx context.Context
	var cancel context.CancelFunc
	if cfg.DialTimeout != 0 {
		opCtx, cancel = context.WithTimeout(ctx, cfg.DialTimeout)
	} else {
		opCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	}
	defer cancel()

	// Check if the root node exists
	resp, err := client.Get(opCtx, etcdRoot)
	if err != nil {
		return errors.Wrap(err, "failed to get the root node from etcd")
	}

	// If root node does not exist and createIfNotExist is true, create the root node with a lease
	if len(resp.Kvs) == 0 {
		if createIfNotExist {
			var leaseTTL int64 = 10
			var leaseResp *clientv3.LeaseGrantResponse
			if leaseTTL > 0 {
				// Create a lease
				leaseResp, err = client.Grant(opCtx, leaseTTL)
				if err != nil {
					return errors.Wrap(err, "failed to create lease in etcd")
				}
			}

			// Put the key with the lease
			putOpts := []clientv3.OpOption{}
			if leaseResp != nil {
				putOpts = append(putOpts, clientv3.WithLease(leaseResp.ID))
			}

			_, err := client.Put(opCtx, etcdRoot, "", putOpts...)
			if err != nil {
				return errors.Wrap(err, "failed to create the root node in etcd")
			}
			fmt.Printf("Root node %s did not exist, but has been created.\n", etcdRoot)
		} else {
			return fmt.Errorf("root node %s does not exist in etcd", etcdRoot)
		}
	} else {
		fmt.Printf("Etcd is running and the root node %s exists.\n", etcdRoot)
	}

	return nil
}

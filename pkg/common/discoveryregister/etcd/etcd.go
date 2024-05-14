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
	"strings"
	"sync"
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

	mu      sync.RWMutex
	connMap map[string][]*grpc.ClientConn
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

	s := &SvcDiscoveryRegistryImpl{
		client:        client,
		resolver:      r,
		rootDirectory: rootDirectory,
		connMap:       make(map[string][]*grpc.ClientConn),
	}
	go s.watchServiceChanges()
	return s, nil
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
	fullServiceKey := fmt.Sprintf("%s/%s", r.rootDirectory, serviceName)
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.connMap[fullServiceKey], nil
}

// GetConn returns a single gRPC client connection for a given service name
func (r *SvcDiscoveryRegistryImpl) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("etcd:///%s/%s", r.rootDirectory, serviceName)
	return grpc.DialContext(ctx, target, append(append(r.dialOptions, opts...), grpc.WithResolvers(r.resolver))...)
}

// GetSelfConnTarget returns the connection target for the current service
func (r *SvcDiscoveryRegistryImpl) GetSelfConnTarget() string {
	return r.rpcRegisterTarget
}

// AddOption appends gRPC dial options to the existing options
func (r *SvcDiscoveryRegistryImpl) AddOption(opts ...grpc.DialOption) {
	r.dialOptions = append(r.dialOptions, opts...)
}

// CloseConn closes a given gRPC client connection
func (r *SvcDiscoveryRegistryImpl) CloseConn(conn *grpc.ClientConn) {
	conn.Close()
}

// Register registers a new service endpoint with etcd
func (r *SvcDiscoveryRegistryImpl) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	r.serviceKey = fmt.Sprintf("%s/%s/%s:%d", r.rootDirectory, serviceName, host, port)
	em, err := endpoints.NewManager(r.client, r.rootDirectory+"/"+serviceName)
	if err != nil {
		return err
	}
	r.endpointMgr = em

	leaseResp, err := r.client.Grant(context.Background(), 30) //
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
		return
	}
	for ka := range ch {
		if ka != nil {
		} else {
			return
		}
	}
}

// watchServiceChanges watches for changes in the service directory
func (r *SvcDiscoveryRegistryImpl) watchServiceChanges() {
	watchChan := r.client.Watch(context.Background(), r.rootDirectory, clientv3.WithPrefix())
	for watchResp := range watchChan {
		updatedPrefixes := make(map[string]struct{}) // Create a set to track updated prefixes

		for _, event := range watchResp.Events {
			prefix, _ := r.splitEndpoint(string(event.Kv.Key))
			if _, alreadyUpdated := updatedPrefixes[prefix]; !alreadyUpdated {
				updatedPrefixes[prefix] = struct{}{} // Mark this prefix as updated
				r.refreshConnMap(prefix)
			}
		}
	}
}

// refreshConnMap fetches the latest endpoints and updates the local map
func (r *SvcDiscoveryRegistryImpl) refreshConnMap(prefix string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fullPrefix := fmt.Sprintf("%s/", prefix)
	resp, err := r.client.Get(context.Background(), fullPrefix, clientv3.WithPrefix())
	if err != nil {
		return
	}
	r.connMap[prefix] = []*grpc.ClientConn{} // Update the connMap with new connections
	for _, kv := range resp.Kvs {
		_, addr := r.splitEndpoint(string(kv.Key))
		conn, err := grpc.DialContext(context.Background(), addr, append(r.dialOptions, grpc.WithResolvers(r.resolver))...)
		if err != nil {
			continue
		}
		r.connMap[prefix] = append(r.connMap[prefix], conn)
	}
}

// splitEndpoint splits the endpoint string into prefix and address
func (r *SvcDiscoveryRegistryImpl) splitEndpoint(input string) (string, string) {
	lastSlashIndex := strings.LastIndex(input, "/")
	if lastSlashIndex != -1 {
		part1 := input[:lastSlashIndex]
		part2 := input[lastSlashIndex+1:]
		return part1, part2
	}
	return input, ""
}

// UnRegister removes the service endpoint from etcd
func (r *SvcDiscoveryRegistryImpl) UnRegister() error {
	if r.endpointMgr == nil {
		return fmt.Errorf("endpoint manager is not initialized")
	}
	err := r.endpointMgr.DeleteEndpoint(context.TODO(), r.serviceKey)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the etcd client connection
func (r *SvcDiscoveryRegistryImpl) Close() {
	if r.client != nil {
		_ = r.client.Close()
	}

	r.mu.Lock()
	defer r.mu.Unlock()
}

// Check verifies if etcd is running by checking the existence of the root node and optionally creates it with a lease
func Check(ctx context.Context, etcdServers []string, etcdRoot string, createIfNotExist bool, options ...ZkOption) error {
	cfg := clientv3.Config{
		Endpoints: etcdServers,
	}
	for _, opt := range options {
		opt(&cfg)
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to connect to etcd")
	}
	defer client.Close()

	var opCtx context.Context
	var cancel context.CancelFunc
	if cfg.DialTimeout != 0 {
		opCtx, cancel = context.WithTimeout(ctx, cfg.DialTimeout)
	} else {
		opCtx, cancel = context.WithTimeout(ctx, 10*time.Second)
	}
	defer cancel()

	resp, err := client.Get(opCtx, etcdRoot)
	if err != nil {
		return errors.Wrap(err, "failed to get the root node from etcd")
	}

	if len(resp.Kvs) == 0 {
		if createIfNotExist {
			var leaseTTL int64 = 10
			var leaseResp *clientv3.LeaseGrantResponse
			if leaseTTL > 0 {
				leaseResp, err = client.Grant(opCtx, leaseTTL)
				if err != nil {
					return errors.Wrap(err, "failed to create lease in etcd")
				}
			}
			putOpts := []clientv3.OpOption{}
			if leaseResp != nil {
				putOpts = append(putOpts, clientv3.WithLease(leaseResp.ID))
			}

			_, err := client.Put(opCtx, etcdRoot, "", putOpts...)
			if err != nil {
				return errors.Wrap(err, "failed to create the root node in etcd")
			}
		} else {
			return fmt.Errorf("root node %s does not exist in etcd", etcdRoot)
		}
	}
	return nil
}

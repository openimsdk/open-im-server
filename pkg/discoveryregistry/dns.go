package discoveryregistry

// type DnsDiscoveryRegistry struct {
// 	opts      []grpc.DialOption
// 	namespace string
// 	clientset *kubernetes.Clientset
// }

// func NewDnsDiscoveryRegistry(namespace string, opts []grpc.DialOption) (*DnsDiscoveryRegistry, error) {
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		return nil, err
// 	}
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &DnsDiscoveryRegistry{
// 		clientset: clientset,
// 		namespace: namespace,
// 		opts:      opts,
// 	}, nil
// }

// func (d DnsDiscoveryRegistry) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
// 	endpoints, err := d.clientset.CoreV1().Endpoints(d.namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	var conns []*grpc.ClientConn
// 	for _, subset := range endpoints.Subsets {
// 		for _, address := range subset.Addresses {
// 			for _, port := range subset.Ports {
// 				conn, err := grpc.DialContext(ctx, net.JoinHostPort(address.IP, string(port.Port)), append(d.opts, opts...)...)
// 				if err != nil {
// 					return nil, err
// 				}
// 				conns = append(conns, conn)
// 			}
// 		}
// 	}
// 	return conns, nil
// }

// func (d DnsDiscoveryRegistry) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
// 	return grpc.DialContext(ctx, fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, d.namespace), append(d.opts, opts...)...)
// }

// func (d *DnsDiscoveryRegistry) AddOption(opts ...grpc.DialOption) {
// 	d.opts = append(d.opts, opts...)
// }

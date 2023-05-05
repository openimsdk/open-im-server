package discoveryregistry

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type DnsDiscoveryRegistry struct {
	opts      []grpc.DialOption
	namespace string
	config    *rest.Config
}

func NewDnsDiscoveryRegistry(namespace string, opts []grpc.DialOption) (*DnsDiscoveryRegistry, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return &DnsDiscoveryRegistry{
		config:    config,
		namespace: namespace,
		opts:      opts,
	}, nil
}

func (d DnsDiscoveryRegistry) GetConns(serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	clientset, err := kubernetes.NewForConfig(d.config)
	if err != nil {
		return nil, err
	}
	endpoints, err := clientset.CoreV1().Endpoints(d.namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	var conns []*grpc.ClientConn
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			for _, port := range subset.Ports {
				conn, err := grpc.Dial(net.JoinHostPort(address.IP, string(port.Port)), append(d.opts, opts...)...)
				if err != nil {
					return nil, err
				}
				conns = append(conns, conn)
			}
		}
	}
	return conns, nil
}

func (d DnsDiscoveryRegistry) GetConn(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, d.namespace), append(d.opts, opts...)...)
}

func (d *DnsDiscoveryRegistry) AddOption(opts ...grpc.DialOption) {
	d.opts = append(d.opts, opts...)
}

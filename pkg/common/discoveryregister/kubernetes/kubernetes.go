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

package kubernetes

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type KubernetesConnManager struct {
	clientset   *kubernetes.Clientset
	namespace   string
	dialOptions []grpc.DialOption

	selfTarget string

	mu      sync.RWMutex
	connMap map[string][]*grpc.ClientConn
}

// NewKubernetesConnManager creates a new connection manager that uses Kubernetes services for service discovery.
func NewKubernetesConnManager(namespace string, options ...grpc.DialOption) (*KubernetesConnManager, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	k := &KubernetesConnManager{
		clientset:   clientset,
		namespace:   namespace,
		dialOptions: options,
		connMap:     make(map[string][]*grpc.ClientConn),
	}

	go k.watchEndpoints()

	return k, nil
}

func (k *KubernetesConnManager) initializeConns(serviceName string) error {
	port, err := k.getServicePort(serviceName)
	if err != nil {
		return err
	}

	endpoints, err := k.clientset.CoreV1().Endpoints(k.namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get endpoints for service %s: %v", serviceName, err)
	}

	var conns []*grpc.ClientConn
	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			target := fmt.Sprintf("%s:%d", address.IP, port)
			conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return fmt.Errorf("failed to dial endpoint %s: %v", target, err)
			}
			conns = append(conns, conn)
		}
	}

	k.mu.Lock()
	defer k.mu.Unlock()
	k.connMap[serviceName] = conns

	// go k.watchEndpoints(serviceName)

	return nil
}

// GetConns returns gRPC client connections for a given Kubernetes service name.
func (k *KubernetesConnManager) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	k.mu.RLock()
	conns, exists := k.connMap[serviceName]
	defer k.mu.RUnlock()

	if exists {
		return conns, nil
	}

	k.mu.Lock()
	defer k.mu.Unlock()

	// Check if another goroutine has already initialized the connections when we released the read lock
	conns, exists = k.connMap[serviceName]
	if exists {
		return conns, nil
	}

	if err := k.initializeConns(serviceName); err != nil {
		return nil, fmt.Errorf("failed to initialize connections for service %s: %v", serviceName, err)
	}

	return k.connMap[serviceName], nil
}

// GetConn returns a single gRPC client connection for a given Kubernetes service name.
func (k *KubernetesConnManager) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	port, err := k.getServicePort(serviceName)
	if err != nil {
		return nil, err
	}

	fmt.Println("SVC port:", port)

	target := fmt.Sprintf("%s.%s.svc.cluster.local:%d", serviceName, k.namespace, port)

	fmt.Println("SVC target:", target)

	return grpc.DialContext(
		ctx,
		target,
		append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, k.dialOptions...)...,
	)
}

// GetSelfConnTarget returns the connection target for the current service.
func (k *KubernetesConnManager) GetSelfConnTarget() string {
	return k.selfTarget
}

// AddOption appends gRPC dial options to the existing options.
func (k *KubernetesConnManager) AddOption(opts ...grpc.DialOption) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.dialOptions = append(k.dialOptions, opts...)
}

// CloseConn closes a given gRPC client connection.
func (k *KubernetesConnManager) CloseConn(conn *grpc.ClientConn) {
	conn.Close()
}

// Close closes all gRPC connections managed by KubernetesConnManager.
func (k *KubernetesConnManager) Close() {
	k.mu.Lock()
	defer k.mu.Unlock()
	for _, conns := range k.connMap {
		for _, conn := range conns {
			_ = conn.Close()
		}
	}
	k.connMap = make(map[string][]*grpc.ClientConn)
}

func (k *KubernetesConnManager) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	return nil
}
func (k *KubernetesConnManager) UnRegister() error {
	return nil
}

func (k *KubernetesConnManager) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	return "", nil
}

func (k *KubernetesConnManager) getServicePort(serviceName string) (int32, error) {
	svc, err := k.clientset.CoreV1().Services(k.namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
	if err != nil {
		fmt.Print("namespace:", k.namespace)
		return 0, fmt.Errorf("failed to get service %s: %v", serviceName, err)
	}

	if len(svc.Spec.Ports) == 0 {
		return 0, fmt.Errorf("service %s has no ports defined", serviceName)
	}

	return svc.Spec.Ports[0].Port, nil
}

// watchEndpoints listens for changes in Pod resources.
func (k *KubernetesConnManager) watchEndpoints() {
	informerFactory := informers.NewSharedInformerFactory(k.clientset, time.Minute*10)
	informer := informerFactory.Core().V1().Pods().Informer()

	// Watch for Pod changes (add, update, delete)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			k.handleEndpointChange(obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			k.handleEndpointChange(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			k.handleEndpointChange(obj)
		},
	})

	informerFactory.Start(context.Background().Done())
	<-context.Background().Done() // Block forever
}

func (k *KubernetesConnManager) handleEndpointChange(obj interface{}) {
	endpoint, ok := obj.(*v1.Endpoints)
	if !ok {
		return
	}
	serviceName := endpoint.Name
	if err := k.initializeConns(serviceName); err != nil {
		fmt.Printf("Error initializing connections for %s: %v\n", serviceName, err)
	}
}

// =================

// initEndpoints initializes connections by fetching all available endpoints in the specified namespace.

// func (k *KubernetesConnManager) initEndpoints() error {
// 	k.mu.Lock()
// 	defer k.mu.Unlock()

// 	pods, err := k.clientset.CoreV1().Pods(k.namespace).List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		return fmt.Errorf("failed to list pods: %v", err)
// 	}

// 	for _, pod := range pods.Items {
// 		if pod.Status.Phase == v1.PodRunning {
// 			target := fmt.Sprintf("%s:%d", address.IP, port)
// 			conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 			conn, err := k.createGRPCConnection(pod)
// 			if err != nil {
// 				return fmt.Errorf("failed to create GRPC connection for pod %s: %v", pod.Name, err)
// 			}
// 			k.connMap[pod.Name] = append(k.connMap[pod.Name], conn)
// 		}
// 	}

// 	return nil
// }

// -----

// func (k *KubernetesConnManager) watchEndpoints1(serviceName string) {
// 	// watch for changes to the service's endpoints
// 	informerFactory := informers.NewSharedInformerFactory(k.clientset, time.Minute)
// 	endpointsInformer := informerFactory.Core().V1().Endpoints().Informer()

// 	endpointsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
// 		AddFunc: func(obj interface{}) {
// 			eps := obj.(*v1.Endpoints)
// 			if eps.Name == serviceName {
// 				k.initializeConns(serviceName)
// 			}
// 		},
// 		UpdateFunc: func(oldObj, newObj interface{}) {
// 			eps := newObj.(*v1.Endpoints)
// 			if eps.Name == serviceName {
// 				k.initializeConns(serviceName)
// 			}
// 		},
// 		DeleteFunc: func(obj interface{}) {
// 			eps := obj.(*v1.Endpoints)
// 			if eps.Name == serviceName {
// 				k.mu.Lock()
// 				defer k.mu.Unlock()
// 				for _, conn := range k.connMap[serviceName] {
// 					_ = conn.Close()
// 				}
// 				delete(k.connMap, serviceName)
// 			}
// 		},
// 	})

// 	informerFactory.Start(wait.NeverStop)
// 	informerFactory.WaitForCacheSync(wait.NeverStop)
// }

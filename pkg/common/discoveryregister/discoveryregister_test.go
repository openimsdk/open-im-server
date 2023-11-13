package discoveryregister

import (
	"context"
	"reflect"
	"testing"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"google.golang.org/grpc"
)

func TestNewDiscoveryRegister(t *testing.T) {
	type args struct {
		envType string
	}
	tests := []struct {
		name    string
		args    args
		want    discoveryregistry.SvcDiscoveryRegistry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDiscoveryRegister(tt.args.envType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDiscoveryRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDiscoveryRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewK8sDiscoveryRegister(t *testing.T) {
	tests := []struct {
		name    string
		want    discoveryregistry.SvcDiscoveryRegistry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewK8sDiscoveryRegister()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewK8sDiscoveryRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewK8sDiscoveryRegister() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sDR_Register(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	type args struct {
		serviceName string
		host        string
		port        int
		opts        []grpc.DialOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			if err := cli.Register(tt.args.serviceName, tt.args.host, tt.args.port, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("K8sDR.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestK8sDR_UnRegister(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			if err := cli.UnRegister(); (err != nil) != tt.wantErr {
				t.Errorf("K8sDR.UnRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestK8sDR_CreateRpcRootNodes(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	type args struct {
		serviceNames []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			if err := cli.CreateRpcRootNodes(tt.args.serviceNames); (err != nil) != tt.wantErr {
				t.Errorf("K8sDR.CreateRpcRootNodes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestK8sDR_RegisterConf2Registry(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	type args struct {
		key  string
		conf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			if err := cli.RegisterConf2Registry(tt.args.key, tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("K8sDR.RegisterConf2Registry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestK8sDR_GetConfFromRegistry(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			got, err := cli.GetConfFromRegistry(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("K8sDR.GetConfFromRegistry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("K8sDR.GetConfFromRegistry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sDR_GetConns(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	type args struct {
		ctx         context.Context
		serviceName string
		opts        []grpc.DialOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			got, err := cli.GetConns(tt.args.ctx, tt.args.serviceName, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("K8sDR.GetConns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("K8sDR.GetConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sDR_GetConn(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	type args struct {
		ctx         context.Context
		serviceName string
		opts        []grpc.DialOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			got, err := cli.GetConn(tt.args.ctx, tt.args.serviceName, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("K8sDR.GetConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("K8sDR.GetConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sDR_GetSelfConnTarget(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			if got := cli.GetSelfConnTarget(); got != tt.want {
				t.Errorf("K8sDR.GetSelfConnTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sDR_AddOption(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	type args struct {
		opts []grpc.DialOption
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			cli.AddOption(tt.args.opts...)
		})
	}
}

func TestK8sDR_CloseConn(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	type args struct {
		conn *grpc.ClientConn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			cli.CloseConn(tt.args.conn)
		})
	}
}

func TestK8sDR_GetClientLocalConns(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string][]*grpc.ClientConn
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			if got := cli.GetClientLocalConns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("K8sDR.GetClientLocalConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestK8sDR_Close(t *testing.T) {
	type fields struct {
		options         []grpc.DialOption
		rpcRegisterAddr string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &K8sDR{
				options:         tt.fields.options,
				rpcRegisterAddr: tt.fields.rpcRegisterAddr,
			}
			cli.Close()
		})
	}
}

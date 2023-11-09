package prom_metrics

import (
	"reflect"
	"testing"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/openimsdk/open-im-server/v3/pkg/common/ginprom"
	"github.com/prometheus/client_golang/prometheus"
)

func TestNewGrpcPromObj(t *testing.T) {
	type args struct {
		cusMetrics []prometheus.Collector
	}
	tests := []struct {
		name    string
		args    args
		want    *prometheus.Registry
		want1   *grpc_prometheus.ServerMetrics
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := NewGrpcPromObj(tt.args.cusMetrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGrpcPromObj() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGrpcPromObj() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NewGrpcPromObj() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetGrpcCusMetrics(t *testing.T) {
	type args struct {
		registerName string
	}
	tests := []struct {
		name string
		args args
		want []prometheus.Collector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGrpcCusMetrics(tt.args.registerName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGrpcCusMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGinCusMetrics(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want []*ginprom.Metric
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGinCusMetrics(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGinCusMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

package base_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"gitlab.slade360emr.com/go/base"
)

func TestRecordGraphqlResolverMetrics(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		startTime time.Time
		name      string
		e         error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success: ok status",
			args: args{
				ctx:       ctx,
				startTime: time.Now(),
				name:      "success_resolver",
				e:         nil,
			},
		},
		{
			name: "success: failed status",
			args: args{
				ctx:       ctx,
				startTime: time.Now(),
				name:      "failed_resolver",
				e:         fmt.Errorf("this resolver will fail"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base.RecordGraphqlResolverMetrics(tt.args.ctx, tt.args.startTime, tt.args.name, tt.args.e)
		})
	}
}

func TestMetricsCollectorService(t *testing.T) {

	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success:staging_env",
			args: args{
				serviceName: "base",
			},
			want: "base-staging",
		},
		{
			name: "success:testing_env",
			args: args{
				serviceName: "base",
			},
			want: "base-testing",
		},
		{
			name: "success:demo_env",
			args: args{
				serviceName: "base",
			},
			want: "base-demo",
		},
		{
			name: "success:prod_env",
			args: args{
				serviceName: "base",
			},
			want: "base-prod",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialEnvironment := os.Getenv("ENVIRONMENT")

			if tt.name == "success:staging_env" {
				os.Setenv("ENVIRONMENT", "staging")
				if got := base.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}
			if tt.name == "success:staging_env" {
				os.Setenv("ENVIRONMENT", "staging")
				if got := base.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}
			if tt.name == "success:testing_env" {
				os.Setenv("ENVIRONMENT", "testing")
				if got := base.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}
			if tt.name == "success:demo_env" {
				os.Setenv("ENVIRONMENT", "demo")
				if got := base.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}
			if tt.name == "success:prod_env" {
				os.Setenv("ENVIRONMENT", "prod")
				if got := base.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}

			os.Setenv("ENVIRONMENT", initialEnvironment)
		})
	}
}

func TestEnableStatsAndTraceExporters(t *testing.T) {
	type args struct {
		ctx     context.Context
		service string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success:enable exporters",
			args: args{
				ctx:     context.Background(),
				service: "test-service",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := base.EnableStatsAndTraceExporters(tt.args.ctx, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnableStatsAndTraceExporters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

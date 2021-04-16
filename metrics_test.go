package base_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
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

func TestMetricsResponseWriter_Header(t *testing.T) {
	rw := httptest.NewRecorder()
	m := base.NewMetricsResponseWriter(rw)

	emptyHeader := map[string][]string{}

	tests := []struct {
		name string
		want http.Header
	}{
		{
			name: "success:empty headers",
			want: emptyHeader,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MetricsResponseWriter.Header() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricsResponseWriter_WriteHeader(t *testing.T) {
	rw := httptest.NewRecorder()
	m := base.NewMetricsResponseWriter(rw)

	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success:default status code",
			args: args{
				code: http.StatusOK,
			},
		},
		{
			name: "success:change default status code",
			args: args{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.WriteHeader(tt.args.code)
			if m.StatusCode != tt.args.code {
				t.Errorf("base.MetricsResponseWriter.WriteHeader() = %v, want %v", rw.Code, tt.args.code)
			}
		})
	}
}

func TestMetricsResponseWriter_Write(t *testing.T) {
	rw := httptest.NewRecorder()
	m := base.NewMetricsResponseWriter(rw)

	sample := []byte("four")

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				b: sample,
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := m.Write(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("MetricsResponseWriter.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MetricsResponseWriter.Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomRequestMetricsMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	mw := base.CustomHTTPRequestMetricsMiddleware()
	h := mw(next)

	rw := httptest.NewRecorder()
	reader := bytes.NewBuffer([]byte("sample"))

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	h.ServeHTTP(rw, req)

}

func TestRecordStats(t *testing.T) {
	rw := httptest.NewRecorder()
	w := base.NewMetricsResponseWriter(rw)
	reader := bytes.NewBuffer([]byte("sample"))
	req := httptest.NewRequest(http.MethodPost, "/", reader)

	type args struct {
		w *base.MetricsResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "collect http metrics",
			args: args{
				w: w,
				r: req,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base.RecordHTTPStats(tt.args.w, tt.args.r)
		})
	}
}

func TestGenerateLatencyBounds(t *testing.T) {
	type args struct {
		max  int
		step int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "Happy case",
			args: args{
				max:  1000,
				step: 100,
			},
			want: []float64{0, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.GenerateLatencyBounds(tt.args.max, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateLatencyBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

package serverutils_test

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

	"github.com/savannahghi/serverutils"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
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
		t.Run(tt.name, func(_ *testing.T) {
			serverutils.RecordGraphqlResolverMetrics(tt.args.ctx, tt.args.startTime, tt.args.name, tt.args.e)
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

				if got := serverutils.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}

			if tt.name == "success:staging_env" {
				os.Setenv("ENVIRONMENT", "staging")

				if got := serverutils.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}

			if tt.name == "success:testing_env" {
				os.Setenv("ENVIRONMENT", "testing")

				if got := serverutils.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}

			if tt.name == "success:demo_env" {
				os.Setenv("ENVIRONMENT", "demo")

				if got := serverutils.MetricsCollectorService(tt.args.serviceName); got != tt.want {
					t.Errorf("MetricsCollectorService() = %v, want %v", got, tt.want)
				}
			}

			if tt.name == "success:prod_env" {
				os.Setenv("ENVIRONMENT", "prod")

				if got := serverutils.MetricsCollectorService(tt.args.serviceName); got != tt.want {
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
			_, err := serverutils.EnableStatsAndTraceExporters(tt.args.ctx, tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnableStatsAndTraceExporters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMetricsResponseWriter_Header(t *testing.T) {
	rw := httptest.NewRecorder()
	m := serverutils.NewMetricsResponseWriter(rw)

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
	m := serverutils.NewMetricsResponseWriter(rw)

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
				t.Errorf("serverutils.MetricsResponseWriter.WriteHeader() = %v, want %v", rw.Code, tt.args.code)
			}
		})
	}
}

func TestMetricsResponseWriter_Write(t *testing.T) {
	rw := httptest.NewRecorder()
	m := serverutils.NewMetricsResponseWriter(rw)

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
		t.Run(tt.name, func(_ *testing.T) {
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

func TestCustomRequestMetricsMiddleware(_ *testing.T) {
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	mw := serverutils.CustomHTTPRequestMetricsMiddleware()
	h := mw(next)

	rw := httptest.NewRecorder()
	reader := bytes.NewBuffer([]byte("sample"))

	req := httptest.NewRequest(http.MethodPost, "/", reader)
	h.ServeHTTP(rw, req)
}

func TestRecordStats(t *testing.T) {
	rw := httptest.NewRecorder()
	w := serverutils.NewMetricsResponseWriter(rw)
	reader := bytes.NewBuffer([]byte("sample"))
	req := httptest.NewRequest(http.MethodPost, "/", reader)

	type args struct {
		w *serverutils.MetricsResponseWriter
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
		t.Run(tt.name, func(_ *testing.T) {
			serverutils.RecordHTTPStats(tt.args.w, tt.args.r)
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
			if got := serverutils.GenerateLatencyBounds(tt.args.max, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateLatencyBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitOtelSDK(t *testing.T) {
	type args struct {
		ctx         context.Context
		serviceName string
	}

	tests := []struct {
		name      string
		args      args
		want      *tracesdk.TracerProvider
		wantErr   bool
		wantPanic bool
	}{
		{
			name: "success:initialize otel sdk",
			args: args{
				ctx:         context.Background(),
				serviceName: "test",
			},
			wantErr:   false,
			wantPanic: false,
		},
		{
			name: "fail:initialize otel sdk missing environment",
			args: args{
				ctx:         context.Background(),
				serviceName: "test",
			},
			wantErr:   true,
			wantPanic: true,
		},
		{
			name: "fail:initialize otel sdk missing jaeger env",
			args: args{
				ctx:         context.Background(),
				serviceName: "test",
			},
			wantErr:   true,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("InitOtelSDK() recover = %v, wantPanic = %v", r, tt.wantPanic)
				}
			}()

			initialEnvironment := os.Getenv("ENVIRONMENT")
			initialJaegerEnv := os.Getenv("JAEGER_URL")

			if tt.name == "success:initialize otel sdk" {
				os.Setenv("JAEGER_URL", "http://jaeger")
				os.Setenv("ENVIRONMENT", "staging")
			}

			if tt.name == "fail:initialize otel sdk missing environment" {
				os.Setenv("JAEGER_URL", "http://jaeger")
				os.Setenv("ENVIRONMENT", "")
			}

			if tt.name == "fail:initialize otel sdk missing jaeger env" {
				os.Setenv("JAEGER_URL", "")
			}

			_, err := serverutils.InitOtelSDK(tt.args.ctx, tt.args.serviceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitOtelSDK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			os.Setenv("ENVIRONMENT", initialEnvironment)
			os.Setenv("JAEGER_URL", initialJaegerEnv)
		})
	}
}

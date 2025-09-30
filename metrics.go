package serverutils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/api/option"
)

// GenerateLatencyBounds is used in generating latency bounds
// The arguments provided should be in millisecond format e.g 1s == 1000ms
// interval will be used as an increment value
// [>=0ms, >=100ms, >=200ms, >=300ms,...., >=1000ms]
func GenerateLatencyBounds(maxVal, interval int) []float64 {
	bounds := []float64{}
	for j := 0; j <= maxVal; j += interval {
		bounds = append(bounds, float64(j))
	}

	return bounds
}

var (
	// LatencyBounds used in aggregating latency
	// should be in ms i.e seconds written in ms eg 1s --> 1000ms
	// [>=0ms, >=10ms, >=20ms, >=30ms,...., >=4s, >=5s, >=6s >=7s]
	//
	// Disclaimer: The interval value should be reasonable so as to avoid many
	// buckets. If the distribution metrics has many buckets, it will not export
	// the metrics.
	LatencyBounds = GenerateLatencyBounds(60000, 200) // 1 min in intervals of 200ms

	// Measure

	GraphqlResolverLatency = stats.Float64(
		"graphql_resolver_latency",
		"The Latency in milliseconds per graphql resolver execution",
		"ms",
	)

	// Tags

	// Resolver is the Graphql resolver used when making a GraphQl request
	ResolverName = tag.MustNewKey("resolver.name")

	// Error is the error recorded if an error occurs
	ResolverErrorMessage = tag.MustNewKey("resolver.error")

	// ResolverStatus is used to tag whether passed or failed
	// it is either pass/fail...make constants
	ResolverStatus = tag.MustNewKey("resolver.status")

	// Views

	GraphqlResolverLatencyView = &view.View{
		Name:        "graphql_resolver_latency_distribution",
		Description: "Time taken by a graphql resolver",
		Measure:     GraphqlResolverLatency,
		// Latency in buckets:
		// [>=0ms, >=10ms, >=20ms, >=30ms,...., >=4s]
		Aggregation: view.Distribution(LatencyBounds...),
		TagKeys:     []tag.Key{ResolverName, ResolverErrorMessage, ResolverStatus},
	}

	GraphqlResolverCountView = &view.View{
		Name:        "graphql_resolver_request_count",
		Description: "The number of times a graphql resolver is executed",
		Measure:     GraphqlResolverLatency,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{ResolverName, ResolverErrorMessage, ResolverStatus},
	}

	// Following are views for the collected metrics i.e how they are exported to the various backends
	HTTPRequestLatency = stats.Float64(
		"http_request_latency",
		"The Latency in milliseconds per http request execution",
		"ms",
	)

	// Path is the URL path (not including query string) in the request.
	HTTPPath = tag.MustNewKey("http.path")

	// StatusCode is the numeric HTTP response status code.
	HTTPStatusCode = tag.MustNewKey("http.status")

	// Method is the HTTP method of the request.
	HTTPMethod = tag.MustNewKey("http.method")

	ServerRequestLatencyView = &view.View{
		Name:        "http_request_latency_distribution",
		Description: "Time taken to process a http request",
		Measure:     HTTPRequestLatency,
		// Latency in buckets:
		// [>=0ms, >=10ms, >=20ms, >=30ms,...., >=4s]
		Aggregation: view.Distribution(LatencyBounds...),
		TagKeys:     []tag.Key{HTTPPath, HTTPStatusCode, HTTPMethod},
	}

	ServerRequestCountView = &view.View{
		Name:        "http_request_count",
		Description: "The number of HTTP requests",
		Measure:     HTTPRequestLatency,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{HTTPPath, HTTPStatusCode, HTTPMethod},
	}

	// DefaultServiceViews are the default/common server views provided by base package
	// The views can be used by the various services
	DefaultServiceViews = []*view.View{GraphqlResolverLatencyView, GraphqlResolverCountView, ServerRequestLatencyView, ServerRequestCountView}
)

// Resolver status values
const (
	ResolverSuccessValue = "OK"
	ResolverFailureValue = "FAILED"
)

// GetRunningEnvironment returns the environment where the service is running. Important
// so as to point to the correct deps
func GetRunningEnvironment() string {
	return MustGetEnvVar(Environment)
}

// MetricsCollectorService returns name of service suffixed by it's running environment
// this helps identify metrics from different services at the backend/metrics viewer.
// e.g namespace in prometheus exporter
func MetricsCollectorService(serviceName string) string {
	var environment string

	if GetRunningEnvironment() == StagingEnv {
		environment = StagingEnv
	}

	if GetRunningEnvironment() == TestingEnv {
		environment = TestingEnv
	}

	if GetRunningEnvironment() == DemoEnv {
		environment = DemoEnv
	}

	if GetRunningEnvironment() == ProdEnv {
		environment = ProdEnv
	}

	return fmt.Sprintf("%s-%s", serviceName, environment)
}

// EnableStatsAndTraceExporters creates Google Cloud Monitoring exporter
func EnableStatsAndTraceExporters(ctx context.Context, service string) (func(), error) {
	exporter, err := metric.New(
		metric.WithProjectID(os.Getenv("GOOGLE_CLOUD_PROJECT")),
		metric.WithMonitoringClientOptions(
			option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create the Google Cloud Monitoring exporter: %w", err)
	}

	k8sService, _ := GetEnvVar("K_SERVICE")
	k8sRevision, _ := GetEnvVar("K_REVISION")
	k8sConfiguration, _ := GetEnvVar("K_CONFIGURATION")

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(service),
		attribute.String("cloudrun.service", k8sService),
		attribute.String("cloudrun.revision", k8sRevision),
		attribute.String("cloudrun.configuration", k8sConfiguration),
	)

	meterProvider := metricsdk.NewMeterProvider(
		metricsdk.WithReader(
			metricsdk.NewPeriodicReader(
				exporter,
				metricsdk.WithInterval(60*time.Second),
			),
		),
		metricsdk.WithResource(res),
	)

	otel.SetMeterProvider(meterProvider)

	deferFuncs := func() {
		if err := meterProvider.ForceFlush(ctx); err != nil {
			fmt.Printf("Error flushing metrics: %v\n", err)
		}

		if err := meterProvider.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down meter provider: %v\n", err)
		}
	}

	return deferFuncs, nil
}

// RecordGraphqlResolverMetrics records the metrics for a specific graphql resolver
// It should be deferred until the execution of the resolver function is completed
func RecordGraphqlResolverMetrics(ctx context.Context, startTime time.Time, name string, e error) {
	// check if there's an error
	if e != nil {
		ctx, _ = tag.New(ctx,
			tag.Insert(ResolverStatus, ResolverFailureValue),
			tag.Insert(ResolverErrorMessage, e.Error()),
		)
	}

	ctx, _ = tag.New(ctx,
		tag.Insert(ResolverName, name),
		tag.Insert(ResolverStatus, ResolverSuccessValue),
	)

	// returns a duration - time elapsed
	duration := time.Since(startTime)

	// duration is in nanoseconds (ns)
	// 1ms = 1000000 ns
	latency := float64(duration / 1000000)

	// Record the starts
	stats.Record(ctx, GraphqlResolverLatency.M(latency))
}

// CustomHTTPRequestMetricsMiddleware is used to implement custom metrics for our http requests
// The custom middleware used to collect any custom http request stats
// It will also be used to capture distributed trace requests and propagate them through context
func CustomHTTPRequestMetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				newResponseWriter := NewMetricsResponseWriter(w)

				next.ServeHTTP(newResponseWriter, r)

				RecordHTTPStats(newResponseWriter, r)
			},
		)
	}
}

// RecordHTTPStats adds tags and records the metrics for a request
func RecordHTTPStats(w *MetricsResponseWriter, r *http.Request) {
	ctx, _ := tag.New(r.Context(),
		tag.Insert(HTTPPath, r.URL.Path),
		tag.Insert(HTTPMethod, r.Method),
		tag.Insert(HTTPStatusCode, fmt.Sprint(w.StatusCode)))

	duration := time.Since(w.StartTime)

	// duration is in nanoseconds (ns)
	// 1ms = 1000000 ns
	latency := float64(duration / 1000000)

	// Record the starts
	stats.Record(ctx, HTTPRequestLatency.M(latency))
}

// MetricsResponseWriter implements the http.ResponseWriter Interface
// it is a wrapper of http.ResponseWriter and enables obtaining measures
type MetricsResponseWriter struct {
	w          http.ResponseWriter
	StatusCode int
	StartTime  time.Time
}

// NewMetricsResponseWriter new http.ResponseWriter wrapper
func NewMetricsResponseWriter(w http.ResponseWriter) *MetricsResponseWriter {
	return &MetricsResponseWriter{w, http.StatusOK, time.Now()}
}

// Header ...
func (m *MetricsResponseWriter) Header() http.Header {
	return m.w.Header()
}

// WriteHeader ...
func (m *MetricsResponseWriter) WriteHeader(code int) {
	m.StatusCode = code
	m.w.WriteHeader(code)
}

// Write ...
func (m *MetricsResponseWriter) Write(b []byte) (int, error) {
	size, err := m.w.Write(b)
	return size, err
}

// InitOtelSDK returns an OpenTelemetry TracerProvider configured to use
// the OTLP HTTP exporter. The returned
// TracerProvider will also use a Resource configured with all the information
// about the service deployment.
func InitOtelSDK(ctx context.Context, serviceName string) (*tracesdk.TracerProvider, error) {
	otlpEndpoint := MustGetEnvVar("JAEGER_URL")

	if !strings.Contains(otlpEndpoint, "/v1/traces") {
		otlpEndpoint = strings.TrimSuffix(otlpEndpoint, "/") + "/v1/traces"
	}

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(otlpEndpoint))
	if err != nil {
		return nil, err
	}

	// Environment variables automatically added to the running containers in cloud run
	// K_SERVICE	The name of the Cloud Run service being run.
	service, _ := GetEnvVar("K_SERVICE")

	// K_REVISION	The name of the Cloud Run revision being run.
	revision, _ := GetEnvVar("K_REVISION")

	// K_CONFIGURATION	The name of the Cloud Run configuration that created the revision.
	configuration, _ := GetEnvVar("K_CONFIGURATION")

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(MetricsCollectorService(serviceName)),
		attribute.String("cloudrun.service", service),
		attribute.String("cloudrun.revision", revision),
		attribute.String("cloudrun.configuration", configuration),
	)

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource),
	)

	// Register the TracerProvider as the global so any imported
	// instrumentation will default to using it
	otel.SetTracerProvider(tp)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	)
	otel.SetTextMapPropagator(propagator)

	return tp, nil
}

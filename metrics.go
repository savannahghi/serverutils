package base

import (
	"context"
	"fmt"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Server HTTP measures used to record metrics
var (
	GraphqlResolverLatency = stats.Float64(
		"graphql_resolver_latency",
		"The Latency in milliseconds per graphql resolver execution",
		"ms",
	)
)

// tags associated with colleted metrics
var (

	// Resolver is the Graphql resolver used when making a GraphQl request
	ResolverName = tag.MustNewKey("resolver.name")

	// Error is the error recorded if an error occurs
	ResolverErrorMessage = tag.MustNewKey("resolver.error")

	// ResolverStatus is used to tag whether passed or failed
	// it is either pass/fail...make constants
	ResolverStatus = tag.MustNewKey("resolver.status")
)

// Resolver status values
const (
	ResolverSuccessValue = "OK"
	ResolverFailureValue = "FAILED"
)

// Views for the collected metrics i.e how they are exported to the various backends
var (
	GraphqlResolverLatencyView = &view.View{
		Name:        "graphql_resolver_latency_distribution",
		Description: "Time taken by a graphql resolver",
		Measure:     GraphqlResolverLatency,
		// Latency in buckets:
		// [>=0ms, >=10ms, >=20ms, >=30ms,...]
		Aggregation: view.Distribution(0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100),
		TagKeys:     []tag.Key{ResolverName, ResolverErrorMessage, ResolverStatus},
	}

	GraphqlResolverCountView = &view.View{
		Name:        "graphql_resolver_request_count",
		Description: "The number of times a graphql resolver is executed",
		Measure:     GraphqlResolverLatency,
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{ResolverName, ResolverErrorMessage, ResolverStatus},
	}
)

// DefaultServiceViews are the default/common server views provided by base package
// The views can be used by the various services
var DefaultServiceViews = []*view.View{GraphqlResolverLatencyView, GraphqlResolverCountView}

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

// EnableStatsAndTraceExporters a wrapper for initializing metrics exporters
// TODO:Look into improvements
func EnableStatsAndTraceExporters(ctx context.Context, service string) (func(), error) {
	// Enable OpenCensus exporters to export metrics
	// to Stackdriver Monitoring.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		MetricPrefix: service,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create the stackdriver exporter: %v", err)
	}

	// Start the metrics exporter
	if err := exporter.StartMetricsExporter(); err != nil {
		return nil, fmt.Errorf("error starting metric exporter: %v", err)
	}

	deferFuncs := func() {
		exporter.Flush()
		exporter.StopMetricsExporter()
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

	duration := time.Since(startTime)
	latency := float64(duration / 1000000)

	// Record the starts
	stats.Record(ctx, GraphqlResolverLatency.M(latency))
}

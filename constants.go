package serverutils

const (
	// AppName is the name of "this server"
	AppName = "api-gateway"

	// DSNEnvVarName is the Sentry reporting config
	DSNEnvVarName = "SENTRY_DSN"

	// AppVersion is the app version (used for StackDriver error reporting)
	AppVersion = "0.0.1"

	// PortEnvVarName is the name of the environment variable that defines the
	// server port
	PortEnvVarName = "PORT"

	// DefaultPort is the default port at which the server listens if the port
	// environment variable is not set
	DefaultPort = "8080"

	// GoogleCloudProjectIDEnvVarName is used to determine the ID of the GCP project e.g for setting up StackDriver client
	GoogleCloudProjectIDEnvVarName = "GOOGLE_CLOUD_PROJECT"

	// DebugEnvVarName is used to determine if we should print extended tracing / logging (debugging aids)
	// to the console
	DebugEnvVarName = "DEBUG"

	// IsRunningTestsEnvVarName is used to determine if we are running in a test environment
	IsRunningTestsEnvVarName = "IS_RUNNING_TESTS"

	// Environment points to where this service is running e.g staging, testing, prod
	Environment = "ENVIRONMENT"

	// StagingEnv runs the service under staging
	StagingEnv = "staging"

	// DemoEnv runs the service under demo
	DemoEnv = "demo"

	// TestingEnv runs the service under testing
	TestingEnv = "testing"

	// ProdEnv runs the service under production
	ProdEnv = "prod"

	// TraceSampleRateEnvVarName indicates the percentage of transactions to be captured when doing performance monitoring
	TraceSampleRateEnvVarName = "SENTRY_TRACE_SAMPLE_RATE"
)

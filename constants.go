package server_utils

const (
	// TestUserPhoneNumber is used by integration tests
	TestUserPhoneNumber        = "+254711223344"
	TestUserPhoneNumberWithPin = "+254778990088"
	TestUserPin                = "1234"
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

	// TestsEnvVarName is used to determine if we are running in a test environment
	IsRunningTestsEnvVarName = "IS_RUNNING_TESTS"
)

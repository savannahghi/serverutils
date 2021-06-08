package base

import "time"

/* #nosec */
const (
	tokenMinLength = 12

	apiPasswordMinLength = 3

	TokenExpiryRatio = 0.95 // Refresh access tokens after 95% of the time is spent

	meURLFragment = "v1/user/me/?format=json"

	DateLayout = "2006-01-02"

	DateTimeFormatLayout = "2006-01-02T15:04:05+03:00"

	defaultRegion = "KE"

	// Sep is a separator, used to create "opaque" IDs
	Sep = "|"

	// DefaultPageSize is used to paginate records (e.g those fetched from Firebase)
	// if there is no user specified page size
	DefaultPageSize = 100

	// FirebaseWebAPIKeyEnvVarName is the name of the env var that holds a Firebase web API key
	// for this project
	FirebaseWebAPIKeyEnvVarName = "FIREBASE_WEB_API_KEY"

	// FirebaseCustomTokenSigninURL is the Google Identity Toolkit API for signing in over REST
	FirebaseCustomTokenSigninURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key="

	// FirebaseRefreshTokenURL is used to request Firebase refresh tokens from Google APIs
	FirebaseRefreshTokenURL = "https://securetoken.googleapis.com/v1/token?key="

	// GoogleApplicationCredentialsEnvVarName is used to obtain service account details from the
	// local server when necessary e.g when running tests on CI or a local developer setup
	GoogleApplicationCredentialsEnvVarName = "GOOGLE_APPLICATION_CREDENTIALS"

	// GoogleProjectNumberEnvVarName is a numeric project number that
	GoogleProjectNumberEnvVarName = "GOOGLE_PROJECT_NUMBER"

	// FDLDomainEnvironmentVariableName is the name of the domain used for short
	// links.
	//
	// e.g https://healthcloud.page.link or https://bwl.page.link
	FDLDomainEnvironmentVariableName = "FIREBASE_DYNAMIC_LINKS_DOMAIN"

	// ServerPublicDomainEnvironmentVariableName is the name of the environment
	// variable at which the server is deployed. It is used to generate long
	// links for shortening
	ServerPublicDomainEnvironmentVariableName = "SERVER_PUBLIC_DOMAIN"

	// TestUserEmail is used by integration tests
	TestUserEmail = "be.well@bewell.co.ke"

	// OTPCollectionName is the name of the collection used to persist single
	// use verification codes on Firebase
	OTPCollectionName         = "otps"
	EmailOptInCollectionName  = "email_opt_ins"
	PhoneOptInCollectionName  = "phone_opt_ins"
	USSDSessionCollectionName = "ussd_signup_sessions"

	// IdentifierCollectionName is used to record randomly generated identifiers so that they
	// are not re-issued
	IdentifierCollectionName = "identifiers"

	// DefaultCalendarEmail is the email address that "owns" the calendar by default
	DefaultCalendarEmail = "be.well@bewell.co.ke"

	// BeWellVirtualPayerSladeCode is the Slade Code for the virtual provider used by the Be.Well app for e.g telemedicine
	BeWellVirtualPayerSladeCode = 2019 // PRO-4683

	// BeWellVirtualProviderSladeCode is the Slade Code for the virtual payer used by the Be.Well app for e.g healthcare lending
	BeWellVirtualProviderSladeCode = 4683 // PAY-2019

	// DefaultRESTAPIPageSize is the page size to use when calling Slade REST API services if the
	// client does not specify a page size
	DefaultRESTAPIPageSize = 100

	// MaxRestAPIPageSize is the largest page size we'll request
	MaxRestAPIPageSize = 250

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

	// CIEnvVarName is set to "true" in CI environments e.g Gitlab CI, Github actions etc.
	// It can be used to opt in to / out of tests in such environments
	CIEnvVarName = "CI"

	// AuthTokenContextKey is used to add/retrieve the Firebase UID on the context
	AuthTokenContextKey = ContextKey("UID")

	// HTTP client settings
	HTTPClientTimeoutSecs = 10

	//RootCollectionSuffix ...
	RootCollectionSuffix = "ROOT_COLLECTION_SUFFIX"

	//Anonymous user identifier
	AnonymousUser = "anonymous user"

	// TestUserPhoneNumber is used by integration tests
	TestUserPhoneNumber        = "+254711223344"
	TestUserPhoneNumberWithPin = "+254778990088"
	TestUserPin                = "1234"

	// Pins collection name
	PINCollectionName = "pins"

	// Secret Key for signing json web tokens
	JWTSecretKey = "JWT_KEY"

	// The file that contains dependency definition. Each service which depends on other service
	// via REST, need to have this file in their root
	DepsFileName = "deps.yaml"

	// env variable pointing to where this service is running e.g staging, testing, prod
	Environment = "ENVIRONMENT"

	// running the service under staging
	StagingEnv = "staging"

	// running the service under demo
	DemoEnv = "demo"

	// running the service under testing
	TestingEnv = "testing"

	// running the service under production
	ProdEnv = "prod"

	// running the service under e2e
	E2eEnv = "e2e"

	// Default login client settings (env var names)
	ClientIDEnvVarName     = "CLIENT_ID"
	ClientSecretEnvVarName = "CLIENT_SECRET"
	UsernameEnvVarName     = "USERNAME"
	PasswordEnvVarName     = "PASSWORD"
	GrantTypeEnvVarName    = "GRANT_TYPE"
	APISchemeEnvVarName    = "API_SCHEME"
	TokenURLEnvVarName     = "TOKEN_URL"
	APIHostEnvVarName      = "HOST"
	WorkstationEnvVarName  = "DEFAULT_WORKSTATION_ID"
	WorkstationHeaderName  = "X-WORKSTATION"
)

var (

	// TimeLocation default timezone
	TimeLocation, _ = time.LoadLocation("Africa/Nairobi")
)

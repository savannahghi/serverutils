package base

import (
	"io"
	"net/http"
	"time"
)

// MockClient creates a mock EDI client for testing
type MockClient struct {

	// Scheme is used to set up the EDI Scheme
	Scheme string

	// Host is used to set up the EDI host
	Host string

	// Initialized is used to mark the mock client as initialized / not
	Initialized bool

	// ServerResp is used to set a mock response from the target server
	ServerResp *http.Response

	// ServerErr is used to simulate an error from communicating with a server
	ServerErr error

	// MeURLErr is used to simulate an error from calling an auth server /me endpoint
	MeURLErr error

	// Grant is used to set the mock OAuth 2 grant type
	Grant string

	domain           string
	clientID         string
	clientSecret     string
	apiTokenURL      string
	authServerDomain string
	username         string
	password         string

	// these are set after successful "authentication"
	accessToken  string
	tokenType    string
	refreshToken string
	accessScope  string
	expiresIn    int
	refreshAt    time.Time

	// these can be used to simulate different behavior in handling of auth server responses
	oauthRespUpdateFunc func(authResp *OAUTHResponse)
	meURL               string

	// use this to affect HTTP behavior e.g to cause a HTTP "error" by replacing the HTTP client transport
	httpClient *http.Client

	// set these to easily simulate errors
	refreshErr error
	authErr    error
}

// MeURL returns the address from which the user's profile should be retrieved
func (c *MockClient) MeURL() (string, error) {
	return c.meURL, c.MeURLErr
}

// HTTPClient returns the mocks HTTP client
func (c *MockClient) HTTPClient() *http.Client {
	return c.httpClient
}

// AccessToken returns the mocks access token
func (c *MockClient) AccessToken() string {
	return c.accessToken
}

// TokenType returns the mocks token type
func (c *MockClient) TokenType() string {
	return c.tokenType
}

// RefreshToken returns the mock refresh token
func (c *MockClient) RefreshToken() string {
	return c.refreshToken
}

// AccessScope returns the mock access scope
func (c *MockClient) AccessScope() string {
	return c.accessScope
}

// ExpiresIn returns the mock expires in
func (c *MockClient) ExpiresIn() int {
	return c.expiresIn
}

// RefreshAt returns the mock refresh at
func (c *MockClient) RefreshAt() time.Time {
	return c.refreshAt
}

// IsInitialized reports the value of the isInitialized boolean on the mock
func (c *MockClient) IsInitialized() bool {
	return c.Initialized
}

// Refresh returns an error (or nil) as set up on the mock
func (c *MockClient) Refresh() error {
	return c.refreshErr
}

// Authenticate returns an error or nil, as configured on the mock
func (c *MockClient) Authenticate() error {
	return c.authErr
}

// MakeRequest returns a configured mock response and error
// method, url and body params are ignored
func (c *MockClient) MakeRequest(_ string, _ string, _ io.Reader) (*http.Response, error) {
	return c.ServerResp, c.ServerErr
}

// APIScheme returns the configured API scheme from the mock
func (c *MockClient) APIScheme() string {
	return c.Scheme
}

// APIHost returns the configured API host from the mock
func (c *MockClient) APIHost() string {
	return c.Host
}

// Domain returns the configured domain from the mock
func (c *MockClient) Domain() string {
	return c.domain
}

// ClientID returns the configured clientID from the mock
func (c *MockClient) ClientID() string {
	return c.clientID
}

// ClientSecret returns the configured clientSecret from the mock
func (c *MockClient) ClientSecret() string {
	return c.clientSecret
}

// APITokenURL returns the configured auth server API token URL from the mock
func (c *MockClient) APITokenURL() string {
	return c.apiTokenURL
}

// AuthServerDomain returns the configured auth server domain from the mock
func (c *MockClient) AuthServerDomain() string {
	return c.authServerDomain
}

// GrantType returns the configured grant type from the mock
func (c *MockClient) GrantType() string {
	return c.Grant
}

// Username returns the configured Username from the mock
func (c *MockClient) Username() string {
	return c.username
}

// Password returns the configured Password from the mock
func (c *MockClient) Password() string {
	return c.password
}

// UpdateAuth applies values from the supplied OAUTHResponse to the mock client through a configured mock function
func (c *MockClient) UpdateAuth(authResp *OAUTHResponse) {
	c.oauthRespUpdateFunc(authResp)
}

// SetInitialized sets the value of the isInitialized boolean on the mock
func (c *MockClient) SetInitialized(isInitialized bool) {
	c.Initialized = isInitialized
}

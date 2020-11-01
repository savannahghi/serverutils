package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GetJWTKey returns a byte slice of the JWT secret key
func GetJWTKey() []byte {
	key := MustGetEnvVar(JWTSecretKey)
	return []byte(key)
}

// Claims a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide public claims
// Provides way for adding private claims
type Claims struct {
	jwt.StandardClaims
}

// GetServiceEnvironmentSuffix get the env suffix where the app is running
// e.g testing, staging, prod, local
func GetServiceEnvironmentSuffix() string {
	environment := MustGetEnvVar(ServiceEnvironmentSuffix)

	return environment
}

// InterServiceClient defines a client for use in interservice communication
type InterServiceClient struct {
	ServiceName string

	environment string
	apiScheme   string
	domain      string
	httpClient  *http.Client
	accessToken string
}

// NewInterserviceClientImpl takes an implementation of the `base.Service`
// interface and uses it to initialize a client for service to service
// communication.
//
// It should be used in preference to the deprecated `NewInterserviceClient`.
func NewInterserviceClientImpl(
	serviceName string,
	paths map[string]string,
) (*InterServiceClient, error) {
	env := GetServiceEnvironmentSuffix()
	return &InterServiceClient{
		ServiceName: serviceName,
		environment: env,
		apiScheme:   "https",
		domain:      "healthcloud.co.ke",
		httpClient: &http.Client{
			Timeout: time.Duration(1 * time.Minute),
		},
	}, nil
}

// NewInterserviceClient initializes a new interservice client
//
// Deprecated: a proper constructor should take service parameters rather than
// hard code specific services within itself.
func NewInterserviceClient(service string) (*InterServiceClient, error) {
	env := GetServiceEnvironmentSuffix()
	return &InterServiceClient{
		ServiceName: service,
		environment: env,
		apiScheme:   "https",
		domain:      "healthcloud.co.ke",
		httpClient: &http.Client{
			Timeout: time.Duration(1 * time.Minute),
		},
	}, nil
}

// CreateAuthToken returns a signed JWT for use in authentication.
func (c InterServiceClient) CreateAuthToken() (string, error) {
	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    c.GenerateBaseURL(c.ServiceName),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTKey())
	if err != nil {
		return "", fmt.Errorf("failed to create token with err: %v", err)
	}

	c.accessToken = tokenString
	return tokenString, nil
}

// GenerateBaseURL generates a URL depending on the environment
func (c InterServiceClient) GenerateBaseURL(service string) string {
	var address string
	if c.environment == "local" {
		port := MustGetEnvVar("PORT")
		address = "http://localhost:" + port
	} else {

		subdomain := fmt.Sprintf("%v-%v", service, c.environment)
		address = fmt.Sprintf("%v://%v.%v", c.apiScheme, subdomain, c.domain)
	}

	return address
}

// GenerateRequestURL generate a url with path for requested resource.
func (c InterServiceClient) GenerateRequestURL(service string, path string) string {

	address := c.GenerateBaseURL(service)

	return fmt.Sprintf("%v/%v", address, path)
}

// MakeRequest performs an inter service http request and returns a response
func (c InterServiceClient) MakeRequest(method string, url string, body interface{}) (*http.Response, error) {

	token, tknErr := c.CreateAuthToken()
	if tknErr != nil {
		return nil, tknErr
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(encoded)

	req, reqErr := http.NewRequest(method, url, payload)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, respErr := c.httpClient.Do(req)
	if respErr != nil {
		return nil, respErr
	}

	if resp.StatusCode > 201 {
		return nil, fmt.Errorf("bad response got: %v", resp.StatusCode)
	}

	return resp, nil
}

// jwtCheckFn is a function type for authorization and authentication checks
// there can be several e.g an authentication check runs first then an authorization
// check runs next if the authentication passes etc
type jwtCheckFn = func(r *http.Request) (bool, map[string]string, *jwt.Token)

// InterServiceAuthenticationMiddleware handles jwt authentication
func InterServiceAuthenticationMiddleware() func(http.Handler) http.Handler {
	// multiple checks can be run in sequence
	jwtCheckFuncs := []jwtCheckFn{HasValidJWTBearerToken}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				errs := []map[string]string{}

				for _, checkFunc := range jwtCheckFuncs {
					shouldContinue, errMap, _ := checkFunc(r)
					if shouldContinue {

						next.ServeHTTP(w, r)
						return
					}
					errs = append(errs, errMap)
				}

				WriteJSONResponse(w, errs, http.StatusUnauthorized)
			})
	}
}

// HasValidJWTBearerToken returns true with no errors if the request has a valid bearer token in the authorization header.
// Otherwise, it returns false and the error in a map with the key "error"
func HasValidJWTBearerToken(r *http.Request) (bool, map[string]string, *jwt.Token) {
	bearerToken, err := ExtractBearerToken(r)
	if err != nil {

		return false, ErrorMap(err), nil
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJWTKey(), nil
	})

	if err != nil {
		return false, ErrorMap(err), nil
	}

	return true, nil, token
}

package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ISCService defines the blueprint of a dependency service. This struct is here to maintain
// uniform structure definitions
type ISCService struct {
	// The name of the service that is been depended upon e.g mailgun, mpesa
	Name string

	// The endpoint where the service serves requests. The dependant should know forehand where to
	// this services lives
	RootDomain string
}

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

// InterServiceClient defines a client for use in interservice communication
type InterServiceClient struct {
	Name              string
	RequestRootDomain string
	httpClient        *http.Client
}

// NewInterserviceClient initializes a new interservice client
func NewInterserviceClient(s ISCService) (*InterServiceClient, error) {
	return &InterServiceClient{
		Name:              s.Name,
		RequestRootDomain: s.RootDomain,
		httpClient: &http.Client{
			Timeout: time.Duration(1 * time.Minute),
		},
	}, nil
}

// CreateAuthToken returns a signed JWT for use in authentication.
func (c InterServiceClient) CreateAuthToken() (string, error) {
	claims := &Claims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTKey())
	if err != nil {
		return "", fmt.Errorf("failed to create token with err: %v", err)
	}

	return tokenString, nil
}

// GenerateRequestURL generate a url with path for requested resource.
func (c InterServiceClient) generateRequestURL(path string) string {
	return fmt.Sprintf("%v/%v", c.RequestRootDomain, path)
}

// MakeRequest performs an inter service http request and returns a response
func (c InterServiceClient) MakeRequest(method string, path string, body interface{}) (*http.Response, error) {

	url := c.generateRequestURL(path)

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

	return c.httpClient.Do(req)
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
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected sigining method")
		}
		return GetJWTKey(), nil
	})

	if err != nil {
		return false, ErrorMap(err), nil
	}

	return true, nil, token
}

// Dep is the dependency definition
type Dep struct {
	DepName       string `yaml:"depName"`
	DepRootDomain string `yaml:"depRootDomain"`
}

//DepsConfig is the config for dependencies of a particular service
type DepsConfig struct {
	Staging    []Dep `yaml:"staging"`
	Testing    []Dep `yaml:"testing"`
	Production []Dep `yaml:"production"`
}

// PathToDepsFile return the path to deps.yaml file
func PathToDepsFile() string {
	cwd, _ := os.Getwd()
	return getDepsPath(filepath.Join(cwd, DepsFileName))
}

// recursively get the path to the deps.yaml file
func getDepsPath(path string) string {
	_, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		n := strings.Split(filepath.Dir(path), "/")
		m := n[:len(n)-1]
		p := filepath.Join(strings.Join(m, "/"), DepsFileName)
		return getDepsPath(p)
	}
	return path
}

// GetRunningEnvironment returns the environment wheere the service is running. Importannt
// so as to point to the correct deps
func GetRunningEnvironment() string {
	return MustGetEnvVar(Environment)
}

// GetDepFromConfig retrives a specific config from config slice
func GetDepFromConfig(name string, config []Dep) *Dep {
	var d Dep
	for _, dep := range config {
		if dep.DepName == name {
			d = dep
		}
	}
	return &d
}

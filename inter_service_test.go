package base_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestGetJWTKey(t *testing.T) {
	existingJWTKey, err := base.GetEnvVar("JWT_KEY")
	if err != nil {
		existingJWTKey = "an open secret"
	}
	os.Setenv("JWT_KEY", "an open secret")
	tests := []struct {
		name string
		want string
	}{
		{
			name: "JWT key environment variable",
			want: "an open secret",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.GetJWTKey(); string(got) != tt.want {
				t.Errorf("base.GetJWTKey() = %v, want %v", got, tt.want)
			}
		})
		os.Setenv("JWT_KEY", existingJWTKey)
	}
}

func TestNewInterserviceClient(t *testing.T) {
	srv, err := base.NewInterserviceClient(base.ISCService{Name: "otp", RootDomain: "https://example.com"})
	assert.Nil(t, err)
	assert.NotNil(t, srv)
	assert.Equal(t, "otp", srv.Name)
	assert.Equal(t, "https://example.com", srv.RequestRootDomain)
}

func TestInterServiceClient_CreateAuthToken(t *testing.T) {
	service, _ := base.NewInterserviceClient(base.ISCService{Name: "otp", RootDomain: "https://example.com"})
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "success create token",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := service
			got, err := c.CreateAuthToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("InterServiceClient.CreateAuthToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.NotEmpty(t, got) {
				t.Errorf("InterServiceClient.CreateAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterServiceClient_MakeRequest(t *testing.T) {
	type args struct {
		name     string
		endpoint string
		method   string
		path     string
		body     interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "Valid request",
			args: args{
				name:     "otp",
				endpoint: "https://example.com",
				method:   http.MethodPost,
				path:     "",
				body:     nil,
			},
			wantErr: false,
		},
		{
			name: "Invalid request",
			args: args{
				name:     "otp",
				endpoint: "https://google.com",
				method:   http.MethodPost,
				path:     "",
				body:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := base.NewInterserviceClient(base.ISCService{Name: tt.args.name, RootDomain: tt.args.endpoint})

			got, err := c.MakeRequest(tt.args.method, tt.args.path, tt.args.body)

			if err != nil && !tt.wantErr {
				t.Errorf("InterServiceClient.MakeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.Nil(t, err)
			}

			if tt.wantErr {
				assert.NotEqual(t, http.StatusOK, got.StatusCode)
			}

		})
	}

}

func TestHasValidJWTBearerToken(t *testing.T) {
	service, _ := base.NewInterserviceClient(base.ISCService{Name: "otp", RootDomain: "https://example.com"})

	validTokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	validToken, _ := service.CreateAuthToken()
	validTokenRequest.Header.Set("Authorization", "Bearer "+validToken)

	emptyHeaderRequest := httptest.NewRequest(http.MethodGet, "/", nil)

	invalidSignatureTokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	invalidSignatureToken, _ := createInvalidAuthToken()
	invalidSignatureTokenRequest.Header.Set("Authorization", "Bearer "+invalidSignatureToken)

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 map[string]string
	}{
		{
			name: "Request with valid JWT",
			args: args{
				r: validTokenRequest,
			},
			want:  true,
			want1: nil,
		},
		{
			name: "Request without a JWT",
			args: args{
				r: emptyHeaderRequest,
			},
			want:  false,
			want1: map[string]string{"error": "expected an `Authorization` request header"},
		},
		{
			name: "Request with an invalid signature",
			args: args{
				r: invalidSignatureTokenRequest,
			},
			want:  false,
			want1: map[string]string{"error": "signature is invalid"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, _ := base.HasValidJWTBearerToken(tt.args.r)
			if got != tt.want {
				t.Errorf("HasValidJWTBearerToken() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("HasValidJWTBearerToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestInterServiceAuthenticationMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	mw := base.InterServiceAuthenticationMiddleware()
	h := mw(next)
	rw := httptest.NewRecorder()
	reader := bytes.NewBuffer([]byte("sample"))

	service, _ := base.NewInterserviceClient(base.ISCService{Name: "otp", RootDomain: "https://example.com"})
	token, _ := service.CreateAuthToken()
	authHeader := fmt.Sprintf("Bearer %s", token)
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Add("Authorization", authHeader)
	h.ServeHTTP(rw, req)

	rw1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodPost, "/", reader)
	h.ServeHTTP(rw1, req1)
}

func createInvalidAuthToken() (string, error) {
	claims := &base.Claims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("Just bad"))
	if err != nil {
		return "", fmt.Errorf("failed to create token with err: %v", err)
	}

	return tokenString, nil
}

func TestGetDepFromConfig(t *testing.T) {
	deps := []base.Dep{
		{
			DepName:       "one",
			DepRootDomain: "https://one.com",
		},
		{
			DepName:       "two",
			DepRootDomain: "https://two.com",
		},
		{
			DepName:       "three",
			DepRootDomain: "https://three.com",
		},
		{
			DepName:       "four",
			DepRootDomain: "https://four.com",
		},
	}

	one := base.GetDepFromConfig("one", deps)
	assert.NotNil(t, one)
	assert.Equal(t, "one", one.DepName)
	assert.Equal(t, "https://one.com", one.DepRootDomain)

	two := base.GetDepFromConfig("two", deps)
	assert.NotNil(t, two)
	assert.Equal(t, "two", two.DepName)
	assert.Equal(t, "https://two.com", two.DepRootDomain)

	three := base.GetDepFromConfig("three", deps)
	assert.NotNil(t, three)
	assert.Equal(t, "three", three.DepName)
	assert.Equal(t, "https://three.com", three.DepRootDomain)

	four := base.GetDepFromConfig("four", deps)
	assert.NotNil(t, four)
	assert.Equal(t, "four", four.DepName)
	assert.Equal(t, "https://four.com", four.DepRootDomain)
}

func TestGetPathToDepsFile(t *testing.T) {
	p := base.PathToDepsFile()
	assert.NotEmpty(t, p)
	assert.Equal(t, true, strings.HasSuffix(p, base.DepsFileName))
}

func TestLoadDepsFromYAML(t *testing.T) {
	got, err := base.LoadDepsFromYAML()
	if err != nil {
		t.Errorf("can't load deps from YAML: %v", err)
		return
	}

	if got == nil {
		t.Errorf("got back nil deps")
		return
	}
}

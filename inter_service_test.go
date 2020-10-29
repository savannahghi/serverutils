package base_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestGetServiceEnvirionmentSuffix(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "service environment variable",
			want: "testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.GetServiceEnvirionmentSuffix(); got != tt.want {
				t.Errorf("base.GetServiceEnvirionmentSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetJWTKey(t *testing.T) {
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
	}
}

func TestNewInterserviceClient(t *testing.T) {
	srv, _ := base.NewInterserviceClient("base")
	type args struct {
		service string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.InterServiceClient
		wantErr bool
	}{
		{
			name: "create inter service client success",
			args: args{
				service: "base",
			},
			want:    srv,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.NewInterserviceClient(tt.args.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("base.NewInterserviceClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("base.NewInterserviceClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterServiceClient_CreateAuthToken(t *testing.T) {
	service, _ := base.NewInterserviceClient("base")
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

func TestInterServiceClient_GenerateBaseURL(t *testing.T) {
	service, _ := base.NewInterserviceClient("base")
	type args struct {
		service string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Generate example url",
			args: args{
				service: "example",
			},
			want: "https://example-testing.healthcloud.co.ke",
		},
		{
			name: "Generate initialized service url",
			args: args{
				service: service.Mailgun.Name,
			},
			want: "https://mailgun-testing.healthcloud.co.ke",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := service
			if got := c.GenerateBaseURL(tt.args.service); got != tt.want {
				t.Errorf("InterServiceClient.GenerateBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterServiceClient_GenerateRequestURL(t *testing.T) {
	service, _ := base.NewInterserviceClient("base")
	type args struct {
		service string
		path    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Generate example url path",
			args: args{
				service: "example",
				path:    "example_path",
			},
			want: "https://example-testing.healthcloud.co.ke/example_path",
		},
		{
			name: "Mailgun send email url path",
			args: args{
				service: service.Mailgun.Name,
				path:    service.Mailgun.Paths["sendEmail"],
			},
			want: "https://mailgun-testing.healthcloud.co.ke/communication/send_email",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := service
			if got := c.GenerateRequestURL(tt.args.service, tt.args.path); got != tt.want {
				t.Errorf("InterServiceClient.GenerateRequestURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterServiceClient_MakeRequest(t *testing.T) {
	service, _ := base.NewInterserviceClient("base")
	type args struct {
		method string
		url    string
		body   interface{}
	}

	message := base.MailgunEMailMessage{
		Subject: "Hello inter service email",
		Text:    "Test Email",
		To:      []string{"ngure.nyaga@healthcloud.co.ke"},
	}

	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "Example url bad request",
			args: args{
				method: http.MethodPost,
				url:    service.GenerateRequestURL("example", "example_path"),
				body: map[string]string{
					"example": "example_request",
				},
			},
			wantErr: true,
		},
		{
			name: "Example mailgun request",
			args: args{
				method: http.MethodPost,
				url: service.GenerateRequestURL(
					service.Mailgun.Name,
					service.Mailgun.Paths["sendEmail"],
				),
				body: message,
			},
			// TODO:Path not yet set up
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := service
			got, err := c.MakeRequest(tt.args.method, tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("InterServiceClient.MakeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterServiceClient.MakeRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasValidJWTBearerToken(t *testing.T) {
	service, _ := base.NewInterserviceClient("base")

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

func TestInterServiceAuthenticationMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	mw := base.InterServiceAuthenticationMiddleware()
	h := mw(next)
	rw := httptest.NewRecorder()
	reader := bytes.NewBuffer([]byte("sample"))

	service, _ := base.NewInterserviceClient("base")
	token, _ := service.CreateAuthToken()
	authHeader := fmt.Sprintf("Bearer %s", token)
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Add("Authorization", authHeader)
	h.ServeHTTP(rw, req)

	rw1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodPost, "/", reader)
	h.ServeHTTP(rw1, req1)
}

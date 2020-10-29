package base_test

import (
	"net/http"
	"os"
	"reflect"
	"testing"

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

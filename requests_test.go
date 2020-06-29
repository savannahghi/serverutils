package base

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRequestToTarget(t *testing.T) {
	clientID := MustGetEnvVar("CLIENT_ID")
	clientSecret := MustGetEnvVar("CLIENT_SECRET")
	username := MustGetEnvVar("USERNAME")
	password := MustGetEnvVar("PASSWORD")
	grantType := MustGetEnvVar("GRANT_TYPE")
	apiScheme := MustGetEnvVar("API_SCHEME")
	apiTokenURL := MustGetEnvVar("TOKEN_URL")
	apiHost := MustGetEnvVar("HOST")
	customHeader := MustGetEnvVar("DEFAULT_WORKSTATION_ID")
	extraHeaders := map[string]string{
		"X-WORKSTATION": customHeader,
	}
	target := map[string]interface{}{}

	client, err := NewServerClient(
		clientID,
		clientSecret,
		apiTokenURL,
		apiHost,
		apiScheme,
		grantType,
		username,
		password,
		extraHeaders,
	)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	type args struct {
		apiClient Client
		method    string
		path      string
		query     string
		content   []byte
		target    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case - nil content",
			args: args{
				apiClient: client,
				method:    http.MethodGet,
				path:      "/api/common/bp_registry/",
				content:   nil,
				target:    &target,
			},
			wantErr: false,
		},
		{
			name: "good case - non nil content",
			args: args{
				apiClient: client,
				method:    http.MethodPost,
				path:      "/api/common/bp_registry/",
				content:   []byte("some content"),
				target:    &target,
			},
			wantErr: false,
		},
		{
			name: "invalid URL",
			args: args{
				apiClient: client,
				method:    http.MethodPost,
				path:      "/this/is/not/a/real/path",
				content:   []byte("some content"),
				target:    &target,
			},
			wantErr: true,
		},
		{
			name: "bad target that will not decode to JSON",
			args: args{
				apiClient: client,
				method:    http.MethodGet,
				path:      "/api/common/bp_registry/",
				content:   nil,
				target:    target, // not a pointer, will trigger an error
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadRequestToTarget(tt.args.apiClient, tt.args.method, tt.args.path, tt.args.query, tt.args.content, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("ReadRequestToTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadAuthServerRequestToTarget(t *testing.T) {
	clientID := MustGetEnvVar("CLIENT_ID")
	clientSecret := MustGetEnvVar("CLIENT_SECRET")
	username := MustGetEnvVar("USERNAME")
	password := MustGetEnvVar("PASSWORD")
	grantType := MustGetEnvVar("GRANT_TYPE")
	apiScheme := MustGetEnvVar("API_SCHEME")
	apiTokenURL := MustGetEnvVar("TOKEN_URL")
	apiHost := MustGetEnvVar("HOST")
	customHeader := MustGetEnvVar("DEFAULT_WORKSTATION_ID")
	extraHeaders := map[string]string{
		"X-WORKSTATION": customHeader,
	}
	target := map[string]interface{}{}

	client, err := NewServerClient(
		clientID,
		clientSecret,
		apiTokenURL,
		apiHost,
		apiScheme,
		grantType,
		username,
		password,
		extraHeaders,
	)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	type args struct {
		client  Client
		method  string
		url     string
		s       string
		content []byte
		target  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid URL",
			args: args{
				client:  client,
				method:  http.MethodPost,
				url:     "/this/is/not/a/real/path",
				content: []byte("some content"),
				target:  &target,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadAuthServerRequestToTarget(tt.args.client, tt.args.method, tt.args.url, tt.args.s, tt.args.content, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("ReadAuthServerRequestToTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

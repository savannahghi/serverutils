package base_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestReadRequestToTarget(t *testing.T) {
	clientID := base.MustGetEnvVar(base.ClientIDEnvVarName)
	clientSecret := base.MustGetEnvVar(base.ClientSecretEnvVarName)
	username := base.MustGetEnvVar(base.UsernameEnvVarName)
	password := base.MustGetEnvVar(base.PasswordEnvVarName)
	grantType := base.MustGetEnvVar(base.GrantTypeEnvVarName)
	apiScheme := base.MustGetEnvVar(base.APISchemeEnvVarName)
	apiTokenURL := base.MustGetEnvVar(base.TokenURLEnvVarName)
	apiHost := base.MustGetEnvVar(base.APIHostEnvVarName)
	customHeader := base.MustGetEnvVar(base.WorkstationEnvVarName)
	extraHeaders := map[string]string{
		base.WorkstationHeaderName: customHeader,
	}
	target := map[string]interface{}{}

	client, err := base.NewServerClient(
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
		apiClient base.Client
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
			if err := base.ReadRequestToTarget(tt.args.apiClient, tt.args.method, tt.args.path, tt.args.query, tt.args.content, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("ReadRequestToTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadAuthServerRequestToTarget(t *testing.T) {
	clientID := base.MustGetEnvVar(base.ClientIDEnvVarName)
	clientSecret := base.MustGetEnvVar(base.ClientSecretEnvVarName)
	username := base.MustGetEnvVar(base.UsernameEnvVarName)
	password := base.MustGetEnvVar(base.PasswordEnvVarName)
	grantType := base.MustGetEnvVar(base.GrantTypeEnvVarName)
	apiScheme := base.MustGetEnvVar(base.APISchemeEnvVarName)
	apiTokenURL := base.MustGetEnvVar(base.TokenURLEnvVarName)
	apiHost := base.MustGetEnvVar(base.APIHostEnvVarName)
	workstationID := base.MustGetEnvVar(base.WorkstationEnvVarName)
	extraHeaders := map[string]string{
		base.WorkstationHeaderName: workstationID,
	}
	target := map[string]interface{}{}

	client, err := base.NewServerClient(
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
		client  base.Client
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
			if err := base.ReadAuthServerRequestToTarget(tt.args.client, tt.args.method, tt.args.url, tt.args.s, tt.args.content, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("ReadAuthServerRequestToTarget() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadWriteRequestToTarget(t *testing.T) {
	clientID := base.MustGetEnvVar(base.ClientIDEnvVarName)
	clientSecret := base.MustGetEnvVar(base.ClientSecretEnvVarName)
	username := base.MustGetEnvVar(base.UsernameEnvVarName)
	password := base.MustGetEnvVar(base.PasswordEnvVarName)
	grantType := base.MustGetEnvVar(base.GrantTypeEnvVarName)
	apiScheme := base.MustGetEnvVar(base.APISchemeEnvVarName)
	apiTokenURL := base.MustGetEnvVar(base.TokenURLEnvVarName)
	apiHost := base.MustGetEnvVar(base.APIHostEnvVarName)
	customHeader := base.MustGetEnvVar(base.WorkstationEnvVarName)
	extraHeaders := map[string]string{
		base.WorkstationHeaderName: customHeader,
	}
	target := map[string]interface{}{}

	client, err := base.NewServerClient(
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
		apiClient base.Client
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
				content:   nil,
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
			_, err := base.ReadWriteRequestToTarget(tt.args.apiClient, tt.args.method, tt.args.path, tt.args.query, tt.args.content, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadWriteRequestToTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err == nil {
				t.Errorf("expected an error to occur")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("did not expected an error to occur")
				return
			}
		})
	}
}

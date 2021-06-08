package base_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestMergeURLValues(t *testing.T) {
	val1 := url.Values{}
	val1.Add("a", "1")

	val2 := url.Values{}
	val2.Add("b", "2")

	expected := url.Values{}
	expected.Add("a", "1")
	expected.Add("b", "2")

	type args struct {
		values []url.Values
	}
	tests := []struct {
		name string
		args args
		want url.Values
	}{
		{
			name: "empty values",
			args: args{
				values: nil,
			},
			want: url.Values{},
		},
		{
			name: "non empty values",
			args: args{
				values: []url.Values{val1, val2},
			},
			want: expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.MergeURLValues(tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeURLValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAPIPaginationParams(t *testing.T) {
	lastPage := url.Values{}
	lastPage.Add("page_size", "10")
	lastPage.Add("page", "1")

	firstPage := url.Values{}
	firstPage.Add("page_size", "10")
	firstPage.Add("page", "1")

	beforePage := url.Values{}
	beforePage.Add("page_size", "100")
	beforePage.Add("page", "1")

	afterPage := url.Values{}
	afterPage.Add("page_size", "100")
	afterPage.Add("page", "1")

	type args struct {
		pagination *base.PaginationInput
	}
	tests := []struct {
		name    string
		args    args
		want    url.Values
		wantErr bool
	}{
		{
			name: "nil pagination params",
			args: args{
				pagination: nil,
			},
			want:    url.Values{},
			wantErr: false,
		},
		{
			name: "pagination last set",
			args: args{
				pagination: &base.PaginationInput{
					Last: 10,
				},
			},
			want:    lastPage,
			wantErr: false,
		},
		{
			name: "pagination first set",
			args: args{
				pagination: &base.PaginationInput{
					First: 10,
				},
			},
			want:    firstPage,
			wantErr: false,
		},
		{
			name: "pagination before set",
			args: args{
				pagination: &base.PaginationInput{
					Before: "12",
				},
			},
			want:    beforePage,
			wantErr: false,
		},
		{
			name: "pagination after set",
			args: args{
				pagination: &base.PaginationInput{
					After: "12",
				},
			},
			want:    afterPage,
			wantErr: false,
		},
		{
			name: "pagination - wrong after format",
			args: args{
				pagination: &base.PaginationInput{
					After: "this is not an int",
				},
			},
			want:    url.Values{},
			wantErr: true,
		},
		{
			name: "pagination - wrong before format",
			args: args{
				pagination: &base.PaginationInput{
					Before: "this is not an int",
				},
			},
			want:    url.Values{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetAPIPaginationParams(tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIPaginationParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAPIPaginationParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServerClient(t *testing.T) {
	// see the README for more guidance on these env vars
	clientID := base.MustGetEnvVar(base.ClientIDEnvVarName)
	clientSecret := base.MustGetEnvVar(base.ClientSecretEnvVarName)
	username := base.MustGetEnvVar(base.UsernameEnvVarName)
	password := base.MustGetEnvVar(base.PasswordEnvVarName)
	grantType := base.MustGetEnvVar(base.GrantTypeEnvVarName)
	apiScheme := base.MustGetEnvVar(base.APISchemeEnvVarName)
	apiTokenURL := base.MustGetEnvVar(base.TokenURLEnvVarName)
	apiHost := base.MustGetEnvVar(base.APIHostEnvVarName)
	workstationID := base.MustGetEnvVar(base.WorkstationEnvVarName)

	if base.IsDebug() {
		log.Printf(
			"Test Client Creds:\nclientID: %s\nclientSecret: %s\nusername: %s\npassword: %s\ngrantType: %s\napiScheme: %s\napiTokenURL: %s\napiHost: %s\nworkstationID: %s\n",
			clientID,
			clientSecret,
			username,
			password,
			grantType,
			apiScheme,
			apiTokenURL,
			apiHost,
			workstationID,
		)
	}

	type args struct {
		clientID     string
		clientSecret string
		apiTokenURL  string
		apiHost      string
		apiScheme    string
		grantType    string
		username     string
		password     string
		extraHeaders map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "invalid credentials (missing)",
			wantErr: true,
		},
		{
			name: "valid credentials, NO custom header",
			args: args{
				clientID:     clientID,
				clientSecret: clientSecret,
				apiTokenURL:  apiTokenURL,
				apiHost:      apiHost,
				apiScheme:    apiScheme,
				grantType:    grantType,
				username:     username,
				password:     password,
				extraHeaders: nil,
			},
			wantErr: false,
		},
		{
			name: "valid credentials, WITH custom header",
			args: args{
				clientID:     clientID,
				clientSecret: clientSecret,
				apiTokenURL:  apiTokenURL,
				apiHost:      apiHost,
				apiScheme:    apiScheme,
				grantType:    grantType,
				username:     username,
				password:     password,
				extraHeaders: map[string]string{
					base.WorkstationHeaderName: workstationID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.NewServerClient(tt.args.clientID, tt.args.clientSecret, tt.args.apiTokenURL, tt.args.apiHost, tt.args.apiScheme, tt.args.grantType, tt.args.username, tt.args.password, tt.args.extraHeaders)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				assert.True(t, got.IsInitialized())

				if tt.args.extraHeaders != nil {
					url := fmt.Sprintf("%s://%s/api/branches/workstationusers/?format=json", apiScheme, apiHost)
					resp, err := got.MakeRequest("GET", url, nil)
					assert.Nil(t, err)
					assert.NotNil(t, resp)
				}
			}
		})
	}
}

func Test_boolEnv(t *testing.T) {
	type args struct {
		envVarName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "env var that exists",
			args: args{
				envVarName: "DEBUG",
			},
			want: true,
		},
		{
			name: "env var that does not exist",
			args: args{
				envVarName: "this is not a real env var name",
			},
			want: false,
		},
		{
			name: "env var in the wrong format",
			args: args{
				envVarName: "GRANT_TYPE",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.BoolEnv(tt.args.envVarName); got != tt.want {
				t.Errorf("boolEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRunningTests(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "default case",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.IsRunningTests(); got != tt.want {
				t.Errorf("IsRunningTests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvVar(t *testing.T) {
	type args struct {
		envVarName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "existing var",
			args: args{
				envVarName: "GRANT_TYPE",
			},
			want:    "password",
			wantErr: false,
		},
		{
			name: "non existent var",
			args: args{
				envVarName: "this is not a valid env var name",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetEnvVar(tt.args.envVarName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnvVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustGetEnvVar(t *testing.T) {
	type args struct {
		envVarName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "existing var",
			args: args{
				envVarName: "GRANT_TYPE",
			},
			want: "password",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.MustGetEnvVar(tt.args.envVarName); got != tt.want {
				t.Errorf("MustGetEnvVar() = %v, want %v", got, tt.want)
			}

			assert.Panics(t, func() {
				base.MustGetEnvVar("this does not exist as an env var")
			})
		})
	}
}

func TestCloseRespBody(t *testing.T) {
	// no-op for nil response
	base.CloseRespBody(nil)

	// closes the resp body when supported
	resp := &http.Response{Body: ioutil.NopCloser(strings.NewReader("some contents"))}
	base.CloseRespBody(resp)

	// does not blow up on unclosable bodies
	req := httptest.NewRequest("POST", "/", strings.NewReader("some contents"))
	blowResp := &http.Response{Body: base.BlowUpOnClose{}, Request: req}
	base.CloseRespBody(blowResp)
}

func TestCheckAPIInitialization(t *testing.T) {
	msg := "the API httpClient is not correctly initialized. Please use the `.Initialize` constructor"
	nilErr := base.CheckAPIInitialization(nil)
	assert.NotNil(t, nilErr)
	assert.Equal(t, msg, nilErr.Error())

	badCl := &base.MockClient{Initialized: false}
	initErr := base.CheckAPIInitialization(badCl)
	assert.NotNil(t, initErr)
	assert.Equal(t, msg, initErr.Error())

	goodCl := &base.MockClient{Initialized: true}
	noErr := base.CheckAPIInitialization(goodCl)
	assert.Nil(t, noErr)
}

func Test_decodeOauthResponseFromJSON_Unreadable_Body(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader("some more contents"))
	blowResp := &http.Response{Body: base.BlowUpOnRead{}, Request: req}
	oauth, err := base.DecodeOAUTHResponseFromJSON(blowResp)
	assert.NotNil(t, err)
	assert.Nil(t, oauth)
	assert.Equal(t, "boom", err.Error())
}

func Test_decodeOauthResponseFromJSON_Unmarshallable_Body(t *testing.T) {
	reader := ioutil.NopCloser(strings.NewReader("very bad not valid JSON"))
	resp := &http.Response{Body: reader}
	oauth, err := base.DecodeOAUTHResponseFromJSON(resp)
	assert.NotNil(t, err)
	assert.Nil(t, oauth)
	assert.Equal(t, "invalid character 'v' looking for beginning of value", err.Error())
}

func Test_decodeOauthResponseFromJSON_Unmarshallable_Success(t *testing.T) {
	jsonResp := `
		{
			"access_token": "GJJGFDGJJGFGJHHJF",
			"scope": "this.is.some.dummy.scope",
			"token_type": "Bearer",
			"expires_in": 3600,
			"refresh_token": "YHGFDSETGJKHFDD"
		}`
	reader := ioutil.NopCloser(strings.NewReader(jsonResp))
	resp := &http.Response{Body: reader}
	oauth, err := base.DecodeOAUTHResponseFromJSON(resp)
	assert.NotNil(t, oauth)
	assert.Nil(t, err)
	assert.Equal(t, oauth.AccessToken, "GJJGFDGJJGFGJHHJF")
	assert.Equal(t, oauth.Scope, "this.is.some.dummy.scope")
	assert.Equal(t, oauth.TokenType, "Bearer")
	assert.Equal(t, oauth.ExpiresIn, 3600)
	assert.Equal(t, oauth.RefreshToken, "YHGFDSETGJKHFDD")

	// update auth
	c := &base.ServerClient{}
	c.UpdateAuth(oauth)
	assert.Equal(t, oauth.AccessToken, c.AccessToken())
	assert.Equal(t, oauth.TokenType, c.TokenType())
	assert.Equal(t, oauth.Scope, c.AccessScope())
	assert.Equal(t, oauth.RefreshToken, c.RefreshToken())
	assert.Equal(t, oauth.ExpiresIn, c.ExpiresIn())

	// wait out most of the token's duration to expiry before attempting to Refresh
	secondsToRefresh := int(float64(c.ExpiresIn()) * base.TokenExpiryRatio)
	refreshTime := time.Now().Add(time.Second * time.Duration(secondsToRefresh))
	assert.InDelta(t, refreshTime.UnixNano(), c.RefreshAt().UnixNano(), 1_000_000_000, "allow up to one second of delta")
	assert.Equal(t, true, c.IsInitialized())
}

func TestComposeAPIURL(t *testing.T) {
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

	if client == nil {
		return
	}

	type args struct {
		client base.Client
		path   string
		query  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "good case",
			args: args{
				client: client,
				path:   "/api/branches/workstations",
				query:  "format=json",
			},
			want: "https://erp-api-testing.healthcloud.co.ke/api/branches/workstations?format=json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.ComposeAPIURL(tt.args.client, tt.args.path, tt.args.query); got != tt.want {
				t.Errorf("ComposeAPIURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccessToken(t *testing.T) {
	// see the README for more guidance on these env vars
	clientID := base.MustGetEnvVar(base.ClientIDEnvVarName)
	clientSecret := base.MustGetEnvVar(base.ClientSecretEnvVarName)
	username := base.MustGetEnvVar(base.UsernameEnvVarName)
	password := base.MustGetEnvVar(base.PasswordEnvVarName)
	grantType := base.MustGetEnvVar(base.GrantTypeEnvVarName)
	apiScheme := base.MustGetEnvVar(base.APISchemeEnvVarName)
	apiTokenURL := base.MustGetEnvVar(base.TokenURLEnvVarName)
	apiHost := base.MustGetEnvVar(base.APIHostEnvVarName)

	tests := []struct {
		name    string
		args    *base.ClientServerOptions
		wantErr bool
	}{
		{
			name: "valid credentials",
			args: &base.ClientServerOptions{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				APITokenURL:  apiTokenURL,
				APIHost:      apiHost,
				APIScheme:    apiScheme,
				GrantType:    grantType,
				Username:     username,
				Password:     password,
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			args: &base.ClientServerOptions{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				APITokenURL:  apiTokenURL,
				APIHost:      apiHost,
				APIScheme:    apiScheme,
				GrantType:    grantType,
				Username:     "username",
				Password:     password,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := base.GetAccessToken(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, token)
			}
		})
	}
}

func TestNewPostRequest(t *testing.T) {
	type args struct {
		url             string
		values          url.Values
		headers         map[string]string
		timeoutDuration int
	}
	tests := []struct {
		name string
		args args
		want *http.Response
	}{
		{
			name: "Test post request returns an error",
			args: args{
				url:             "/some url",
				values:          url.Values{"data": []string{"dummy data"}},
				headers:         map[string]string{"Content-Type": "application/json"},
				timeoutDuration: 200,
			},
		},
		{
			name: "invalid URL",
			args: args{
				url:             "this is not a real URL",
				values:          url.Values{"data": []string{"dummy data"}},
				headers:         map[string]string{"Content-Type": "application/json"},
				timeoutDuration: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := base.NewPostRequest(tt.args.url, tt.args.values, tt.args.headers, tt.args.timeoutDuration)
			assert.Nil(t, payload)
			assert.NotNil(t, err)
		})
	}
}

func TestServerClient_MeURL(t *testing.T) {
	// badly initialized client
	badClient, err := base.NewServerClient(
		base.MustGetEnvVar(base.ClientIDEnvVarName),
		base.MustGetEnvVar(base.ClientSecretEnvVarName),
		"this is not a valid url",
		base.MustGetEnvVar(base.APIHostEnvVarName),
		base.MustGetEnvVar(base.APISchemeEnvVarName),
		base.MustGetEnvVar(base.GrantTypeEnvVarName),
		base.MustGetEnvVar(base.UsernameEnvVarName),
		base.MustGetEnvVar(base.PasswordEnvVarName),
		map[string]string{
			base.WorkstationHeaderName: base.MustGetEnvVar(base.WorkstationEnvVarName),
		},
	)
	assert.NotNil(t, err)
	assert.Nil(t, badClient)

	// properly initialized client
	goodClient, err := base.NewServerClient(
		base.MustGetEnvVar(base.ClientIDEnvVarName),
		base.MustGetEnvVar(base.ClientSecretEnvVarName),
		base.MustGetEnvVar(base.TokenURLEnvVarName),
		base.MustGetEnvVar(base.APIHostEnvVarName),
		base.MustGetEnvVar(base.APISchemeEnvVarName),
		base.MustGetEnvVar(base.GrantTypeEnvVarName),
		base.MustGetEnvVar(base.UsernameEnvVarName),
		base.MustGetEnvVar(base.PasswordEnvVarName),
		map[string]string{
			base.WorkstationHeaderName: base.MustGetEnvVar(base.WorkstationEnvVarName),
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, goodClient)
	if goodClient != nil {
		got, err := goodClient.MeURL()
		assert.Nil(t, err)
		assert.Equal(t, "https://auth.healthcloud.co.ke/v1/user/me/?format=json", got)
	}
}

func TestDefaultServerClient(t *testing.T) {
	client, err := base.DefaultServerClient()
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestCheckAPIClientPostConditions(t *testing.T) {
	goodClient, err := base.DefaultServerClient()
	assert.Nil(t, err)
	assert.NotNil(t, goodClient)

	type args struct {
		client base.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid client",
			args: args{
				client: goodClient,
			},
			wantErr: false,
		},
		{
			name: "invalid client - no access token",
			args: args{
				client: &base.ServerClient{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := base.CheckAPIClientPostConditions(tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("CheckAPIClientPostConditions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerClient_HTTPClient(t *testing.T) {
	c, err := base.DefaultServerClient()
	assert.Nil(t, err)
	assert.NotNil(t, c)

	assert.NotNil(t, c.HTTPClient())
}

func TestServerClient_Refresh(t *testing.T) {
	c, err := base.DefaultServerClient()
	assert.Nil(t, err)
	assert.NotNil(t, c)

	err = c.Refresh()
	assert.Nil(t, err)
}

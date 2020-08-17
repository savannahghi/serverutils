package base

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
			if got := MergeURLValues(tt.args.values...); !reflect.DeepEqual(got, tt.want) {
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
		pagination *PaginationInput
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
				pagination: &PaginationInput{
					Last: 10,
				},
			},
			want:    lastPage,
			wantErr: false,
		},
		{
			name: "pagination first set",
			args: args{
				pagination: &PaginationInput{
					First: 10,
				},
			},
			want:    firstPage,
			wantErr: false,
		},
		{
			name: "pagination before set",
			args: args{
				pagination: &PaginationInput{
					Before: "12",
				},
			},
			want:    beforePage,
			wantErr: false,
		},
		{
			name: "pagination after set",
			args: args{
				pagination: &PaginationInput{
					After: "12",
				},
			},
			want:    afterPage,
			wantErr: false,
		},
		{
			name: "pagination - wrong after format",
			args: args{
				pagination: &PaginationInput{
					After: "this is not an int",
				},
			},
			want:    url.Values{},
			wantErr: true,
		},
		{
			name: "pagination - wrong before format",
			args: args{
				pagination: &PaginationInput{
					Before: "this is not an int",
				},
			},
			want:    url.Values{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAPIPaginationParams(tt.args.pagination)
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

func TestServerClient_IsInitialized(t *testing.T) {
	type fields struct {
		isInitialized bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "initialized client",
			fields: fields{
				isInitialized: true,
			},
			want: true,
		},
		{
			name: "uninitialized client",
			fields: fields{
				isInitialized: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ServerClient{
				isInitialized: tt.fields.isInitialized,
			}
			if got := c.IsInitialized(); got != tt.want {
				t.Errorf("ServerClient.IsInitialized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerClient_AuthServerDomain(t *testing.T) {
	type fields struct {
		authServerDomain string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "auth server domain set",
			fields: fields{
				authServerDomain: "https://auth.healthcloud.co.ke",
			},
			want: "https://auth.healthcloud.co.ke",
		},
		{
			name: "auth server domain not set",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ServerClient{
				authServerDomain: tt.fields.authServerDomain,
			}
			if got := c.AuthServerDomain(); got != tt.want {
				t.Errorf("ServerClient.IsInitialized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServerClient(t *testing.T) {
	// see the README for more guidance on these env vars
	clientID := MustGetEnvVar("CLIENT_ID")
	clientSecret := MustGetEnvVar("CLIENT_SECRET")
	username := MustGetEnvVar("USERNAME")
	password := MustGetEnvVar("PASSWORD")
	grantType := MustGetEnvVar("GRANT_TYPE")
	apiScheme := MustGetEnvVar("API_SCHEME")
	apiTokenURL := MustGetEnvVar("TOKEN_URL")
	apiHost := MustGetEnvVar("HOST")
	customHeader := MustGetEnvVar("DEFAULT_WORKSTATION_ID")

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
					"X-WORKSTATION": customHeader,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServerClient(tt.args.clientID, tt.args.clientSecret, tt.args.apiTokenURL, tt.args.apiHost, tt.args.apiScheme, tt.args.grantType, tt.args.username, tt.args.password, tt.args.extraHeaders)
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
			if got := boolEnv(tt.args.envVarName); got != tt.want {
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
			if got := IsRunningTests(); got != tt.want {
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
			got, err := GetEnvVar(tt.args.envVarName)
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
			if got := MustGetEnvVar(tt.args.envVarName); got != tt.want {
				t.Errorf("MustGetEnvVar() = %v, want %v", got, tt.want)
			}

			assert.Panics(t, func() {
				MustGetEnvVar("this does not exist as an env var")
			})
		})
	}
}

func TestNewEDIClient_Invalid_Domain(t *testing.T) {
	c := ServerClient{
		clientID:         "",
		clientSecret:     "",
		apiTokenURL:      "",
		authServerDomain: "",
		apiHost:          "",
		apiScheme:        "",
		grantType:        "",
		username:         "",
		password:         "",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, " is not a valid clientId, expected a non-blank alphanumeric string of at least 12 characters", err.Error())
}

func TestNewEDIClient_Invalid_ClientID(t *testing.T) {
	c := ServerClient{
		clientID:         "",
		clientSecret:     "",
		apiTokenURL:      "",
		authServerDomain: "",
		apiHost:          "",
		apiScheme:        "",
		grantType:        "",
		username:         "",
		password:         "",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, " is not a valid clientId, expected a non-blank alphanumeric string of at least 12 characters", err.Error())
}

func TestNewEDIClient_Invalid_ClientSecret(t *testing.T) {
	c := ServerClient{
		clientID:         "SDFGHJKLOIUYG",
		clientSecret:     "",
		apiTokenURL:      "",
		authServerDomain: "",
		apiHost:          "",
		apiScheme:        "",
		grantType:        "",
		username:         "",
		password:         "",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, " is not a valid clientSecret, expected a non-blank alphanumeric string of at least 12 characters", err.Error())
}

func TestNewEDIClient_Invalid_APITokenURL(t *testing.T) {
	c := ServerClient{
		clientID:         "SDFGHJKLOIUYG",
		clientSecret:     "RTFFfghjfkjkhjkghkgfhadciuy",
		apiTokenURL:      "",
		authServerDomain: "",
		apiHost:          "",
		apiScheme:        "",
		grantType:        "",
		username:         "",
		password:         "",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, " is not a valid apiTokenURL, expected an http(s) URL", err.Error())
}

func TestNewEDIClient_Invalid_APIHost(t *testing.T) {
	c := ServerClient{
		clientID:         "SDFGHJKLOIUYG",
		clientSecret:     "RTFFfghjfkjkhjkghkgfhadciuy",
		apiTokenURL:      "http://localhost:9000/o/token/",
		authServerDomain: "localhost",
		apiHost:          "",
		apiScheme:        "",
		grantType:        "",
		username:         "",
		password:         "",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, " is not a valid apiHost, expected a valid IP or domain name", err.Error())
}

func TestNewEDIClient_Invalid_APIScheme(t *testing.T) {
	c := ServerClient{
		clientID:         "SDFGHJKLOIUYG",
		clientSecret:     "RTFFfghjfkjkhjkghkgfhadciuy",
		apiTokenURL:      "http://localhost:9000/o/token/",
		authServerDomain: "localhost",
		apiHost:          "localhost",
		apiScheme:        "ftp",
		grantType:        "",
		username:         "",
		password:         "",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, "ftp is not a valid apiScheme, expected http or https", err.Error())
}

func TestNewEDIClient_Invalid_GrantType(t *testing.T) {
	c := ServerClient{
		clientID:         "SDFGHJKLOIUYG",
		clientSecret:     "RTFFfghjfkjkhjkghkgfhadciuy",
		apiTokenURL:      "http://localhost:9000/o/token/",
		authServerDomain: "localhost",
		apiHost:          "localhost",
		apiScheme:        "http",
		grantType:        "token_refresh",
		username:         "",
		password:         "",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, "the only supported OAuth grant type for now is 'password'", err.Error())
}

func TestNewEDIClient_Invalid_Username(t *testing.T) {
	c := ServerClient{
		clientID:         "SDFGHJKLOIUYG",
		clientSecret:     "RTFFfghjfkjkhjkghkgfhadciuy",
		apiTokenURL:      "http://localhost:9000/o/token/",
		authServerDomain: "localhost",
		apiHost:          "localhost",
		apiScheme:        "http",
		grantType:        "password",
		username:         "not_an_email",
		password:         "",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, "the Username should be a valid email address", err.Error())
}

func TestNewEDIClient_Invalid_Password(t *testing.T) {
	c := ServerClient{
		clientID:         "SDFGHJKLOIUYG",
		clientSecret:     "RTFFfghjfkjkhjkghkgfhadciuy",
		apiTokenURL:      "http://localhost:9000/o/token/",
		authServerDomain: "localhost",
		apiHost:          "localhost",
		apiScheme:        "http",
		grantType:        "password",
		username:         "user@mail.com",
		password:         "x",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("the Password should be a string of at least %d characters", apiPasswordMinLength), err.Error())
}

func TestCloseRespBody(t *testing.T) {
	// no-op for nil response
	CloseRespBody(nil)

	// closes the resp body when supported
	resp := &http.Response{Body: ioutil.NopCloser(strings.NewReader("some contents"))}
	CloseRespBody(resp)

	// does not blow up on unclosable bodies
	req := httptest.NewRequest("POST", "/", strings.NewReader("some contents"))
	blowResp := &http.Response{Body: BlowUpOnClose{}, Request: req}
	CloseRespBody(blowResp)
}

func TestCheckEDIAPIInitialization(t *testing.T) {
	msg := "the EDI httpClient is not correctly initialized. Please use the `.Initialize` constructor"
	nilErr := CheckAPIInitialization(nil)
	assert.NotNil(t, nilErr)
	assert.Equal(t, msg, nilErr.Error())

	badCl := &MockClient{Initialized: false}
	initErr := CheckAPIInitialization(badCl)
	assert.NotNil(t, initErr)
	assert.Equal(t, msg, initErr.Error())

	goodCl := &MockClient{Initialized: true}
	noErr := CheckAPIInitialization(goodCl)
	assert.Nil(t, noErr)
}

func TestCheckEDIClientPostConditions(t *testing.T) {
	tests := map[string]struct {
		input Client
		want  error
	}{
		"invalid_access_tokens": {
			input: &ServerClient{
				accessToken: "invalid",
			},
			want: fmt.Errorf("invalid access token after EDIAPIClient initialization"),
		},
		"invalid_token_type": {
			input: &ServerClient{
				accessToken: "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:   "Bogus",
			},
			want: fmt.Errorf("invalid token type after EDIAPIClient initialization, expected 'Bearer'"),
		},
		"invalid_refresh_token": {
			input: &ServerClient{
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "bad",
			},
			want: fmt.Errorf("invalid Refresh token after EDIAPIClient initialization"),
		},
		"invalid_access_scope": {
			input: &ServerClient{
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "jfdahfdjafhdjfhdalkfjdhkfhasdk",
				accessScope:  "bad",
			},
			want: fmt.Errorf("invalid access scope text after EDIAPIClient initialization"),
		},
		"invalid_expires_in": {
			input: &ServerClient{
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "jfdahfdjafhdjfhdalkfjdhkfhasdk",
				accessScope:  "scope blah blah blah more blah blah blah",
				expiresIn:    -1,
			},
			want: fmt.Errorf("invalid expiresIn after EDIAPIClient initialization"),
		},
		"invalid_refresh_at": {
			input: &ServerClient{
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "jfdahfdjafhdjfhdalkfjdhkfhasdk",
				accessScope:  "scope blah blah blah more blah blah blah",
				expiresIn:    3600,
				refreshAt:    time.Unix(0, 0),
			},
			want: fmt.Errorf("invalid past refreshAt after EDIAPIClient initialization"),
		},
		"passing_case": {
			input: &ServerClient{
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "jfdahfdjafhdjfhdalkfjdhkfhasdk",
				accessScope:  "scope blah blah blah more blah blah blah",
				expiresIn:    3600,
				refreshAt:    time.Now().Add(time.Second * 3600),
			},
			want: nil,
		},
		"invalid_credentials": {
			input: &MockClient{
				authErr:      fmt.Errorf("mock auth error"),
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "jfdahfdjafhdjfhdalkfjdhkfhasdk",
				accessScope:  "scope blah blah blah more blah blah blah",
				expiresIn:    3600,
				refreshAt:    time.Now().Add(time.Second * 3600),
			},
			want: fmt.Errorf("mock auth error"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := CheckAPIClientPostConditions(tc.input)

			// if the validator passes, err will be nil
			if err != nil {
				assert.Equalf(t, tc.want.Error(), err.Error(), "expected error '%s' not found")
			}
		})
	}
}

func Test_decodeOauthResponseFromJSON_Unreadable_Body(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader("some more contents"))
	blowResp := &http.Response{Body: BlowUpOnRead{}, Request: req}
	oauth, err := decodeOauthResponseFromJSON(blowResp)
	assert.NotNil(t, err)
	assert.Nil(t, oauth)
	assert.Equal(t, "boom", err.Error())
}

func Test_decodeOauthResponseFromJSON_Unmarshallable_Body(t *testing.T) {
	reader := ioutil.NopCloser(strings.NewReader("very bad not valid JSON"))
	resp := &http.Response{Body: reader}
	oauth, err := decodeOauthResponseFromJSON(resp)
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
	oauth, err := decodeOauthResponseFromJSON(resp)
	assert.NotNil(t, oauth)
	assert.Nil(t, err)
	assert.Equal(t, oauth.AccessToken, "GJJGFDGJJGFGJHHJF")
	assert.Equal(t, oauth.Scope, "this.is.some.dummy.scope")
	assert.Equal(t, oauth.TokenType, "Bearer")
	assert.Equal(t, oauth.ExpiresIn, 3600)
	assert.Equal(t, oauth.RefreshToken, "YHGFDSETGJKHFDD")

	// update auth
	c := &ServerClient{}
	c.updateAuth(oauth)
	assert.Equal(t, oauth.AccessToken, c.accessToken)
	assert.Equal(t, oauth.TokenType, c.tokenType)
	assert.Equal(t, oauth.Scope, c.accessScope)
	assert.Equal(t, oauth.RefreshToken, c.refreshToken)
	assert.Equal(t, oauth.ExpiresIn, c.expiresIn)

	// wait out most of the token's duration to expiry before attempting to Refresh
	secondsToRefresh := int(float64(c.expiresIn) * tokenExpiryRatio)
	refreshTime := time.Now().Add(time.Second * time.Duration(secondsToRefresh))
	assert.InDelta(t, refreshTime.UnixNano(), c.refreshAt.UnixNano(), 1_000_000_000, "allow up to one second of delta")
	assert.Equal(t, true, c.isInitialized)
}

func TestEDIAPIClient_Authenticate_Non_20x_Response(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)
		assert.NotNil(t, body)
		w.WriteHeader(http.StatusInternalServerError) // 500 error response
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	cl := srv.Client()
	client := ServerClient{
		httpClient:       cl,
		clientID:         "cidhhhhhhhhhhh",
		clientSecret:     "csffffffffffff",
		grantType:        "password",
		username:         "yusa@yusa.io",
		password:         "the greatest Password in the world",
		isInitialized:    true,
		apiTokenURL:      srv.URL, // important, this invokes our test server!
		apiHost:          "localhost",
		authServerDomain: "localhost",
		apiScheme:        "http",
	}
	err := client.Authenticate()
	assert.NotNil(t, err)
	assert.Equal(t, "server error status: 500", err.Error())

	refreshErr := client.Refresh()
	assert.NotNil(t, refreshErr)
	assert.Equal(t, "server error status: 500", refreshErr.Error())
}

func TestEDIAPIClient_Authenticate_HTTP_Error(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)
		assert.NotNil(t, body)
		panic(nil) // abort the connection prematurely
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	cl := srv.Client()
	client := ServerClient{
		httpClient:       cl,
		clientID:         "cidxxxxxxxxxxx",
		clientSecret:     "csxxxxxxxxxxx",
		grantType:        "password",
		username:         "yusa@yusa.io",
		password:         "the greatest Password in the world",
		isInitialized:    true,
		apiTokenURL:      srv.URL, // important, this invokes our test server!
		apiHost:          "localhost",
		authServerDomain: "localhost",
		apiScheme:        "http",
	}
	err := client.Authenticate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), ": EOF")

	refreshErr := client.Refresh()
	assert.NotNil(t, refreshErr)
	assert.Contains(t, refreshErr.Error(), ": EOF")
}

func TestEDIAPIClient_Authenticate_Unmarshalable_Response(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)
		assert.NotNil(t, body)

		respBytes := []byte("some stuff")
		bytes, writeErr := w.Write(respBytes)
		assert.Nil(t, writeErr)
		assert.Equal(t, 10, bytes)
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	cl := srv.Client()
	client := ServerClient{
		httpClient:       cl,
		clientID:         "cidxxxxxxxxxx",
		clientSecret:     "csxxxxxxxxxxx",
		grantType:        "password",
		username:         "yusa@yusa.io",
		password:         "the greatest Password in the world",
		isInitialized:    true,
		apiTokenURL:      srv.URL, // important, this invokes our test server!
		apiHost:          "localhost",
		apiScheme:        "http",
		authServerDomain: "localhost",
	}
	err := client.Authenticate()
	assert.NotNil(t, err)
	assert.Equal(t, "invalid character 's' looking for beginning of value", err.Error())

	refreshErr := client.Refresh()
	assert.NotNil(t, refreshErr)
	assert.Equal(t, "invalid character 's' looking for beginning of value", refreshErr.Error())
}

func TestEDIAPIClient_Authenticate_Success(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)
		assert.NotNil(t, body)

		jsonResp := `
		{
			"access_token": "GJJGFDGJJGFGJHHJF",
			"scope": "this.is.some.dummy.scope",
			"token_type": "Bearer",
			"expires_in": 3600,
			"refresh_token": "YHGFDSETGJKHFDD"
		}`
		respBytes := []byte(jsonResp)
		bytes, writeErr := w.Write(respBytes)
		assert.Nil(t, writeErr)
		assert.NotNil(t, bytes)
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	cl := srv.Client()
	client := ServerClient{
		httpClient:       cl,
		clientID:         "cidxxxxxxxxxxx",
		clientSecret:     "csxxxxxxxxxxx",
		grantType:        "password",
		username:         "yusa@yusa.io",
		password:         "the greatest Password in the world",
		isInitialized:    true,
		apiTokenURL:      srv.URL, // important, this invokes our test server!
		apiHost:          "localhost",
		authServerDomain: "localhost",
		apiScheme:        "http",
	}
	err := client.Authenticate()
	assert.Nil(t, err) // successful authentication

	refreshErr := client.Refresh()
	assert.Nil(t, refreshErr)
}

func TestEDIAPIClient_Authenticate_Failing_Preconditions(t *testing.T) {
	client := ServerClient{
		clientID:         "cidxxxxxxxxxxx",
		clientSecret:     "csxxxxxxxxxxx",
		grantType:        "password",
		username:         "yusa@yusa.io",
		password:         "the greatest Password in the world",
		isInitialized:    true,
		apiHost:          "", // bad value, will trigger preconditions check
		authServerDomain: "localhost",
		apiScheme:        "http",
	}
	err := client.Authenticate()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "is not a valid apiTokenURL, expected an http(s) UR")
}

func TestEDIAPIClient_Initialize_Failing_Postconditions(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)
		assert.NotNil(t, body)

		// the access_token is too short, will trigger a post-condition fail
		jsonResp := `
		{
			"access_token": "x",
			"scope": "this.is.some.dummy.scope",
			"token_type": "Bearer",
			"expires_in": 3600,
			"refresh_token": "YHGFDSETGJKHFDD"
		}`
		respBytes := []byte(jsonResp)
		bytes, writeErr := w.Write(respBytes)
		assert.Nil(t, writeErr)
		assert.NotNil(t, bytes)
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	cl := srv.Client()
	client := ServerClient{
		httpClient:       cl,
		clientID:         "cidxxxxxxxxxxx",
		clientSecret:     "csxxxxxxxxxxx",
		grantType:        "password",
		username:         "yusa@yusa.io",
		password:         "the greatest Password in the world",
		isInitialized:    true,
		apiTokenURL:      srv.URL, // important, this invokes our test server!
		apiHost:          "localhost",
		authServerDomain: "localhost",
		apiScheme:        "http",
	}
	err := client.Initialize()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid access token after EDIAPIClient initialization")
}

func TestNewEDIClient_Auth_Fail(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)
		assert.NotNil(t, body)
		w.WriteHeader(http.StatusInternalServerError) // 500 error response
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	c := ServerClient{
		clientID:         "SDFGHJKLOIUYG",
		clientSecret:     "RTFFfghjfkjkhjkghkgfhadciuy",
		apiTokenURL:      srv.URL,
		authServerDomain: "localhost",
		apiHost:          "localhost",
		apiScheme:        "http",
		grantType:        "password",
		username:         "user@mail.com",
		password:         "legit Password",
	}
	err := c.Initialize()
	assert.NotNil(t, err)
	assert.Equal(t, "server error status: 500", err.Error())
}

func TestEDIAPIClient_MakeRequest_Happy_Case(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer JHJKLHLKJHJKHKJHKJHKJ", r.Header.Get("Authorization"))

		body, rErr := ioutil.ReadAll(r.Body)
		assert.Nil(t, rErr)
		assert.NotNil(t, body)
		w.Header().Set("Content-Type", "application/json")

		jsonOutput := `{"a": 1}`
		bytes, wErr := w.Write([]byte(jsonOutput))
		assert.Nil(t, wErr)
		assert.Equal(t, len(jsonOutput), bytes)
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	cl := srv.Client()
	now := time.Now()
	future := now.Add(time.Second * 3600)
	client := ServerClient{
		httpClient:    cl,
		clientID:      "cid",
		clientSecret:  "cs",
		grantType:     "password",
		username:      "yusa@yusa.io",
		password:      "the greatest Password in the world",
		isInitialized: true,
		accessToken:   "JHJKLHLKJHJKHKJHKJHKJ",
		refreshAt:     future,
	}
	// the request is sent to our test server at srv.URL
	resp, err := client.MakeRequest("GET", srv.URL, strings.NewReader("stuff"))
	log.Printf("%#v\n", err)
	assert.Nil(t, err)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	respBody, readErr := ioutil.ReadAll(resp.Body)
	assert.Nil(t, readErr)
	assert.Equal(t, `{"a": 1}`, string(respBody))
}

func TestEDIAPIClient_MakeRequest_Wrong_Content_Type(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer JHJKLHLKJHJKHKJHKJHKJ", r.Header.Get("Authorization"))

		body, rErr := ioutil.ReadAll(r.Body)
		assert.Nil(t, rErr)
		assert.NotNil(t, body)
		w.Header().Set("Content-Type", "text/plain")

		jsonOutput := `{"a": 1}`
		bytes, wErr := w.Write([]byte(jsonOutput))
		assert.Nil(t, wErr)
		assert.Equal(t, len(jsonOutput), bytes)
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	cl := srv.Client()
	now := time.Now()
	future := now.Add(time.Second * 3600)
	client := ServerClient{
		httpClient:    cl,
		clientID:      "cid",
		clientSecret:  "cs",
		grantType:     "password",
		username:      "yusa@yusa.io",
		password:      "the greatest Password in the world",
		isInitialized: true,
		accessToken:   "JHJKLHLKJHJKHKJHKJHKJ",
		refreshAt:     future,
	}
	// the request is sent to our test server at srv.URL
	resp, err := client.MakeRequest("GET", srv.URL, strings.NewReader("stuff"))
	log.Printf("%#v\n", err)
	assert.NotNil(t, err)
	assert.Equal(t, "expected application/json Content-Type, got text/plain", err.Error())
	assert.Nil(t, resp)
}

func TestEDIAPIClient_MakeRequest_Client_Error(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer JHJKLHLKJHJKHKJHKJHKJ", r.Header.Get("Authorization"))

		body, rErr := ioutil.ReadAll(r.Body)
		assert.Nil(t, rErr)
		assert.NotNil(t, body)
		w.Header().Set("Content-Type", "text/plain")

		panic(nil) // abort the request
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()

	cl := srv.Client()
	now := time.Now()
	future := now.Add(time.Second * 3600)
	client := ServerClient{
		httpClient:    cl,
		clientID:      "cid",
		clientSecret:  "cs",
		grantType:     "password",
		username:      "yusa@yusa.io",
		password:      "the greatest Password in the world",
		isInitialized: true,
		accessToken:   "JHJKLHLKJHJKHKJHKJHKJ",
		refreshAt:     future,
	}
	// the request is sent to our test server at srv.URL
	resp, err := client.MakeRequest("GET", srv.URL, strings.NewReader("stuff"))
	log.Printf("%#v\n", err)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), ": EOF")
	assert.Nil(t, resp)
}

func TestEDIAPIClient_MakeRequest_Refresh_Needed_Succeeds(t *testing.T) {
	ediHf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer GJJGFDGJJGFGJHHJF", r.Header.Get("Authorization"))

		body, rErr := ioutil.ReadAll(r.Body)
		assert.Nil(t, rErr)
		assert.NotNil(t, body)
		w.Header().Set("Content-Type", "application/json")

		jsonOutput := `{"a": 1}`
		bytes, wErr := w.Write([]byte(jsonOutput))
		assert.Nil(t, wErr)
		assert.Equal(t, len(jsonOutput), bytes)
	})
	ediSrv := httptest.NewServer(ediHf)
	defer ediSrv.Close()

	// authentication and refresh should succeed
	authHF := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		body, _ := ioutil.ReadAll(r.Body)
		assert.NotNil(t, body)

		jsonResp := `
		{
			"access_token": "GJJGFDGJJGFGJHHJF",
			"scope": "this.is.some.dummy.scope",
			"token_type": "Bearer",
			"expires_in": 3600,
			"refresh_token": "YHGFDSETGJKHFDD"
		}`
		respBytes := []byte(jsonResp)
		bytes, writeErr := w.Write(respBytes)
		assert.Nil(t, writeErr)
		assert.NotNil(t, bytes)
	})
	authSrv := httptest.NewServer(authHF)
	defer authSrv.Close()

	cl := ediSrv.Client()
	now := time.Now()
	past := now.Add(time.Second * -3600)
	assert.Less(t, past.UnixNano(), now.UnixNano())
	client := ServerClient{
		httpClient:    cl,
		clientID:      "cid",
		clientSecret:  "cs",
		grantType:     "password",
		username:      "yusa@yusa.io",
		password:      "the greatest Password in the world",
		isInitialized: true,
		accessToken:   "GJJGFDGJJGFGJHHJF",
		refreshAt:     past,
		apiTokenURL:   authSrv.URL,
	}
	// the request is sent to our test server at ediSrv.URL
	resp, err := client.MakeRequest("GET", ediSrv.URL, strings.NewReader("stuff"))
	log.Printf("%#v\n", err)
	assert.Nil(t, err)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	respBody, readErr := ioutil.ReadAll(resp.Body)
	assert.Nil(t, readErr)
	assert.Equal(t, `{"a": 1}`, string(respBody))
}

func TestEDIAPIClient_MakeRequest_Refresh_Needed_Error(t *testing.T) {
	ediHf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonOutput := `{"a": 1}`
		bytes, wErr := w.Write([]byte(jsonOutput))
		assert.Nil(t, wErr)
		assert.Equal(t, len(jsonOutput), bytes)
	})
	ediSrv := httptest.NewServer(ediHf)
	defer ediSrv.Close()

	// authentication and refresh should succeed
	authHF := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(nil)
	})
	authSrv := httptest.NewServer(authHF)
	defer authSrv.Close()

	cl := ediSrv.Client()
	now := time.Now()
	past := now.Add(time.Second * -3600)
	assert.Less(t, past.UnixNano(), now.UnixNano())
	client := ServerClient{
		httpClient:    cl,
		clientID:      "cid",
		clientSecret:  "cs",
		grantType:     "password",
		username:      "yusa@yusa.io",
		password:      "the greatest Password in the world",
		isInitialized: true,
		accessToken:   "GJJGFDGJJGFGJHHJF",
		refreshAt:     past,
		apiTokenURL:   authSrv.URL,
	}
	// the request is sent to our test server at ediSrv.URL
	resp, err := client.MakeRequest("GET", ediSrv.URL, strings.NewReader("stuff"))
	log.Printf("%#v\n", err)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), ": EOF")
}

func TestEDIAPIClient_MeURL(t *testing.T) {
	// invalid URL test case first
	badURL := "\t[@]<>{}^|\n"
	c := ServerClient{apiTokenURL: badURL}
	url, err := c.MeURL()
	assert.NotNil(t, err)
	assert.Equal(t, "", url)
	assert.Contains(t, err.Error(), "invalid control character in URL")

	// valid URL
	validURL := "https://accounts-core.release.slade360.co.ke/oauth2/token/"
	c1 := ServerClient{apiTokenURL: validURL}
	goodURL, goodERR := c1.MeURL()
	assert.Nil(t, goodERR)
	assert.Equal(t, "https://accounts-core.release.slade360.co.ke/v1/user/me/?format=json", goodURL)
}

func TestEDIAPIClient_RefreshAt(t *testing.T) {
	now := time.Now()
	c := ServerClient{refreshAt: now}
	assert.Equal(t, now, c.RefreshAt())
}

func TestEDIAPIClient_HTTPClient(t *testing.T) {
	c := ServerClient{httpClient: http.DefaultClient}
	assert.Equal(t, http.DefaultClient, c.HTTPClient())
}

func TestEDIAPIClient_Refresh_Uninitialized_Client(t *testing.T) {
	c := ServerClient{isInitialized: false}
	err := c.Refresh()
	assert.NotNil(t, err)
	assert.Equal(t, "cannot Refresh API tokens on an uninitialized client", err.Error())
}

func TestEDIAPIClient_Refresh_Client_Error(t *testing.T) {
	hf := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(nil) // abort the request, causing a refresh error
	})
	srv := httptest.NewServer(hf)
	defer srv.Close()
	cl := srv.Client()
	c := ServerClient{
		httpClient:    cl,
		clientID:      "cid",
		clientSecret:  "cs",
		grantType:     "password",
		username:      "yusa@yusa.io",
		password:      "the greatest Password in the world",
		isInitialized: true,
		apiTokenURL:   srv.URL, // important, this invokes our test server!
	}
	err := c.Refresh()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), ": EOF")
}

func TestEDIAPIClient_MakeRequest_Invalid_URL(t *testing.T) {
	c := ServerClient{
		isInitialized: true,
		refreshAt:     time.Now().Add(time.Second * 3600),
	}
	badURL := "\n\t\r"
	resp, err := c.MakeRequest("GET", badURL, strings.NewReader("dummy"))
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "net/url: invalid control character in URL")
}

func TestComposeAPIURL(t *testing.T) {

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
		client Client
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
			want: "https://erp-api-staging.healthcloud.co.ke/api/branches/workstations?format=json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposeAPIURL(tt.args.client, tt.args.path, tt.args.query); got != tt.want {
				t.Errorf("ComposeAPIURL() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestGetAccessToken(t *testing.T) {
	// see the README for more guidance on these env vars
	clientID := MustGetEnvVar("CLIENT_ID")
	clientSecret := MustGetEnvVar("CLIENT_SECRET")
	username := MustGetEnvVar("USERNAME")
	password := MustGetEnvVar("PASSWORD")
	grantType := MustGetEnvVar("GRANT_TYPE")
	apiScheme := MustGetEnvVar("API_SCHEME")
	apiTokenURL := MustGetEnvVar("TOKEN_URL")
	apiHost := MustGetEnvVar("HOST")

	tests := []struct {
		name    string
		args    *ServerClient
		wantErr bool
	}{
		{
			name: "valid credentials",
			args: &ServerClient{
				clientID:     clientID,
				clientSecret: clientSecret,
				apiTokenURL:  apiTokenURL,
				apiHost:      apiHost,
				apiScheme:    apiScheme,
				grantType:    grantType,
				username:     username,
				password:     password,
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			args: &ServerClient{
				clientID:     clientID,
				clientSecret: clientSecret,
				apiTokenURL:  apiTokenURL,
				apiHost:      apiHost,
				apiScheme:    apiScheme,
				grantType:    grantType,
				username:     "username",
				password:     password,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetAccessToken(tt.args)
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
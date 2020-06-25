package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
			want: errors.New("invalid access token after EDIAPIClient initialization"),
		},
		"invalid_token_type": {
			input: &ServerClient{
				accessToken: "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:   "Bogus",
			},
			want: errors.New("invalid token type after EDIAPIClient initialization, expected 'Bearer'"),
		},
		"invalid_refresh_token": {
			input: &ServerClient{
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "bad",
			},
			want: errors.New("invalid Refresh token after EDIAPIClient initialization"),
		},
		"invalid_access_scope": {
			input: &ServerClient{
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "jfdahfdjafhdjfhdalkfjdhkfhasdk",
				accessScope:  "bad",
			},
			want: errors.New("invalid access scope text after EDIAPIClient initialization"),
		},
		"invalid_expires_in": {
			input: &ServerClient{
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "jfdahfdjafhdjfhdalkfjdhkfhasdk",
				accessScope:  "scope blah blah blah more blah blah blah",
				expiresIn:    -1,
			},
			want: errors.New("invalid expiresIn after EDIAPIClient initialization"),
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
			want: errors.New("invalid past refreshAt after EDIAPIClient initialization"),
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
				authErr:      errors.New("mock auth error"),
				accessToken:  "hjhjkhkjhjkhklhkhkhjhkjh",
				tokenType:    "Bearer",
				refreshToken: "jfdahfdjafhdjfhdalkfjdhkfhasdk",
				accessScope:  "scope blah blah blah more blah blah blah",
				expiresIn:    3600,
				refreshAt:    time.Now().Add(time.Second * 3600),
			},
			want: errors.New("mock auth error"),
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

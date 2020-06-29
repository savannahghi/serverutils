package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/errorreporting"
	"github.com/stretchr/testify/assert"
)

const (
	testUserEmail         = "automated.test.bewell.ediproxy.main@healthcloud.co.ke"
	listPayersQuery       = `{"query": "query ListPayers() { payers { id, name, slade_code } }"}`
	listContactsQuery     = `{"query": "query ListContacts() { contacts { id, contact, contact_type, is_home, is_personal, is_work } }"}`
	createContactMutation = `{"query": "mutation {createContact(input: {contact: \"+254722000000\", contact_type: PHO}){id, contact, contact_type, is_home, is_personal, is_work}}"}`
)

func TestListenAddress(t *testing.T) {
	addr := ListenAddress()
	if addr != ":8080" {
		t.Errorf("unexpected listen address, got %s", addr)
	}
}

func TestFirebaseAuthenticationMiddlewareExtractBearerToken(t *testing.T) {
	fc := &FirebaseClient{}
	ctx := context.Background()
	user, userErr := GetOrCreateFirebaseUser(ctx, testUserEmail, fc)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := CreateFirebaseCustomToken(ctx, user.UID, fc)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	client := &http.Client{Timeout: time.Second * 10}
	idTokens, idErr := fc.AuthenticateCustomFirebaseToken(customToken, client)
	if idErr != nil {
		t.Errorf("unable to get ID token for user %#v, custom token %s", user, idErr)
		return
	}

	// make the necessary requests
	validRequest, _ := http.NewRequest(http.MethodPost, "/graphql", nil)
	validRequest.Header.Set("Authorization", "Bearer "+idTokens.IDToken)

	emptyRequest, _ := http.NewRequest(http.MethodPost, "/graphql", nil)

	invalidRequest, _ := http.NewRequest(http.MethodPost, "/graphql", nil)
	invalidRequest.Header.Set("Authorization", "bogus")

	badTokenRequest, _ := http.NewRequest(http.MethodPost, "/graphql", nil)
	badTokenRequest.Header.Set("Authorization", "Bearer bogus")

	// test extraction and validation
	tests := []struct {
		req *http.Request

		expectedOutput string
		expectedErrMsg string
	}{
		{
			req:            validRequest,
			expectedOutput: idTokens.IDToken,
		},
		{
			req:            emptyRequest,
			expectedOutput: "",
			expectedErrMsg: "expected an `Authorization` request header",
		},
		{
			req:            invalidRequest,
			expectedOutput: "",
			expectedErrMsg: "the `Authorization` header contents should start with `Bearer`",
		},
		{
			req:            badTokenRequest,
			expectedOutput: "bogus",
			expectedErrMsg: "invalid auth token: incorrect number of segments; see https://firebase.google.com/docs/auth/admin/verify-id-tokens for details on how to retrieve a valid ID token",
		},
	}
	for _, tc := range tests {
		token, err := ExtractBearerToken(tc.req)
		if token != tc.expectedOutput {
			t.Errorf("expected to get '%s' as the token, got '%s'", tc.expectedOutput, token)
		}
		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Errorf("expected to get error message '%s', got '%s'", tc.expectedErrMsg, err.Error())
		}

		// validate the bearer token
		if err == nil {
			app, fErr := fc.InitFirebase()
			if fErr != nil {
				t.Errorf("unable to initialize Firebase, err %s", fErr)
			}
			authToken, authErr := ValidateBearerToken(ctx, token, app)
			if authErr != nil && tc.req != badTokenRequest {
				t.Errorf("unable to validate extracted bearer token, err %s", authErr)
			}
			if authToken != nil && authToken.UID != user.UID {
				t.Errorf("unexpected auth token UID, not the same as that of the user we created")
			}
			if authErr != nil && tc.req == badTokenRequest {
				if authErr.Error() != tc.expectedErrMsg {
					t.Errorf("unexpected error msg, expected\n\t%s\n, got \n\t%s\n", tc.expectedErrMsg, authErr.Error())
				}
			}
		}
	}
}

func TestErrorMap(t *testing.T) {
	err := fmt.Errorf("test error")
	errMap := ErrorMap(err)
	if errMap["error"] == "" {
		t.Errorf("empty error key in errMap")
	}
	if errMap["error"] != "test error" {
		t.Errorf("expected the error value to be '%s', got '%s'", "test error", errMap["error"])
	}
}

func TestInitFirebase(t *testing.T) {
	fc := FirebaseClient{}
	fb, err := fc.InitFirebase()
	assert.Nil(t, err)
	assert.NotNil(t, fb)
}

func TestHasValidBearerToken(t *testing.T) {
	invalidTokenMsg := "invalid auth token: incorrect number of segments; see https://firebase.google.com/docs/auth/admin/verify-id-tokens for details on how to retrieve a valid ID token"

	fc := &FirebaseClient{}
	ctx := context.Background()
	user, userErr := GetOrCreateFirebaseUser(ctx, testUserEmail, fc)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := CreateFirebaseCustomToken(ctx, user.UID, fc)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	client := &http.Client{Timeout: time.Second * 10}
	idTokens, idErr := fc.AuthenticateCustomFirebaseToken(customToken, client)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	app, fErr := fc.InitFirebase()
	assert.Nil(t, fErr)
	assert.NotNil(t, app)

	// make the necessary request contexts
	validRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", nil)
	validRequest.Header.Set("Authorization", "Bearer "+idTokens.IDToken)

	emptyRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", nil)

	invalidRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", nil)
	invalidRequest.Header.Set("Authorization", "bogus")

	badTokenRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", nil)
	badTokenRequest.Header.Set("Authorization", "Bearer bogus")

	// listContactsContext should succeed with no error
	listPayersRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", strings.NewReader(listPayersQuery))
	listPayersRequest.Header.Set("Authorization", "Bearer "+idTokens.IDToken)

	// listContactsContext has no auth header, an error is expected
	listContactsRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", strings.NewReader(listContactsQuery))

	// test extraction and validation
	tests := map[string]struct {
		req        *http.Request
		successful bool
		errMap     map[string]string
	}{
		"valid_auth_header": {
			successful: true,
			errMap:     nil,
			req:        validRequest,
		},
		"bad auth token (not valid on Firebase auth but format is correct)": {
			successful: false,
			errMap:     map[string]string{"error": invalidTokenMsg},
			req:        badTokenRequest,
		},
		"missing_auth_header": {
			successful: false,
			errMap:     map[string]string{"error": "expected an `Authorization` request header"},
			req:        emptyRequest,
		},
		"invalid_auth_header": {
			successful: false,
			errMap:     map[string]string{"error": "the `Authorization` header contents should start with `Bearer`"},
			req:        invalidRequest,
		},
		"authenticated_payer_list_query": {
			successful: true,
			errMap:     nil,
			req:        listPayersRequest,
		},
		"unauthenticated_contact_list_query": {
			req:        listContactsRequest,
			successful: false,
			errMap:     map[string]string{"error": "expected an `Authorization` request header"}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			successful, errMap, authToken := hasValidFirebaseBearerToken(tc.req, app)
			assert.Equal(t, tc.successful, successful)
			assert.Equal(t, tc.errMap, errMap)
			if tc.successful {
				assert.NotNil(t, authToken)
			}
		})
	}
}

func Test_getFirebaseUser(t *testing.T) {
	username, unameErr := GetEnvVar("USERNAME")
	assert.Nil(t, unameErr)
	password, passwordErr := GetEnvVar("PASSWORD")
	assert.Nil(t, passwordErr)

	ctx := context.Background()
	tests := map[string]struct {
		creds            *LoginCreds
		wr               *httptest.ResponseRecorder
		expectedStatus   int
		expectedResponse string
	}{
		"invalid_creds": {
			creds: &LoginCreds{
				Username: "completely",
				Password: "bogus",
			},
			wr:               httptest.NewRecorder(),
			expectedStatus:   500,
			expectedResponse: "{\"error\":\"malformed email string: \\\"completely\\\"\"}",
		},
		"valid_creds": {
			creds: &LoginCreds{
				Username: username,
				Password: password,
			},
			wr:               httptest.NewRecorder(),
			expectedStatus:   200,
			expectedResponse: "",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			GetFirebaseUser(ctx, tc.wr, tc.creds)
			assert.Equal(t, tc.expectedStatus, tc.wr.Code)

			respBody := tc.wr.Body.Bytes()
			assert.Equal(t, tc.expectedResponse, string(respBody))
		})
	}
}

func getFirebaseTokens(t *testing.T) (string, *FirebaseUserTokens) {
	fc := &FirebaseClient{}
	ctx := context.Background()

	firebaseUser, userErr := GetOrCreateFirebaseUser(ctx, testUserEmail, fc)
	assert.Nil(t, userErr)

	customToken, tokenErr := CreateFirebaseCustomToken(ctx, firebaseUser.UID, fc)
	assert.Nil(t, tokenErr)

	idTokens, authErr := fc.AuthenticateCustomFirebaseToken(customToken, http.DefaultClient)
	assert.Nil(t, authErr)

	return customToken, idTokens
}

func Test_authenticateCustomToken(t *testing.T) {
	validCustomToken, _ := getFirebaseTokens(t)
	tests := map[string]struct {
		customToken        string
		httpClient         *http.Client
		firebaseClient     IFirebaseClient
		wr                 *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		"successful_authentication": {
			customToken:        validCustomToken,
			httpClient:         http.DefaultClient,
			firebaseClient:     &FirebaseClient{},
			wr:                 httptest.NewRecorder(),
			expectedStatusCode: 200,
			expectedResponse:   "",
		},
		"unsuccessful_authentication": {
			customToken:        "invalid token will trigger err",
			httpClient:         http.DefaultClient,
			firebaseClient:     &FirebaseClient{},
			wr:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "firebase HTTP error, status code 400",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			AuthenticateCustomToken(
				tc.wr, tc.customToken, tc.httpClient, tc.firebaseClient)
			assert.Equal(t, tc.expectedStatusCode, tc.wr.Code)
			assert.Contains(t, tc.wr.Body.String(), tc.expectedResponse)
		})
	}
}

func Test_convertStringToInt(t *testing.T) {
	tests := map[string]struct {
		val                string
		rw                 *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		"successful_conversion": {
			val:                "768",
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: 200,
		},
		"failed_conversion": {
			val:                "not an int",
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "{\"error\":\"strconv.Atoi: parsing \\\"not an int\\\": invalid syntax\"}",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ConvertStringToInt(tc.rw, tc.val)
			assert.Equal(t, tc.expectedStatusCode, tc.rw.Code)
			assert.Equal(t, tc.expectedResponse, tc.rw.Body.String())
		})
	}
}

func Test_decodeRefreshResponse(t *testing.T) {
	refreshJSONBytes, jsonErr := json.Marshal(firebaseRefreshResponse{
		ExpiresIn:    "3600",
		TokenType:    "bearer",
		RefreshToken: "this is a refresh token",
		IDToken:      "this is an ID token",
		UserID:       "some uid",
		ProjectID:    "some project id",
	})
	assert.Nil(t, jsonErr)

	tests := map[string]struct {
		resp               *http.Response
		wr                 *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		"successful_decode": {
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(string(refreshJSONBytes))),
			},
			wr:                 httptest.NewRecorder(),
			expectedStatusCode: 200,
			expectedResponse:   "",
		},
		"failed_decode": {
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader("invalid response won't decode")),
			},
			wr:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "{\"error\":\"invalid character 'i' looking for beginning of value\"}",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			decodeRefreshResponse(tc.wr, tc.resp)
			assert.Equal(t, tc.expectedStatusCode, tc.wr.Code)
			assert.Equal(t, tc.expectedResponse, tc.wr.Body.String())
		})
	}
}

func Test_composeRefreshRequest(t *testing.T) {
	apiKey, apiKeyErr := GetEnvVar(FirebaseWebAPIKeyEnvVarName)
	assert.Nil(t, apiKeyErr)

	tests := map[string]struct {
		creds *refreshCreds
	}{
		"successful": {
			creds: &refreshCreds{RefreshToken: "mock refresh token"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			refreshURL, refreshData := composeRefreshRequest(tc.creds)
			assert.Equal(t, FirebaseRefreshTokenURL+apiKey, refreshURL)

			payloadBytes, readErr := ioutil.ReadAll(refreshData)
			assert.Nil(t, readErr)
			assert.Equal(
				t,
				"grant_type=refresh_token&refresh_token=mock+refresh+token", string(payloadBytes),
			)
		})
	}
}

func Test_postRefreshRequest(t *testing.T) {
	postErrHTTPClient := MockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("ka-boom")
	})
	clientErrorResponseHTTPClient := MockHTTPClient(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			Body:       ioutil.NopCloser(strings.NewReader("dummy response content")),
			Status:     "400 Bad Request",
			StatusCode: http.StatusBadRequest,
		}
		return resp, nil
	})
	okResponseHTTPClient := MockHTTPClient(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			Body:       ioutil.NopCloser(strings.NewReader("dummy response content")),
			Status:     "200 OK",
			StatusCode: http.StatusOK,
		}
		return resp, nil
	})

	tests := map[string]struct {
		httpClient         *http.Client
		refreshURL         string
		encodedRefreshData io.Reader
		expectNilResp      bool
		wr                 *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		"successful": {
			httpClient:         okResponseHTTPClient,
			refreshURL:         "http://r.frsh.url",
			encodedRefreshData: strings.NewReader("wont be used"),
			expectNilResp:      false,
			wr:                 httptest.NewRecorder(),
			expectedStatusCode: 200,
			expectedResponse:   "",
		},
		"post_error": {
			httpClient:         postErrHTTPClient,
			refreshURL:         "http://r.frsh.url",
			encodedRefreshData: strings.NewReader("wont be used"),
			expectNilResp:      true,
			wr:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "ka-boom",
		},
		"bad_status_code": {
			httpClient:         clientErrorResponseHTTPClient,
			refreshURL:         "http://r.frsh.url",
			encodedRefreshData: strings.NewReader("wont be used"),
			expectNilResp:      true,
			wr:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "{\"error\":\"refresh auth server error: 400\"}",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp := postRefreshRequest(
				tc.wr, tc.httpClient, tc.refreshURL, tc.encodedRefreshData)
			assert.Equal(t, tc.expectedStatusCode, tc.wr.Code)
			assert.Contains(t, tc.wr.Body.String(), tc.expectedResponse)

			if tc.expectNilResp {
				assert.Nil(t, resp)
			}
		})
	}
}

func Test_logoutFirebase(t *testing.T) {
	ctx := context.Background()

	goodFc := &FirebaseClient{}
	firebaseApp, faErr := goodFc.InitFirebase()
	assert.Nil(t, faErr)
	authClient, clErr := firebaseApp.Auth(ctx)

	assert.Nil(t, clErr)
	assert.NotNil(t, authClient)
	assert.NotNil(t, authClient.RevokeRefreshTokens)

	initErrFc := &MockFirebaseClient{
		MockAppInitErr: fmt.Errorf("mock firebase init err"),
	}
	initApp, initErr := initErrFc.InitFirebase()
	assert.NotNil(t, initErr)
	assert.Nil(t, initApp)

	authErrFc := &MockFirebaseClient{
		MockApp: &MockFirebaseApp{
			MockAuthErr: fmt.Errorf("can't get an auth client"),
		},
		MockFirebaseAuthError: fmt.Errorf("mock firebase auth err"),
	}
	revokeErrFc := &MockFirebaseClient{
		MockApp: &MockFirebaseApp{
			MockAuthClient: authClient,
		},
	}

	tests := map[string]struct {
		ctx                context.Context
		fc                 IFirebaseClient
		req                logoutRequest
		rw                 *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		"real_call_to_firebase_invalid_uid": {
			ctx:                ctx,
			fc:                 goodFc,
			req:                logoutRequest{UID: "dummy uid"},
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "http error status: 400",
		},
		"init_error": {
			ctx:                ctx,
			fc:                 initErrFc,
			req:                logoutRequest{UID: "dummy uid"},
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "",
		},
		"auth_error": {
			ctx:                ctx,
			fc:                 authErrFc,
			req:                logoutRequest{UID: "dummy uid"},
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "",
		},
		"revoke_error": {
			ctx:                ctx,
			fc:                 revokeErrFc,
			req:                logoutRequest{UID: "dummy uid"},
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: 500,
			expectedResponse:   "",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fmt.Println(name)
			logoutFirebase(tc.ctx, tc.rw, tc.fc, &tc.req)
			assert.Equal(t, tc.expectedStatusCode, tc.rw.Code)
			assert.Contains(t, tc.rw.Body.String(), tc.expectedResponse)
		})
	}
}

func Test_Refresh(t *testing.T) {
	ctx := context.Background()
	fc := &FirebaseClient{}
	refreshFunc := GetRefreshFunc(http.DefaultClient)

	user, userErr := GetOrCreateFirebaseUser(ctx, testUserEmail, fc)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := CreateFirebaseCustomToken(ctx, user.UID, fc)
	assert.Nil(t, tokenErr)

	idTokens, authErr := fc.AuthenticateCustomFirebaseToken(customToken, http.DefaultClient)
	assert.Nil(t, authErr)

	validRefreshCreds := refreshCreds{RefreshToken: idTokens.RefreshToken}
	validRefreshCredsBytes, vMarshalErr := json.Marshal(validRefreshCreds)
	assert.Nil(t, vMarshalErr)
	validReq, err := http.NewRequestWithContext(ctx, "POST", "/login", bytes.NewReader(validRefreshCredsBytes))
	assert.Nil(t, err)

	invalidRefreshCreds := refreshCreds{RefreshToken: "this token is not valid"}
	invalidRefreshCredsBytes, iMarshalErr := json.Marshal(invalidRefreshCreds)
	assert.Nil(t, iMarshalErr)
	invalidReq, err := http.NewRequestWithContext(ctx, "POST", "/login", bytes.NewReader(invalidRefreshCredsBytes))
	assert.Nil(t, err)

	tests := map[string]struct {
		req                *http.Request
		rw                 *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		"valid_refresh_request_token_valid": {
			expectedStatusCode: http.StatusOK,
			req:                validReq,
			rw:                 httptest.NewRecorder(),
			expectedResponse:   "{\"expires_in\":3600,\"id_token\":\"",
		},
		"valid_refresh_format_token_invalid": {
			expectedStatusCode: 500,
			req:                invalidReq,
			rw:                 httptest.NewRecorder(),
			expectedResponse:   "{\"error\":\"refresh auth server error: 400\"}",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			refreshFunc(tc.rw, tc.req)
			assert.Equal(t, tc.expectedStatusCode, tc.rw.Code)
			assert.Contains(t, tc.rw.Body.String(), tc.expectedResponse)
		})
	}
}

func Test_Logout(t *testing.T) {
	ctx := context.Background()
	fc := &FirebaseClient{}
	logoutFunc := GetLogoutFunc(ctx, fc)

	user, userErr := GetOrCreateFirebaseUser(ctx, testUserEmail, fc)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)
	goodLogoutRequest := logoutRequest{UID: user.UID}
	goodLogoutRequestBytes, gMarshalErr := json.Marshal(goodLogoutRequest)
	assert.Nil(t, gMarshalErr)
	goodLogoutReq, err := http.NewRequestWithContext(ctx, "POST", "/logout", bytes.NewReader(goodLogoutRequestBytes))
	assert.Nil(t, err)

	invalidLogoutRequest := logoutRequest{UID: "this uid does not exist"}
	invalidLogoutRequestBytes, marshalErr := json.Marshal(invalidLogoutRequest)
	assert.Nil(t, marshalErr)
	invalidLogoutReq, err := http.NewRequestWithContext(ctx, "POST", "/logout", bytes.NewReader(invalidLogoutRequestBytes))
	assert.Nil(t, err)

	tests := map[string]struct {
		expectedStatusCode int
		expectedResponse   string
		req                *http.Request
		rw                 *httptest.ResponseRecorder
	}{
		"valid_logout_request_uid_exists": {
			expectedStatusCode: 200,
			req:                goodLogoutReq,
			rw:                 httptest.NewRecorder(),
		},
		"valid_logout_format_uid_does_not_exist": {
			expectedStatusCode: 500,
			expectedResponse:   "http error status: 400",
			req:                invalidLogoutReq,
			rw:                 httptest.NewRecorder(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			logoutFunc(tc.rw, tc.req)
			assert.Equal(t, tc.expectedStatusCode, tc.rw.Code)
			assert.Contains(t, tc.rw.Body.String(), tc.expectedResponse)
		})
	}
}

func Test_StackDriver_Setup(t *testing.T) {
	errorClient := StackDriver(context.Background())
	err := fmt.Errorf("test error")
	if errorClient != nil {
		errorClient.Report(errorreporting.Entry{
			Error: err,
		})
	}
}

func TestSentry(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "default case",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Sentry(); (err != nil) != tt.wantErr {
				t.Errorf("Sentry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractBearerToken(t *testing.T) {

	fc := &FirebaseClient{}
	ctx := context.Background()
	user, userErr := GetOrCreateFirebaseUser(ctx, testUserEmail, fc)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := CreateFirebaseCustomToken(ctx, user.UID, fc)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	client := &http.Client{Timeout: time.Second * 10}
	idTokens, idErr := fc.AuthenticateCustomFirebaseToken(customToken, client)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	app, fErr := fc.InitFirebase()
	assert.Nil(t, fErr)
	assert.NotNil(t, app)

	// make the necessary request contexts
	validRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", nil)
	validRequest.Header.Set("Authorization", "Bearer "+idTokens.IDToken)

	nilHeaderReq := httptest.NewRequest(
		http.MethodPost, "/", nil)
	nilHeaderReq.Header = nil

	emptyRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", nil)

	invalidRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", nil)
	invalidRequest.Header.Set("Authorization", "bogus")

	badTokenRequest, _ := http.NewRequest(
		http.MethodPost, "/graphql", nil)
	badTokenRequest.Header.Set("Authorization", "Bearer bogus")

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid request",
			args: args{
				r: validRequest,
			},
			want:    idTokens.IDToken,
			wantErr: false,
		},
		{
			name: "request with no header",
			args: args{
				r: emptyRequest,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "nil request",
			args: args{
				r: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "badly formatted request",
			args: args{
				r: invalidRequest,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "request with a nil header",
			args: args{
				r: nilHeaderReq,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "bad auth token",
			args: args{
				r: badTokenRequest,
			},
			want:    "bogus", // ExtractBearerToken does not validate, it simply extracts
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractBearerToken(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStackDriver(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy case",
			args: args{
				ctx: ctx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StackDriver(tt.args.ctx)
			assert.NotNil(t, got)
		})
	}
}

func TestFetchUserProfile(t *testing.T) {
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
		w         http.ResponseWriter
		ediClient Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				w:         httptest.NewRecorder(),
				ediClient: client,
			},
		},
		{
			name: "error case",
			args: args{
				w:         httptest.NewRecorder(),
				ediClient: &ServerClient{}, // not initialized
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FetchUserProfile(tt.args.w, tt.args.ediClient)
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestLogStartupError(t *testing.T) {
	type args struct {
		ctx context.Context
		err error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
				err: fmt.Errorf("this is a test error"),
			},
		},
		{
			name: "nil error",
			args: args{
				ctx: context.Background(),
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LogStartupError(tt.args.ctx, tt.args.err)
		})
	}
}

func TestDecodeJSONToTargetStruct(t *testing.T) {
	type target struct {
		A string `json:"a,omitempty"`
	}
	targetStruct := target{}

	type args struct {
		w            http.ResponseWriter
		r            *http.Request
		targetStruct interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "good case",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/",
					bytes.NewBuffer([]byte(
						"{\"a\":\"1\"}",
					)),
				),
				targetStruct: &targetStruct,
			},
		},
		{
			name: "invalid / failed decode",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodGet,
					"/",
					nil,
				),
				targetStruct: &targetStruct,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DecodeJSONToTargetStruct(tt.args.w, tt.args.r, tt.args.targetStruct)
		})
	}
}

func TestRequestDebugMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	mw := RequestDebugMiddleware()
	h := mw(next)

	rw := httptest.NewRecorder()
	reader := bytes.NewBuffer([]byte("sample"))
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	h.ServeHTTP(rw, req)

	rw1 := httptest.NewRecorder()
	reader1 := ioutil.NopCloser(bytes.NewBuffer([]byte("will be closed")))
	err := reader1.Close()
	assert.Nil(t, err)
	req1 := httptest.NewRequest(http.MethodPost, "/", reader1)
	h.ServeHTTP(rw1, req1)
}

func TestAuthenticationMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fc := &FirebaseClient{}
	fa, err := fc.InitFirebase()
	assert.Nil(t, err)
	assert.NotNil(t, fa)

	mw := AuthenticationMiddleware(fa)
	h := mw(next)
	rw := httptest.NewRecorder()
	reader := bytes.NewBuffer([]byte("sample"))
	idToken := GetIDToken(t)
	authHeader := fmt.Sprintf("Bearer %s", idToken)
	req := httptest.NewRequest(http.MethodPost, "/", reader)
	req.Header.Add("Authorization", authHeader)
	h.ServeHTTP(rw, req)

	rw1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodPost, "/", reader)
	h.ServeHTTP(rw1, req1)

}

package base_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func getFirebaseTokens(t *testing.T) (string, *base.FirebaseUserTokens) {
	ctx := context.Background()

	firebaseUser, userErr := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
	assert.Nil(t, userErr)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, firebaseUser.UID)
	assert.Nil(t, tokenErr)

	idTokens, authErr := base.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, authErr)

	return customToken, idTokens
}

func TestLoginClient(t *testing.T) {
	username, err := base.GetEnvVar(base.UsernameEnvVarName)
	assert.Nil(t, err)

	password, err := base.GetEnvVar(base.PasswordEnvVarName)
	assert.Nil(t, err)

	client, err := base.LoginClient(username, password)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.True(t, client.IsInitialized())
}

func TestLoginClient_Error(t *testing.T) {
	badUsername := "u" // not a valid email
	badPassword := "p" // too short
	client, err := base.LoginClient(badUsername, badPassword)
	assert.Nil(t, client)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the username `u` is not a valid email")
}

func TestAPIClient(t *testing.T) {
	got, err := base.APIClient() // initialize with default params
	assert.Nil(t, err)
	assert.NotNil(t, got)
}

func TestExtractBearerToken(t *testing.T) {
	fc := &base.FirebaseClient{}
	ctx := context.Background()
	user, userErr := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := base.AuthenticateCustomFirebaseToken(customToken)
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
			got, err := base.ExtractBearerToken(tt.args.r)
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

func TestFirebaseAuthenticationMiddlewareExtractBearerToken(t *testing.T) {
	ctx := context.Background()
	user, userErr := base.GetOrCreateFirebaseUser(ctx, testUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := base.AuthenticateCustomFirebaseToken(customToken)
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
		token, err := base.ExtractBearerToken(tc.req)
		if token != tc.expectedOutput {
			t.Errorf("expected to get '%s' as the token, got '%s'", tc.expectedOutput, token)
		}
		if err != nil && err.Error() != tc.expectedErrMsg {
			t.Errorf("expected to get error message '%s', got '%s'", tc.expectedErrMsg, err.Error())
		}

		// validate the bearer token
		if err == nil {
			authToken, err := base.ValidateBearerToken(ctx, token)
			if err != nil && tc.req != badTokenRequest {
				t.Errorf("unable to validate extracted bearer token, err %s", err)
			}
			if authToken != nil && authToken.UID != user.UID {
				t.Errorf("unexpected auth token UID, not the same as that of the user we created")
			}
			if err != nil && tc.req == badTokenRequest {
				if err.Error() != tc.expectedErrMsg {
					t.Errorf("unexpected error msg, expected\n\t%s\n, got \n\t%s\n", tc.expectedErrMsg, err.Error())
				}
			}
		}
	}
}

func TestAuthenticationMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	fc := &base.FirebaseClient{}
	fa, err := fc.InitFirebase()
	assert.Nil(t, err)
	assert.NotNil(t, fa)

	mw := base.AuthenticationMiddleware(fa)
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

func TestFetchUserProfile(t *testing.T) {
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

	type args struct {
		ediClient base.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ediClient: client,
			},
		},
		{
			name: "error case",
			args: args{
				ediClient: &base.ServerClient{}, // not initialized
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.FetchUserProfile(tt.args.ediClient)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_authenticateCustomToken(t *testing.T) {
	validCustomToken, _ := getFirebaseTokens(t)
	tests := map[string]struct {
		customToken        string
		httpClient         *http.Client
		expectedStatusCode int
		expectedResponse   string
		wantErr            bool
	}{
		"successful_authentication": {
			customToken:        validCustomToken,
			httpClient:         http.DefaultClient,
			expectedStatusCode: 200,
			expectedResponse:   "",
			wantErr:            false,
		},
		"unsuccessful_authentication": {
			customToken:        "invalid token will trigger err",
			httpClient:         http.DefaultClient,
			expectedStatusCode: 500,
			expectedResponse:   "firebase HTTP error, status code 400",
			wantErr:            true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tokens, err := base.AuthenticateCustomToken(tc.customToken, tc.httpClient)
			if !tc.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, tokens)
			}
		})
	}
}

func Test_Refresh(t *testing.T) {
	ctx := context.Background()
	refreshFunc := base.GetRefreshFunc()

	user, userErr := base.GetOrCreateFirebaseUser(ctx, testUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)

	idTokens, authErr := base.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, authErr)

	validRefreshCreds := base.RefreshCreds{RefreshToken: idTokens.RefreshToken}
	validRefreshCredsBytes, vMarshalErr := json.Marshal(validRefreshCreds)
	assert.Nil(t, vMarshalErr)
	validReq, err := http.NewRequestWithContext(ctx, "POST", "/login", bytes.NewReader(validRefreshCredsBytes))
	assert.Nil(t, err)

	invalidRefreshCreds := base.RefreshCreds{RefreshToken: "this token is not valid"}
	invalidRefreshCredsBytes, iMarshalErr := json.Marshal(invalidRefreshCreds)
	assert.Nil(t, iMarshalErr)
	invalidReq, err := http.NewRequestWithContext(ctx, "POST", "/login", bytes.NewReader(invalidRefreshCredsBytes))
	assert.Nil(t, err)

	tests := map[string]struct {
		req                *http.Request
		rw                 *httptest.ResponseRecorder
		expectedStatusCode int
	}{
		"valid_refresh_request_token_valid": {
			expectedStatusCode: http.StatusOK,
			req:                validReq,
			rw:                 httptest.NewRecorder(),
		},
		"valid_refresh_format_token_invalid": {
			expectedStatusCode: http.StatusInternalServerError,
			req:                invalidReq,
			rw:                 httptest.NewRecorder(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			refreshFunc(tc.rw, tc.req)
			assert.Equal(t, tc.expectedStatusCode, tc.rw.Code)
		})
	}
}

func TestGetLoginFunc(t *testing.T) {
	ctx := context.Background()
	fc := &base.FirebaseClient{}
	loginFunc := base.GetLoginFunc(ctx, fc)

	goodLoginCredsJSONBytes, err := json.Marshal(&base.LoginCreds{
		Username: base.MustGetEnvVar(base.UsernameEnvVarName),
		Password: base.MustGetEnvVar(base.PasswordEnvVarName),
	})
	assert.Nil(t, err)
	assert.NotNil(t, goodLoginCredsJSONBytes)
	goodLoginCredsReq := httptest.NewRequest(http.MethodGet, "/", nil)
	goodLoginCredsReq.Body = ioutil.NopCloser(bytes.NewReader(goodLoginCredsJSONBytes))

	incorrectLoginCredsJSONBytes, err := json.Marshal(&base.LoginCreds{
		Username: "not a real username",
		Password: "not a real password",
	})
	assert.Nil(t, err)
	assert.NotNil(t, incorrectLoginCredsJSONBytes)
	incorrectLoginCredsReq := httptest.NewRequest(http.MethodGet, "/", nil)
	incorrectLoginCredsReq.Body = ioutil.NopCloser(bytes.NewReader(incorrectLoginCredsJSONBytes))

	wrongFormatLoginCredsJSONBytes, err := json.Marshal(&base.AccessTokenPayload{})
	assert.Nil(t, err)
	assert.NotNil(t, wrongFormatLoginCredsJSONBytes)
	wrongFormatLoginCredsReq := httptest.NewRequest(http.MethodGet, "/", nil)
	wrongFormatLoginCredsReq.Body = ioutil.NopCloser(bytes.NewReader(wrongFormatLoginCredsJSONBytes))

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "invalid login credentials - format",
			args: args{
				w: httptest.NewRecorder(),
				r: wrongFormatLoginCredsReq,
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "incorrect login credentials - good format but won't login",
			args: args{
				w: httptest.NewRecorder(),
				r: incorrectLoginCredsReq,
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "correct login credentials",
			args: args{
				w: httptest.NewRecorder(),
				r: goodLoginCredsReq,
			},
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginFunc(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			assert.True(t, ok)
			assert.NotNil(t, rec)

			assert.Equal(t, rec.Code, tt.wantStatusCode)
		})
	}
}

func TestGetVerifyTokenFunc(t *testing.T) {
	ctx := context.Background()
	fc := &base.FirebaseClient{}
	verifyFunc := base.GetVerifyTokenFunc(ctx, fc)

	invalidAccessTokenJSONBytes, err := json.Marshal(
		&base.AccessTokenPayload{
			AccessToken: "not a valid access token",
		},
	)
	assert.Nil(t, err)

	invalidTokenReq := httptest.NewRequest(http.MethodGet, "/", nil)
	invalidTokenReq.Body = ioutil.NopCloser(bytes.NewReader(invalidAccessTokenJSONBytes))
	invalidTokenReq.Header.Add("X-Authorization", "Bearer not a real token")

	client, err := base.APIClient()
	assert.NotNil(t, client)
	assert.Nil(t, err)
	validTokenReq := httptest.NewRequest(http.MethodGet, "/", nil)
	validAccessTokenJSONBytes, err := json.Marshal(
		&base.AccessTokenPayload{
			AccessToken: client.AccessToken(),
		},
	)
	assert.Nil(t, err)
	validTokenReq.Body = ioutil.NopCloser(bytes.NewReader(validAccessTokenJSONBytes))
	validTokenReq.Header.Add("X-Authorization", "Bearer not a real token")

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "invalid access token - request with no header",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "valid header format with invalid token",
			args: args{
				w: httptest.NewRecorder(),
				r: invalidTokenReq,
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "valid header format with valid token",
			args: args{
				w: httptest.NewRecorder(),
				r: validTokenReq,
			},
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifyFunc(tt.args.w, tt.args.r)

			rec, ok := tt.args.w.(*httptest.ResponseRecorder)
			assert.True(t, ok)
			assert.NotNil(t, rec)

			assert.Equal(t, rec.Code, tt.wantStatusCode)
		})
	}
}

func Test_decodeRefreshResponse(t *testing.T) {
	refreshJSONBytes, jsonErr := json.Marshal(base.FirebaseRefreshResponse{
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
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "",
		},
		"failed_decode": {
			resp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader("invalid response won't decode")),
			},
			wr:                 httptest.NewRecorder(),
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "{\"error\":\"invalid character 'i' looking for beginning of value\"}",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			base.DecodeRefreshResponse(tc.wr, tc.resp)
			assert.Equal(t, tc.expectedStatusCode, tc.wr.Code)
			assert.Equal(t, tc.expectedResponse, tc.wr.Body.String())
		})
	}
}

func Test_composeRefreshRequest(t *testing.T) {
	apiKey, apiKeyErr := base.GetEnvVar(base.FirebaseWebAPIKeyEnvVarName)
	assert.Nil(t, apiKeyErr)

	tests := map[string]struct {
		creds *base.RefreshCreds
	}{
		"successful": {
			creds: &base.RefreshCreds{RefreshToken: "mock refresh token"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			refreshURL, refreshData := base.ComposeRefreshRequest(tc.creds)
			assert.Equal(t, base.FirebaseRefreshTokenURL+apiKey, refreshURL)

			payloadBytes, readErr := ioutil.ReadAll(refreshData)
			assert.Nil(t, readErr)
			assert.Equal(
				t,
				"grant_type=refresh_token&refresh_token=mock+refresh+token", string(payloadBytes),
			)
		})
	}
}

func Test_logoutFirebase(t *testing.T) {
	ctx := context.Background()

	goodFc := &base.FirebaseClient{}
	firebaseApp, faErr := goodFc.InitFirebase()
	assert.Nil(t, faErr)
	authClient, clErr := firebaseApp.Auth(ctx)

	assert.Nil(t, clErr)
	assert.NotNil(t, authClient)
	assert.NotNil(t, authClient.RevokeRefreshTokens)

	initErrFc := &base.MockFirebaseClient{
		MockAppInitErr: fmt.Errorf("mock firebase init err"),
	}
	initApp, initErr := initErrFc.InitFirebase()
	assert.NotNil(t, initErr)
	assert.Nil(t, initApp)

	authErrFc := &base.MockFirebaseClient{
		MockApp: &base.MockFirebaseApp{
			MockAuthErr: fmt.Errorf("can't get an auth client"),
		},
		MockFirebaseAuthError: fmt.Errorf("mock firebase auth err"),
	}
	revokeErrFc := &base.MockFirebaseClient{
		MockApp: &base.MockFirebaseApp{
			MockAuthClient: authClient,
		},
	}

	tests := map[string]struct {
		ctx                context.Context
		fc                 base.IFirebaseClient
		req                base.LogoutRequest
		rw                 *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		"real_call_to_firebase_invalid_uid": {
			ctx:                ctx,
			fc:                 goodFc,
			req:                base.LogoutRequest{UID: "dummy uid"},
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "http error status: 400",
		},
		"init_error": {
			ctx:                ctx,
			fc:                 initErrFc,
			req:                base.LogoutRequest{UID: "dummy uid"},
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "",
		},
		"auth_error": {
			ctx:                ctx,
			fc:                 authErrFc,
			req:                base.LogoutRequest{UID: "dummy uid"},
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "",
		},
		"revoke_error": {
			ctx:                ctx,
			fc:                 revokeErrFc,
			req:                base.LogoutRequest{UID: "dummy uid"},
			rw:                 httptest.NewRecorder(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			base.LogoutFirebase(tc.ctx, tc.rw, tc.fc, &tc.req)
			assert.Equal(t, tc.expectedStatusCode, tc.rw.Code)
			assert.Contains(t, tc.rw.Body.String(), tc.expectedResponse)
		})
	}
}

func Test_Logout(t *testing.T) {
	ctx := context.Background()
	fc := &base.FirebaseClient{}
	logoutFunc := base.GetLogoutFunc(ctx, fc)

	user, userErr := base.GetOrCreateFirebaseUser(ctx, testUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)
	goodLogoutRequest := base.LogoutRequest{UID: user.UID}
	goodLogoutRequestBytes, gMarshalErr := json.Marshal(goodLogoutRequest)
	assert.Nil(t, gMarshalErr)
	goodLogoutReq, err := http.NewRequestWithContext(ctx, "POST", "/logout", bytes.NewReader(goodLogoutRequestBytes))
	assert.Nil(t, err)

	invalidLogoutRequest := base.LogoutRequest{UID: "this uid does not exist"}
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
			expectedStatusCode: http.StatusOK,
			req:                goodLogoutReq,
			rw:                 httptest.NewRecorder(),
		},
		"valid_logout_format_uid_does_not_exist": {
			expectedStatusCode: http.StatusInternalServerError,
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

func TestHasValidBearerToken(t *testing.T) {
	invalidTokenMsg := "invalid auth token: incorrect number of segments; see https://firebase.google.com/docs/auth/admin/verify-id-tokens for details on how to retrieve a valid ID token"

	fc := &base.FirebaseClient{}
	ctx := context.Background()
	user, userErr := base.GetOrCreateFirebaseUser(ctx, testUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := base.AuthenticateCustomFirebaseToken(customToken)
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
			successful, errMap, authToken := base.HasValidFirebaseBearerToken(tc.req, app)
			assert.Equal(t, tc.successful, successful)
			assert.Equal(t, tc.errMap, errMap)
			if tc.successful {
				assert.NotNil(t, authToken)
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	nilHeaderRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	nilHeaderRequest.Header = nil

	noHeaderRequest := httptest.NewRequest(http.MethodGet, "/", nil)

	headerWithNoPrefixRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	headerWithNoPrefixRequest.Header.Add("Authorization", "has no prefix")

	goodRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	goodRequest.Header.Add("Authorization", "Bearer something something")

	type args struct {
		r      *http.Request
		header string
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				r:      goodRequest,
				header: "Authorization",
				prefix: "Bearer",
			},
			want:    "something something",
			wantErr: false,
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
			name: "nil header",
			args: args{
				r: nilHeaderRequest,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no header",
			args: args{
				r: noHeaderRequest,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "header with no prefix",
			args: args{
				r:      headerWithNoPrefixRequest,
				header: "Authorization",
				prefix: "Bearer",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.ExtractToken(tt.args.r, tt.args.header, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAccessToken(t *testing.T) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)
	assert.NotNil(t, firebaseApp)

	client, err := base.APIClient()
	assert.Nil(t, err)
	assert.NotNil(t, client)

	validBearerToken := fmt.Sprintf("Bearer %s", client.AccessToken())

	noXBearerTokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	noXBearerTokenRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

	invalidXBearerTokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	invalidXBearerTokenRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
	invalidXBearerTokenRequest.Header.Add("X-Authorization", "Bearer not a valid token")

	validXBearerTokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	validXBearerTokenRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
	validXBearerTokenRequest.Header.Add("X-Authorization", validBearerToken)

	type args struct {
		accessToken string
	}
	tests := []struct {
		name        string
		args        args
		wantIsValid bool
		wantErr     bool
	}{
		{
			name: "valid access token",
			args: args{
				accessToken: client.AccessToken(),
			},
			wantIsValid: true,
			wantErr:     false,
		},
		{
			name: "invalid access token",
			args: args{
				accessToken: "this is not a valid access token",
			},
			wantIsValid: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, creds, err := base.ValidateAccessToken(tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if valid != tt.wantIsValid {
				t.Errorf("ValidateAccessToken() got = %v, want %v", valid, tt.wantIsValid)
			}
			if !tt.wantErr {
				assert.NotNil(t, creds)
			}
		})
	}
}

func TestHasValidSlade360AccessToken(t *testing.T) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)
	assert.NotNil(t, firebaseApp)

	client, err := base.APIClient()
	assert.Nil(t, err)
	assert.NotNil(t, client)

	validBearerToken := fmt.Sprintf("Bearer %s", client.AccessToken())

	noXBearerTokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	noXBearerTokenRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

	invalidXBearerTokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	invalidXBearerTokenRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
	invalidXBearerTokenRequest.Header.Add("X-Authorization", "Bearer not a valid token")

	validXBearerTokenRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	validXBearerTokenRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
	validXBearerTokenRequest.Header.Add("X-Authorization", validBearerToken)

	type args struct {
		r           *http.Request
		firebaseApp base.IFirebaseApp
	}
	tests := []struct {
		name              string
		args              args
		wantHasValidToken bool
	}{
		{
			name: "bad bearer token",
			args: args{
				r:           noXBearerTokenRequest,
				firebaseApp: firebaseApp,
			},
			wantHasValidToken: false,
		},
		{
			name: "invalid bearer token",
			args: args{
				r:           invalidXBearerTokenRequest,
				firebaseApp: firebaseApp,
			},
			wantHasValidToken: false,
		},
		{
			name: "valid access token",
			args: args{
				r:           validXBearerTokenRequest,
				firebaseApp: firebaseApp,
			},
			wantHasValidToken: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasValidToken, errMap, authToken := base.HasValidSlade360AccessToken(tt.args.r, tt.args.firebaseApp)
			if hasValidToken != tt.wantHasValidToken {
				t.Errorf("HasValidSlade360AccessToken() got = %v, want %v", hasValidToken, tt.wantHasValidToken)
			}
			if tt.wantHasValidToken {
				assert.NotNil(t, authToken)
				assert.Nil(t, errMap)
			}
			if !tt.wantHasValidToken {
				assert.NotNil(t, errMap)
			}
		})
	}
}

func TestGetFirebaseUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		creds *base.LoginCreds
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: context.Background(),
				creds: &base.LoginCreds{
					Username: base.MustGetEnvVar(base.UsernameEnvVarName),
					Password: base.MustGetEnvVar(base.PasswordEnvVarName),
				},
			},
			wantErr: false,
		},
		{
			name: "nil creds",
			args: args{
				ctx:   context.Background(),
				creds: nil,
			},
			wantErr: true,
		},
		{
			name: "bad creds",
			args: args{
				ctx: context.Background(),
				creds: &base.LoginCreds{
					Username: "completely",
					Password: "bad",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetFirebaseUser(tt.args.ctx, tt.args.creds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFirebaseUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestGetDefaultUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		creds *base.LoginCreds
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: context.Background(),
				creds: &base.LoginCreds{
					Username: base.MustGetEnvVar(base.UsernameEnvVarName),
					Password: base.MustGetEnvVar(base.PasswordEnvVarName),
				},
			},
			wantErr: false,
		},
		{
			name: "bad creds",
			args: args{
				ctx: context.Background(),
				creds: &base.LoginCreds{
					Username: "completely",
					Password: "bad",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, profile, err := base.GetDefaultUser(
				tt.args.ctx,
				tt.args.creds,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, tokens)
				assert.NotNil(t, profile)
			}
		})
	}
}

func TestValidateLoginCreds(t *testing.T) {
	goodCreds := &base.LoginCreds{
		Username: "yusa",
		Password: "pass",
	}
	goodCredsJSONBytes, err := json.Marshal(goodCreds)
	assert.Nil(t, err)
	assert.NotNil(t, goodCredsJSONBytes)

	validRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	validRequest.Body = ioutil.NopCloser(bytes.NewReader(goodCredsJSONBytes))

	emptyCredsRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	emptyCredsRequest.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *base.LoginCreds
		wantErr bool
	}{
		{
			name: "valid creds",
			args: args{
				w: httptest.NewRecorder(),
				r: validRequest,
			},
			want: &base.LoginCreds{
				Username: "yusa",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "invalid creds",
			args: args{
				w: httptest.NewRecorder(),
				r: emptyCredsRequest,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.ValidateLoginCreds(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLoginCreds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateLoginCreds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReportErr(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		err    error
		status int
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "setting the status from the error reporter",
			args: args{
				w:      httptest.NewRecorder(),
				err:    fmt.Errorf("an error for testing"),
				status: http.StatusBadRequest,
			},
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base.ReportErr(tt.args.w, tt.args.err, tt.args.status)

			rw, ok := tt.args.w.(*httptest.ResponseRecorder)
			assert.True(t, ok)
			assert.NotNil(t, rw)

			assert.Equal(t, tt.wantStatus, rw.Code)
		})
	}
}

func TestPostFirebaseRefreshRequest(t *testing.T) {
	ctx := context.Background()
	fc := &base.FirebaseClient{}
	loginFunc := base.GetLoginFunc(ctx, fc)
	loginCreds := &base.LoginCreds{
		Username: base.MustGetEnvVar(base.UsernameEnvVarName),
		Password: base.MustGetEnvVar(base.PasswordEnvVarName),
	}
	loginCredsJSONBytes, err := json.Marshal(loginCreds)
	assert.Nil(t, err)
	assert.NotNil(t, loginCredsJSONBytes)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginCredsJSONBytes))
	loginFunc(w, r)
	assert.Equal(t, http.StatusOK, w.Code) // login succeeded

	var loginResp base.LoginResponse
	err = json.NewDecoder(w.Result().Body).Decode(&loginResp)
	assert.Nil(t, err)
	assert.NotZero(t, loginResp.IDToken)
	assert.NotZero(t, loginResp.CustomToken)
	assert.NotZero(t, loginResp.RefreshToken)
	assert.NotZero(t, loginResp.ExpiresIn)
	assert.NotZero(t, loginResp.UID)
	assert.NotZero(t, loginResp.Email)

	goodRefreshCreds := &base.RefreshCreds{RefreshToken: loginResp.RefreshToken}
	goodRefreshURL, goodRefreshPayload := base.ComposeRefreshRequest(goodRefreshCreds)
	assert.NotNil(t, goodRefreshPayload)
	assert.NotZero(t, goodRefreshURL)

	badRefreshCreds := &base.RefreshCreds{RefreshToken: "not a real token"}
	badRefreshURL, badRefreshPayload := base.ComposeRefreshRequest(badRefreshCreds)
	assert.NotNil(t, badRefreshPayload)
	assert.NotZero(t, badRefreshURL)

	type args struct {
		w                  http.ResponseWriter
		refreshURL         string
		encodedRefreshData io.Reader
	}
	tests := []struct {
		name          string
		args          args
		wantStatus    int
		shouldSucceed bool
	}{
		{
			name: "valid refresh creds - should work",
			args: args{
				w:                  httptest.NewRecorder(),
				refreshURL:         goodRefreshURL,
				encodedRefreshData: goodRefreshPayload,
			},
			wantStatus:    http.StatusOK,
			shouldSucceed: false,
		},
		{
			name: "invalid refresh creds - should get an error",
			args: args{
				w:                  httptest.NewRecorder(),
				refreshURL:         badRefreshURL,
				encodedRefreshData: badRefreshPayload,
			},
			wantStatus:    http.StatusInternalServerError,
			shouldSucceed: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := base.PostFirebaseRefreshRequest(
				tt.args.w,
				tt.args.refreshURL,
				tt.args.encodedRefreshData,
			)
			if tt.shouldSucceed {
				assert.Equal(t, tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestNewERPClient(t *testing.T) {
	got, err := base.NewERPClient()
	assert.Nil(t, err)
	assert.NotNil(t, got)
}

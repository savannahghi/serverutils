package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"firebase.google.com/go/auth"
)

// authCheckFn is a function type for authorization and authentication checks
// there can be several e.g an authentication check runs first then an authorization
// check runs next if the authentication passes etc
type authCheckFn = func(
	r *http.Request,
	firebaseApp IFirebaseApp,
) (bool, map[string]string, *auth.Token)

// LoginClient returns an API client that is logged in with the supplied username and password
func LoginClient(username string, password string) (Client, error) {
	clientID, clientIDErr := GetEnvVar(ClientIDEnvVarName)
	if clientIDErr != nil {
		return nil, clientIDErr
	}

	clientSecret, clientSecretErr := GetEnvVar(ClientSecretEnvVarName)
	if clientSecretErr != nil {
		return nil, clientSecretErr
	}

	apiTokenURL, apiTokenURLErr := GetEnvVar(TokenURLEnvVarName)
	if apiTokenURLErr != nil {
		return nil, apiTokenURLErr
	}

	apiHost, apiHostErr := GetEnvVar(APIHostEnvVarName)
	if apiHostErr != nil {
		return nil, apiHostErr
	}

	apiScheme, apiSchemeErr := GetEnvVar(APISchemeEnvVarName)
	if apiSchemeErr != nil {
		return nil, apiSchemeErr
	}

	grantType, grantTypeErr := GetEnvVar(GrantTypeEnvVarName)
	if grantTypeErr != nil {
		return nil, grantTypeErr
	}
	extraHeaders := make(map[string]string)
	return NewServerClient(
		clientID, clientSecret, apiTokenURL, apiHost, apiScheme, grantType, username, password, extraHeaders)
}

// APIClient retrieves an EDI username and password from the environment,
// logs in to Slade 360 EDI and returns an initialized
//
// If any error is encountered, a nil client and error are returned.
func APIClient() (Client, error) {
	username, err := GetEnvVar(UsernameEnvVarName)
	if err != nil {
		return nil, err
	}
	password, err := GetEnvVar(PasswordEnvVarName)
	if err != nil {
		return nil, err
	}
	return LoginClient(username, password)
}

// ExtractToken extracts a token with the specified prefix from the specified header
func ExtractToken(r *http.Request, header string, prefix string) (string, error) {
	if r == nil {
		return "", fmt.Errorf("nil request")
	}
	if r.Header == nil {
		return "", fmt.Errorf("no headers, can't extract bearer token")
	}
	authHeader := r.Header.Get(header)
	if authHeader == "" {
		return "", fmt.Errorf("expected an `%s` request header", header)
	}
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("the `Authorization` header contents should start with `Bearer`")
	}
	tokenOnly := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	return tokenOnly, nil
}

// ExtractBearerToken gets a bearer token from an Authorization header.
//
// This is expected to contain a Firebase idToken prefixed with "Bearer "
func ExtractBearerToken(r *http.Request) (string, error) {
	return ExtractToken(r, "Authorization", "Bearer")
}

// ExtractXBearerToken gets a bearer token from an X-Authorization header.
//
// This is expected to contain a Slade 360 auth server access token prefixed with "Bearer "
func ExtractXBearerToken(r *http.Request) (string, error) {
	return ExtractToken(r, "X-Authorization", "Bearer")
}

// AuthenticationMiddleware decodes the share session cookie and packs the session into context
func AuthenticationMiddleware(firebaseApp IFirebaseApp) func(http.Handler) http.Handler {
	// multiple checks will be run in sequence (order matters)
	// the first check to succeed will call `c.Next()` and `return`
	// this means that more permissive checks (e.g exceptions) should come first
	checkFuncs := []authCheckFn{HasValidFirebaseBearerToken, HasValidSlade360AccessToken}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				errs := []map[string]string{}
				// in case authorization does not succeed, accumulated errors
				// are returned to the client
				for _, checkFunc := range checkFuncs {
					shouldContinue, errMap, authToken := checkFunc(r, firebaseApp)
					if shouldContinue {
						// put the auth token in the context
						ctx := context.WithValue(r.Context(), AuthTokenContextKey, authToken)

						// and call the next with our new context
						r = r.WithContext(ctx)
						next.ServeHTTP(w, r)
						return
					}
					errs = append(errs, errMap)
				}

				// if we got here, it is because we have errors.
				// write an error response)
				WriteJSONResponse(w, errs, http.StatusUnauthorized)
			},
		)
	}
}

// GetFirebaseUser logs in the user with the supplied credentials and retursn their
// Firebase auth user record
func GetFirebaseUser(ctx context.Context, creds *LoginCreds) (*auth.UserRecord, error) {
	if creds == nil {
		return nil, fmt.Errorf("nil creds, can't get firebase user")
	}
	user, err := GetOrCreateFirebaseUser(ctx, creds.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FetchUserProfile gets and returns the Slade 360 auth server profile of the
// user logged in to the supplied EDI client
func FetchUserProfile(authClient Client) (*EDIUserProfile, error) {
	meURL, err := authClient.MeURL()
	if err != nil {
		return nil, err
	}

	var userProfile EDIUserProfile
	err = ReadAuthServerRequestToTarget(authClient, "GET", meURL, "", nil, &userProfile)
	if err != nil {
		return nil, err
	}

	return &userProfile, nil
}

// AuthenticateCustomToken verifies the identity of a user on the basis of a
// Firebase custom token and writes the result to a HTTP response writer.
func AuthenticateCustomToken(customToken string, httpClient *http.Client) (*FirebaseUserTokens, error) {
	userTokens, err := AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, err
	}
	return userTokens, nil
}

// GetRefreshFunc is used to refresh OAuth tokens
func GetRefreshFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds := &RefreshCreds{}
		DecodeJSONToTargetStruct(w, r, creds)

		refreshURL, encodedRefreshData := ComposeRefreshRequest(creds)
		resp := PostFirebaseRefreshRequest(w, refreshURL, encodedRefreshData)
		if resp == nil {
			ReportErr(w, fmt.Errorf("nil response from Firebase for refresh request"), http.StatusInternalServerError)
			return
		}
		firebaseRefreshResp := DecodeRefreshResponse(w, resp)
		if firebaseRefreshResp == nil {
			ReportErr(w, fmt.Errorf("unable to decode response from Firebase for refresh request"), http.StatusInternalServerError)
			return
		}

		refreshResponse := RefreshResponse{
			ExpiresIn:    ConvertStringToInt(w, firebaseRefreshResp.ExpiresIn),
			IDToken:      firebaseRefreshResp.IDToken,
			RefreshToken: firebaseRefreshResp.RefreshToken,
			TokenType:    firebaseRefreshResp.TokenType,
		}
		WriteJSONResponse(w, refreshResponse, http.StatusOK)
	}
}

// GetLogoutFunc logs the user out of Firebase
func GetLogoutFunc(ctx context.Context, fc IFirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logoutReq := &LogoutRequest{}
		DecodeJSONToTargetStruct(w, r, logoutReq)
		LogoutFirebase(ctx, w, fc, logoutReq)
	}
}

// GetLoginFunc returns a function that can authenticate against both Slade 360 and Firebase
func GetLoginFunc(ctx context.Context, fc IFirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds, err := ValidateLoginCreds(w, r)
		if err != nil {
			ReportErr(w, err, http.StatusBadRequest)
			return
		}

		ediClient, err := LoginClient(creds.Username, creds.Password)
		if err != nil {
			ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		// having gotten here (user logged into Slade 360 successfully), we go nuclear
		// an panic for any other errors. Any errors that occur downstream from here are
		// symptomatic of a very big external problem. We are happy to crash and set off
		// alarms. Dead programs tell tales!
		firebaseUser, err := GetFirebaseUser(ctx, creds)
		if err != nil {
			ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		customToken, err := CreateFirebaseCustomToken(ctx, firebaseUser.UID)
		if err != nil {
			ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		userTokens, err := AuthenticateCustomToken(customToken, ediClient.HTTPClient())
		if err != nil {
			ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		userProfile, err := FetchUserProfile(ediClient)
		if err != nil {
			ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		loginResp := LoginResponse{
			CustomToken:   customToken,
			Scope:         ediClient.AccessScope(),
			ExpiresIn:     ConvertStringToInt(w, userTokens.ExpiresIn),
			IDToken:       userTokens.IDToken,
			RefreshToken:  userTokens.RefreshToken,
			TokenType:     ediClient.TokenType(),
			UserProfile:   userProfile,
			UID:           firebaseUser.UID,
			Email:         firebaseUser.Email,
			DisplayName:   firebaseUser.DisplayName,
			EmailVerified: firebaseUser.EmailVerified,
			PhoneNumber:   firebaseUser.PhoneNumber,
			PhotoURL:      firebaseUser.PhotoURL,
			Disabled:      firebaseUser.Disabled,
			TenantID:      firebaseUser.TenantID,
			ProviderID:    firebaseUser.ProviderID,
		}
		WriteJSONResponse(w, loginResp, http.StatusOK)
	}
}

// ValidateAccessToken determines if a Slade 360 auth server token is valid when checked against
// a Slade 360 auth server that is configured via environment variables.
func ValidateAccessToken(accessToken string) (bool, *LoginCreds, error) {
	if accessToken == "" {
		msg := "unable to get `accessToken` from the input"
		return false, nil, fmt.Errorf(msg)
	}
	// We'll need a configured EDI client later.
	// Here, it is configured using values from the environment.
	username, err := GetEnvVar(UsernameEnvVarName)
	if err != nil {
		return false, nil, err
	}

	password, err := GetEnvVar(PasswordEnvVarName)
	if err != nil {
		return false, nil, err
	}

	creds := &LoginCreds{
		Username: username,
		Password: password,
	}

	ediClient, err := LoginClient(creds.Username, creds.Password)
	if err != nil {
		return false, creds, err
	}

	url, err := GetEnvVar(introspectionURLEnvVarName)
	if err != nil {
		return false, creds, err
	}
	log.Printf("Access token URL: %s", url)

	payload := map[string]string{
		"token":      accessToken,
		"token_type": "access_token",
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return false, creds, err
	}

	httpClient := ediClient.HTTPClient()
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(payloadJSON))
	if err != nil {
		return false, creds, err
	}

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		msg := fmt.Sprintf("error from token endpoint, status %d", resp.StatusCode)
		return false, creds, fmt.Errorf(msg)
	}

	var respDict map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respDict)
	if err != nil {
		return false, creds, err
	}

	isValid, ok := respDict["is_valid"]
	if !ok {
		msg := "expected is_valid key in the response, did not find one"
		return false, creds, fmt.Errorf(msg)
	}

	isValidBool, ok := isValid.(bool)
	if !ok {
		msg := fmt.Sprintf(
			"expected is_valid to be a bool, got %T", isValid)
		return false, creds, fmt.Errorf(msg)
	}

	if !isValidBool {
		msg := "the supplied access token is not valid"
		return false, creds, fmt.Errorf(msg)
	}

	return true, creds, nil
}

// GetDefaultUser returns the service's default user, configured via env vars
func GetDefaultUser(
	ctx context.Context,
	creds *LoginCreds,
) (
	*FirebaseUserTokens,
	*EDIUserProfile,
	error,
) {
	firebaseUser, err := GetFirebaseUser(ctx, creds)
	if err != nil {
		return nil, nil, err
	}

	customToken, err := CreateFirebaseCustomToken(ctx, firebaseUser.UID)
	if err != nil {
		return nil, nil, err
	}

	ediClient, err := APIClient()
	if err != nil {
		return nil, nil, err
	}

	userTokens, err := AuthenticateCustomToken(customToken, ediClient.HTTPClient())
	if err != nil {
		return nil, nil, err
	}

	authServerProfile, err := FetchUserProfile(ediClient)
	if err != nil {
		return nil, nil, err
	}
	return userTokens, authServerProfile, nil
}

// GetVerifyTokenFunc confirms that an EDI access token (supplied) is valid.
//
// If it is valid, it exchanges it for a Firebase ID token.
func GetVerifyTokenFunc(ctx context.Context, fc IFirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We expect the user to supply JSON that includes a valid Slade 360 auth server access token
		// under the key `accessToken`. If this key is not supplied, we cannot continue.
		inp := &AccessTokenPayload{}
		DecodeJSONToTargetStruct(w, r, inp)

		isValid, creds, err := ValidateAccessToken(inp.AccessToken)
		if err != nil {
			ReportErr(w, err, http.StatusBadRequest)
			return
		}
		if !isValid {
			ReportErr(w, fmt.Errorf("the supplied access token is not valid"), http.StatusBadRequest)
			return
		}

		firebaseUser, err := GetFirebaseUser(ctx, creds)
		if err != nil {
			ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		customToken, err := CreateFirebaseCustomToken(ctx, firebaseUser.UID)
		if err != nil {
			ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		userTokens, authServerProfile, err := GetDefaultUser(ctx, creds)
		if err != nil {
			ReportErr(w, err, http.StatusInternalServerError)
			return
		}

		loginResp := LoginResponse{
			CustomToken:   customToken,
			ExpiresIn:     ConvertStringToInt(w, userTokens.ExpiresIn),
			IDToken:       userTokens.IDToken,
			RefreshToken:  userTokens.RefreshToken,
			TokenType:     "Bearer",
			UserProfile:   authServerProfile,
			UID:           firebaseUser.UID,
			Email:         firebaseUser.Email,
			DisplayName:   firebaseUser.DisplayName,
			EmailVerified: firebaseUser.EmailVerified,
			PhoneNumber:   firebaseUser.PhoneNumber,
			PhotoURL:      firebaseUser.PhotoURL,
			Disabled:      firebaseUser.Disabled,
			TenantID:      firebaseUser.TenantID,
			ProviderID:    firebaseUser.ProviderID,
		}
		WriteJSONResponse(w, loginResp, http.StatusOK)
	}
}

// DecodeRefreshResponse reads from a HTTP token refresh response and decodes into a struct
func DecodeRefreshResponse(w http.ResponseWriter, resp *http.Response) *FirebaseRefreshResponse {
	var refreshResp FirebaseRefreshResponse
	err := json.NewDecoder(resp.Body).Decode(&refreshResp)
	if err != nil {
		ReportErr(w, err, http.StatusBadRequest)
		return nil
	}
	return &refreshResp
}

// ComposeRefreshRequest composes a URL and params to request a token refresh with
func ComposeRefreshRequest(creds *RefreshCreds) (string, io.Reader) {
	key, err := GetEnvVar(FirebaseWebAPIKeyEnvVarName)
	if err != nil {
		log.Panic(err)
	}
	refreshURL := FirebaseRefreshTokenURL + key
	refreshData := url.Values{}
	refreshData.Set("grant_type", "refresh_token")
	refreshData.Set("refresh_token", creds.RefreshToken)
	encodedRefreshData := strings.NewReader(refreshData.Encode())
	return refreshURL, encodedRefreshData
}

// PostFirebaseRefreshRequest sends a refresh request to the relevant auth server
func PostFirebaseRefreshRequest(
	w http.ResponseWriter,
	refreshURL string,
	encodedRefreshData io.Reader,
) *http.Response {
	httpClient := http.Client{
		Timeout: time.Second * HTTPClientTimeoutSecs,
	}
	resp, err := httpClient.Post(
		refreshURL,
		"application/x-www-form-urlencoded",
		encodedRefreshData,
	)
	if err != nil {
		ReportErr(w, err, http.StatusInternalServerError)
		return nil
	}
	if resp != nil && resp.StatusCode > http.StatusPartialContent {
		rawResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Unable to read Firebase refresh request HTTP response")
		}
		log.Printf("Firebase refresh request raw response: %s", string(rawResp))
		ReportErr(w, fmt.Errorf("firebase auth refresh error: %d", resp.StatusCode), http.StatusInternalServerError)
		return nil
	}
	return resp
}

// LogoutFirebase logs out a Firebase session
func LogoutFirebase(ctx context.Context, w http.ResponseWriter, fc IFirebaseClient, logoutReq *LogoutRequest) {
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		ReportErr(w, err, http.StatusInternalServerError)
		return
	}

	authClient, err := firebaseApp.Auth(ctx)
	if err != nil {
		ReportErr(w, err, http.StatusInternalServerError)
		return
	}

	err = authClient.RevokeRefreshTokens(ctx, logoutReq.UID)
	if err != nil {
		ReportErr(w, err, http.StatusInternalServerError)
		return
	}
}

// HasValidFirebaseBearerToken returns true with no errors if the request has a valid bearer token in the authorization header.
// Otherwise, it returns false and the error in a map with the key "error"
func HasValidFirebaseBearerToken(r *http.Request, firebaseApp IFirebaseApp) (bool, map[string]string, *auth.Token) {
	bearerToken, err := ExtractBearerToken(r)
	if err != nil {
		// this error here will only be returned to the user if all the verification functions in the chain fail
		return false, ErrorMap(err), nil
	}

	validToken, err := ValidateBearerToken(r.Context(), bearerToken)
	if err != nil {
		return false, ErrorMap(err), nil
	}

	return true, nil, validToken
}

// HasValidSlade360AccessToken checks that a request has a valid Slade 360 auth server X-Authorization bearer token.
//
// The header name should be "X-Authorization". The token should be prefixed with "Bearer ".
func HasValidSlade360AccessToken(
	r *http.Request,
	firebaseApp IFirebaseApp,
) (bool, map[string]string, *auth.Token) {
	accessToken, err := ExtractXBearerToken(r)
	if err != nil {
		// this error here will only be returned to the user if all the verification functions in the chain fail
		return false, ErrorMap(err), nil
	}
	isValid, creds, err := ValidateAccessToken(accessToken)
	if err != nil {
		// this error here will only be returned to the user if all the verification functions in the chain fail
		return false, ErrorMap(err), nil
	}
	if !isValid {
		return false, ErrorMap(fmt.Errorf("the supplied access token is not valid")), nil
	}
	userTokens, _, err := GetDefaultUser(r.Context(), creds)
	if err != nil {
		return false, ErrorMap(err), nil
	}

	auth, err := firebaseApp.Auth(r.Context())
	if err != nil {
		return false, ErrorMap(err), nil
	}

	authToken, err := auth.VerifyIDToken(r.Context(), userTokens.IDToken)
	if err != nil {
		return false, ErrorMap(err), nil
	}

	return true, nil, authToken
}

// ValidateLoginCreds checks that the credentials supplied in the indicated request are valid
func ValidateLoginCreds(w http.ResponseWriter, r *http.Request) (*LoginCreds, error) {
	creds := &LoginCreds{}
	DecodeJSONToTargetStruct(w, r, creds)
	if creds.Username == "" || creds.Password == "" {
		err := fmt.Errorf("invalid credentials, expected a username AND password")
		ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return creds, nil
}

// ReportErr writes the indicated error to supplied response writer and also logs it
func ReportErr(w http.ResponseWriter, err error, status int) {
	log.Printf("%s", err)
	WriteJSONResponse(w, ErrorMap(err), status)
}

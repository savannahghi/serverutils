package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/profiler"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"firebase.google.com/go/auth"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// LoginCreds is used to (de)serialize the login username and password
type LoginCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse is used to (de)serialize the result of a successful login
type LoginResponse struct {
	CustomToken   string          `json:"custom_token"`
	Scope         string          `json:"scope"`
	ExpiresIn     int             `json:"expires_in"`
	IDToken       string          `json:"id_token"`
	RefreshToken  string          `json:"refresh_token"`
	TokenType     string          `json:"token_type"`
	UserProfile   *EDIUserProfile `json:"user_profile"`
	UID           string          `json:"uid"`
	Email         string          `json:"email"`
	DisplayName   string          `json:"display_name"`
	EmailVerified bool            `json:"email_verified"`
	PhoneNumber   string          `json:"phone_number"`
	PhotoURL      string          `json:"photo_url"`
	Disabled      bool            `json:"disabled"`
	TenantID      string          `json:"tenant_id"`
	ProviderID    string          `json:"provider_id"`
	Setup         *setupProcess   `json:"setup,omitempty"`
}

type refreshCreds struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	UID string `json:"uid"`
}

// EDIUserProfile is used to (de)serialialize the Slade 360 auth server
// profile of the logged in user.
type EDIUserProfile struct {
	ID              int      `json:"id"`
	GUID            string   `json:"guid"`
	Email           string   `json:"email"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	OtherNames      string   `json:"other_names"`
	IsStaff         bool     `json:"is_staff"`
	IsActive        bool     `json:"is_active"`
	Organisation    int      `json:"organisation"`
	BusinessPartner string   `json:"business_partner"`
	Roles           []string `json:"roles"`
	BPType          string   `json:"bp_type"`
}

type setupProcess struct {
	Progress        int           `json:"progress"`
	CompletedSteps  []interface{} `json:"completedSteps"`
	IncompleteSteps []interface{} `json:"incompleteSteps"`
	NextStep        interface{}   `json:"nextStep"`
}

type refreshResponse struct {
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type firebaseRefreshResponse struct {
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
	ProjectID    string `json:"project_id"`
}

// Sentry initializes Sentry, for error reporting
func Sentry() error {
	dsn, err := GetEnvVar(DSNEnvVarName)
	if err != nil {
		return err
	}
	return sentry.Init(sentry.ClientOptions{Dsn: dsn})
}

// ListenAddress determines what port to listen on and falls back to a default
// port if the environment does not supply a port
func ListenAddress() string {
	port := os.Getenv(PortEnvVarName)
	if port == "" {
		port = DefaultPort
	}
	address := fmt.Sprintf(":%s", port)
	return address
}

// ExtractBearerToken gets a bearer token from an Authorization header
func ExtractBearerToken(r *http.Request) (string, error) {
	if r == nil {
		return "", fmt.Errorf("nil request")
	}
	if r.Header == nil {
		return "", fmt.Errorf("no headers, can't extract bearer token")
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("expected an `Authorization` request header")
	}
	if !strings.HasPrefix(authHeader, BearerTokenPrefix) {
		return "", fmt.Errorf("the `Authorization` header contents should start with `Bearer`")
	}
	tokenOnly := strings.TrimPrefix(authHeader, BearerTokenPrefix)
	return tokenOnly, nil
}

// ErrorMap turns the supplied error into a map with "error" as the key
func ErrorMap(err error) map[string]string {
	errMap := make(map[string]string)
	errMap["error"] = err.Error()
	return errMap
}

// authCheckFn is a function type for authorization and authentication checks
// there can be several e.g an authentication check runs first then an authorization
// check runs next if the authentication passes etc
type authCheckFn = func(r *http.Request, firebaseApp IFirebaseApp) (bool, map[string]string, *auth.Token)

// hasValidFirebaseBearerToken returns true with no errors if the request has a valid bearer token in the authorization header.
// Otherwise, it returns false and the error in a map with the key "error"
func hasValidFirebaseBearerToken(r *http.Request, firebaseApp IFirebaseApp) (bool, map[string]string, *auth.Token) {
	bearerToken, tokenExtractErr := ExtractBearerToken(r)
	if tokenExtractErr != nil {
		return false, ErrorMap(tokenExtractErr), nil
	}

	validToken, err := ValidateBearerToken(r.Context(), bearerToken, firebaseApp)
	if err != nil {
		return false, ErrorMap(err), nil
	}

	return true, nil, validToken
}

// RequestDebugMiddleware dumps the incoming HTTP request to the log for inspection
func RequestDebugMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Errorf("Unable to read request body for debugging: error %#v", err)
				}
				req, err := httputil.DumpRequest(r, true)
				if err != nil {
					log.Errorf("Unable to dump cloned request for debugging: error %#v", err)
				}
				if IsDebug() {
					log.Printf("Raw request: %v", string(req))
				}
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				next.ServeHTTP(w, r)
			},
		)
	}
}

// AuthenticationMiddleware decodes the share session cookie and packs the session into context
func AuthenticationMiddleware(firebaseApp IFirebaseApp) func(http.Handler) http.Handler {
	// multiple checks will be run in sequence (order matters)
	// the first check to succeed will call `c.Next()` and `return`
	// this means that more permissive checks (e.g exceptions) should come first
	checkFuncs := []authCheckFn{hasValidFirebaseBearerToken}

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

// LogStartupError is used to e.g log fatal startup errors.
// It logs, attempts to report the error to StackDriver then panics/crashes.
func LogStartupError(ctx context.Context, err error) {
	errorClient := StackDriver(ctx)
	if err != nil {
		if errorClient != nil {
			errorClient.Report(errorreporting.Entry{Error: err})
		}
		log.WithFields(log.Fields{"error": err}).Error("Server startup error")
	}
}

// DecodeJSONToTargetStruct maps JSON from a HTTP request to a struct.
func DecodeJSONToTargetStruct(w http.ResponseWriter, r *http.Request, targetStruct interface{}) {
	decErr := json.NewDecoder(r.Body).Decode(targetStruct)
	if decErr != nil {
		WriteJSONResponse(w, ErrorMap(decErr), http.StatusBadRequest)
		return
	}
}

// GetFirebaseUser logs in the user with the supplied credentials and retursn their
// Firebase auth user record
func GetFirebaseUser(ctx context.Context, w http.ResponseWriter, creds *LoginCreds) *auth.UserRecord {
	fc := &FirebaseClient{}
	user, uErr := GetOrCreateFirebaseUser(ctx, creds.Username, fc)
	if uErr != nil {
		WriteJSONResponse(w, ErrorMap(uErr), http.StatusInternalServerError)
		return nil
	}
	return user
}

// FetchUserProfile gets and returns the Slade 360 auth server profile of the
// user logged in to the supplied EDI client
func FetchUserProfile(w http.ResponseWriter, client Client) *EDIUserProfile {
	meURL, err := client.MeURL()
	if err != nil {
		WriteJSONResponse(w, ErrorMap(err), http.StatusInternalServerError)
		return nil
	}

	var userProfile EDIUserProfile
	err = ReadAuthServerRequestToTarget(client, "GET", meURL, "", nil, &userProfile)
	if err != nil {
		WriteJSONResponse(w, ErrorMap(err), http.StatusInternalServerError)
		return nil
	}

	return &userProfile
}

// ConvertStringToInt converts a supplied string value to an integer.
// It writes an error to the JSON response writer if the conversion fails.
func ConvertStringToInt(w http.ResponseWriter, val string) int {
	converted, convErr := strconv.Atoi(val)
	if convErr != nil {
		WriteJSONResponse(w, ErrorMap(convErr), http.StatusInternalServerError)
		return -1 // sentinel value
	}
	return converted
}

// AuthenticateCustomToken verifies the identity of a user on the basis of a
// Firebase custom token and writes the result to a HTTP response writer.
func AuthenticateCustomToken(w http.ResponseWriter, customToken string, httpClient *http.Client, fc IFirebaseClient) *FirebaseUserTokens {
	userTokens, authErr := fc.AuthenticateCustomFirebaseToken(customToken, httpClient)
	if authErr != nil {
		WriteJSONResponse(w, ErrorMap(authErr), http.StatusInternalServerError)
		return nil
	}
	return userTokens
}

func decodeRefreshResponse(w http.ResponseWriter, resp *http.Response) *firebaseRefreshResponse {
	var refreshResp firebaseRefreshResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&refreshResp)
	if decodeErr != nil {
		WriteJSONResponse(w, ErrorMap(decodeErr), http.StatusInternalServerError)
		return nil
	}
	return &refreshResp
}

func composeRefreshRequest(creds *refreshCreds) (string, io.Reader) {
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

func postRefreshRequest(w http.ResponseWriter, httpClient *http.Client, refreshURL string, encodedRefreshData io.Reader) *http.Response {
	resp, postErr := httpClient.Post(refreshURL, "application/x-www-form-urlencoded", encodedRefreshData)
	if postErr != nil {
		WriteJSONResponse(w, ErrorMap(postErr), http.StatusInternalServerError)
		return nil
	}
	if resp != nil && (resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices) {
		statusErr := fmt.Errorf("refresh auth server error: %d", resp.StatusCode)
		WriteJSONResponse(w, ErrorMap(statusErr), http.StatusInternalServerError)
		return nil
	}
	return resp
}

func logoutFirebase(ctx context.Context, w http.ResponseWriter, fc IFirebaseClient, logoutReq *logoutRequest) {
	firebaseApp, faErr := fc.InitFirebase()
	if faErr != nil {
		WriteJSONResponse(w, ErrorMap(faErr), http.StatusInternalServerError)
		return
	}

	authClient, clErr := firebaseApp.Auth(ctx)
	if clErr != nil {
		WriteJSONResponse(w, ErrorMap(clErr), http.StatusInternalServerError)
		return
	}

	revokeErr := authClient.RevokeRefreshTokens(ctx, logoutReq.UID)
	if revokeErr != nil {
		WriteJSONResponse(w, ErrorMap(revokeErr), http.StatusInternalServerError)
		return
	}
}

// GetRefreshFunc is used to refresh OAuth tokens
func GetRefreshFunc(httpClient *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds := &refreshCreds{}
		DecodeJSONToTargetStruct(w, r, creds)

		refreshURL, encodedRefreshData := composeRefreshRequest(creds)
		resp := postRefreshRequest(w, httpClient, refreshURL, encodedRefreshData)
		if resp == nil {
			msg := "nil response from Firebase for refresh request"
			log.Printf(msg)
			_, _ = w.Write([]byte(msg))
			return
		}
		firebaseRefreshResp := decodeRefreshResponse(w, resp)
		if firebaseRefreshResp == nil {
			msg := "unable to decode response from Firebase for refresh request"
			log.Printf(msg)
			_, _ = w.Write([]byte(msg))
			return
		}

		refreshResponse := refreshResponse{
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
		logoutReq := &logoutRequest{}
		DecodeJSONToTargetStruct(w, r, logoutReq)
		logoutFirebase(ctx, w, fc, logoutReq)
	}
}

// StackDriver initializes StackDriver logging, error reporting, profiling etc
func StackDriver(ctx context.Context) *errorreporting.Client {
	// project setup
	projectID, err := GetEnvVar(GoogleCloudProjectIDEnvVarName)
	if err != nil {
		log.WithFields(log.Fields{
			"environment variable name": GoogleCloudProjectIDEnvVarName,
			"error":                     err,
		}).Error("Unable to determine the Google Cloud Project, can't set up StackDriver")
		return nil
	}

	// logging
	loggingClient, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		log.WithFields(log.Fields{
			"project ID": projectID,
			"error":      err,
		}).Error("Unable to initialize logging client")
		return nil
	}
	defer closeStackDriverLoggingClient(loggingClient)

	// error reporting
	errorClient, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName: AppName,
		OnError: func(err error) {
			log.WithFields(log.Fields{
				"project ID":   projectID,
				"service name": AppName,
				"error":        err,
			}).Info("Unable to initialize error client")
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to initialize error client")
		return nil
	}
	defer closeStackDriverErrorClient(errorClient)

	// tracing
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"project ID": projectID,
			"error":      err,
		}).Info("Unable to initialize tracing")
		return errorClient // the error client is already initialized, return it
	}
	trace.RegisterExporter(exporter)

	// profiler
	err = profiler.Start(profiler.Config{
		Service:        AppName,
		ServiceVersion: AppVersion,
		ProjectID:      projectID,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"project ID":      projectID,
			"service name":    AppName,
			"service version": AppVersion,
			"error":           err,
		}).Info("Unable to initialize profiling")
		return errorClient // the error client is already initialized, return it
	}

	return errorClient
}

// WriteJSONResponse writes the content supplied via the `source` parameter to
// the supplied http ResponseWriter. The response is returned with the indicated
// status.
func WriteJSONResponse(w http.ResponseWriter, source interface{}, status int) {
	w.WriteHeader(status)
	content, err := json.Marshal(source)
	if err != nil {
		msg := fmt.Sprintf("error when marshalling %#v to JSON bytes: %#v", source, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(content)
	if err != nil {
		msg := fmt.Sprintf(
			"error when writing JSON %s to http.ResponseWriter: %#v", string(content), err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func closeStackDriverLoggingClient(loggingClient *logging.Client) {
	err := loggingClient.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Unable to close StackDriver logging client")
	}
}

func closeStackDriverErrorClient(errorClient *errorreporting.Client) {
	err := errorClient.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("Unable to close StackDriver error client")
	}
}

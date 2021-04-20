package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/imroc/req"
	"github.com/stretchr/testify/assert"
)

const (
	anonymousUserUID  = "AgkGYKUsRifO2O9fTLDuVCMr2hb2" // This is an anonymous user
	verifyPhone       = "testing/verify_phone"
	createUserByPhone = "testing/create_user_by_phone"
	loginByPhone      = "testing/login_by_phone"
	removeUserByPhone = "testing/remove_user"
	addAdmin          = "testing/add_admin_permissions"
)

// ContextKey is used as a type for the UID key for the Firebase *auth.Token on context.Context.
// It is a custom type in order to minimize context key collissions on the context
// (.and to shut up golint).
type ContextKey string

// GetAuthenticatedContext returns a logged in context, useful for test purposes
func GetAuthenticatedContext(t *testing.T) context.Context {
	ctx := context.Background()
	authToken := getAuthToken(ctx, t)
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext
}

// GetAuthenticatedContextAndToken returns a logged in context and ID token.
// It is useful for test purposes
func GetAuthenticatedContextAndToken(t *testing.T) (context.Context, *auth.Token) {
	ctx := context.Background()
	authToken := getAuthToken(ctx, t)
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext, authToken
}

// GetAuthenticatedContextAndBearerToken returns a logged in context and bearer token.
// It is useful for test purposes
func GetAuthenticatedContextAndBearerToken(t *testing.T) (context.Context, string) {
	ctx := context.Background()
	authToken, bearerToken := getAuthTokenAndBearerToken(ctx, t)
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext, bearerToken
}

func getAuthToken(ctx context.Context, t *testing.T) *auth.Token {
	authToken, _ := getAuthTokenAndBearerToken(ctx, t)
	return authToken
}

func getAuthTokenAndBearerToken(ctx context.Context, t *testing.T) (*auth.Token, string) {
	user, userErr := GetOrCreateFirebaseUser(ctx, TestUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := ValidateBearerToken(ctx, bearerToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	return authToken, bearerToken
}

// GetOrCreateAnonymousUser creates an anonymous user
// For documentation and test purposes only
func GetOrCreateAnonymousUser(ctx context.Context) (*auth.UserRecord, error) {
	authClient, err := GetFirebaseAuthClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get or create Firebase client: %w", err)
	}
	existingUser, userErr := authClient.GetUser(ctx, anonymousUserUID)

	if userErr == nil {
		return existingUser, nil
	}

	params := (&auth.UserToCreate{})
	newUser, createErr := authClient.CreateUser(ctx, params)
	if createErr != nil {
		return nil, createErr
	}
	return newUser, nil
}

// GetAnonymousContext returns an anonymous logged in context, useful for test purposes
func GetAnonymousContext(t *testing.T) context.Context {
	ctx := context.Background()
	authToken := getAnonymousAuthToken(ctx, t)
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext
}

func getAnonymousAuthToken(ctx context.Context, t *testing.T) *auth.Token {
	user, userErr := GetOrCreateAnonymousUser(ctx)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := ValidateBearerToken(ctx, bearerToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	return authToken
}

// VerifyTestPhoneNumber checks if the test `Phone Number` exists as a primary
// phone number in any user profile record
func VerifyTestPhoneNumber(
	t *testing.T,
	phone string,
	onboardingClient *InterServiceClient,
) (string, error) {
	verifyPhonePayload := map[string]interface{}{
		"phoneNumber": phone,
	}

	resp, err := onboardingClient.MakeRequest(
		http.MethodPost,
		verifyPhone,
		verifyPhonePayload,
	)

	if err != nil {
		return "", fmt.Errorf("unable to make a verify phone number request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to convert response to string: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", string(body))
	}

	var otp OtpResponse
	err = json.Unmarshal(body, &otp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	return otp.OTP, nil
}

// LoginTestPhoneUser returns user response data for a created test user allowing
// them to run and access test resources
func LoginTestPhoneUser(
	t *testing.T,
	phone string,
	PIN string,
	flavour Flavour,
	onboardingClient *InterServiceClient,
) (*UserResponse, error) {
	loginPayload := map[string]interface{}{
		"phoneNumber": phone,
		"pin":         PIN,
		"flavour":     flavour,
	}

	resp, err := onboardingClient.MakeRequest(
		http.MethodPost,
		loginByPhone,
		loginPayload,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to make a login request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to login : %s, with status code %v",
			phone,
			resp.StatusCode,
		)
	}
	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to convert response to string: %v", err)
	}

	var response *UserResponse
	err = json.Unmarshal(code, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	return response, nil
}

// AddAdminPermissions adds ADMIN permissions to our test user
func AddAdminPermissions(
	t *testing.T,
	onboardingClient *InterServiceClient,
	phone string,
) error {
	phonePayload := map[string]interface{}{
		"phoneNumber": phone,
	}

	resp, err := onboardingClient.MakeRequest(
		http.MethodPost,
		addAdmin,
		phonePayload,
	)

	if err != nil {
		return fmt.Errorf("unable to make add admin request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to convert response to string: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %v with resp body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// CreateOrLoginTestPhoneNumberUser creates an phone number test user if they
// do not exist or `Logs them in` if the test user exists to retrieve
// authenticated user response
// For documentation and test purposes only
func CreateOrLoginTestPhoneNumberUser(
	t *testing.T,
	onboardingClient *InterServiceClient,
) (*UserResponse, error) {
	phone := TestUserPhoneNumber
	PIN := TestUserPin
	flavour := FlavourConsumer

	if onboardingClient == nil {
		return nil, fmt.Errorf("nil ISC client")
	}

	otp, err := VerifyTestPhoneNumber(t, phone, onboardingClient)
	if err != nil {
		if strings.Contains(
			err.Error(),
			strconv.Itoa(int(PhoneNumberInUse)),
		) {
			// enforce perms
			err = AddAdminPermissions(t, onboardingClient, phone)
			if err != nil {
				return nil, fmt.Errorf("unable to add admin permissions: %v", err)
			}

			return LoginTestPhoneUser(
				t,
				phone,
				PIN,
				flavour,
				onboardingClient,
			)
		}
		return nil, fmt.Errorf("failed to verify test phone number: %v", err)
	}
	createUserPayload := map[string]interface{}{
		"phoneNumber": phone,
		"pin":         PIN,
		"flavour":     flavour,
		"otp":         otp,
	}

	resp, err := onboardingClient.MakeRequest(
		http.MethodPost,
		createUserByPhone,
		createUserPayload,
	)

	if err != nil {
		return nil, fmt.Errorf("unable to make a sign up request: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unable to sign up : %s, with status code %v",
			phone,
			resp.StatusCode,
		)
	}
	signUpResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to convert response to string: %v", err)
	}

	var response *UserResponse
	err = json.Unmarshal(signUpResp, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	err = AddAdminPermissions(t, onboardingClient, phone)
	if err != nil {
		return nil, fmt.Errorf("unable to add admin permissions: %v", err)
	}

	return response, nil

}

// GetPhoneNumberAuthenticatedContextAndToken returns a phone number logged in context
// and an auth Token that contains the the test user UID useful for test purposes
func GetPhoneNumberAuthenticatedContextAndToken(
	t *testing.T,
	onboardingClient *InterServiceClient,
) (context.Context, *auth.Token, error) {
	ctx := context.Background()
	userResponse, err := CreateOrLoginTestPhoneNumberUser(t, onboardingClient)
	if err != nil {
		return nil, nil, err
	}
	authToken := &auth.Token{
		UID: userResponse.Auth.UID,
	}
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext, authToken, nil
}

// GetDefaultHeaders returns headers used in inter service communication acceptance tests
func GetDefaultHeaders(t *testing.T, rootDomain string, serviceName string) map[string]string {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": GetInterserviceBearerTokenHeader(t, rootDomain, serviceName),
	}
}

// GetInterserviceBearerTokenHeader returns a valid isc bearer token header
func GetInterserviceBearerTokenHeader(t *testing.T, rootDomain string, serviceName string) string {
	isc := GetInterserviceClient(t, rootDomain, serviceName)
	authToken, err := isc.CreateAuthToken()
	assert.Nil(t, err)
	assert.NotZero(t, authToken)
	bearerHeader := fmt.Sprintf("Bearer %s", authToken)
	return bearerHeader
}

// GetInterserviceClient returns an isc client used in acceptance testing
func GetInterserviceClient(t *testing.T, rootDomain string, serviceName string) *InterServiceClient {
	service := ISCService{
		Name:       serviceName,
		RootDomain: rootDomain,
	}
	isc, err := NewInterserviceClient(service)
	assert.Nil(t, err)
	assert.NotNil(t, isc)
	return isc
}

// GetAuthenticatedContextFromUID creates an auth.Token given a valid uid
func GetAuthenticatedContextFromUID(ctx context.Context, uid string) (*auth.Token, error) {
	customToken, err := CreateFirebaseCustomToken(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to create an authenticated token: %w", err)
	}

	idTokens, err := AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticated custom token: %w", err)
	}

	authToken, err := ValidateBearerToken(ctx, idTokens.IDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to validate bearer token: %w", err)
	}

	return authToken, nil
}

// RemoveTestPhoneNumberUser removes the records created by the
// test phonenumber user
func RemoveTestPhoneNumberUser(
	t *testing.T,
	onboardingClient *InterServiceClient,
) error {
	if onboardingClient == nil {
		return fmt.Errorf("nil ISC client")
	}

	payload := map[string]interface{}{
		"phoneNumber": TestUserPhoneNumber,
	}
	resp, err := onboardingClient.MakeRequest(
		http.MethodPost,
		removeUserByPhone,
		payload,
	)
	if err != nil {
		return fmt.Errorf("unable to make a request to remove test user: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil // This is a test utility. Do not block if the user is not found
	}

	return nil
}

package base

import (
	"context"
	"fmt"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/imroc/req"
	"github.com/stretchr/testify/assert"
)

const anonymousUserUID = "AgkGYKUsRifO2O9fTLDuVCMr2hb2" // This is an anonymous user

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

// GetOrCreatePhoneNumberUser creates an phone number user
// For documentation and test purposes only
func GetOrCreatePhoneNumberUser(ctx context.Context, msisdn string) (*auth.UserRecord, error) {
	authClient, err := GetFirebaseAuthClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get or create Firebase client: %w", err)
	}
	existingUser, userErr := authClient.GetUserByPhoneNumber(ctx, msisdn)

	if userErr == nil {
		return existingUser, nil
	}

	params := (&auth.UserToCreate{}).
		PhoneNumber(msisdn)
	newUser, createErr := authClient.CreateUser(ctx, params)
	if createErr != nil {
		return nil, createErr
	}
	return newUser, nil
}

// GetPhoneNumberAuthenticatedContext returns a phone number logged in context, useful for test purposes
func GetPhoneNumberAuthenticatedContext(t *testing.T) context.Context {
	ctx := context.Background()
	authToken := getPhoneNumberAuthToken(ctx, t)
	authenticatedContext := context.WithValue(ctx, AuthTokenContextKey, authToken)
	return authenticatedContext
}

func getPhoneNumberAuthToken(ctx context.Context, t *testing.T) *auth.Token {
	user, userErr := GetOrCreatePhoneNumberUser(ctx, TestUserPhoneNumber)
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

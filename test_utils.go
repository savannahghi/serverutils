package base

import (
	"context"
	"net/http"
	"testing"
	"time"

	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
)

// ContextKey is used as a type for the UID key for the Firebase *auth.Token on context.Context.
// It is a custom type in order to minimize context key collissions on the context
// (.and to shut up golint).
type ContextKey string

// AuthTokenContextKey is used to add/retrieve the Firebase UID on the context
const AuthTokenContextKey = ContextKey("UID")

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

func getAuthToken(ctx context.Context, t *testing.T) *auth.Token {
	fc := &FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)

	user, userErr := GetOrCreateFirebaseUser(ctx, TestUserEmail, fc)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := CreateFirebaseCustomToken(ctx, user.UID, fc)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	client := &http.Client{Timeout: time.Second * 10}
	idTokens, idErr := fc.AuthenticateCustomFirebaseToken(customToken, client)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := ValidateBearerToken(ctx, bearerToken, firebaseApp)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	return authToken
}
package base

import (
	"context"
	"fmt"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func getJSONGoogleApplicationCredentials() ([]byte, error) {
	// a future iteration of this needs to decrypt GPG encoded creds using the
	// secret pass phrase
	return GPGEncryptedJSONGoogleApplicationCredentials, nil
}

// GetGSUITEDelegatedAuthorityTokenSource gets a token source to be used in Google Cloud APIs that
// require impersonation of a user e.g Google Calendar.
//
// This uses G-Suite domain wide delegation. See https://developers.google.com/calendar/auth#perform-g-suite-domain-wide-delegation-of-authority .
func GetGSUITEDelegatedAuthorityTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	jsonCredentials, err := getJSONGoogleApplicationCredentials()
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(
		jsonCredentials,
		calendar.CalendarScope,
	)
	if err != nil {
		return nil, fmt.Errorf("JWTConfigFromJSON: %v", err)
	}
	config.Subject = DefaultCalendarEmail

	ts := config.TokenSource(ctx)
	return ts, nil
}

// GetLoggedInUserUID retrieves the logged in user's Firebase UID from the
// supplied context and returns an error if it does not succeed
func GetLoggedInUserUID(ctx context.Context) (string, error) {
	authToken, err := GetUserTokenFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("auth token not found in context: %w", err)
	}
	return authToken.UID, nil
}

// MustGetLoggedInUserUID retrieves the logged in user's Firebase UID from the
// supplied context and panics if it does not succeed
func MustGetLoggedInUserUID(ctx context.Context) string {
	authToken, err := GetUserTokenFromContext(ctx)
	if err != nil {
		log.Panicf("unable to get auth token from context: %s", err)
	}
	return authToken.UID
}

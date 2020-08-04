package base

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// DefaultCalendarEmail is the email address that "owns" the calendar by default
const DefaultCalendarEmail = "be.well@healthcloud.co.ke"

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

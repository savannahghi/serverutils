package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"firebase.google.com/go/messaging"
	"github.com/lithammer/shortuuid"
	"google.golang.org/api/option"
)

/* #nosec */
const (
	// FirebaseWebAPIKeyEnvVarName is the name of the env var that holds a Firebase web API key
	// for this project
	FirebaseWebAPIKeyEnvVarName = "FIREBASE_WEB_API_KEY"

	// FirebaseCustomTokenSigninURL is the Google Identity Toolkit API for signing in over REST
	FirebaseCustomTokenSigninURL = "https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key="

	// FirebaseRefreshTokenURL is used to request Firebase refresh tokens from Google APIs
	FirebaseRefreshTokenURL = "https://securetoken.googleapis.com/v1/token?key="

	// GoogleApplicationCredentialsEnvVarName is used to obtain service account details from the
	// local server when necessary e.g when running tests on CI or a local developer setup
	GoogleApplicationCredentialsEnvVarName = "GOOGLE_APPLICATION_CREDENTIALS"

	// TestUserEmail is used by integration tests
	TestUserEmail = "automated.test.user.bewell-app-ci@healthcloud.co.ke"

	// OTPCollectionName is the name of the collection used to persist single
	// use verification codes on Firebase
	OTPCollectionName = "otps"

	// IdentifierCollectionName is used to record randomly generated identifiers so that they
	// are not re-issued
	IdentifierCollectionName = "identifiers"
)

// FirebaseTokenExchangePayload is marshalled into JSON and sent to the Firebase Auth REST API
// when exchanging a custom token for an ID token that can be used to make API calls
type FirebaseTokenExchangePayload struct {
	Token             string `json:"token"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

// FirebaseUserTokens is the unmarshalling target for the JSON response received from the Firebase Auth REST API
// when exchanging a custom token for an ID token that can be used to make API calls
type FirebaseUserTokens struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

// IFirebaseApp is an interface that has been extracted in order to support mocking of Firebase auth in tests
type IFirebaseApp interface {
	Auth(ctx context.Context) (*auth.Client, error)
	Firestore(ctx context.Context) (*firestore.Client, error)
	Messaging(ctx context.Context) (*messaging.Client, error)
}

// IFirebaseClient defines the Firebase methods that we depend on
// It has been defined in order to facilitate mocking for tests
type IFirebaseClient interface {
	InitFirebase() (IFirebaseApp, error)
	AuthenticateCustomFirebaseToken(customAuthToken string, client *http.Client) (*FirebaseUserTokens, error)
}

// FirebaseClient is an implementation of the FirebaseClient interface
type FirebaseClient struct{}

// InitFirebase ensures that we have a working Firebase configuration
func (fc *FirebaseClient) InitFirebase() (IFirebaseApp, error) {
	appCreds, err := GetEnvVar(GoogleApplicationCredentialsEnvVarName)
	if err != nil {
		return firebase.NewApp(
			context.Background(),
			nil,
			option.WithCredentialsFile(appCreds),
		)
	}
	return firebase.NewApp(
		context.Background(),
		nil,
	)
}

// AuthenticateCustomFirebaseToken takes a custom Firebase auth token and tries to fetch an ID token
// If successful, a pointer to the ID token is returned
// Otherwise, an error is returned
func (fc *FirebaseClient) AuthenticateCustomFirebaseToken(customAuthToken string, httpClient *http.Client) (*FirebaseUserTokens, error) {
	apiKey, apiKeyErr := GetEnvVar(FirebaseWebAPIKeyEnvVarName)
	if apiKeyErr != nil {
		return nil, apiKeyErr
	}

	payload := FirebaseTokenExchangePayload{
		Token:             customAuthToken,
		ReturnSecureToken: true,
	}
	payloadBytes, _ := json.Marshal(payload) // err intentionally ignored, static typing makes it very hard to get this error

	url := FirebaseCustomTokenSigninURL + apiKey
	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(payloadBytes))
	defer CloseRespBody(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bs, err := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"firebase HTTP error, status code %d\nBody: %s\nBody read error: %s", resp.StatusCode, string(bs), err)
	}
	var tokenResp FirebaseUserTokens
	unmarshalErr := json.NewDecoder(resp.Body).Decode(&tokenResp)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return &tokenResp, nil
}

// CreateFirebaseCustomToken creates a custom auth token for the user with the
// indicated UID
func CreateFirebaseCustomToken(ctx context.Context, uid string, fc IFirebaseClient) (string, error) {
	authClient, err := GetFirebaseAuthClient(ctx, fc)
	if err != nil {
		return "", fmt.Errorf("unable to create custom Firebase token: %w", err)
	}
	return authClient.CustomToken(ctx, uid)
}

// GetOrCreateFirebaseUser retrieves the user record of the user with the given email
// or creates a new one if no user has the specified email
func GetOrCreateFirebaseUser(ctx context.Context, email string, fc IFirebaseClient) (*auth.UserRecord, error) {
	authClient, err := GetFirebaseAuthClient(ctx, fc)
	if err != nil {
		return nil, fmt.Errorf("unable to get or create Firebase user: %w", err)
	}
	existingUser, userErr := authClient.GetUserByEmail(ctx, email)
	if userErr == nil {
		return existingUser, nil
	}

	// try creating, assume the user could not be found
	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(false).
		Disabled(false)
	newUser, createErr := authClient.CreateUser(ctx, params)
	if createErr != nil {
		return nil, createErr
	}
	return newUser, nil
}

// GetFirebaseAuthClient initializes a Firebase Authentication client
func GetFirebaseAuthClient(ctx context.Context, fc IFirebaseClient) (*auth.Client, error) {
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Firebase app: %w", err)
	}
	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Firebase auth client: %w", err)
	}
	return client, nil
}

// ValidateBearerToken checks the bearer token for validity against Firebase
func ValidateBearerToken(ctx context.Context, token string, firebaseApp IFirebaseApp) (*auth.Token, error) {
	client, clientErr := firebaseApp.Auth(ctx)
	if clientErr != nil {
		return nil, fmt.Errorf("error getting Auth client: " + clientErr.Error())
	}
	verifiedToken, verifyErr := client.VerifyIDToken(ctx, token)
	if verifyErr != nil {
		return nil, fmt.Errorf("invalid auth token: " + verifyErr.Error())
	}
	return verifiedToken, nil
}

// SaveDataToFirestore takes the supplied data (which can be a map of string to
// interface{} or a struct with json/firestore tags), a collection name and an
// intialized firestore client then tries to save the data to that collection.
func SaveDataToFirestore(firestoreClient *firestore.Client, collection string,
	data interface{}) (string, error) {
	ctx := context.Background()
	docRef, _, err := firestoreClient.Collection(collection).Add(ctx, data)
	if err != nil {
		return "", err
	}
	return docRef.ID, err
}

// UpdateRecordOnFirestore takes the supplied data (which can be a map of string to
// interface{} or a struct with json/firestore tags), a collection name and an
// intialized firestore client then tries to update the data in that object
func UpdateRecordOnFirestore(
	firestoreClient *firestore.Client,
	collection string,
	id string,
	data interface{},
) error {
	ctx := context.Background()
	_, err := firestoreClient.Collection(collection).Doc(id).Set(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// GenerateSafeIdentifier generates a shortened alphanumeric identifier.
func GenerateSafeIdentifier() string {
	return shortuuid.New()
}

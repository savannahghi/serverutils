package base

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func GetIDToken(t *testing.T) string {
	ctx := context.Background()
	fc := &FirebaseClient{}
	user, err := GetOrCreateFirebaseUser(ctx, TestUserEmail, fc)
	if err != nil {
		t.Errorf("unable to create Firebase user for email %v, error %v", TestUserEmail, err)
	}

	// test custom token generation
	customToken, err := CreateFirebaseCustomToken(ctx, user.UID, fc)
	if err != nil {
		t.Errorf("unable to get custom token for %#v", user)
	}

	// test authentication of custom Firebase tokens
	client := &http.Client{Timeout: time.Second * 10}
	idTokens, err := fc.AuthenticateCustomFirebaseToken(customToken, client)
	if err != nil {
		t.Errorf("unable to exchange custom token for ID tokens, error %s", err)
	}
	if idTokens.IDToken == "" {
		t.Errorf("got blank ID token")
	}
	return idTokens.IDToken
}

func TestGetOrCreateFirebaseUser(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		email string
	}{
		{email: TestUserEmail},
		{email: fmt.Sprintf("test.bewell.%s@healthcloud.co.ke", uuid.New().String())},
	}
	fc := &FirebaseClient{}
	for _, tc := range tests {
		user, err := GetOrCreateFirebaseUser(ctx, tc.email, fc)
		if err != nil {
			t.Errorf("unable to create Firebase user for email %v, error %v", tc.email, err)
		}

		// sanity check
		if user.Email != tc.email {
			t.Errorf("expected to get back a user with email %s, got %s", tc.email, user.Email)
		}

		// test custom token generation
		customToken, err := CreateFirebaseCustomToken(ctx, user.UID, fc)
		if err != nil {
			t.Errorf("unable to get custom token for %#v", user)
		}

		// test authentication of custom Firebase tokens
		client := &http.Client{Timeout: time.Second * 10}
		idTokens, err := fc.AuthenticateCustomFirebaseToken(customToken, client)
		if err != nil {
			t.Errorf("unable to exchange custom token for ID tokens, error %s", err)
		}
		if idTokens.IDToken == "" {
			t.Errorf("got blank ID token")
		}
	}
}

func TestAuthenticateCustomFirebaseToken_Invalid_Token(t *testing.T) {
	invalidToken := uuid.New().String()
	client := &http.Client{Timeout: time.Second * 10}
	fc := FirebaseClient{}
	returnToken, err := fc.AuthenticateCustomFirebaseToken(invalidToken, client)
	assert.Errorf(t, err, "expected invalid token to fail authentication with message %s")
	var nilToken *FirebaseUserTokens
	assert.Equal(t, nilToken, returnToken)
}

func TestAuthenticateCustomFirebaseToken_Valid_Token(t *testing.T) {
	fc := &FirebaseClient{}
	ctx := context.Background()
	user, err := GetOrCreateFirebaseUser(ctx, TestUserEmail, fc)
	assert.Nilf(t, err, "unexpected user retrieval error '%s'")
	validToken, tokenErr := CreateFirebaseCustomToken(ctx, user.UID, fc)
	assert.Nilf(t, tokenErr, "unexpected custom token creation error '%s'")
	client := &http.Client{Timeout: time.Second * 10}
	idTokens, validateErr := fc.AuthenticateCustomFirebaseToken(validToken, client)
	assert.Nilf(t, validateErr, "unexpected custom token validation/exchange error '%s'")
	assert.NotNilf(t, idTokens.IDToken, "expected ID token to be non nil")
}

func TestFirebaseClient_AuthenticateCustomFirebaseToken_HTTP_Error(t *testing.T) {
	client := MockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`Error!`)),
			Header:     make(http.Header), // Must be set to non-nil value or it panics
		}, nil
	})
	fc := FirebaseClient{}
	userTokens, err := fc.AuthenticateCustomFirebaseToken("fake_token", client)
	assert.Nil(t, userTokens)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "firebase HTTP error, status code 500")
}

func TestFirebaseClient_AuthenticateCustomFirebaseToken_JSON_Error(t *testing.T) {
	client := MockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`Bad json, will not parse`)),
			Header:     make(http.Header), // Must be set to non-nil value or it panics
		}, nil
	})
	fc := FirebaseClient{}
	userTokens, err := fc.AuthenticateCustomFirebaseToken("fake_token", client)
	assert.Nil(t, userTokens)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid character 'B' looking for beginning of value")
}

func TestFirebaseClient_AuthenticateCustomFirebaseToken_Request_Error(t *testing.T) {
	client := MockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("ka-boom")
	})
	fc := FirebaseClient{}
	userTokens, err := fc.AuthenticateCustomFirebaseToken("fake_token", client)
	assert.Nil(t, userTokens)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ka-boom")
}

func TestCreateFirebaseCustomToken_Init_Err(t *testing.T) {
	mockFc := &MockFirebaseClient{MockAppInitErr: fmt.Errorf("mock init error")}
	ctx := context.Background()
	token, err := CreateFirebaseCustomToken(ctx, "uid", mockFc)
	assert.Equal(t, "", token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "mock init error")
}

func TestCreateFirebaseCustomToken_Auth_Err(t *testing.T) {
	mockFc := &MockFirebaseClient{
		MockApp: &MockFirebaseApp{
			MockAuthErr: fmt.Errorf("mock auth error"),
		},
	}
	ctx := context.Background()
	token, err := CreateFirebaseCustomToken(ctx, "uid", mockFc)
	assert.Equal(t, "", token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "mock auth error")
}

func TestGetOrCreateFirebaseUser_Init_Err(t *testing.T) {
	mockFc := &MockFirebaseClient{MockAppInitErr: fmt.Errorf("mock init error")}
	ctx := context.Background()
	ur, err := GetOrCreateFirebaseUser(ctx, "user@mail.com", mockFc)
	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "mock init error")
}

func TestGetOrCreateFirebaseUser_Auth_Err(t *testing.T) {
	mockFc := &MockFirebaseClient{
		MockApp: &MockFirebaseApp{
			MockAuthErr: fmt.Errorf("mock auth error"),
		},
	}
	ctx := context.Background()
	ur, err := GetOrCreateFirebaseUser(ctx, "user@mail.com", mockFc)
	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "mock auth error")
}

func TestGetOrCreateFirebaseUser_Invalid_Email(t *testing.T) {
	fc := &FirebaseClient{}
	ctx := context.Background()
	ur, err := GetOrCreateFirebaseUser(ctx, "invalid_email_should_blow_up", fc)
	assert.Nil(t, ur)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "malformed email string: ")
}

func TestValidateBearerToken_Client_Err(t *testing.T) {
	auth := &MockFirebaseApp{MockAuthErr: fmt.Errorf("mock error from initializing firebase auth")}
	ctx := context.Background()
	token, err := ValidateBearerToken(ctx, "invalid token", auth)
	assert.Nil(t, token)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "error getting Auth client: mock error from initializing firebase auth")
}

func TestGenerateSafeIdentifier(t *testing.T) {
	id := GenerateSafeIdentifier()
	assert.NotZero(t, id)
}

func TestUpdateRecordOnFirestore(t *testing.T) {
	firestoreClient := getFirestoreClient(t)
	collection := "integration_test_collection"
	data := map[string]string{
		"a_key_for_testing": uuid.New().String(),
	}
	id, err := SaveDataToFirestore(firestoreClient, collection, data)
	assert.Nil(t, err)
	assert.NotZero(t, id)

	type args struct {
		firestoreClient *firestore.Client
		collection      string
		id              string
		data            interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				firestoreClient: firestoreClient,
				collection:      collection,
				id:              id,
				data: map[string]string{
					"a_key_for_testing": uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "invalid id",
			args: args{
				firestoreClient: firestoreClient,
				collection:      collection,
				id:              "this is a fake id that should not be found",
				data: map[string]string{
					"a_key_for_testing": uuid.New().String(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateRecordOnFirestore(tt.args.firestoreClient, tt.args.collection, tt.args.id, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRecordOnFirestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserTokenFromContext(t *testing.T) {
	authenticatedContext, authToken := GetAuthenticatedContextAndToken(t)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *auth.Token
		wantErr bool
	}{
		{
			name: "good case - authenticated context",
			args: args{
				ctx: authenticatedContext,
			},
			want:    authToken,
			wantErr: false,
		},
		{
			name: "unauthenticated context",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserTokenFromContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserTokenFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserTokenFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

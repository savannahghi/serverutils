package base_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestGetOrCreateAnonymousUser(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Anonymous user happy case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetOrCreateAnonymousUser(tt.args.ctx)
			assert.NotNil(t, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAnonymousUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAuthenticatedContextFromUID(t *testing.T) {
	ctx := context.Background()

	// create a valid uid

	type args struct {
		uid string
	}
	tests := []struct {
		name      string
		args      args
		changeEnv bool
		wantErr   bool
	}{
		{
			name: "valid case",
			args: args{
				uid: "some invalid uid",
			},
			changeEnv: false,
			wantErr:   false,
		},
		{
			name: "invalid: wrong uid",
			args: args{
				uid: "some invalid uid",
			},
			changeEnv: true,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			initialKey := os.Getenv("FIREBASE_WEB_API_KEY")

			if tt.changeEnv {
				os.Setenv("FIREBASE_WEB_API_KEY", "invalidkey")
			}

			got, err := base.GetAuthenticatedContextFromUID(ctx, tt.args.uid)
			if got == nil && !tt.wantErr {
				t.Errorf("invalid auth token")
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAuthenticatedContextFromUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			os.Setenv("FIREBASE_WEB_API_KEY", initialKey)

		})
	}
}

func onboardingISCClient() (*base.InterServiceClient, error) {
	onboardingClient, err := base.NewInterserviceClient(
		base.ISCService{
			Name:       base.OnboardingName,
			RootDomain: base.OnboardingRootDomain,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize onboarding ISC client: %v", err)
	}
	return onboardingClient, nil
}

func TestVerifyTestPhoneNumber(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
	}
	type args struct {
		phone            string
		onboardingClient *base.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: verify a phone number does not exist",
			args: args{
				phone:            "+254711999888",
				onboardingClient: onboardingClient,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := base.VerifyTestPhoneNumber(
				t,
				tt.args.phone,
				tt.args.onboardingClient,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyTestPhoneNumber() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if tt.wantErr && otp != "" {
				t.Errorf("expected no otp to be sent but got %v, since the error %v occurred",
					otp,
					err,
				)
				return
			}

			if !tt.wantErr && otp == "" {
				t.Errorf("expected an otp to be sent, since no error occurred")
				return
			}
		})
	}
}

func TestCreateOrLoginTestPhoneNumberUser(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
	}
	type args struct {
		onboardingClient *base.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: create a test user successfully",
			args: args{
				onboardingClient: onboardingClient,
			},
			wantErr: false,
		},
		{
			name: "failure: failed to create a test user successfully",
			args: args{
				onboardingClient: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userResponse, err := base.CreateOrLoginTestPhoneNumberUser(t, tt.args.onboardingClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrLoginTestPhoneNumberUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && userResponse != nil {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					userResponse,
					err,
				)
				return
			}

			if !tt.wantErr && userResponse == nil {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
		})
	}
}

func TestRemoveTestPhoneNumberUser(t *testing.T) {
	onboardingClient, err := onboardingISCClient()
	if err != nil {
		t.Errorf("failed to initialize onboarding test ISC client")
		return
	}

	TestCreateOrLoginTestPhoneNumberUser(t)

	type args struct {
		onboardingClient *base.InterServiceClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: remove the created test user",
			args: args{
				onboardingClient: onboardingClient,
			},
			wantErr: false,
		},
		{
			name: "failure: failed to remove the created test user",
			args: args{
				onboardingClient: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := base.RemoveTestPhoneNumberUser(
				t,
				tt.args.onboardingClient,
			); (err != nil) != tt.wantErr {
				t.Errorf("RemoveTestPhoneNumberUser() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

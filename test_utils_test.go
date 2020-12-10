package base

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
			got, err := GetOrCreateAnonymousUser(tt.args.ctx)
			assert.NotNil(t, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAnonymousUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
func TestGetOrCreatePhoneNumberUser(t *testing.T) {
	type args struct {
		ctx    context.Context
		msisdn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Create phone number user happy case",
			args: args{
				ctx:    context.Background(),
				msisdn: "+254711223344",
			},
			wantErr: false,
		},
		{
			name: "Create phone number user invalid case",
			args: args{
				ctx:    context.Background(),
				msisdn: "not a phone number",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOrCreatePhoneNumberUser(tt.args.ctx, tt.args.msisdn)
			if err == nil {
				assert.NotNil(t, got)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrCreatePhoneNumberUser() error = %v, wantErr %v", err, tt.wantErr)
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

			got, err := GetAuthenticatedContextFromUID(ctx, tt.args.uid)
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

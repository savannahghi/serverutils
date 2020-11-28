package base_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestGetTokenSource(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetGSUITEDelegatedAuthorityTokenSource(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTokenSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}

func TestGetLoggedInUserUID(t *testing.T) {
	authenticatedContext, token := base.GetAuthenticatedContextAndToken(t)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx: authenticatedContext,
			},
			want:    token.UID,
			wantErr: false,
		},
		{
			name: "bad case",
			args: args{
				ctx: context.Background(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetLoggedInUserUID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLoggedInUserUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLoggedInUserUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustGetLoggedInUserUID(t *testing.T) {
	authenticatedContext, token := base.GetAuthenticatedContextAndToken(t)
	unauthenticatedContext := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "good case",
			args: args{
				ctx: authenticatedContext,
			},
			want: token.UID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.MustGetLoggedInUserUID(tt.args.ctx); got != tt.want {
				t.Errorf("MustGetLoggedInUserUID() = %v, want %v", got, tt.want)
			}
			assert.Panics(t, func() {
				base.MustGetLoggedInUserUID(unauthenticatedContext)
			})
		})
	}
}

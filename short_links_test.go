package base_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestShortenLink(t *testing.T) {
	dynamicLinkDomain, err := base.GetEnvVar(base.FDLDomainEnvironmentVariableName)
	assert.Nil(t, err)

	type args struct {
		ctx      context.Context
		longLink string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:      context.Background(),
				longLink: "https://console.cloud.google.com/run/detail/europe-west1/api-gateway/revisions?project=bewell-app-testing",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.ShortenLink(tt.args.ctx, tt.args.longLink)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShortenLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Contains(t, got, dynamicLinkDomain)
		})
	}
}

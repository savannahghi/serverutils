package go_utils_test

import (
	"testing"

	base "github.com/savannahghi/go_utils"
)

func TestGetConsumerTerms(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "Good case",
			want:    base.ConsumerTerms,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetConsumerTerms()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConsumerTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetConsumerTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProviderTerms(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "Good case",
			want:    base.ProviderTerms,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetProviderTerms()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProviderTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetProviderTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

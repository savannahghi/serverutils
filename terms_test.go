package base

import (
	"testing"
)

func TestGetConsumerTerms(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "Good case",
			want:    ConsumerTerms,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConsumerTerms()
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
			want:    ProviderTerms,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProviderTerms()
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

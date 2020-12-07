package base_test

import (
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestUserProfile_IsEntity(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "default case - just checking that the profile is marked as an entity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := base.UserProfile{}
			u.IsEntity()
		})
	}
}

func TestCover_IsEntity(t *testing.T) {
	type fields struct {
		PayerName      string
		PayerSladeCode int
		MemberNumber   string
		MemberName     string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "default case - just checking that the cover is marked as an entity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := base.Cover{}
			c.IsEntity()
		})
	}
}

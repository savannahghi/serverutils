package base_test

import (
	"reflect"
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

func TestRoleType_Permissions(t *testing.T) {
	employeePermissions := []base.PermissionType{
		base.PermissionTypeCreateConsumer,
		base.PermissionTypeUpdateConsumer,
		base.PermissionTypeDeleteConsumer,
		base.PermissionTypeCreatePatient,
		base.PermissionTypeUpdatePatient,
		base.PermissionTypeDeletePatient,
		base.PermissionTypeIdentifyPatient,
	}
	agentPermissions := []base.PermissionType{
		base.PermissionTypeCreatePartner,
		base.PermissionTypeUpdatePartner,
		base.PermissionTypeCreateConsumer,
		base.PermissionTypeUpdateConsumer,
	}

	tests := []struct {
		name    string
		r       base.RoleType
		want    []base.PermissionType
		wantErr bool
	}{
		{
			name: "valid role type permissions",
			r:    base.RoleTypeEmployee,
			want: employeePermissions,
		},
		{
			name: "valid role type permissions",
			r:    base.RoleTypeAgent,
			want: agentPermissions,
		},
		{
			name: "invalid role type permissions",
			r:    "IMPOSTER",
			want: []base.PermissionType{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.r.Permissions()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoleType.Permissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		r    base.RoleType
		want bool
	}{
		{
			name: "valid employee role type",
			r:    base.RoleTypeEmployee,
			want: true,
		},
		{
			name: "valid agent role type",
			r:    base.RoleTypeAgent,
			want: true,
		},
		{
			name: "invalid role type",
			r:    "IMPOSTER",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.IsValid(); got != tt.want {
				t.Errorf("RoleType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

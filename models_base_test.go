package base_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestTypeof(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{
				v: "this is a string",
			},
			want: "string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.Typeof(tt.args.v); got != tt.want {
				t.Errorf("Typeof() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIDValue_String(t *testing.T) {
	tests := []struct {
		name string
		val  base.IDValue
		want string
	}{
		{
			name: "happy case",
			val:  base.IDValue("mimi ni id"),
			want: "mimi ni id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.String(); got != tt.want {
				t.Errorf("IDValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalID(t *testing.T) {
	type args struct {
		id string
		n  base.Node
	}
	tests := []struct {
		name string
		args args
		want base.ID
	}{
		{
			name: "good case",
			args: args{
				id: "1",
				n:  &base.Model{},
			},
			want: base.IDValue("MXwqYmFzZS5Nb2RlbA=="),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.MarshalID(tt.args.id, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewString(t *testing.T) {
	ns := base.NewString("a string")
	assert.Equal(t, "a string", *ns)
}

func TestModel_GetID(t *testing.T) {
	type fields struct {
		ID           string
		Name         string
		Description  string
		Deleted      bool
		CreatedByUID string
		UpdatedByUID string
	}
	tests := []struct {
		name   string
		fields fields
		want   base.ID
	}{
		{
			name: "good case",
			fields: fields{
				ID: "an ID",
			},
			want: base.IDValue("an ID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &base.Model{
				ID:           tt.fields.ID,
				Name:         tt.fields.Name,
				Description:  tt.fields.Description,
				Deleted:      tt.fields.Deleted,
				CreatedByUID: tt.fields.CreatedByUID,
				UpdatedByUID: tt.fields.UpdatedByUID,
			}
			if got := c.GetID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_SetID(t *testing.T) {
	type fields struct {
		ID           string
		Name         string
		Description  string
		Deleted      bool
		CreatedByUID string
		UpdatedByUID string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "good case",
			args: args{
				id: "an ID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &base.Model{
				ID:           tt.fields.ID,
				Name:         tt.fields.Name,
				Description:  tt.fields.Description,
				Deleted:      tt.fields.Deleted,
				CreatedByUID: tt.fields.CreatedByUID,
				UpdatedByUID: tt.fields.UpdatedByUID,
			}
			c.SetID(tt.args.id)
			assert.Equal(t, c.GetID(), base.IDValue(tt.args.id))
		})
	}
}

func TestModel_IsNode(t *testing.T) {
	type fields struct {
		ID           string
		Name         string
		Description  string
		Deleted      bool
		CreatedByUID string
		UpdatedByUID string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "default case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &base.Model{
				ID:           tt.fields.ID,
				Name:         tt.fields.Name,
				Description:  tt.fields.Description,
				Deleted:      tt.fields.Deleted,
				CreatedByUID: tt.fields.CreatedByUID,
				UpdatedByUID: tt.fields.UpdatedByUID,
			}
			c.IsNode()
		})
	}
}

func TestModelsIsEntity(t *testing.T) {

	t5 := base.Attachment{}
	t5.IsEntity()

	t6 := base.EDIUserProfile{}
	t6.IsEntity()

	t7 := base.LogoutRequest{}
	t7.IsEntity()

	t8 := base.RefreshCreds{}
	t8.IsEntity()

	t9 := base.LoginResponse{}
	t9.IsEntity()

	t10 := base.LoginCreds{}
	t10.IsEntity()

	t11 := base.EmailOptIn{}
	t11.IsEntity()

	t12 := base.USSDSessionLog{}
	t12.IsEntity()

	t13 := base.PhoneOptIn{}
	t13.IsEntity()

	t14 := base.FilterParam{}
	t14.IsEntity()

	t15 := base.FilterInput{}
	t15.IsEntity()

	t16 := base.SortInput{}
	t16.IsEntity()

	t17 := base.PaginationInput{}
	t17.IsEntity()

	t19 := base.Attachment{}
	t19.IsEntity()
}

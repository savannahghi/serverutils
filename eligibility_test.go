package base_test

import (
	"bytes"
	"strconv"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestBenefitType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.BenefitType
		want bool
	}{
		{
			name: "valid benefit type",
			e:    base.BenefitTypeDental,
			want: true,
		},
		{
			name: "unknown benefit type",
			e:    base.BenefitType("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("BenefitType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBenefitType_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.BenefitType
		want string
	}{
		{
			name: "dental benefit",
			e:    base.BenefitTypeDental,
			want: "DENTAL",
		},
		{
			name: "OP benefit",
			e:    base.BenefitTypeOutpatient,
			want: "OUTPATIENT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("BenefitType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBenefitType_UnmarshalGQL(t *testing.T) {
	target := base.BenefitType("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.BenefitType
		args    args
		wantErr bool
	}{
		{
			name: "valid benefit type",
			e:    &target,
			args: args{
				v: "OPTICAL",
			},
			wantErr: false,
		},
		{
			name: "invalid benefit type",
			e:    &target,
			args: args{
				v: "not a valid benefit",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("BenefitType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBenefitType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.BenefitType
		wantW string
	}{
		{
			name:  "maternity benefit",
			e:     base.BenefitTypeMaternity,
			wantW: strconv.Quote("MATERNITY"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("BenefitType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestRelationship_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Relationship
		want bool
	}{
		{
			name: "valid relationship",
			e:    base.RelationshipSpouse,
			want: true,
		},
		{
			name: "invalid relationship",
			e:    base.Relationship("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Relationship.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationship_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Relationship
		want string
	}{
		{
			name: "child",
			e:    base.RelationshipChild,
			want: "CHILD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Relationship.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelationship_UnmarshalGQL(t *testing.T) {
	target := base.Relationship("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Relationship
		args    args
		wantErr bool
	}{
		{
			name: "valid case - child",
			e:    &target,
			args: args{
				v: "CHILD",
			},
			wantErr: false,
		},
		{
			name: "invalid case",
			e:    &target,
			args: args{
				v: "not a real relationship type",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Relationship.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRelationship_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Relationship
		wantW string
	}{
		{
			name:  "valid marshal of father relationship",
			e:     base.RelationshipFather,
			wantW: strconv.Quote("FATHER"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Relationship.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

package base_test

import (
	"bytes"
	"strconv"
	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestGender_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Gender
		want string
	}{
		{
			name: "male",
			e:    base.GenderMale,
			want: "male",
		},
		{
			name: "female",
			e:    base.GenderFemale,
			want: "female",
		},
		{
			name: "unknown",
			e:    base.GenderUnknown,
			want: "unknown",
		},
		{
			name: "other",
			e:    base.GenderOther,
			want: "other",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Gender.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGender_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Gender
		want bool
	}{
		{
			name: "valid male",
			e:    base.GenderMale,
			want: true,
		},
		{
			name: "invalid gender",
			e:    base.Gender("this is not a real gender"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Gender.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGender_UnmarshalGQL(t *testing.T) {
	female := base.GenderFemale
	invalid := base.Gender("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Gender
		args    args
		wantErr bool
	}{
		{
			name: "valid female gender",
			e:    &female,
			args: args{
				v: "female",
			},
			wantErr: false,
		},
		{
			name: "invalid gender",
			e:    &invalid,
			args: args{
				v: "this is not a real gender",
			},
			wantErr: true,
		},
		{
			name: "non string gender",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Gender.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGender_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Gender
		wantW string
	}{
		{
			name:  "valid unknown gender enum",
			e:     base.GenderUnknown,
			wantW: strconv.Quote("unknown"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Gender.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFieldType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.FieldType
		want bool
	}{
		{
			name: "valid string field type",
			e:    base.FieldTypeString,
			want: true,
		},
		{
			name: "invalid field type",
			e:    base.FieldType("this is not a real field type"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("FieldType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldType_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.FieldType
		want string
	}{
		{
			name: "valid boolean field type as string",
			e:    base.FieldTypeBoolean,
			want: "BOOLEAN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("FieldType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldType_UnmarshalGQL(t *testing.T) {
	intEnum := base.FieldType("")
	invalid := base.FieldType("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.FieldType
		args    args
		wantErr bool
	}{
		{
			name: "valid integer enum",
			e:    &intEnum,
			args: args{
				v: "INTEGER",
			},
			wantErr: false,
		},
		{
			name: "invalid enum",
			e:    &invalid,
			args: args{
				v: "NOT A VALID ENUM",
			},
			wantErr: true,
		},
		{
			name: "wrong type -int",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("FieldType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFieldType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.FieldType
		wantW string
	}{
		{
			name:  "number field type",
			e:     base.FieldTypeNumber,
			wantW: strconv.Quote("NUMBER"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FieldType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestOperation_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Operation
		want bool
	}{
		{
			name: "valid operation",
			e:    base.OperationEqual,
			want: true,
		},
		{
			name: "invalid operation",
			e:    base.Operation("hii sio valid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Operation.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Operation
		want string
	}{
		{
			name: "valid case - contains",
			e:    base.OperationContains,
			want: "CONTAINS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Operation.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperation_UnmarshalGQL(t *testing.T) {
	valid := base.Operation("")
	invalid := base.Operation("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Operation
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			e:    &valid,
			args: args{
				v: "CONTAINS",
			},
			wantErr: false,
		},
		{
			name: "invalid string value",
			e:    &invalid,
			args: args{
				v: "NOT A REAL OPERATION",
			},
			wantErr: true,
		},
		{
			name: "invalid non string value",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Operation.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperation_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Operation
		wantW string
	}{
		{
			name:  "good case",
			e:     base.OperationContains,
			wantW: strconv.Quote("CONTAINS"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Operation.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestSortOrder_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.SortOrder
		want string
	}{
		{
			name: "good case",
			e:    base.SortOrderAsc,
			want: "ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("SortOrder.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortOrder_UnmarshalGQL(t *testing.T) {
	so := base.SortOrder("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.SortOrder
		args    args
		wantErr bool
	}{
		{
			name: "valid sort order",
			e:    &so,
			args: args{
				v: "ASC",
			},
			wantErr: false,
		},
		{
			name: "invalid sort order string",
			e:    &so,
			args: args{
				v: "not a valid sort order",
			},
			wantErr: true,
		},
		{
			name: "invalid sort order - non string",
			e:    &so,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SortOrder.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSortOrder_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.SortOrder
		wantW string
	}{
		{
			name:  "good case",
			e:     base.SortOrderDesc,
			wantW: strconv.Quote("DESC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SortOrder.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestContentType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.ContentType
		want bool
	}{
		{
			name: "good case",
			e:    base.ContentTypeJpg,
			want: true,
		},
		{
			name: "bad case",
			e:    base.ContentType("not a real content type"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ContentType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentType_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.ContentType
		want string
	}{
		{
			name: "default case",
			e:    base.ContentTypePdf,
			want: "PDF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ContentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContentType_UnmarshalGQL(t *testing.T) {
	var sc base.ContentType
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.ContentType
		args    args
		wantErr bool
	}{
		{
			name: "valid unmarshal",
			e:    &sc,
			args: args{
				v: "PDF",
			},
			wantErr: false,
		},
		{
			name: "invalid unmarshal",
			e:    &sc,
			args: args{
				v: "this is not a valid scalar value",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ContentType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContentType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.ContentType
		wantW string
	}{
		{
			name:  "default case",
			e:     base.ContentTypePdf,
			wantW: strconv.Quote("PDF"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ContentType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestLanguage_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Language
		want bool
	}{
		{
			name: "good case",
			e:    base.LanguageEn,
			want: true,
		},
		{
			name: "bad case",
			e:    base.Language("not a real language"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Language.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanguage_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Language
		want string
	}{
		{
			name: "default case",
			e:    base.LanguageEn,
			want: "en",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Language.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLanguage_UnmarshalGQL(t *testing.T) {
	var sc base.Language

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Language
		args    args
		wantErr bool
	}{
		{
			name: "valid unmarshal",
			e:    &sc,
			args: args{
				v: "en",
			},
			wantErr: false,
		},
		{
			name: "invalid unmarshal",
			e:    &sc,
			args: args{
				v: "this is not a valid scalar value",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Language.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLanguage_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Language
		wantW string
	}{
		{
			name:  "default case",
			e:     base.LanguageEn,
			wantW: strconv.Quote("en"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Language.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPractitionerSpecialty_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.PractitionerSpecialty
		want bool
	}{
		{
			name: "good case",
			e:    base.PractitionerSpecialtyAnaesthesia,
			want: true,
		},
		{
			name: "bad case",
			e:    base.PractitionerSpecialty("not a real specialty"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("PractitionerSpecialty.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPractitionerSpecialty_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.PractitionerSpecialty
		want string
	}{
		{
			name: "default case",
			e:    base.PractitionerSpecialtyAnaesthesia,
			want: "ANAESTHESIA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("PractitionerSpecialty.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPractitionerSpecialty_UnmarshalGQL(t *testing.T) {
	var sc base.PractitionerSpecialty
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.PractitionerSpecialty
		args    args
		wantErr bool
	}{
		{
			name: "valid unmarshal",
			e:    &sc,
			args: args{
				v: "ANAESTHESIA",
			},
			wantErr: false,
		},
		{
			name: "invalid unmarshal",
			e:    &sc,
			args: args{
				v: "this is not a valid scalar value",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("PractitionerSpecialty.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPractitionerSpecialty_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.PractitionerSpecialty
		wantW string
	}{
		{
			name:  "default case",
			e:     base.PractitionerSpecialtyAnaesthesia,
			wantW: strconv.Quote("ANAESTHESIA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PractitionerSpecialty.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

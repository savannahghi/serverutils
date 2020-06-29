package base

import (
	"bytes"
	"strconv"
	"testing"
)

func TestGender_String(t *testing.T) {
	tests := []struct {
		name string
		e    Gender
		want string
	}{
		{
			name: "male",
			e:    GenderMale,
			want: "male",
		},
		{
			name: "female",
			e:    GenderFemale,
			want: "female",
		},
		{
			name: "unknown",
			e:    GenderUnknown,
			want: "unknown",
		},
		{
			name: "other",
			e:    GenderOther,
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
		e    Gender
		want bool
	}{
		{
			name: "valid male",
			e:    GenderMale,
			want: true,
		},
		{
			name: "invalid gender",
			e:    Gender("this is not a real gender"),
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
	female := GenderFemale
	invalid := Gender("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *Gender
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
		e     Gender
		wantW string
	}{
		{
			name:  "valid unknown gender enum",
			e:     GenderUnknown,
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
		e    FieldType
		want bool
	}{
		{
			name: "valid string field type",
			e:    FieldTypeString,
			want: true,
		},
		{
			name: "invalid field type",
			e:    FieldType("this is not a real field type"),
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
		e    FieldType
		want string
	}{
		{
			name: "valid boolean field type as string",
			e:    FieldTypeBoolean,
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
	intEnum := FieldType("")
	invalid := FieldType("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *FieldType
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
		e     FieldType
		wantW string
	}{
		{
			name:  "number field type",
			e:     FieldTypeNumber,
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
		e    Operation
		want bool
	}{
		{
			name: "valid operation",
			e:    OperationEqual,
			want: true,
		},
		{
			name: "invalid operation",
			e:    Operation("hii sio valid"),
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
		e    Operation
		want string
	}{
		{
			name: "valid case - contains",
			e:    OperationContains,
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
	valid := Operation("")
	invalid := Operation("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *Operation
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
		e     Operation
		wantW string
	}{
		{
			name:  "good case",
			e:     OperationContains,
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
		e    SortOrder
		want string
	}{
		{
			name: "good case",
			e:    SortOrderAsc,
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
	so := SortOrder("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *SortOrder
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
		e     SortOrder
		wantW string
	}{
		{
			name:  "good case",
			e:     SortOrderDesc,
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

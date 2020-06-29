package base

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCursor(t *testing.T) {
	type args struct {
		offset int
	}
	tests := []struct {
		name string
		args args
		want *Cursor
	}{
		{
			name: "offset one",
			args: args{
				offset: 1,
			},
			want: &Cursor{Offset: 1},
		},
		{
			name: "default offset",
			args: args{
				offset: 0,
			},
			want: &Cursor{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCursor(tt.args.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeCursor(t *testing.T) {
	type args struct {
		cursor *Cursor
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default zero offset cursor",
			args: args{
				cursor: &Cursor{},
			},
			want: "gaZPZmZzZXTTAAAAAAAAAAA=",
		},
		{
			name: "negative cursor",
			args: args{
				cursor: &Cursor{Offset: -1},
			},
			want: "gaZPZmZzZXTT//////////8=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeCursor(tt.args.cursor); got != tt.want {
				t.Errorf("EncodeCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateAndEncodeCursor(t *testing.T) {
	zeroCur := "gaZPZmZzZXTTAAAAAAAAAAA="
	negCur := "gaZPZmZzZXTT//////////8="
	type args struct {
		offset int
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "default zero offset cursor",
			args: args{
				offset: 0,
			},
			want: &zeroCur,
		},
		{
			name: "negative cursor",
			args: args{
				offset: -1,
			},
			want: &negCur,
		},
	}
	for _, tt := range tests {
		got := CreateAndEncodeCursor(tt.args.offset)
		assert.NotNil(t, got)
		assert.Equal(t, *got, *tt.want)
	}
}

func TestDecodeCursor(t *testing.T) {
	zeroCur := "gaZPZmZzZXTTAAAAAAAAAAA="
	negCur := "gaZPZmZzZXTT//////////8="
	type args struct {
		cursor string
	}
	tests := []struct {
		name    string
		args    args
		want    *Cursor
		wantErr bool
	}{
		{
			name: "zero cursor",
			args: args{
				cursor: zeroCur,
			},
			want:    &Cursor{Offset: 0},
			wantErr: false,
		},
		{
			name: "minus one cursor",
			args: args{
				cursor: negCur,
			},
			want:    &Cursor{Offset: -1},
			wantErr: false,
		},
		{
			name: "invalid cursor",
			args: args{
				cursor: "this is not a valid cursor",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeCursor(tt.args.cursor)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeCursor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecodeCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustDecodeCursor(t *testing.T) {
	zeroCur := "gaZPZmZzZXTTAAAAAAAAAAA="
	negCur := "gaZPZmZzZXTT//////////8="
	type args struct {
		cursor string
	}
	tests := []struct {
		name string
		args args
		want *Cursor
	}{
		{
			name: "zero cursor",
			args: args{
				cursor: zeroCur,
			},
			want: &Cursor{Offset: 0},
		},
		{
			name: "minus one cursor",
			args: args{
				cursor: negCur,
			},
			want: &Cursor{Offset: -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustDecodeCursor(tt.args.cursor); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustDecodeCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

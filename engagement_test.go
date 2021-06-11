package go_utils_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	base "github.com/savannahghi/go_utils"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

const (
	sampleVideoURL = "https://www.youtube.com/watch?v=bPiofmZGb8o"
	intMax         = 9007199254740990
)

func getBlankActionType() *base.ActionType {
	at := base.ActionType("")
	return &at
}

func getEmptyJson(t *testing.T) []byte {
	emptyJSONBytes, err := json.Marshal(map[string]string{})
	assert.Nil(t, err)
	assert.NotNil(t, emptyJSONBytes)
	return emptyJSONBytes
}

func TestActionType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.ActionType
		want bool
	}{
		{
			name: "valid case",
			e:    base.ActionTypeFloating,
			want: true,
		},
		{
			name: "invalid case",
			e:    base.ActionType("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ActionType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionType_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.ActionType
		want string
	}{
		{
			name: "primary",
			e:    base.ActionTypePrimary,
			want: "PRIMARY",
		},
		{
			name: "secondary",
			e:    base.ActionTypeSecondary,
			want: "SECONDARY",
		},
		{
			name: "overflow",
			e:    base.ActionTypeOverflow,
			want: "OVERFLOW",
		},
		{
			name: "floating",
			e:    base.ActionTypeFloating,
			want: "FLOATING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ActionType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionType_UnmarshalGQL(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.ActionType
		args    args
		wantErr bool
	}{
		{
			name: "primary",
			e:    getBlankActionType(),
			args: args{
				v: "PRIMARY",
			},
			wantErr: false,
		},
		{
			name: "secondary",
			e:    getBlankActionType(),
			args: args{
				v: "SECONDARY",
			},
			wantErr: false,
		},
		{
			name: "overflow",
			e:    getBlankActionType(),
			args: args{
				v: "OVERFLOW",
			},
			wantErr: false,
		},
		{
			name: "floating",
			e:    getBlankActionType(),
			args: args{
				v: "FLOATING",
			},
			wantErr: false,
		},
		{
			name: "invalid - should error",
			e:    getBlankActionType(),
			args: args{
				v: "bogus bonoko",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"ActionType.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestActionType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.ActionType
		wantW string
	}{
		{
			name:  "floating",
			e:     base.ActionTypeFloating,
			wantW: `"FLOATING"`,
		},
		{
			name:  "primary",
			e:     base.ActionTypePrimary,
			wantW: `"PRIMARY"`,
		},
		{
			name:  "secondary",
			e:     base.ActionTypeSecondary,
			wantW: `"SECONDARY"`,
		},
		{
			name:  "overflow",
			e:     base.ActionTypeOverflow,
			wantW: `"OVERFLOW"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf(
					"ActionType.MarshalGQL() = %v, want %v",
					gotW,
					tt.wantW,
				)
			}
		})
	}
}

func TestHandling_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Handling
		want bool
	}{
		{
			name: "valid case",
			e:    base.HandlingFullPage,
			want: true,
		},
		{
			name: "invalid case",
			e:    base.Handling("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Handling.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandling_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Handling
		want string
	}{
		{
			name: "simple case",
			e:    base.HandlingInline,
			want: "INLINE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Handling.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandling_UnmarshalGQL(t *testing.T) {
	target := base.Handling("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Handling
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "INLINE",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Handling.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestHandling_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Handling
		wantW string
	}{
		{
			name:  "simple case",
			e:     base.HandlingFullPage,
			wantW: `"FULL_PAGE"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf(
					"Handling.MarshalGQL() = %v, want %v",
					gotW,
					tt.wantW,
				)
			}
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Status
		want bool
	}{
		{
			name: "valid case",
			e:    base.StatusDone,
			want: true,
		},
		{
			name: "invalid case",
			e:    base.Status("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Status.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Status
		want string
	}{
		{
			name: "simple case",
			e:    base.StatusDone,
			want: "DONE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_UnmarshalGQL(t *testing.T) {
	target := base.Status("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Status
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "DONE",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Status.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestStatus_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Status
		wantW string
	}{
		{
			name:  "simple case",
			e:     base.StatusDone,
			wantW: `"DONE"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Status.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestVisibility_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Visibility
		want bool
	}{
		{
			name: "valid case",
			e:    base.VisibilityHide,
			want: true,
		},
		{
			name: "invalid case",
			e:    base.Visibility("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Visibility.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisibility_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Visibility
		want string
	}{

		{
			name: "simple case",
			e:    base.VisibilityShow,
			want: "SHOW",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Visibility.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisibility_UnmarshalGQL(t *testing.T) {
	target := base.Visibility("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Visibility
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "SHOW",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Visibility.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestVisibility_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Visibility
		wantW string
	}{
		{
			name:  "simple case",
			e:     base.VisibilityHide,
			wantW: `"HIDE"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf(
					"Visibility.MarshalGQL() = %v, want %v",
					gotW,
					tt.wantW,
				)
			}
		})
	}
}

func TestChannel_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Channel
		want bool
	}{
		{
			name: "valid case",
			e:    base.ChannelEmail,
			want: true,
		},
		{
			name: "invalid case",
			e:    base.Channel("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Channel.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Channel
		want string
	}{
		{
			name: "simple case",
			e:    base.ChannelEmail,
			want: "EMAIL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Channel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChannel_UnmarshalGQL(t *testing.T) {
	target := base.Channel("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Channel
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "EMAIL",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Channel.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestChannel_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Channel
		wantW string
	}{
		{
			name:  "simple case",
			e:     base.ChannelEmail,
			wantW: `"EMAIL"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Channel.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFlavour_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Flavour
		want bool
	}{
		{
			name: "valid case",
			e:    base.FlavourConsumer,
			want: true,
		},
		{
			name: "invalid case",
			e:    base.Flavour("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Flavour.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlavour_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Flavour
		want string
	}{
		{
			name: "simple case",
			e:    base.FlavourConsumer,
			want: "CONSUMER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Flavour.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlavour_UnmarshalGQL(t *testing.T) {
	target := base.Flavour("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Flavour
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "PRO",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Flavour.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestFlavour_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Flavour
		wantW string
	}{
		{
			name:  "simple case",
			e:     base.FlavourPro,
			wantW: `"PRO"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Flavour.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestKeys_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.Keys
		want bool
	}{
		{
			name: "valid case",
			e:    base.KeysActions,
			want: true,
		},
		{
			name: "invalid case",
			e:    base.Keys("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Keys.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeys_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.Keys
		want string
	}{
		{
			name: "simple case",
			e:    base.KeysActions,
			want: "actions",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Keys.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeys_UnmarshalGQL(t *testing.T) {
	target := base.Keys("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.Keys
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "actions",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(
				tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"Keys.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestKeys_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.Keys
		wantW string
	}{
		{
			name:  "simple case",
			e:     base.KeysActions,
			wantW: `"actions"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Keys.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestBooleanFilter_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.BooleanFilter
		want bool
	}{
		{
			name: "valid case",
			e:    base.BooleanFilterBoth,
			want: true,
		},
		{
			name: "invalid case",
			e:    base.BooleanFilter("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("BooleanFilter.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBooleanFilter_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.BooleanFilter
		want string
	}{
		{
			name: "simple case",
			e:    base.BooleanFilterFalse,
			want: "FALSE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("BooleanFilter.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBooleanFilter_UnmarshalGQL(t *testing.T) {
	target := base.BooleanFilter("")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.BooleanFilter
		args    args
		wantErr bool
	}{
		{
			name: "successful case",
			e:    &target,
			args: args{
				v: "BOTH",
			},
			wantErr: false,
		},
		{
			name: "failure case",
			e:    &target,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf(
					"BooleanFilter.UnmarshalGQL() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestBooleanFilter_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.BooleanFilter
		wantW string
	}{
		{
			name:  "simple case",
			e:     base.BooleanFilterBoth,
			wantW: `"BOTH"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf(
					"BooleanFilter.MarshalGQL() = %v, want %v",
					gotW,
					tt.wantW,
				)
			}
		})
	}
}

func TestLinkType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     base.LinkType
		wantW string
	}{
		{
			name:  "PDF document",
			e:     base.LinkTypePdfDocument,
			wantW: `"PDF_DOCUMENT"`,
		},
		{
			name:  "PNG Image",
			e:     base.LinkTypePngImage,
			wantW: `"PNG_IMAGE"`,
		},
		{
			name:  "YouTube Video",
			e:     base.LinkTypeYoutubeVideo,
			wantW: `"YOUTUBE_VIDEO"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("LinkType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestLinkType_UnmarshalGQL(t *testing.T) {
	l := base.LinkType("")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *base.LinkType
		args    args
		wantErr bool
	}{
		{
			name: "invalid link type",
			e:    &l,
			args: args{
				v: "bogus",
			},
			wantErr: true,
		},
		{
			name: "valid - pdf",
			e:    &l,
			args: args{
				v: "PDF_DOCUMENT",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("LinkType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLinkType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    base.LinkType
		want bool
	}{
		{
			name: "PDF document",
			e:    base.LinkTypePdfDocument,
			want: true,
		},
		{
			name: "PNG Image",
			e:    base.LinkTypePngImage,
			want: true,
		},
		{
			name: "YouTube Video",
			e:    base.LinkTypeYoutubeVideo,
			want: true,
		},
		{
			name: "invalid link type",
			e:    base.LinkType("bogus"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("LinkType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinkType_String(t *testing.T) {
	tests := []struct {
		name string
		e    base.LinkType
		want string
	}{
		{
			name: "YouTube video",
			e:    base.LinkTypeYoutubeVideo,
			want: "YOUTUBE_VIDEO",
		},
		{
			name: "PDF document",
			e:    base.LinkTypePdfDocument,
			want: "PDF_DOCUMENT",
		},
		{
			name: "PNG image",
			e:    base.LinkTypePngImage,
			want: "PNG_IMAGE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("LinkType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Text:           "some message text",
		PostedByName:   ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		Timestamp:      time.Now(),
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &base.Message{}
			if err := msg.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Message.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestItem_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Item{
		ID:             "item-1",
		SequenceNumber: 1,
		Expiry:         time.Now(),
		Persistent:     true,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		Author:         "Bot 1",
		Tagline:        "Bot speaks...",
		Label:          "DRUGS",
		Timestamp:      time.Now(),
		Summary:        "I am a bot...",
		Text:           "This bot can speak",
		TextType:       base.TextTypeMarkdown,
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		},
		Actions: []base.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: 1,
				Name:           "ACTION_NAME",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypeSecondary,
				Handling:       base.HandlingFullPage,
			},
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingInline,
			},
		},
		Conversations: []base.Message{
			{
				ID:             "msg-2",
				SequenceNumber: 1,
				Text:           "hii ni reply",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
				Timestamp:      time.Now(),
			},
		},
		Users: []string{
			"user-1",
			"user-2",
		},
		Groups: []string{
			"group-1",
			"group-2",
		},
		NotificationChannels: []base.Channel{
			base.ChannelFcm,
			base.ChannelEmail,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Expiry         time.Time
		Persistent     bool
		Status         base.Status
		Visibility     base.Visibility
		Icon           base.Link
		Author         string
		Tagline        string
		Label          string
		Timestamp      time.Time
		Summary        string
		Text           string
		Links          []base.Link
		Actions        []base.Action
		Conversations  []base.Message
		Users          []string
		Groups         []string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &base.Item{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Expiry:         tt.fields.Expiry,
				Persistent:     tt.fields.Persistent,
				Status:         tt.fields.Status,
				Visibility:     tt.fields.Visibility,
				Icon:           tt.fields.Icon,
				Author:         tt.fields.Author,
				Tagline:        tt.fields.Tagline,
				Label:          tt.fields.Label,
				Timestamp:      tt.fields.Timestamp,
				Summary:        tt.fields.Summary,
				Text:           tt.fields.Text,
				Links:          tt.fields.Links,
				Actions:        tt.fields.Actions,
				Conversations:  tt.fields.Conversations,
				Users:          tt.fields.Users,
				Groups:         tt.fields.Groups,
			}
			if err := it.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Item.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestNudge_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Nudge{
		ID:             "nudge-1",
		SequenceNumber: 1,
		Visibility:     base.VisibilityShow,
		Status:         base.StatusPending,
		Title:          "Update your profile!",
		Links: []base.Link{
			base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		},
		Text: "An up to date profile will help us serve you better!",
		Actions: []base.Action{
			{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingInline,
			},
		},
		Groups: []string{
			"group-1",
			"group-2",
		},
		Users: []string{
			"user-1",
			"user-2",
		},
		NotificationChannels: []base.Channel{
			base.ChannelFcm,
			base.ChannelEmail,
			base.ChannelSms,
			base.ChannelWhatsapp,
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Visibility     base.Visibility
		Status         base.Status
		Title          string
		Text           string
		Links          []base.Link
		Actions        []base.Action
		Groups         []string
		Users          []string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nu := &base.Nudge{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Visibility:     tt.fields.Visibility,
				Status:         tt.fields.Status,
				Title:          tt.fields.Title,
				Links:          tt.fields.Links,
				Text:           tt.fields.Text,
				Actions:        tt.fields.Actions,
				Groups:         tt.fields.Groups,
				Users:          tt.fields.Users,
			}
			if err := nu.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Nudge.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestAction_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Action{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Name:           "ACTION_NAME",
		Icon:           base.GetPNGImageLink(base.LogoURL, "title", "description", base.BlankImageURL),
		ActionType:     base.ActionTypeSecondary,
		Handling:       base.HandlingFullPage,
		AllowAnonymous: false,
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)
	assert.Greater(t, len(validBytes), 3)

	type fields struct {
		ID             string
		SequenceNumber int
		Name           string
		Icon           base.Link
		ActionType     base.ActionType
		Handling       base.Handling
		AllowAnonymous bool
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &base.Action{}
			if err := ac.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Action.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestAction_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID             string
		SequenceNumber int
		Name           string
		Icon           base.Link
		ActionType     base.ActionType
		Handling       base.Handling
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid action",
			fields: fields{
				ID:             "action-1",
				SequenceNumber: 1,
				Name:           "First action",
				Icon: base.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				ActionType: base.ActionTypePrimary,
				Handling:   base.HandlingInline,
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &base.Action{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Name:           tt.fields.Name,
				Icon:           tt.fields.Icon,
				ActionType:     tt.fields.ActionType,
				Handling:       tt.fields.Handling,
			}
			got, err := ac.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Action.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestNudge_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID                   string
		SequenceNumber       int
		Visibility           base.Visibility
		Status               base.Status
		Title                string
		Links                []base.Link
		Text                 string
		Actions              []base.Action
		Groups               []string
		Users                []string
		NotificationChannels []base.Channel
		NotificationBody     base.NotificationBody
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case - valid nudge",
			fields: fields{
				ID:             "nudge-1",
				SequenceNumber: 1,
				Visibility:     base.VisibilityShow,
				Status:         base.StatusPending,
				Title:          "Update your profile!",
				Links: []base.Link{
					base.GetPNGImageLink(
						base.LogoURL, "title", "description", base.BlankImageURL),
				},
				Text: "An up to date profile will help us serve you better!",
				Actions: []base.Action{
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypePrimary,
						Handling:       base.HandlingInline,
						AllowAnonymous: false,
					},
				},
				Groups: []string{
					"group-1",
					"group-2",
				},
				Users: []string{
					"user-1",
					"user-2",
				},
				NotificationChannels: []base.Channel{
					base.ChannelFcm,
					base.ChannelEmail,
					base.ChannelSms,
					base.ChannelWhatsapp,
				},
				NotificationBody: base.NotificationBody{
					PublishMessage:   "Your nudge has been successfully published",
					ResolveMessage:   "Your nudge has been successfully resolved",
					DeleteMessage:    "Your nudge has been successfully deleted",
					UnresolveMessage: "Your nudge has been successfully unresolved",
					ShowMessage:      "Your nudge has been successfully shown",
					HideMessage:      "Your nudge has been successfully hidden",
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nu := &base.Nudge{
				ID:                   tt.fields.ID,
				SequenceNumber:       tt.fields.SequenceNumber,
				Visibility:           tt.fields.Visibility,
				Status:               tt.fields.Status,
				Title:                tt.fields.Title,
				Links:                tt.fields.Links,
				Text:                 tt.fields.Text,
				Actions:              tt.fields.Actions,
				Groups:               tt.fields.Groups,
				Users:                tt.fields.Users,
				NotificationChannels: tt.fields.NotificationChannels,
				NotificationBody:     tt.fields.NotificationBody,
			}
			got, err := nu.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Nudge.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestItem_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID                   string
		SequenceNumber       int
		Expiry               time.Time
		Persistent           bool
		Status               base.Status
		Visibility           base.Visibility
		Icon                 base.Link
		Author               string
		Tagline              string
		Label                string
		Timestamp            time.Time
		Summary              string
		Text                 string
		TextType             base.TextType
		Links                []base.Link
		Actions              []base.Action
		Conversations        []base.Message
		Users                []string
		Groups               []string
		NotificationChannels []base.Channel
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case - valid item",
			fields: fields{
				ID:             "item-1",
				SequenceNumber: 1,
				Expiry:         time.Now(),
				Persistent:     true,
				Status:         base.StatusPending,
				Visibility:     base.VisibilityShow,
				Icon: base.GetPNGImageLink(
					base.LogoURL, "title", "description", base.BlankImageURL),
				Author:    "Bot 1",
				Tagline:   "Bot speaks...",
				Label:     "DRUGS",
				Timestamp: time.Now(),
				Summary:   "I am a bot...",
				Text:      "This bot can speak",
				TextType:  base.TextTypeMarkdown,
				Links: []base.Link{
					base.GetPNGImageLink(
						base.LogoURL, "title", "description", base.BlankImageURL),
				},
				Actions: []base.Action{
					{
						ID:             ksuid.New().String(),
						SequenceNumber: 1,
						Name:           "ACTION_NAME",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypeSecondary,
						Handling:       base.HandlingFullPage,
						AllowAnonymous: false,
					},
					{
						ID:             "action-1",
						SequenceNumber: 1,
						Name:           "First action",
						Icon: base.GetPNGImageLink(
							base.LogoURL, "title", "description", base.BlankImageURL),
						ActionType:     base.ActionTypePrimary,
						Handling:       base.HandlingInline,
						AllowAnonymous: false,
					},
				},
				Conversations: []base.Message{
					{
						ID:             "msg-2",
						SequenceNumber: 1,
						Text:           "hii ni reply",
						ReplyTo:        "msg-1",
						PostedByName:   ksuid.New().String(),
						PostedByUID:    ksuid.New().String(),
						Timestamp:      time.Now(),
					},
				},
				Users: []string{
					"user-1",
					"user-2",
				},
				Groups: []string{
					"group-1",
					"group-2",
				},
				NotificationChannels: []base.Channel{
					base.ChannelFcm,
					base.ChannelEmail,
					base.ChannelSms,
					base.ChannelWhatsapp,
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			it := &base.Item{
				ID:                   tt.fields.ID,
				SequenceNumber:       tt.fields.SequenceNumber,
				Expiry:               tt.fields.Expiry,
				Persistent:           tt.fields.Persistent,
				Status:               tt.fields.Status,
				Visibility:           tt.fields.Visibility,
				Icon:                 tt.fields.Icon,
				Author:               tt.fields.Author,
				Tagline:              tt.fields.Tagline,
				Label:                tt.fields.Label,
				Timestamp:            tt.fields.Timestamp,
				Summary:              tt.fields.Summary,
				Text:                 tt.fields.Text,
				TextType:             tt.fields.TextType,
				Links:                tt.fields.Links,
				Actions:              tt.fields.Actions,
				Conversations:        tt.fields.Conversations,
				Users:                tt.fields.Users,
				Groups:               tt.fields.Groups,
				NotificationChannels: tt.fields.NotificationChannels,
			}
			got, err := it.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Item.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestMessage_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID             string
		SequenceNumber int
		Text           string
		ReplyTo        string
		PostedByName   string
		PostedByUID    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				ID:             "msg-2",
				SequenceNumber: 1,
				Text:           "this is a message",
				ReplyTo:        "msg-1",
				PostedByName:   ksuid.New().String(),
				PostedByUID:    ksuid.New().String(),
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &base.Message{
				ID:             tt.fields.ID,
				SequenceNumber: tt.fields.SequenceNumber,
				Text:           tt.fields.Text,
				ReplyTo:        tt.fields.ReplyTo,
				PostedByName:   tt.fields.PostedByName,
				PostedByUID:    tt.fields.PostedByUID,
				Timestamp:      time.Now(),
			}
			got, err := msg.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Message.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestContext_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Context{
		UserID:         "uid-1",
		Flavour:        base.FlavourConsumer,
		OrganizationID: "org-1",
		LocationID:     "loc-1",
		Timestamp:      time.Now(),
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)

	type fields struct {
		UserID         string
		OrganizationID string
		LocationID     string
		Timestamp      time.Time
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &base.Context{
				UserID:         tt.fields.UserID,
				OrganizationID: tt.fields.OrganizationID,
				LocationID:     tt.fields.LocationID,
				Timestamp:      tt.fields.Timestamp,
			}
			if err := ct.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Context.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestContext_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		UserID         string
		Flavour        base.Flavour
		OrganizationID string
		LocationID     string
		Timestamp      time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				UserID:         "uid-1",
				Flavour:        base.FlavourConsumer,
				OrganizationID: "org-1",
				LocationID:     "loc-1",
				Timestamp:      time.Now(),
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &base.Context{
				UserID:         tt.fields.UserID,
				Flavour:        tt.fields.Flavour,
				OrganizationID: tt.fields.OrganizationID,
				LocationID:     tt.fields.LocationID,
				Timestamp:      tt.fields.Timestamp,
			}
			got, err := ct.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Context.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestPayload_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Payload{
		Data: map[string]interface{}{"a": 1},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)

	type fields struct {
		Data map[string]interface{}
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := &base.Payload{
				Data: tt.fields.Data,
			}
			if err := pl.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Payload.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestPayload_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		Data map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				Data: map[string]interface{}{"a": 1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := &base.Payload{
				Data: tt.fields.Data,
			}
			got, err := pl.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Payload.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestEvent_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)

	validElement := base.Event{
		ID:   "event-1",
		Name: "THIS_EVENT",
		Context: base.Context{
			UserID:         "user-1",
			Flavour:        base.FlavourConsumer,
			OrganizationID: "org-1",
			LocationID:     "loc-1",
			Timestamp:      time.Now(),
		},
		Payload: base.Payload{
			Data: map[string]interface{}{"a": 1},
		},
	}
	validBytes, err := json.Marshal(validElement)
	assert.Nil(t, err)
	assert.NotNil(t, validBytes)

	type fields struct {
		ID      string
		Name    string
		Context base.Context
		Payload base.Payload
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid JSON",
			args: args{
				b: validBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &base.Event{
				ID:      tt.fields.ID,
				Name:    tt.fields.Name,
				Context: tt.fields.Context,
				Payload: tt.fields.Payload,
			}
			if err := ev.ValidateAndUnmarshal(
				tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf(
					"Event.ValidateAndUnmarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestEvent_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID      string
		Name    string
		Context base.Context
		Payload base.Payload
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid case",
			fields: fields{
				ID:   "event-1",
				Name: "THIS_EVENT",
				Context: base.Context{
					UserID:         "user-1",
					Flavour:        base.FlavourConsumer,
					OrganizationID: "org-1",
					LocationID:     "loc-1",
					Timestamp:      time.Now(),
				},
				Payload: base.Payload{
					Data: map[string]interface{}{"a": 1},
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid case - empty",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := &base.Event{
				ID:      tt.fields.ID,
				Name:    tt.fields.Name,
				Context: tt.fields.Context,
				Payload: tt.fields.Payload,
			}
			got, err := ev.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Event.ValidateAndMarshal() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr {
				assert.NotZero(t, got)
			}
		})
	}
}

func TestLink_ValidateAndUnmarshal(t *testing.T) {
	emptyJSONBytes := getEmptyJson(t)
	validLink := base.Link{
		ID:          ksuid.New().String(),
		URL:         sampleVideoURL,
		LinkType:    base.LinkTypeYoutubeVideo,
		Title:       "title",
		Description: "description",
		Thumbnail:   base.BlankImageURL,
	}
	validLinkJSONBytes, err := json.Marshal(validLink)
	assert.Nil(t, err)
	assert.NotNil(t, validLinkJSONBytes)
	assert.Greater(t, len(validLinkJSONBytes), 3)

	invalidVideoLink := base.Link{
		ID:          ksuid.New().String(),
		URL:         "www.example.com/not_a_youtube_video",
		LinkType:    base.LinkTypeYoutubeVideo,
		Title:       "title",
		Description: "description",
		Thumbnail:   base.BlankImageURL,
	}
	invalidLinkJSONBytes, err := json.Marshal(invalidVideoLink)
	assert.Nil(t, err)
	assert.NotNil(t, invalidLinkJSONBytes)
	assert.Greater(t, len(invalidLinkJSONBytes), 3)

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid link",
			args: args{
				b: validLinkJSONBytes,
			},
			wantErr: false,
		},
		{
			name: "invalid link",
			args: args{
				b: invalidLinkJSONBytes,
			},
			wantErr: true,
		},
		{
			name: "empty JSON - invalid",
			args: args{
				b: emptyJSONBytes,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &base.Link{}
			if err := l.ValidateAndUnmarshal(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Link.ValidateAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLink_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		ID          string
		URL         string
		Type        base.LinkType
		Title       string
		Description string
		Thumbnail   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid link",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         sampleVideoURL,
				Type:        base.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: false,
		},
		{
			name: "invalid URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "not a valid URL",
				Type:        base.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid YouTube URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_video",
				Type:        base.LinkTypeYoutubeVideo,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid PNG URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_png",
				Type:        base.LinkTypePngImage,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: true,
		},
		{
			name: "invalid PDF URL",
			fields: fields{
				ID:          ksuid.New().String(),
				URL:         "www.example.com/not_a_pdf",
				Type:        base.LinkTypePdfDocument,
				Title:       "title",
				Description: "description",
				Thumbnail:   base.BlankImageURL,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &base.Link{
				ID:          tt.fields.ID,
				URL:         tt.fields.URL,
				LinkType:    tt.fields.Type,
				Title:       tt.fields.Title,
				Description: tt.fields.Description,
				Thumbnail:   tt.fields.Thumbnail,
			}
			got, err := l.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Link.ValidateAndMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}
}

func TestValidateAndUnmarshal(t *testing.T) {
	msg := base.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Errorf("can't marshal message to JSON: %w", err)
		return
	}
	type args struct {
		sch string
		b   []byte
		el  base.Element
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				sch: base.MessageSchemaFile,
				b:   msgBytes,
				el:  &base.Message{},
			},
			wantErr: false,
		},
		{
			name: "invalid case",
			args: args{
				sch: base.MessageSchemaFile,
				b:   []byte("this should not pass validation"),
				el:  &base.Message{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := base.ValidateAndUnmarshal(tt.args.sch, tt.args.b, tt.args.el); (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAndMarshal(t *testing.T) {
	msg := base.Message{
		ID:             ksuid.New().String(),
		SequenceNumber: 1,
		Text:           ksuid.New().String(),
		ReplyTo:        ksuid.New().String(),
		PostedByUID:    ksuid.New().String(),
		PostedByName:   ksuid.New().String(),
		Timestamp:      time.Now(),
	}
	type args struct {
		sch string
		el  base.Element
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				sch: base.MessageSchemaFile,
				el:  &msg,
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "invalid case",
			args: args{
				sch: base.MessageSchemaFile,
				el:  &base.Message{},
			},
			wantNil: true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.ValidateAndMarshal(tt.args.sch, tt.args.el)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil && got == nil {
				t.Errorf("ValidateAndMarshal() got unexpected nil")
				return
			}
		})
	}
}

func TestNotificationBody_ValidateAndUnmarshal(t *testing.T) {
	validNotificationBody := base.NotificationBody{
		PublishMessage:   "publish message",
		DeleteMessage:    "delete message",
		ResolveMessage:   "resolve message",
		UnresolveMessage: "unresolve message",
		ShowMessage:      "show message",
		HideMessage:      "hide message",
	}
	validJSONBytes, err := json.Marshal(validNotificationBody)
	if err != nil {
		t.Errorf("an error occurred: %v", err)
		return
	}

	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy notification body",
			args: args{
				b: validJSONBytes,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nb := &base.NotificationBody{}
			if err := nb.ValidateAndUnmarshal(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("NotificationBody.ValidateAndUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotificationBody_ValidateAndMarshal(t *testing.T) {
	type fields struct {
		PublishMessage   string
		DeleteMessage    string
		ResolveMessage   string
		UnresolveMessage string
		ShowMessage      string
		HideMessage      string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "happy case",
			fields: fields{
				PublishMessage:   "publish message",
				DeleteMessage:    "delete message",
				ResolveMessage:   "resolve message",
				UnresolveMessage: "unresolve message",
				ShowMessage:      "show message",
				HideMessage:      "hide message",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nb := &base.NotificationBody{
				PublishMessage:   tt.fields.PublishMessage,
				DeleteMessage:    tt.fields.DeleteMessage,
				ResolveMessage:   tt.fields.ResolveMessage,
				UnresolveMessage: tt.fields.UnresolveMessage,
				ShowMessage:      tt.fields.ShowMessage,
				HideMessage:      tt.fields.HideMessage,
			}
			notificationBody, err := nb.ValidateAndMarshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("NotificationBody.ValidateAndMarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if notificationBody == nil {
				t.Errorf("expected notificationBody")
				return
			}
		})
	}
}

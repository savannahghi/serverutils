package base_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/segmentio/ksuid"
	"gitlab.slade360emr.com/go/base"
	"google.golang.org/api/idtoken"
)

func TestGetPubSubTopic(t *testing.T) {
	type args struct {
		m *base.PubSubPayload
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "nil payload",
			args: args{
				m: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "payload with no topic key",
			args: args{
				m: &base.PubSubPayload{
					Message: base.PubSubMessage{
						Attributes: map[string]string{
							"bad": "key",
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "payload withcorrect topic key",
			args: args{
				m: &base.PubSubPayload{
					Message: base.PubSubMessage{
						Attributes: map[string]string{
							"topicID": "some-topic",
						},
					},
				},
			},
			want:    "some-topic",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.GetPubSubTopic(tt.args.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPubSubTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPubSubTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyPubSubJWTAndDecodePayload(t *testing.T) {
	invalidHeaderReq := httptest.NewRequest(http.MethodPost, "/", nil)
	invalidHeaderReq.Header.Add("Authorization", "Bearer stuff")

	validHeaderReq := httptest.NewRequest(http.MethodPost, "/", nil)
	validHeaderReq.Header.Add("Authorization", "Bearer stuff")

	ctx := context.Background()
	goodIDTokenSource, err := idtoken.NewTokenSource(ctx, base.Aud)
	if err != nil {
		t.Errorf("can't initialize ID token source: %v", err)
		return
	}
	goodToken, err := goodIDTokenSource.Token()
	if err != nil {
		t.Errorf("error getting ID token: %v", err)
		return
	}
	goodTokenHeader := fmt.Sprintf(
		"%s %s", goodToken.TokenType, goodToken.AccessToken)
	testPayload := &base.PubSubPayload{
		Subscription: "test",
		Message: base.PubSubMessage{
			MessageID: ksuid.New().String(),
			Data:      []byte(ksuid.New().String()),
			Attributes: map[string]string{
				"topicID": ksuid.New().String(),
			},
		},
	}
	testPayloadBytes, err := json.Marshal(testPayload)
	if err != nil {
		t.Errorf("error marshalling test payload: %v", err)
		return
	}
	goodTokenReq := httptest.NewRequest(
		http.MethodPost, "/", bytes.NewBuffer(testPayloadBytes))
	goodTokenReq.Header.Add("Authorization", goodTokenHeader)

	badIssuerTokenSource, err := idtoken.NewTokenSource(ctx, "bad audience")
	if err != nil {
		t.Errorf("can't initialize ID token source: %v", err)
		return
	}
	badToken, err := badIssuerTokenSource.Token()
	if err != nil {
		t.Errorf("error getting ID token: %v", err)
		return
	}
	badTokenHeader := fmt.Sprintf(
		"%s %s", badToken.TokenType, badToken.AccessToken)
	badTokenReq := httptest.NewRequest(
		http.MethodPost, "/", bytes.NewBuffer(testPayloadBytes))
	badTokenReq.Header.Add("Authorization", badTokenHeader)

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *base.PubSubPayload
		wantErr bool
	}{
		{
			name: "no authorization header",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(http.MethodPost, "/", nil),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid authorization header",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: invalidHeaderReq,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid authorization header",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: goodTokenReq,
			},
			want:    testPayload,
			wantErr: false,
		},
		{
			name: "valid payload with wrong audience",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: badTokenReq,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.VerifyPubSubJWTAndDecodePayload(tt.args.w, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyPubSubJWTAndDecodePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyPubSubJWTAndDecodePayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

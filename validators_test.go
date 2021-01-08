package base_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func GetFirestoreClient(t *testing.T) *firestore.Client {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)

	ctx := base.GetAuthenticatedContext(t)
	firestoreClient, err := firebaseApp.Firestore(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, firestoreClient)
	return firestoreClient
}

func TestValidateCoordinates(t *testing.T) {
	type args struct {
		coordinates string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		want1   float64
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				coordinates: "-1.2881361,36.7815616",
			},
			want:    -1.2881361,
			want1:   36.7815616,
			wantErr: false,
		},
		{
			name: "wrong input format - no comma separator",
			args: args{
				coordinates: "not a valid input format",
			},
			wantErr: true,
		},
		{
			name: "wrong input format - unparseable lat",
			args: args{
				coordinates: "a,1",
			},
			wantErr: true,
		},
		{
			name: "wrong input format - unparseable long",
			args: args{
				coordinates: "1,b",
			},
			wantErr: true,
		},
		{
			name: "wrong input format - lat out of range",
			args: args{
				coordinates: "-98,2.3",
			},
			wantErr: true,
		},
		{
			name: "wrong input format - long out of range",
			args: args{
				coordinates: "-80,201",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := base.ValidateCoordinates(tt.args.coordinates)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCoordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateCoordinates() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateCoordinates() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestValidateMSISDN(t *testing.T) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)

	ctx := base.GetAuthenticatedContext(t)
	firestoreClient, err := firebaseApp.Firestore(ctx)
	assert.Nil(t, err)

	otpMsisdn := "+254722000000"
	normalized, err := base.NormalizeMSISDN(otpMsisdn)
	assert.Nil(t, err)

	validOtpCode := rand.Int()
	validOtpData := map[string]interface{}{
		"authorizationCode": strconv.Itoa(validOtpCode),
		"isValid":           true,
		"message":           "testing OTP message",
		"msisdn":            normalized,
		"timestamp":         time.Now(),
	}
	_, err = base.SaveDataToFirestore(firestoreClient, base.SuffixCollection(base.OTPCollectionName), validOtpData)
	assert.Nil(t, err)

	invalidOtpCode := rand.Int()
	invalidOtpData := map[string]interface{}{
		"authorizationCode": strconv.Itoa(invalidOtpCode),
		"isValid":           false,
		"message":           "testing OTP message",
		"msisdn":            normalized,
		"timestamp":         time.Now(),
	}
	_, err = base.SaveDataToFirestore(firestoreClient, base.SuffixCollection(base.OTPCollectionName), invalidOtpData)
	assert.Nil(t, err)

	type args struct {
		msisdn           string
		verificationCode string
		isUSSD           bool
		firestoreClient  *firestore.Client
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "invalid phone format",
			args: args{
				msisdn: "not a valid phone format",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "ussd session validation",
			args: args{
				msisdn:           "0722000000",
				verificationCode: uuid.NewV1().String(),
				isUSSD:           true,
				firestoreClient:  firestoreClient,
			},
			want:    "+254722000000",
			wantErr: false,
		},
		{
			name: "non existent verification code for non USSD",
			args: args{
				msisdn:           "0722000000",
				verificationCode: uuid.NewV1().String(),
				isUSSD:           false,
				firestoreClient:  firestoreClient,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "valid verification code for non USSD",
			args: args{
				msisdn:           "0722000000",
				verificationCode: strconv.Itoa(validOtpCode),
				isUSSD:           false,
				firestoreClient:  firestoreClient,
			},
			want:    "+254722000000",
			wantErr: false,
		},
		{
			name: "used (invalid) verification code for non USSD",
			args: args{
				msisdn:           "0722000000",
				verificationCode: strconv.Itoa(invalidOtpCode),
				isUSSD:           false,
				firestoreClient:  firestoreClient,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.ValidateMSISDN(tt.args.msisdn, tt.args.verificationCode, tt.args.isUSSD, tt.args.firestoreClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateMSISDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)

	ctx := base.GetAuthenticatedContext(t)
	firestoreClient, err := firebaseApp.Firestore(ctx)
	assert.Nil(t, err)

	type args struct {
		email           string
		optIn           bool
		firestoreClient *firestore.Client
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "first valid email, opted in",
			args: args{
				email:           "ngure.nyaga@savannahinformatics.com",
				optIn:           true,
				firestoreClient: firestoreClient,
			},
			wantErr: false,
		},
		{
			name: "second valid email, opted in",
			args: args{
				email:           "ngure.nyaga@healthcloud.com",
				optIn:           true,
				firestoreClient: firestoreClient,
			},
			wantErr: false,
		},
		{
			name: "third valid email,  notopted in",
			args: args{
				email:           "ngurenyaga@gmail.com",
				optIn:           true,
				firestoreClient: firestoreClient,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			args: args{
				email:           "not a valid email",
				optIn:           true,
				firestoreClient: firestoreClient,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := base.ValidateEmail(tt.args.email, tt.args.optIn, tt.args.firestoreClient); (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsMSISDNValid(t *testing.T) {

	tests := []struct {
		name   string
		msisdn string
		want   bool
	}{
		{
			name:   "valid : kenyan with code",
			msisdn: "+254722000000",
			want:   true,
		},
		{
			name:   "valid : kenyan with code and spaces",
			msisdn: "+254 722 000 000",
			want:   true,
		},
		{
			name:   "valid : kenyan without code",
			msisdn: "0722000000",
			want:   true,
		},
		{
			name:   "valid : kenyan without code and spaces",
			msisdn: "0722 000 000",
			want:   true,
		},
		{
			name:   "valid : kenyan without plus sign",
			msisdn: "+254722000000",
			want:   true,
		},
		{
			name:   "valid : kenyan without plus sign and spaces",
			msisdn: "+254 722 000 000",
			want:   true,
		},
		{
			name:   "invalid : kenyan with alphanumeric1",
			msisdn: "+25472abc0000",
			want:   false,
		},
		{
			name:   "invalid : kenyan with alphanumeric2",
			msisdn: "072abc0000",
			want:   false,
		},
		{
			name:   "invalid : kenyan with unwanted characters : asterisk",
			msisdn: "072*120000",
			want:   false,
		},
		{
			name:   "invalid : kenyan without code with plus sign as prefix",
			msisdn: "+0722000000",
			want:   false,
		},
		{
			name:   "valid : international with code",
			msisdn: "(+351)282435050",
			want:   true,
		},
		{
			name:   "valid : international with code and spaces",
			msisdn: "(+351) 282 43 50 50",
			want:   true,
		},
		{
			name:   "valid : international with extension",
			msisdn: "1-234-567-8901 ext1234",
			want:   true,
		},
		{
			name:   "ivalid : international with alphanumeric",
			msisdn: "90191919qwe",
			want:   false,
		},
		{
			name:   "invalid : international with unwanted characters : asterisk",
			msisdn: "(+351) 282 *3 50 50",
			want:   false,
		},
		{
			name:   "invalid : international with unwanted characters : assorted",
			msisdn: "(+351) $82 *3 50 50",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.IsMSISDNValid(tt.msisdn); got != tt.want {
				t.Errorf("IsMSISDNValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustNormalizeMSISDN(t *testing.T) {
	type args struct {
		msisdn string
	}
	validOutput := "+254722000000"
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{
			name: "valid : E164 input",
			args: args{
				msisdn: "+254722000000",
			},
			want:    &validOutput,
			wantErr: false,
		},
		{
			name: "valid: non E164 input",
			args: args{
				msisdn: "0722000000",
			},
			want:    &validOutput,
			wantErr: false,
		},
		{
			name: "invalid: non E164 input with alphanumeric",
			args: args{
				msisdn: "0722abc000",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid: E164 input with alphanumeric",
			args: args{
				msisdn: "+254722wer0000",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.MustNormalizeMSISDN(tt.args.msisdn)
			if tt.wantErr {
				assert.Nil(t, got)
				assert.NotNil(t, err)
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, *tt.want, *got)
			}
		})
	}
}

func TestIntSliceContains(t *testing.T) {
	type args struct {
		s []int
		e int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "slice which contains the int",
			args: args{
				s: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
				e: 7,
			},
			want: true,
		},
		{
			name: "slice which does NOT contain the int",
			args: args{
				s: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
				e: 79,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.IntSliceContains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("IntSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAndSaveMSISDN(t *testing.T) {
	fc := GetFirestoreClient(t)

	type args struct {
		msisdn           string
		verificationCode string
		isUSSD           bool
		optIn            bool
		firestoreClient  *firestore.Client
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "invalid phone number/OTP, non USSD",
			args: args{
				msisdn:           "0722000000",
				verificationCode: "not a real one",
				isUSSD:           false,
				optIn:            true,
				firestoreClient:  fc,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "valid phone number, USSD, opt in true",
			args: args{
				msisdn:           "0722000000",
				verificationCode: "this is a ussd session ID from the telco",
				isUSSD:           true,
				optIn:            true,
				firestoreClient:  fc,
			},
			want:    "+254722000000",
			wantErr: false,
		},
		{
			name: "valid phone number, USSD, opt in false",
			args: args{
				msisdn:           "0722000000",
				verificationCode: "this is a ussd session ID from the telco",
				isUSSD:           true,
				optIn:            false,
				firestoreClient:  fc,
			},
			want:    "+254722000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.ValidateAndSaveMSISDN(tt.args.msisdn, tt.args.verificationCode, tt.args.isUSSD, tt.args.optIn, tt.args.firestoreClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndSaveMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateAndSaveMSISDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceContains(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string found in slice",
			args: args{
				s: []string{"a", "b", "c", "d", "e"},
				e: "a",
			},
			want: true,
		},
		{
			name: "string not found in slice",
			args: args{
				s: []string{"a", "b", "c", "d", "e"},
				e: "z",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := base.StringSliceContains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalizeMSISDN(t *testing.T) {
	type args struct {
		msisdn string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "good Kenyan number, full E164 format",
			args: args{
				"+254723002959",
			},
			want:    "+254723002959",
			wantErr: false,
		},
		{
			name: "good Kenyan number, no + prefix",
			args: args{
				"254723002959",
			},
			want:    "+254723002959",
			wantErr: false,
		},
		{
			name: "good Kenyan number, no international dialling code",
			args: args{
				"0723002959",
			},
			want:    "+254723002959",
			wantErr: false,
		},
		{
			name: "good US number, full E164 format",
			args: args{
				"+16125409037",
			},
			want:    "+16125409037",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := base.NormalizeMSISDN(tt.args.msisdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("NormalizeMSISDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

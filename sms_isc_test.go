package base_test

import (
	"testing"

	"gitlab.slade360emr.com/go/base"
)

const testPhone = "+254723002959"

func TestSendSMS(t *testing.T) {

	// Note: This is a very brittle test case.
	// Any change to the service urls would probably lead to a failure
	// There's probably a better way to do this (Mocking *wink wink)
	// But I (Farad) felt this is the best way of doing it i.e. Acceptance Testing
	newSmsIsc, _ := base.NewInterserviceClient(base.ISCService{
		Name:       "sms",
		RootDomain: "https://sms-staging.healthcloud.co.ke",
	})

	newTwilioIsc, _ := base.NewInterserviceClient(base.ISCService{
		Name:       "twilio",
		RootDomain: "https://twilio-staging.healthcloud.co.ke",
	})

	smsEndPoint := "internal/send_sms"

	type args struct {
		phoneNumbers    []string
		message         string
		smsIscClient    base.SmsISC
		twilioIscClient base.SmsISC
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good test case",
			args: args{
				phoneNumbers: []string{testPhone},
				message:      "Test Text Message",
				smsIscClient: base.SmsISC{
					Isc:      newSmsIsc,
					EndPoint: "internal/send_sms",
				},
				twilioIscClient: base.SmsISC{
					Isc:      newTwilioIsc,
					EndPoint: smsEndPoint,
				},
			},
			wantErr: false,
		},
		{
			name: "bad test case: Empty Message",
			args: args{
				phoneNumbers: []string{testPhone},
				message:      "",
				smsIscClient: base.SmsISC{
					Isc:      newSmsIsc,
					EndPoint: smsEndPoint,
				},
				twilioIscClient: base.SmsISC{
					Isc:      newTwilioIsc,
					EndPoint: smsEndPoint,
				},
			},
			wantErr: true,
		},
		{
			name: "bad test case: No Phone Numbers",
			args: args{
				phoneNumbers: []string{},
				message:      "Test Text Message",
				smsIscClient: base.SmsISC{
					Isc:      newSmsIsc,
					EndPoint: smsEndPoint,
				},
				twilioIscClient: base.SmsISC{
					Isc:      newTwilioIsc,
					EndPoint: smsEndPoint,
				},
			},
			wantErr: true,
		},
		{
			name: "bad test case: Invalid Phone Numbers",
			args: args{
				phoneNumbers: []string{"not-a-number"},
				message:      "Test Text Message",
				smsIscClient: base.SmsISC{
					Isc:      newSmsIsc,
					EndPoint: smsEndPoint,
				},
				twilioIscClient: base.SmsISC{
					Isc:      newTwilioIsc,
					EndPoint: smsEndPoint,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := base.SendSMS(tt.args.phoneNumbers, tt.args.message, tt.args.smsIscClient, tt.args.twilioIscClient); (err != nil) != tt.wantErr {
				t.Errorf("SendSMS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

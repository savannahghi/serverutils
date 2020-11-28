package base

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

// SmsISC is a representation of an ISC client
type SmsISC struct {
	Isc      *InterServiceClient
	EndPoint string
}

// SendSMS is send a text message to specified phone No.s both local and foreign
func SendSMS(phoneNumbers []string, message string, smsClient, twilioClient SmsISC) error {

	if message == "" {
		return fmt.Errorf("sms not sent: `message` needs to be supplied")
	}

	foreignPhoneNos := []string{}
	localPhoneNos := []string{}

	for _, phone := range phoneNumbers {
		if IsKenyanNumber(phone) {
			localPhoneNos = append(localPhoneNos, phone)
			continue
		}
		foreignPhoneNos = append(foreignPhoneNos, phone)
	}

	if len(localPhoneNos) < 1 && len(foreignPhoneNos) < 1 {
		return fmt.Errorf("sms not sent: `phone numbers` need to be supplied")
	}

	if len(foreignPhoneNos) >= 1 {
		err := makeRequest(foreignPhoneNos, message, twilioClient.EndPoint, *twilioClient.Isc)
		if err != nil {
			return fmt.Errorf("sms not sent: %v", err)
		}
	}

	if len(localPhoneNos) >= 1 {
		err := makeRequest(localPhoneNos, message, smsClient.EndPoint, *smsClient.Isc)
		if err != nil {
			return fmt.Errorf("sms not sent: %v", err)
		}
	}

	return nil
}

func makeRequest(phoneNumbers []string, message, EndPoint string, client InterServiceClient) error {
	payload := map[string]interface{}{
		"to":      phoneNumbers,
		"message": message,
	}
	resp, err := client.MakeRequest(http.MethodPost, EndPoint, payload)
	if err != nil {
		return err
	}
	if IsDebug() {
		b, _ := httputil.DumpResponse(resp, true)
		log.Println(string(b))
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send SMS : %w, with status code %v", err, resp.StatusCode)
	}
	return nil
}

//IsKenyanNumber checks if phone number belongs to KENYA TELECOM
func IsKenyanNumber(phoneNumber string) bool {
	return strings.HasPrefix(phoneNumber, "+254")
}

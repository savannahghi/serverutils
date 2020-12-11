package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	// VerifyOTPEndPoint ISC endpoint to verify OTP
	VerifyOTPEndPoint = "internal/verify_otp/"
	// SendOTPEndPoint ISC endpoint to sent OTP
	SendOTPEndPoint = "internal/send_otp/"
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

// VerifyOTP confirms a phone number is valid by verifying the code that was sent to the number
func VerifyOTP(msisdn string, otp string, otpClient *InterServiceClient) (bool, error) {
	if otpClient == nil {
		return false, fmt.Errorf("nil OTP client")
	}

	_, err := NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("invalid phone format: %w", err)
	}

	type VerifyOTP struct {
		Msisdn           string `json:"msisdn"`
		VerificationCode string `json:"verificationCode"`
	}

	verifyPayload := VerifyOTP{
		Msisdn:           msisdn,
		VerificationCode: otp,
	}

	resp, err := otpClient.MakeRequest(http.MethodPost, VerifyOTPEndPoint, verifyPayload)
	if err != nil {
		return false, fmt.Errorf(
			"can't complete OTP verification request: %w", err)
	}

	if IsDebug() {
		b, _ := httputil.DumpResponse(resp, true)
		log.Println(string(b))
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unable to verify OTP : %w, with status code %v", err, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("can't read OTP response data: %w", err)
	}

	type otpResponse struct {
		IsVerified bool `json:"IsVerified"`
	}

	var r otpResponse
	err = json.Unmarshal(data, &r)
	if err != nil {
		return false, fmt.Errorf(
			"can't unmarshal OTP response data from JSON: %w", err)
	}

	return r.IsVerified, nil
}

// SendOTPHelper is a helper used in tests to send OTP to a test number
func SendOTPHelper(msisdn string, otpClient *InterServiceClient) (string, error) {
	// we prepare the OTP payload
	payload := map[string]interface{}{
		"msisdn": msisdn,
	}
	// make the request
	resp, err := otpClient.MakeRequest(http.MethodPost, SendOTPEndPoint, payload)
	if err != nil {
		return "", fmt.Errorf("unable to make a send otp request: %w", err)
	}
	// inspect the response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unable to generate otp, with status code %v", resp.StatusCode)
	}
	// read the response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to convert response to string: %v", err)
	}
	// reset the response body to the original unread state so that decode can
	// continue
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// Decode and store result to a variable
	type otpResponse struct {
		OTP string `json:"otp"`
	}
	// store the response in a variable and return
	OTPResp := &otpResponse{}
	if err := json.NewDecoder(resp.Body).Decode(OTPResp); err != nil {
		return "", fmt.Errorf("InternalServerError: unable to decode verify OTP response: %v", err)
	}
	return OTPResp.OTP, nil
}

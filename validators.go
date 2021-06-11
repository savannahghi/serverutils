package go_utils

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/asaskevich/govalidator"
	"github.com/ttacon/libphonenumber"
)

// ValidateEmail returns an error if the supplied string does not have a
// valid format or resolvable host
func ValidateEmail(
	email string, optIn bool, firestoreClient *firestore.Client) error {
	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	if optIn {
		data := EmailOptIn{
			Email:   email,
			OptedIn: optIn,
		}
		_, err := SaveDataToFirestore(
			firestoreClient, EmailOptInCollectionName, data)
		if err != nil {
			return fmt.Errorf("unable to save email opt in: %v", err)
		}
	}
	return nil
}

// IsMSISDNValid uses regular expression to validate the a phone number
func IsMSISDNValid(msisdn string) bool {
	if len(msisdn) < 10 {
		return false
	}
	reKen := regexp.MustCompile(`^(?:254|\+254|0)?((7|1)(?:(?:[129][0-9])|(?:0[0-8])|(4[0-1]))[0-9]{6})$`)
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	if !reKen.MatchString(msisdn) {
		return re.MatchString(msisdn)
	}
	return reKen.MatchString(msisdn)
}

// NormalizeMSISDN validates the input phone number.
// For valid phone numbers, it normalizes them to international format
// e.g +2547........
func NormalizeMSISDN(msisdn string) (*string, error) {
	if !IsMSISDNValid(msisdn) {
		return nil, fmt.Errorf("invalid phone number: %s", msisdn)
	}
	num, err := libphonenumber.Parse(msisdn, defaultRegion)
	if err != nil {
		return nil, err
	}
	formatted := libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
	cleaned := strings.ReplaceAll(formatted, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	return &cleaned, nil
}

// MustNormalizeMSISDN validates the input phone number otherwise it panics
func MustNormalizeMSISDN(msisdn string) (*string, error) {
	if !IsMSISDNValid(msisdn) {
		return nil, fmt.Errorf("invalid phone number: %s", msisdn)
	}
	num, err := libphonenumber.Parse(msisdn, defaultRegion)
	if err != nil {
		log.Panic(fmt.Sprintf("invalid phone number: %s", msisdn))
	}
	formatted := libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
	cleaned := strings.ReplaceAll(formatted, " ", "")
	return &cleaned, nil
}

// ValidateMSISDN returns an error if the MSISDN format is wrong or the
// supplied verification code is not valid
//
// Deprecated: Should implement `VerifyOTP` instead. This helps to confirm if a phonenumber
// is valid by verifying the code sent to it.
func ValidateMSISDN(
	msisdn, verificationCode string,
	isUSSD bool, firestoreClient *firestore.Client) (string, error) {

	// check the format
	normalized, err := NormalizeMSISDN(msisdn)
	if err != nil {
		return "", fmt.Errorf("invalid phone format: %v", err)
	}

	// save a USSD log for USSD registrations
	if isUSSD {
		log := USSDSessionLog{
			MSISDN:    msisdn,
			SessionID: verificationCode,
		}
		_, err = SaveDataToFirestore(
			firestoreClient, SuffixCollection(USSDSessionCollectionName), log)
		if err != nil {
			return "", fmt.Errorf("unable to save USSD session: %v", err)
		}
		return *normalized, nil
	}

	// check if the OTP is on file / known
	query := firestoreClient.Collection(SuffixCollection(OTPCollectionName)).Where(
		"isValid", "==", true,
	).Where(
		"msisdn", "==", normalized,
	).Where(
		"authorizationCode", "==", verificationCode,
	)
	ctx := context.Background()
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve verification codes: %v", err)
	}
	if len(docs) == 0 {
		return "", fmt.Errorf("no matching verification codes found")
	}

	for _, doc := range docs {
		otpData := doc.Data()
		otpData["isValid"] = false
		err = UpdateRecordOnFirestore(
			firestoreClient, SuffixCollection(OTPCollectionName), doc.Ref.ID, otpData)
		if err != nil {
			return "", fmt.Errorf("unable to save updated OTP document: %v", err)
		}
	}

	return *normalized, nil
}

// ValidateAndSaveMSISDN returns an error if the MSISDN format is wrong or the
// supplied verification code is not valid
func ValidateAndSaveMSISDN(
	msisdn, verificationCode string, isUSSD bool, optIn bool,
	firestoreClient *firestore.Client) (string, error) {
	validated, err := ValidateMSISDN(
		msisdn, verificationCode, isUSSD, firestoreClient)
	if err != nil {
		return "", fmt.Errorf("invalid MSISDN: %s", err)
	}
	if optIn {
		data := PhoneOptIn{
			MSISDN:  validated,
			OptedIn: optIn,
		}
		_, err = SaveDataToFirestore(
			firestoreClient, PhoneOptInCollectionName, data)
		if err != nil {
			return "", fmt.Errorf("unable to save email opt in: %v", err)
		}
	}
	return validated, nil
}

// StringSliceContains tests if a string is contained in a slice of strings
func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// IntSliceContains tests if a string is contained in a slice of strings
func IntSliceContains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ValidateCoordinates takes a coordinates string (should be similar to "-1.2881361,36.7815616")
// validates it, parses it and returns it as a pair of floats.
//
// If the validation or parsing fails, an error is returned.
func ValidateCoordinates(coordinates string) (float64, float64, error) {
	latlong := strings.Split(coordinates, ",")
	if len(latlong) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinates; expected two parts separated by a comma")
	}

	latStr := strings.TrimSpace(latlong[0])
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse lat: %w", err)
	}
	if lat < -90 || lat > 90 {
		return 0, 0, fmt.Errorf("latitude out of range, expected a value between -90 and 90")
	}

	longStr := strings.TrimSpace(latlong[1])
	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse long: %w", err)
	}
	if long < -180 || long > 180 {
		return 0, 0, fmt.Errorf("longitude out of range, expected a value between -180 and 180")
	}

	return lat, long, nil
}

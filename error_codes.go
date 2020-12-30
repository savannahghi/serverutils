package base

// ErrorCode are  used to determine the nature of an error, and why it occurred
// both the frontend and backend should be aware of these codes
type ErrorCode int

const (
	// OK is returned on success.
	OK ErrorCode = iota + 1

	// Internal errors means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	// it's value is 2
	Internal

	// UndefinedArguments errors means either one or more arguments to
	// a method have not been specified
	// it's value is 3
	UndefinedArguments

	// PhoneNumberInUse indicates that a phone number has an associated user profile.
	// this error can occur when fetching a user profile using a phone number, to check
	// that the phone number has not already been registered. The check usually runs
	// on both PRIMARY PHONE and SECONDARY PHONE
	// it's value is 4
	PhoneNumberInUse

	// EmailAddressInUse indicates that an email address has an associated user profile.
	// this error can occur when fetching a user profile using an email address, to check
	// that the email address has not already been registered. The check usually runs
	// on both PRIMARY EMAIL ADDRESS and SECONDARY EMAIL ADDRESS.
	// it's value is 5
	EmailAddressInUse

	// UsernameInUse indicates that a username has an associated user profile.
	// this error can occur when trying a update a user's username with a username that already has been taken
	// it's value is 6
	UsernameInUse

	// ProfileNotFound errors means a user profile does not exist with the provided parameters
	// This occures when fetching a user profile either by UID, ID , PHONE NUMBER or EMAIL and no
	// matching record is found
	// it's value is 7
	ProfileNotFound

	// PINMismatch errors means that the provided PINS do not match (are not similar)
	// it's value is 8
	PINMismatch

	// PINNotFound errors means a user PIN does not exist with the provided parameters
	// This occurs when fetching a PIN by the user's user profile ID and no
	// matching record is found. This should never occur and if it does then i means
	// there is a serious issue with our data
	// it's value is 9
	PINNotFound

	// UserNotFound errors means that a user's firebase auth account does not exists. This occurs
	// when fectching a firebase user by either a phone number or an email and their record is not found
	// it's value is 10
	UserNotFound
)

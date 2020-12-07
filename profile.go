package base

// Cover is used to save a user's insurance details.
type Cover struct {
	PayerName      string `json:"payerName,omitempty" firestore:"payerName"`
	PayerSladeCode int    `json:"payerSladeCode,omitempty" firestore:"payerSladeCode"`
	MemberNumber   string `json:"memberNumber,omitempty" firestore:"memberNumber"`
	MemberName     string `json:"memberName,omitempty" firestore:"memberName"`
}

// VerifiedEmail represents an email and it's verification status
type VerifiedEmail struct {
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

// VerifiedMsisdn represents an E164 phone number and it's verification status
type VerifiedMsisdn struct {
	Msisdn   string `json:"msisdn"`
	Verified bool   `json:"verified"`
}

// IsEntity marks this struct as a GraphQL entity
func (c Cover) IsEntity() {}

// UserProfile serializes the profile of the logged in user.
type UserProfile struct {
	// globally unique identifier for a profile
	ID string `json:"id" firestore:"id"`

	// `VerifiedIdentifiers` represent various ways the user has been able to login
	// and these providers point to the same user
	VerifiedIdentifiers []string `json:"verified_identifiers" firestore:"verifiedIdentifiers"`

	// a profile contains a user's contact details
	Msisdns        []string         `json:"msisdns,omitempty" firestore:"msisdns"`
	Emails         []string         `json:"emails,omitempty" firestore:"emails"`
	PushTokens     []string         `json:"pushTokens,omitempty" firestore:"pushTokens"`
	VerifiedEmails []VerifiedEmail  `json:"verifiedEmails,omitempty" firestore:"verifiedEmails"`
	VerifiedPhones []VerifiedMsisdn `json:"verifiedPhones,omitempty" firestore:"verifiedPhones"`

	// we determine if a user is "live" by examining fields on their profile
	TermsAccepted                      bool  `json:"termsAccepted,omitempty" firestore:"termsAccepted"`
	IsApproved                         bool  `json:"isApproved,omitempty" firestore:"isApproved"`
	Active                             bool  `json:"active" firestore:"active"`
	PractitionerApproved               *bool `json:"practitionerApproved,omitempty" firestore:"practitionerApproved"`
	PractitionerTermsOfServiceAccepted *bool `json:"practitionerTermsOfServiceAccepted,omitempty" firestore:"practitionerTermsOfServiceAccepted"`

	// a user's profile photo can be stored as base 64 encoded PNG
	PhotoBase64      string      `json:"photoBase64,omitempty" firestore:"photoBase64"`
	PhotoContentType ContentType `json:"photoContentType,omitempty" firestore:"photoContentType"`

	// a user can have zero or more insurance covers
	Covers []Cover `json:"covers,omitempty" firestore:"covers"`

	// a user's biodata is stored on the profile
	Name        *string  `json:"name,omitempty" firestore:"name"`
	Bio         *string  `json:"bio,omitempty" firestore:"bio"`
	DateOfBirth *Date    `json:"dateOfBirth,omitempty" firestore:"dateOfBirth,omitempty"`
	Gender      *Gender  `json:"gender,omitempty" firestore:"gender,omitempty"`
	Language    Language `json:"language,omitempty" firestore:"language"`
	PatientID   *string  `json:"patientID,omitempty" firestore:"patientID"`

	// testers are whitelisted via their profiles
	IsTester                   bool `json:"isTester,omitempty" firestore:"isTester"`
	CanExperiment              bool `json:"canExperiment,omitempty" firestore:"canExperiment"`
	AskAgainToSetIsTester      bool `json:"askAgainToSetIsTester,omitempty" firestore:"askAgainToSetIsTester"`
	AskAgainToSetCanExperiment bool `json:"askAgainToSetCanExperiment,omitempty" firestore:"askAgainToSetCanExperiment"`

	// these flags are used to determine paths to take/omit in the UI
	HasPin                  bool `json:"hasPin,omitempty" firestore:"hasPin"`
	HasSupplierAccount      bool `json:"hasSupplierAccount,omitempty" firestore:"hasSupplierAccount"`
	HasCustomerAccount      bool `json:"hasCustomerAccount,omitempty" firestore:"hasCustomerAccount"`
	PractitionerHasServices bool `json:"practitionerHasServices,omitempty" firestore:"practitionerHasServices"`
}

// IsEntity marks a profile as a GraphQL entity
func (u UserProfile) IsEntity() {}

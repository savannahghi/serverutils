package base

// Cover is used to save a user's insurance details.
type Cover struct {
	PayerName      string `json:"payer_name,omitempty" firestore:"payerName"`
	PayerSladeCode int    `json:"payer_slade_code,omitempty" firestore:"payerSladeCode"`
	MemberNumber   string `json:"member_number,omitempty" firestore:"memberNumber"`
	MemberName     string `json:"member_name,omitempty" firestore:"memberName"`
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
	PushTokens     []string         `json:"push_tokens,omitempty" firestore:"pushTokens"`
	VerifiedEmails []VerifiedEmail  `json:"verified_emails,omitempty" firestore:"verifiedEmails"`
	VerifiedPhones []VerifiedMsisdn `json:"verified_phones,omitempty" firestore:"verifiedPhones"`

	// should be true for admins
	IsAdmin bool `json:"isAdmin" firestore:"isAdmin"`

	// we determine if a user is "live" by examining fields on their profile
	TermsAccepted                      bool  `json:"terms_accepted,omitempty" firestore:"termsAccepted"`
	IsApproved                         bool  `json:"is_approved" firestore:"isApproved"`
	Active                             bool  `json:"active" firestore:"active"`
	PractitionerApproved               *bool `json:"practitioner_approved,omitempty" firestore:"practitionerApproved"`
	PractitionerTermsOfServiceAccepted *bool `json:"practitioner_term_of_service_accepted,omitempty" firestore:"practitionerTermsOfServiceAccepted"`

	// a user's profile photo can be stored as base 64 encoded PNG
	PhotoBase64      string      `json:"photo_base64,omitempty" firestore:"photoBase64"`
	PhotoContentType ContentType `json:"photo_content_type,omitempty" firestore:"photoContentType"`

	// a user can have zero or more insurance covers
	Covers []Cover `json:"covers,omitempty" firestore:"covers"`

	// a user's biodata is stored on the profile
	Name        *string  `json:"name,omitempty" firestore:"name"`
	Bio         *string  `json:"bio,omitempty" firestore:"bio"`
	DateOfBirth *Date    `json:"date_of_birth,omitempty" firestore:"dateOfBirth,omitempty"`
	Gender      *Gender  `json:"gender,omitempty" firestore:"gender,omitempty"`
	Language    Language `json:"language,omitempty" firestore:"language"`
	PatientID   *string  `json:"patient_id,omitempty" firestore:"patientID"`

	// testers are whitelisted via their profiles
	IsTester                   bool `json:"is_tester" firestore:"isTester"`
	CanExperiment              bool `json:"can_experiment,omitempty" firestore:"canExperiment"`
	AskAgainToSetIsTester      bool `json:"ask_again_to_set_is_tester,omitempty" firestore:"askAgainToSetIsTester"`
	AskAgainToSetCanExperiment bool `json:"ask_again_to_set_can_experiment,omitempty" firestore:"askAgainToSetCanExperiment"`

	// these flags are used to determine paths to take/omit in the UI
	HasPin                  bool `json:"has_pin" firestore:"hasPin"`
	HasSupplierAccount      bool `json:"has_supplier_account" firestore:"hasSupplierAccount"`
	HasCustomerAccount      bool `json:"has_customer_account" firestore:"hasCustomerAccount"`
	PractitionerHasServices bool `json:"practitioner_has_services,omitempty" firestore:"practitionerHasServices"`
}

// IsEntity marks a profile as a GraphQL entity
func (u UserProfile) IsEntity() {}

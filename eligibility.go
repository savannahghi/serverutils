package base

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"time"
)

// Eligibility (de)serializes eligibility information from te Slade 360 EDI
// eligibility v3 base.
//
// The API version matters: a major revision of the API is likely to introduce
// a new benefit structure.
//
// This structure has not mapped all the fields that are in EDI eligibility.
// Those that are not important to/for the mobile app have been omitted.
type Eligibility struct {
	Message         string                   `json:"message"`
	Status          string                   `json:"status"`
	StatusReason    string                   `json:"statusReason"`
	EligibilityTime *time.Time               `json:"eligibilityTime"`
	Benefits        []*EligibilityBenefit    `json:"benefits"`
	Cover           *EligibilityCoverDetails `json:"cover"`
	PayerDetails    *EligibilityPayer        `json:"payerDetails"`
	Member          *EligibilityMember       `json:"member"`
	Inclusions      []*string                `json:"inclusions"`
	Exclusions      []*string                `json:"exclusions"`
	Dependants      []*Dependant             `json:"dependants"`
}

// EligibilityBenefit matches the structure of eligibility information from
// the Slade 360 EDI Eligibility v3 base.
//
// The API version matters: a major revision of the API is likely to introduce
// a new benefit structure.
type EligibilityBenefit struct {
	// these fields identify the exact benefit
	Benefit     string      `json:"benefit"`
	BenefitType BenefitType `json:"benefitType"`
	BenefitCode string      `json:"benefitCode"`
	ParentCode  *string     `json:"parentCode"`

	// the benefit status in human readable form. UIs should use ValiditySummary
	Status              string  `json:"status"`
	AvailabilitySummary *string `json:"availabilitySummary"` // computed
	ValiditySummmary    *string `json:"validitySummmary"`

	// these numbers explain how much is available/spent
	BenefitLimit     float64  `json:"benefitLimit"` // for wellness card only
	AvailableBalance *float64 `json:"availableBalance"`
	PayerBalance     *float64 `json:"payerBalance"`
	ThresholdAmount  *float64 `json:"thresholdAmount"`
	ReservedAmount   *float64 `json:"reservedAmount"`
	SpentAmount      *float64 `json:"spentAmount"`
	PercentageUsed   *float64 `json:"percentageUsed"`

	// these fields determine the guidance that is displayed for pre-auths
	PreauthRequired    bool     `json:"preauthRequired"`
	VisitLimit         float64  `json:"visitLimit"`
	PreauthThreshold   *Decimal `json:"preauthThreshold"`
	PreauthExplanation *string  `json:"preauthExplanation"`

	// these fields determine the guidance that is displayed for copayments
	CopayType      *string   `json:"copayType"`
	CopayAppliesTo []*string `json:"copayAppliesTo"`
	CopayValue     *float64  `json:"copayValue"`

	// should be set for shared benefits only
	SharedWith []*Dependant `json:"sharedWith"`

	// should be set to the relevant provider panel only
	AllowedProviders []*AllowedProvider `json:"allowedProviders"`

	// determine whether to show PROVIDERS the balance or not
	ShowBalance bool `json:"showBalance"`

	// other stuff
	HasExcessProtection *bool     `json:"hasExcessProtection"`
	Errors              []*string `json:"errors"`
}

// EligibilityCoverDetails serializes the cover details from the Slade 360 EDI
// eligibility v3 base.
//
// The API version matters: a major revision of the API is likely to introduce
// a new benefit structure.
//
// This structure has not mapped all EDI eligibility cover details. Those that
// are not important for the mobile app have been omitted.
type EligibilityCoverDetails struct {
	Group                 string    `json:"group"`
	GroupCode             string    `json:"groupCode"`
	GroupType             string    `json:"groupType"`
	IsEmployerGroup       bool      `json:"isEmployerGroup"`
	Category              string    `json:"category"`
	CategoryCode          string    `json:"categoryCode"`
	PolicyNumber          string    `json:"policyNumber"`
	EffectivePolicyNumber string    `json:"effectivePolicyNumber"`
	ValidFrom             time.Time `json:"validFrom"`
	ValidTo               time.Time `json:"validTo"`
}

// EligibilityPayer serializes the payer details from the Slade 360 EDI
// eligibility v3 base.
//
// The API version matters: a major revision of the API is likely to introduce
// a new benefit structure.
type EligibilityPayer struct {
	SladeCode int                    `json:"sladeCode"`
	BpType    string                 `json:"bpType"`
	Name      string                 `json:"name"`
	Contacts  map[string]interface{} `json:"contacts"`
}

// EligibilityMember (de)serializes membership information from the Slade 360
// eligibility v3 base.
//
// The API version matters: a major revision of the API is likely to introduce
// a new benefit structure.
type EligibilityMember struct {
	ID                    int                      `json:"id"`
	GUID                  string                   `json:"guid"`
	FirstName             string                   `json:"firstName"`
	LastName              string                   `json:"lastName"`
	Names                 string                   `json:"names"`
	BeneficiaryCode       string                   `json:"beneficiaryCode"`
	MembershipNumber      string                   `json:"membershipNumber"`
	Gender                string                   `json:"gender"`
	IsPrincipal           bool                     `json:"isPrincipal"`
	IsMinor               bool                     `json:"isMinor"`
	IsActive              bool                     `json:"isActive"`
	IsEnrolled            bool                     `json:"isEnrolled"`
	IsVip                 bool                     `json:"isVip"`
	IsOfFpCaptureAge      bool                     `json:"isOfFpCaptureAge"`
	OtpOnlyAllowed        bool                     `json:"otpOnlyAllowed"`
	DateOfBirth           *time.Time               `json:"dateOfBirth"`
	Title                 *string                  `json:"title"`
	OtherNames            *string                  `json:"otherNames"`
	PhotoURL              *string                  `json:"photoURL"`
	PrincipalRelationship *string                  `json:"principalRelationship"`
	IDNo                  *string                  `json:"idNo"`
	NhifNo                *string                  `json:"nhifNo"`
	HasSladeCard          []interface{}            `json:"hasSladeCard"`
	Identifiers           []*EligibilityIdentifier `json:"identifiers"`
	Contacts              []*EligibilityContact    `json:"contacts"`
	PrincipalMember       *PrincipalMember         `json:"principalMember"`
	HasCard               *bool                    `json:"hasCard"`
	CardTypes             []*string                `json:"cardTypes"`
}

// PrincipalMember (de)serializes principal member information from the Slade 360
// eligibility v3 base.
//
// The API version matters: a major revision of the API is likely to introduce
// a new benefit structure.
type PrincipalMember struct {
	ID                    int                      `json:"id"`
	GUID                  string                   `json:"guid"`
	FirstName             string                   `json:"firstName"`
	LastName              string                   `json:"lastName"`
	Names                 string                   `json:"names"`
	BeneficiaryCode       string                   `json:"beneficiaryCode"`
	MembershipNumber      string                   `json:"membershipNumber"`
	Gender                string                   `json:"gender"`
	IsPrincipal           bool                     `json:"isPrincipal"`
	IsMinor               bool                     `json:"isMinor"`
	IsActive              bool                     `json:"isActive"`
	IsEnrolled            bool                     `json:"isEnrolled"`
	IsVip                 bool                     `json:"isVip"`
	IsOfFpCaptureAge      bool                     `json:"isOfFpCaptureAge"`
	OtpOnlyAllowed        bool                     `json:"otpOnlyAllowed"`
	DateOfBirth           *time.Time               `json:"dateOfBirth"`
	Title                 *string                  `json:"title"`
	OtherNames            *string                  `json:"otherNames"`
	PhotoURL              *string                  `json:"photoURL"`
	PrincipalRelationship *string                  `json:"principalRelationship"`
	IDNo                  *string                  `json:"idNo"`
	NhifNo                *string                  `json:"nhifNo"`
	HasSladeCard          []interface{}            `json:"hasSladeCard"`
	Identifiers           []*EligibilityIdentifier `json:"identifiers"`
	Contacts              []*EligibilityContact    `json:"contacts"`
}

// EligibilityContact serializes beneficiary eligibility contact information
// from the Slade 360 eligibility v3 base.
//
// The API version matters: a major revision of the API is likely to introduce
// a new benefit structure.
type EligibilityContact struct {
	ID            int    `json:"id"`
	GUID          string `json:"guid"`
	IsConfirmed   bool   `json:"isConfirmed"`
	IsMainContact bool   `json:"isMainContact"`
	IsVerified    bool   `json:"isVerified"`
	ContactType   string `json:"contactType"`
	ContactValue  string `json:"contactValue"`
	Active        bool   `json:"active"`
}

// EligibilityIdentifier is used to serialize beneficiaries' identifiers,
// as retrieved from Slade 360 EDI eligibility.
//
// The API version matters: a major revision of the API is likely to introduce
// a new benefit structure.
type EligibilityIdentifier struct {
	ID               int    `json:"id"`
	GUID             string `json:"guid"`
	Identifier       string `json:"identifier"`
	IdentifierType   string `json:"identifierType"`
	IsMainIdentifier bool   `json:"isMainIdentifier"`
}

// AllowedProvider is used to serialize panels for simple lists in the app.
type AllowedProvider struct {
	ProviderName      string `json:"providerName"`
	ProviderSladeCode int    `json:"providerSladeCode"`
	CopayNeeded       bool   `json:"copayNeeded"`
	CopayAmount       int    `json:"copayAmount"`
}

// Provider is used to serialize details from Slade 360 EDI
type Provider struct {
	ID        int    `json:"id"`
	GUID      string `json:"guid"`
	Name      string `json:"name"`
	SladeCode int    `json:"sladeCode"`
}

// Dependant is used to serialize dependants' information
type Dependant struct {
	Name         string       `json:"name"`
	Initials     string       `json:"initials"`
	MemberNumber string       `json:"memberNumber"`
	Relationship Relationship `json:"relationship"`
	PatientID    string       `json:"patientID,omitempty"`
}

// BenefitType is used to classify benefits for the Be.Well app.
type BenefitType string

// standard benefit types
const (
	BenefitTypeOutpatient BenefitType = "OUTPATIENT"
	BenefitTypeInpatient  BenefitType = "INPATIENT"
	BenefitTypeDental     BenefitType = "DENTAL"
	BenefitTypeOptical    BenefitType = "OPTICAL"
	BenefitTypeMaternity  BenefitType = "MATERNITY"
	BenefitTypeOther      BenefitType = "OTHER"
)

// AllBenefitType is a list of known (acceptable) benefit types
var AllBenefitType = []BenefitType{
	BenefitTypeOutpatient,
	BenefitTypeInpatient,
	BenefitTypeDental,
	BenefitTypeOptical,
	BenefitTypeMaternity,
	BenefitTypeOther,
}

// IsValid returns true for valid benefit types
func (e BenefitType) IsValid() bool {
	switch e {
	case BenefitTypeOutpatient,
		BenefitTypeInpatient,
		BenefitTypeDental,
		BenefitTypeOptical,
		BenefitTypeMaternity,
		BenefitTypeOther:
		return true
	}
	return false
}

// String renders a benefit type as a string
func (e BenefitType) String() string {
	return string(e)
}

// UnmarshalGQL reads a benefit type from the supplied input
func (e *BenefitType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BenefitType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BenefitType", str)
	}
	return nil
}

// MarshalGQL writes the benefit type to the supplied writer as a quoted string
func (e BenefitType) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// Relationship encodes the relationship between a dependant and
// their principal.
//
// The values currently present on Slade 360 EDI at the time of writing.
// are:
//
//     Child
//     DAUGHTER
//     PRINCIPAL
//     Brother
//     CHILD
//     SPOUSE
//     Partner
//     Spouse
//     Unmarried Child
//     FATHER
//     UNMARRIED CHILD
//     Father
//     MOTHER
//     Self
//     Sister
//     Mother
//     Married Child
//     Principal
//     SON
//
// These values are normalized as follows:
//
// CHILD
//
//     Child
//     DAUGHTER
//     CHILD
//     Unmarried Child
//     UNMARRIED CHILD
//     Married Child
//     SON
//
// PRINCIPAL
//
//     PRINCIPAL
//     Self
//     Principal
//
// SIBLING
//
//     Brother
//     Sister
//
// SPOUSE
//
//     SPOUSE
//     Partner
//     Spouse
//
// FATHER
//
//     Father
//     FATHER
//
// MOTHER
//
//     MOTHER
//     Mother
type Relationship string

// known EDI relationship types
const (
	RelationshipSpouse    Relationship = "SPOUSE"
	RelationshipChild     Relationship = "CHILD"
	RelationshipPrincipal Relationship = "PRINCIPAL"
	RelationshipSibling   Relationship = "SIBLING"
	RelationshipFather    Relationship = "FATHER"
	RelationshipMother    Relationship = "MOTHER"
)

// AllRelationship is the list of all known EDI relationship types
var AllRelationship = []Relationship{
	RelationshipSpouse,
	RelationshipChild,
	RelationshipPrincipal,
	RelationshipSibling,
	RelationshipFather,
	RelationshipMother,
}

// IsValid returns true for valid relationship types
func (e Relationship) IsValid() bool {
	switch e {
	case RelationshipSpouse,
		RelationshipChild,
		RelationshipPrincipal,
		RelationshipSibling,
		RelationshipFather,
		RelationshipMother:
		return true
	}
	return false
}

// String renders the EDI relationship as a string
func (e Relationship) String() string {
	return string(e)
}

// UnmarshalGQL translates the input into an EDI relationship
func (e *Relationship) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Relationship(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Relationship", str)
	}
	return nil
}

// MarshalGQL writes the relationship to the supplied writer as a quoted string
func (e Relationship) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

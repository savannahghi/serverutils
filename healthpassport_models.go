package base

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// SimplePatientRegistrationInput provides a simplified API to support registration
// of patients.
type SimplePatientRegistrationInput struct {
	ID                      string                    `json:"id"`
	Names                   []*NameInput              `json:"names"`
	IdentificationDocuments []*IdentificationDocument `json:"identificationDocuments"`
	BirthDate               Date                      `json:"birthDate"`
	PhoneNumbers            []*PhoneNumberInput       `json:"phoneNumbers"`
	Photos                  []*PhotoInput             `json:"photos"`
	Emails                  []*EmailInput             `json:"emails"`
	PhysicalAddresses       []*PhysicalAddress        `json:"physicalAddresses"`
	PostalAddresses         []*PostalAddress          `json:"postalAddresses"`
	Gender                  string                    `json:"gender"`
	Active                  bool                      `json:"active"`
	MaritalStatus           MaritalStatus             `json:"maritalStatus"`
	Languages               []Language                `json:"languages"`
	ReplicateUSSD           bool                      `json:"replicate_ussd,omitempty"`
}

// BreakGlassEpisodeCreationInput is used to start emergency episodes via a
// break glass protocol
type BreakGlassEpisodeCreationInput struct {
	PatientID       string `json:"patientID" firestore:"patientID"`
	ProviderCode    string `json:"providerCode" firestore:"providerCode"`
	PractitionerUID string `json:"practitionerUID" firestore:"practitionerUID"`
	Msisdn          string `json:"msisdn" firestore:"msisdn"`
	PatientPhone    string `json:"patientPhone" firestore:"patientPhone"`
	Otp             string `json:"otp" firestore:"otp"`
	FullAccess      bool   `json:"fullAccess" firestore:"fullAccess"`
}

// OTPEpisodeCreationInput is used to start patient visits via OTP
type OTPEpisodeCreationInput struct {
	PatientID    string `json:"patientID"`
	ProviderCode string `json:"providerCode"`
	Msisdn       string `json:"msisdn"`
	Otp          string `json:"otp"`
	FullAccess   bool   `json:"fullAccess"`
}

// OTPEpisodeUpgradeInput is used to upgrade existing open episodes
type OTPEpisodeUpgradeInput struct {
	EpisodeID string `json:"episodeID"`
	Msisdn    string `json:"msisdn"`
	Otp       string `json:"otp"`
}

// SimpleNHIFInput adds NHIF membership details as an extra identifier.
type SimpleNHIFInput struct {
	PatientID             string       `json:"patientID"`
	MembershipNumber      string       `json:"membershipNumber"`
	FrontImageBase64      *string      `json:"frontImageBase64"`
	FrontImageContentType *ContentType `json:"frontImageContentType"`
	RearImageBase64       *string      `json:"rearImageBase64"`
	RearImageContentType  *ContentType `json:"rearImageContentType"`
}

// SimpleNextOfKinInput is used to add next of kin to a patient.
type SimpleNextOfKinInput struct {
	PatientID         string              `json:"patientID"`
	Names             []*NameInput        `json:"names"`
	PhoneNumbers      []*PhoneNumberInput `json:"phoneNumbers"`
	Emails            []*EmailInput       `json:"emails"`
	PhysicalAddresses []*PhysicalAddress  `json:"physicalAddresses"`
	PostalAddresses   []*PostalAddress    `json:"postalAddresses"`
	Gender            string              `json:"gender"`
	Relationship      RelationshipType    `json:"relationship"`
	Active            bool                `json:"active"`
	BirthDate         Date                `json:"birthDate"`
}

// USSDEpisodeCreationInput is used to start episodes via USSD
type USSDEpisodeCreationInput struct {
	PatientID    string `json:"patientID"`
	ProviderCode string `json:"providerCode"`
	SessionID    string `json:"sessionID"`
	Msisdn       string `json:"msisdn"`
	FullAccess   bool   `json:"fullAccess"`
}

// EpisodeOfCare is An association between a patient and an organization /
// healthcare provider(s) during which time encounters may occur. The managing
// organization assumes a level of responsibility for the patient during this
// time.
type EpisodeOfCare struct {
	ID                   string              `json:"id,omitempty"`
	Identifier           []*Identifier       `json:"identifier,omitempty"`
	Status               EpisodeOfCareStatus `json:"status,omitempty"`
	Type                 []*CodeableConcept  `json:"type,omitempty"`
	Patient              Reference           `json:"patient,omitempty"`
	ManagingOrganization *Reference          `json:"managingOrganization,omitempty"`
	Period               *Period             `json:"period,omitempty"`
}

// EpisodeOfCarePayload is used to return the results after creation of
// episodes of care
type EpisodeOfCarePayload struct {
	EpisodeOfCare *EpisodeOfCare `json:"episodeOfCare"`
	TotalVisits   int            `json:"totalVisits"`
}

// Reference defines references to other FHIR resources.
type Reference struct {
	Reference  string     `json:"reference,omitempty"`
	Type       string     `json:"type,omitempty"`
	Identifier Identifier `json:"identifier,omitempty"`
	Display    *string    `json:"display,omitempty"`
}

// Patient description from FHIR: Demographics and other administrative information
// about an individual or animal receiving care or other health-related services.
type Patient struct {
	ID            string            `json:"id"`
	Identifier    []*Identifier     `json:"identifier"`
	Active        bool              `json:"active"`
	Name          []*HumanName      `json:"name"`
	Telecom       []*ContactPoint   `json:"telecom"`
	Gender        string            `json:"gender"`
	BirthDate     Date              `json:"birthDate"`
	Address       []*Address        `json:"address"`
	MaritalStatus *CodeableConcept  `json:"maritalStatus"`
	Photo         []*Attachment     `json:"photo"`
	Contact       []*PatientContact `json:"contact"`
	Communication []*Communication  `json:"communication"`
}

// Names renders the patient's names as a string
func (p Patient) Names() string {
	name := ""
	if p.Name == nil {
		return name
	}

	names := []string{}
	for _, hn := range p.Name {
		if hn == nil {
			continue
		}
		if hn.Text == nil {
			continue
		}
		names = append(names, *hn.Text)
	}
	name = strings.Join(names, " | ")
	return name
}

// IsEntity ...
func (p Patient) IsEntity() {}

// PatientConnection is a Relay style connection for use in listings of FHIR
// patient records.
type PatientConnection struct {
	Edges    []*PatientEdge `json:"edges"`
	PageInfo *PageInfo      `json:"pageInfo"`
}

// PatientContact is a contact party (e.g. guardian, partner, friend) for the
// patient.
type PatientContact struct {
	ID           string             `json:"id"`
	Relationship []*CodeableConcept `json:"relationship"`
	Name         *HumanName         `json:"name"`
	Telecom      []*ContactPoint    `json:"telecom"`
	Address      *Address           `json:"address"`
	Gender       *string            `json:"gender"`
	Period       *Period            `json:"period"`
}

// IsEntity ...
func (p PatientContact) IsEntity() {}

// PatientContactInput is used to create and update patient contacts
type PatientContactInput struct {
	ID           *string                 `json:"id"`
	Relationship []*CodeableConceptInput `json:"relationship"`
	Name         *HumanNameInput         `json:"name"`
	Telecom      []*ContactPointInput    `json:"telecom"`
	Address      *AddressInput           `json:"address"`
	Gender       *string                 `json:"gender"`
	Period       *PeriodInput            `json:"period"`
}

// PatientEdge is a Relay style edge for listings of FHIR patient records.
type PatientEdge struct {
	Cursor          string   `json:"cursor"`
	Node            *Patient `json:"node"`
	HasOpenEpisodes bool     `json:"hasOpenEpisodes"`
}

// PatientInput is used to create patient records.
type PatientInput struct {
	ID            *string                `json:"id"`
	Identifier    []*IdentifierInput     `json:"identifier"`
	Active        bool                   `json:"active"`
	Name          []*HumanNameInput      `json:"name"`
	Telecom       []*ContactPointInput   `json:"telecom"`
	Gender        string                 `json:"gender"`
	BirthDate     Date                   `json:"birthDate"`
	Address       []*AddressInput        `json:"address"`
	MaritalStatus *CodeableConceptInput  `json:"maritalStatus"`
	Photo         []*AttachmentInput     `json:"photo"`
	Contact       []*PatientContactInput `json:"contact"`
	Communication []*CommunicationInput  `json:"communication"`
}

// PatientPayload is used to return patient records and ancillary data after
// mutations.
type PatientPayload struct {
	PatientRecord   *Patient         `json:"patientRecord,omitempty"`
	HasOpenEpisodes bool             `json:"hasOpenEpisodes,omitempty"`
	OpenEpisodes    []*EpisodeOfCare `json:"openEpisodes,omitempty"`
}

// RetirePatientInput is used to retire patient records.
type RetirePatientInput struct {
	ID string `json:"id"`
}

// PatientExtraInformationInput is used to update patient records metadata.
type PatientExtraInformationInput struct {
	PatientID     string         `json:"patientID"`
	MaritalStatus *MaritalStatus `json:"maritalStatus"`
	Languages     []*Language    `json:"languages"`
	Emails        []*EmailInput  `json:"emails"`
}

// USSDNextOfKinCreationInput is used to register next of kin via USSD.
type USSDNextOfKinCreationInput struct {
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	OtherNames string    `json:"otherNames"`
	BirthDate  time.Time `json:"birthDate"`
	Gender     string    `json:"gender"`
	Active     bool      `json:"active"`
	ParentID   string    `json:"parentID"`
}

// EpisodeOfCareStatus is used to record the status of an episode of care.
type EpisodeOfCareStatus string

// Episode of care valueset status values
const (
	// EpisodeOfCareStatusPlanned ...
	EpisodeOfCareStatusPlanned EpisodeOfCareStatus = "planned"
	// EpisodeOfCareStatusWaitlist ...
	EpisodeOfCareStatusWaitlist EpisodeOfCareStatus = "waitlist"
	// EpisodeOfCareStatusActive ...
	EpisodeOfCareStatusActive EpisodeOfCareStatus = "active"
	// EpisodeOfCareStatusOnhold ...
	EpisodeOfCareStatusOnhold EpisodeOfCareStatus = "onhold"
	// EpisodeOfCareStatusFinished ...
	EpisodeOfCareStatusFinished EpisodeOfCareStatus = "finished"
	// EpisodeOfCareStatusCancelled ...
	EpisodeOfCareStatusCancelled EpisodeOfCareStatus = "cancelled"
	// EpisodeOfCareStatusEnteredInError ...
	EpisodeOfCareStatusEnteredInError EpisodeOfCareStatus = "entered_in_error"
	// EpisodeOfCareStatusEnteredInErrorCanonical ...
	EpisodeOfCareStatusEnteredInErrorCanonical EpisodeOfCareStatus = "entered-in-error"
)

// AllEpisodeOfCareStatus is a list of episode of care statuses
var AllEpisodeOfCareStatus = []EpisodeOfCareStatus{
	EpisodeOfCareStatusPlanned,
	EpisodeOfCareStatusWaitlist,
	EpisodeOfCareStatusActive,
	EpisodeOfCareStatusOnhold,
	EpisodeOfCareStatusFinished,
	EpisodeOfCareStatusCancelled,
	EpisodeOfCareStatusEnteredInError,
	EpisodeOfCareStatusEnteredInErrorCanonical,
}

// IsValid validates episode of care status values
func (e EpisodeOfCareStatus) IsValid() bool {
	switch e {
	case EpisodeOfCareStatusPlanned, EpisodeOfCareStatusWaitlist, EpisodeOfCareStatusActive, EpisodeOfCareStatusOnhold, EpisodeOfCareStatusFinished, EpisodeOfCareStatusCancelled, EpisodeOfCareStatusEnteredInError, EpisodeOfCareStatusEnteredInErrorCanonical:
		return true
	}
	return false
}

// String ...
func (e EpisodeOfCareStatus) String() string {
	return string(e)
}

// UnmarshalGQL converts the input value into an episode of care
func (e *EpisodeOfCareStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EpisodeOfCareStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EpisodeOfCareStatus", str)
	}
	return nil
}

// MarshalGQL writes the episode of care status value to the supplied writer
func (e EpisodeOfCareStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))

}

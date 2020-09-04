package base

import (
	"net/url"
	"time"
)

// QueryParam is an interface used for filter and sort parameters
type QueryParam interface {
	ToURLValues() (values url.Values)
}

// PaginationInput represents paging parameters
type PaginationInput struct {
	First  int    `json:"first"`
	Last   int    `json:"last"`
	After  string `json:"after"`
	Before string `json:"before"`
}

// SortInput is a generic container for strongly typed sorting parameters
type SortInput struct {
	SortBy []*SortParam `json:"sortBy"`
}

// SortParam represents a single field sort parameter
type SortParam struct {
	FieldName string    `json:"fieldName"`
	SortOrder SortOrder `json:"sortOrder"`
}

// FilterInput is s generic container for strongly type filter parameters
type FilterInput struct {
	Search   *string        `json:"search"`
	FilterBy []*FilterParam `json:"filterBy"`
}

// FilterParam represents a single field filter parameter
type FilterParam struct {
	FieldName           string      `json:"fieldName"`
	FieldType           FieldType   `json:"fieldType"`
	ComparisonOperation Operation   `json:"comparisonOperation"`
	FieldValue          interface{} `json:"fieldValue"`
}

// PhoneOptIn is used to persist and manage phone communication whitelists
type PhoneOptIn struct {
	MSISDN  string `json:"msisdn" firestore:"msisdn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

// USSDSessionLog is used to persist a log of USSD sessions
type USSDSessionLog struct {
	MSISDN    string `json:"msisdn" firestore:"msisdn"`
	SessionID string `json:"sessionID" firestore:"sessionID"`
}

// EmailOptIn is used to persist and manage email communication whitelists
type EmailOptIn struct {
	Email   string `json:"email" firestore:"optedIn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

// SladeAPIListRespBase defines the fields that are common on list endpoints
// for Slade 360 APIs
type SladeAPIListRespBase struct {
	Count       int    `json:"count,omitempty"`
	Next        string `json:"next,omitempty"`
	Previous    string `json:"previous,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
	CurrentPage int    `json:"current_page,omitempty"`
	TotalPages  int    `json:"total_pages,omitempty"`
	StartIndex  int    `json:"start_index,omitempty"`
	EndIndex    int    `json:"end_index,omitempty"`
}

// LoginCreds is used to (de)serialize the login username and password
type LoginCreds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type setupProcess struct {
	Progress        int           `json:"progress"`
	CompletedSteps  []interface{} `json:"completedSteps"`
	IncompleteSteps []interface{} `json:"incompleteSteps"`
	NextStep        interface{}   `json:"nextStep"`
}

// LoginResponse is used to (de)serialize the result of a successful login
type LoginResponse struct {
	CustomToken   string          `json:"custom_token"`
	Scope         string          `json:"scope"`
	ExpiresIn     int             `json:"expires_in"`
	IDToken       string          `json:"id_token"`
	RefreshToken  string          `json:"refresh_token"`
	TokenType     string          `json:"token_type"`
	UserProfile   *EDIUserProfile `json:"user_profile"`
	UID           string          `json:"uid"`
	Email         string          `json:"email"`
	DisplayName   string          `json:"display_name"`
	EmailVerified bool            `json:"email_verified"`
	PhoneNumber   string          `json:"phone_number"`
	PhotoURL      string          `json:"photo_url"`
	Disabled      bool            `json:"disabled"`
	TenantID      string          `json:"tenant_id"`
	ProviderID    string          `json:"provider_id"`
	Setup         *setupProcess   `json:"setup,omitempty"`
}

// RefreshCreds models the inputs expected from an API client when refreshing tokens
type RefreshCreds struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest models the inputs expected from an API client when requesting a log out
type LogoutRequest struct {
	UID string `json:"uid"`
}

// EDIUserProfile is used to (de)serialialize the Slade 360 auth server
// profile of the logged in user.
type EDIUserProfile struct {
	ID              int      `json:"id"`
	GUID            string   `json:"guid"`
	Email           string   `json:"email"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	OtherNames      string   `json:"other_names"`
	IsStaff         bool     `json:"is_staff"`
	IsActive        bool     `json:"is_active"`
	Organisation    int      `json:"organisation"`
	BusinessPartner string   `json:"business_partner"`
	Roles           []string `json:"roles"`
	BPType          string   `json:"bp_type"`
}

// RefreshResponse is used to return the results of a successful token refresh to an API client
type RefreshResponse struct {
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// FirebaseRefreshResponse is used to (de)serialize the results of a successful Firebase token refresh
type FirebaseRefreshResponse struct {
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	UserID       string `json:"user_id"`
	ProjectID    string `json:"project_id"`
}

// AccessTokenPayload is used to accept access token verification requests from API clients
type AccessTokenPayload struct {
	AccessToken string `json:"accessToken"`
}

// EmailInput is used to register patient emails.
type EmailInput struct {
	Email              string `json:"email"`
	CommunicationOptIn bool   `json:"communicationOptIn"`
}

// IdentificationDocument is used to input e.g National ID or passport document
// numbers at patient registration.
type IdentificationDocument struct {
	DocumentType     IDDocumentType `json:"documentType"`
	DocumentNumber   string         `json:"documentNumber"`
	Title            *string        `json:"title,omitempty"`
	ImageContentType *ContentType   `json:"imageContentType,omitempty"`
	ImageBase64      *string        `json:"imageBase64,omitempty"`
}

// NameInput is used to input patient names.
type NameInput struct {
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	OtherNames *string `json:"otherNames"`
}

// PhoneNumberInput is used to input phone numbers.
type PhoneNumberInput struct {
	Msisdn             string `json:"msisdn"`
	VerificationCode   string `json:"verificationCode"`
	IsUssd             bool   `json:"isUSSD"`
	CommunicationOptIn bool   `json:"communicationOptIn"`
}

// PhotoInput is used to upload patient photos.
type PhotoInput struct {
	PhotoContentType ContentType `json:"photoContentType"`
	PhotoBase64data  string      `json:"photoBase64data"`
	PhotoFilename    string      `json:"photoFilename"`
}

// PhysicalAddress is used to record a precise physical address.
type PhysicalAddress struct {
	MapsCode        string `json:"mapsCode"`
	PhysicalAddress string `json:"physicalAddress"`
}

// PostalAddress is used to record patient's postal addresses
type PostalAddress struct {
	PostalAddress string `json:"postalAddress"`
	PostalCode    string `json:"postalCode"`
}

// Address description from FHIR: An address expressed using postal conventions
// (as opposed to GPS or other location definition formats).
//
// This data type may be used to convey addresses for use in delivering mail as
// well as for visiting locations which might not be valid for mail delivery.
//
// There are a variety of postal address formats defined around the world.
type Address struct {
	ID         string      `json:"id"`
	Use        AddressUse  `json:"use"`
	Type       AddressType `json:"type"`
	Text       string      `json:"text"`
	Line       []*string   `json:"line"`
	City       *string     `json:"city"`
	District   *string     `json:"district"`
	State      *string     `json:"state"`
	PostalCode *string     `json:"postalCode"`
	Country    Country     `json:"country"`
	Period     *Period     `json:"period"`
}

// AddressInput is used to create postal and physical addresses.
//
// IMPORTANT:
//
// For physical addresses, use Google Maps co-ordinates or plus codes.
// See: https://support.google.com/maps/answer/18539?co=GENIE.Platform%3DDesktop&hl=en
type AddressInput struct {
	ID         *string      `json:"id"`
	Use        AddressUse   `json:"use"`
	Type       AddressType  `json:"type"`
	Text       string       `json:"text"`
	Line       []*string    `json:"line"`
	City       *string      `json:"city"`
	District   *string      `json:"district"`
	State      *string      `json:"state"`
	PostalCode *string      `json:"postalCode"`
	Country    Country      `json:"country"`
	Period     *PeriodInput `json:"period"`
}

// Attachment is for containing or referencing attachments - additional data
// content defined in other formats. The most common use of this type is to
// include images or reports in some report format such as PDF. However, it can
// be used for any data that has a MIME type.
type Attachment struct {
	ID          string    `json:"id"`
	ContentType string    `json:"contentType"`
	Language    string    `json:"language"`
	Data        string    `json:"data"`
	URL         string    `json:"url"`
	Size        int64     `json:"size"`
	Hash        string    `json:"hash"`
	Title       string    `json:"title"`
	Creation    time.Time `json:"creation"`
}

// AttachmentInput is used to create attachments.
type AttachmentInput struct {
	ID          string     `json:"id"`
	ContentType string     `json:"contentType"`
	Language    string     `json:"language"`
	Data        string     `json:"data"`
	URL         string     `json:"url"`
	Size        int64      `json:"size"`
	Hash        string     `json:"hash"`
	Title       string     `json:"title"`
	Creation    *time.Time `json:"creation"`
}

// CodeableConcept description from FHIR: A concept that may be defined by a
// formal reference to a terminology or ontology or may be provided by text.
type CodeableConcept struct {
	ID     string    `json:"id,omitempty"`
	Coding []*Coding `json:"coding,omitempty"`
	Text   string    `json:"text,omitempty"`
}

// CodeableConceptInput is used to create codeable concepts.
type CodeableConceptInput struct {
	ID     *string        `json:"id"`
	Coding []*CodingInput `json:"coding"`
	Text   string         `json:"text"`
}

// Coding description from FHIR: A reference to a code defined by a terminology system.
type Coding struct {
	System       string  `json:"system"`
	Version      *string `json:"version"`
	Code         string  `json:"code"`
	Display      *string `json:"display"`
	UserSelected *bool   `json:"userSelected"`
}

// CodingInput is used to set coding.
type CodingInput struct {
	System       string `json:"system"`
	Version      string `json:"version"`
	Code         string `json:"code"`
	Display      string `json:"display"`
	UserSelected bool   `json:"userSelected"`
}

// Communication description from FHIR: Information about a person
// that is involved in the care for a patient, but who is not the target of
// healthcare, nor has a formal responsibility in the care process.
type Communication struct {
	Language  *CodeableConcept `json:"language"`
	Preferred bool             `json:"preferred"`
}

// CommunicationInput is used to create send communication data in GraphQL.
type CommunicationInput struct {
	Language  *CodeableConceptInput `json:"language"`
	Preferred bool                  `json:"preferred"`
}

// ContactPoint description from FHIR: Details for all kinds of technology
// mediated contact points for a person or organization, including telephone,
// email, etc.
type ContactPoint struct {
	ID     string             `json:"id"`
	System ContactPointSystem `json:"system"`
	Value  string             `json:"value"`
	Use    ContactPointUse    `json:"use"`
	Rank   int64              `json:"rank"`
	Period *Period            `json:"period"`
}

// ContactPointInput is used to input contact details e.g phone, email etc.
type ContactPointInput struct {
	ID     *string            `json:"id"`
	System ContactPointSystem `json:"system"`
	Value  string             `json:"value"`
	Use    ContactPointUse    `json:"use"`
	Rank   int64              `json:"rank"`
	Period *PeriodInput       `json:"period"`
}

// HumanName description from FHIR: A human's name with the ability to identify
// parts and usage.
type HumanName struct {
	ID     string    `json:"id"`
	Use    NameUse   `json:"use"`
	Text   *string   `json:"text"`
	Family string    `json:"family"`
	Given  []string  `json:"given"`
	Prefix []*string `json:"prefix"`
	Suffix []*string `json:"suffix"`
	Period *Period   `json:"period"`
}

// HumanNameInput is used to input patient names.
type HumanNameInput struct {
	ID     *string      `json:"id"`
	Use    NameUse      `json:"use"`
	Text   *string      `json:"text"`
	Family string       `json:"family"`
	Given  []string     `json:"given"`
	Prefix []*string    `json:"prefix"`
	Suffix []*string    `json:"suffix"`
	Period *PeriodInput `json:"period"`
}

// Identifier is used to represent a numeric or alphanumeric string that is
// associated with a single object or entity within a given system. Typically,
// identifiers are used to connect content in resources to external content
// available in other frameworks or protocols. Identifiers are associated with
// objects and may be changed or retired due to human or system process and
// errors.
type Identifier struct {
	ID     string           `json:"id,omitempty"`
	Use    IdentifierUse    `json:"use,omitempty"`
	Type   *CodeableConcept `json:"type,omitempty"`
	System *string          `json:"system,omitempty"`
	Value  *string          `json:"value,omitempty"`
	Period *Period          `json:"period,omitempty"`
}

// IdentifierInput is used to create and update identifiers.
type IdentifierInput struct {
	ID     string                `json:"id"`
	Use    IdentifierUse         `json:"use"`
	Type   *CodeableConceptInput `json:"type"`
	System *string               `json:"system"`
	Value  *string               `json:"value"`
	Period *PeriodInput          `json:"period"`
}

// Period is a FHIR https://www.hl7.org/fhir/datatypes.html#Period.
//
// A period should have a start, end or both. It's an error to have a period that
// has null start and end times.
type Period struct {
	ID    string     `json:"id,omitempty"`
	Start *time.Time `json:"start,omitempty"`
	End   *time.Time `json:"end,omitempty"`
}

// PeriodInput is used to set time ranges e.g validity.
//
// A period should have a start, end or both. It's an error to have a period that
// has null start and end times.
type PeriodInput struct {
	ID    *string    `json:"id"`
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

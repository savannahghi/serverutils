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

//IsEntity ...
func (p PaginationInput) IsEntity() {}

// SortInput is a generic container for strongly typed sorting parameters
type SortInput struct {
	SortBy []*SortParam `json:"sortBy"`
}

//IsEntity ...
func (s SortInput) IsEntity() {}

// SortParam represents a single field sort parameter
type SortParam struct {
	FieldName string    `json:"fieldName"`
	SortOrder SortOrder `json:"sortOrder"`
}

//IsEntity ...
func (s SortParam) IsEntity() {}

// FilterInput is s generic container for strongly type filter parameters
type FilterInput struct {
	Search   *string        `json:"search"`
	FilterBy []*FilterParam `json:"filterBy"`
}

//IsEntity ...
func (f FilterInput) IsEntity() {}

// FilterParam represents a single field filter parameter
type FilterParam struct {
	FieldName           string      `json:"fieldName"`
	FieldType           FieldType   `json:"fieldType"`
	ComparisonOperation Operation   `json:"comparisonOperation"`
	FieldValue          interface{} `json:"fieldValue"`
}

//IsEntity ...
func (f FilterParam) IsEntity() {}

// PhoneOptIn is used to persist and manage phone communication whitelists
type PhoneOptIn struct {
	MSISDN  string `json:"msisdn" firestore:"msisdn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

//IsEntity ...
func (p PhoneOptIn) IsEntity() {}

// USSDSessionLog is used to persist a log of USSD sessions
type USSDSessionLog struct {
	MSISDN    string `json:"msisdn" firestore:"msisdn"`
	SessionID string `json:"sessionID" firestore:"sessionID"`
}

//IsEntity ...
func (p USSDSessionLog) IsEntity() {}

// EmailOptIn is used to persist and manage email communication whitelists
type EmailOptIn struct {
	Email   string `json:"email" firestore:"optedIn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
}

//IsEntity ...
func (e EmailOptIn) IsEntity() {}

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

//IsEntity ...
func (l LoginCreds) IsEntity() {}

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

//IsEntity ...
func (l LoginResponse) IsEntity() {}

// RefreshCreds models the inputs expected from an API client when refreshing tokens
type RefreshCreds struct {
	RefreshToken string `json:"refresh_token"`
}

//IsEntity ...
func (r RefreshCreds) IsEntity() {}

// LogoutRequest models the inputs expected from an API client when requesting a log out
type LogoutRequest struct {
	UID string `json:"uid"`
}

//IsEntity ...
func (l LogoutRequest) IsEntity() {}

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

//IsEntity ...
func (e EDIUserProfile) IsEntity() {}

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

//IsEntity ...
func (a Attachment) IsEntity() {}

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

// PIN is used to store a PIN (Personal Identifiation Number) associated
// to a phone number sign up to Firebase
type PIN struct {
	UID     string `json:"uid" firestore:"uid"`
	MSISDN  string `json:"msisdn,omitempty" firestore:"msisdn"`
	PIN     string `json:"pin,omitempty" firestore:"pin"`
	IsValid bool   `json:"isValid,omitempty" firestore:"isValid"`
}

//IsEntity ...
func (p PIN) IsEntity() {}

// FinancialYearAndCurrency generic struct for the default financial year and default currency
type FinancialYearAndCurrency struct {
	ID           *string `json:"id"`
	Active       *bool   `json:"active"`
	Organisation *string `json:"organisation"`
	IsDefault    bool    `json:"is_default,omitempty"`
	ISOCode      *string `json:"iso_code,omitempty"`
}

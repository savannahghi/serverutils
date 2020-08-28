package base

import "net/url"

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

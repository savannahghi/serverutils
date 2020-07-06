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

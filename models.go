package base

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

// SortInput is a generic container for strongly typed sorting parameters
type SortInput struct {
	SortBy []*SortParam `json:"sortBy"`
}

// SortParam represents a single field sort parameter
type SortParam struct {
	FieldName string    `json:"fieldName"`
	SortOrder SortOrder `json:"sortOrder"`
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

// EmailOptIn is used to persist and manage email communication whitelists
type EmailOptIn struct {
	Email   string `json:"email" firestore:"optedIn"`
	OptedIn bool   `json:"optedIn" firestore:"optedIn"`
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

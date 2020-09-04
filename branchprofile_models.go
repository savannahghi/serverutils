package base

import "time"

// BranchProfileRatingInput is used to submit ratings
type BranchProfileRatingInput struct {
	BranchProfileID string `json:"branchProfileID" firestore:"branchProfileID"`
	Review          string `json:"review" firestore:"review"`
	Rating          int    `json:"rating" firestore:"rating"`
}

// BranchProfileRating records the ratings of a provider or doctor by users
type BranchProfileRating struct {
	ID              string    `json:"id" firestore:"id"`
	BranchProfileID string    `json:"branchProfileID" firestore:"branchProfileID"`
	UID             string    `json:"uid" firestore:"uid"`
	Reviewer        string    `json:"reviewer" firestore:"reviewer"`
	Review          string    `json:"review" firestore:"review"`
	Updated         time.Time `json:"updated" firestore:"updated"`
	Rating          int       `json:"rating" firestore:"rating"`
	Deleted         bool      `json:"deleted" firestore:"deleted"`
}

// IsNode is a "label" that marks this struct (and those that embed it) as
// implementations of the "Base" interface defined in our GraphQL schema.
func (bpr *BranchProfileRating) IsNode() {}

// GetID returns the struct's's ID value
func (bpr *BranchProfileRating) GetID() ID {
	return IDValue(bpr.ID)
}

// IsEntity ...
func (bpr BranchProfileRating) IsEntity() {}

// SetID set's the branch profile's ID
func (bpr *BranchProfileRating) SetID(id string) {
	bpr.ID = id
}

// BranchProfile holds additional information about Slade 360 Charge Master
// branches.
type BranchProfile struct {
	ID                           string                  `json:"id" firestore:"id"`
	BranchSladeCode              string                  `json:"branchSladeCode" firestore:"branchSladeCode"`
	BranchName                   string                  `json:"branchName" firestore:"branchName"`
	OrganizationSladeCode        string                  `json:"organizationSladeCode" firestore:"organizationSladeCode"`
	Coordinates                  string                  `json:"coordinates" firestore:"coordinates"`
	Active                       bool                    `json:"active" firestore:"active"`
	Specialties                  []PractitionerSpecialty `json:"specialties,omitempty" firestore:"specialties"`
	SortPriority                 int                     `json:"sortPriority" firestore:"sortPriority"`
	Profile                      Markdown                `json:"profile" firestore:"profile"`
	License                      string                  `json:"license" firestore:"license"`
	PhotoContentType             ContentType             `json:"photoContentType" firestore:"photoContentType"`
	PhotoBase64data              string                  `json:"photoBase64data" firestore:"photoBase64data"`
	AverageConsultationPrice     float64                 `json:"averageConsultationPrice" firestore:"averageConsultationPrice"`
	AverageTeleconsultationPrice float64                 `json:"averageTeleconsultationPrice" firestore:"averageTeleconsultationPrice"`
	Phones                       []string                `json:"phones" firestore:"phones"`
	Emails                       []string                `json:"emails" firestore:"emails"`
	ProfilePages                 []URL                   `json:"profilePages" firestore:"profilePages"`
	CalendarID                   string                  `json:"calendarID" firestore:"calendarID"`
	HasTelehealth                bool                    `json:"hasTelehealth" firestore:"hasTelehealth"`
	Geohash                      string                  `json:"geohash" firestore:"geohash"`
	FHIRReference                *string                 `json:"fhirReference" firestore:"fhirReference"`
	Deleted                      bool                    `json:"-" firestore:"deleted"` // for soft deletes

	AverageRating   float64 `json:"averageRating" firestore:"averageRating"`
	NumberOfRatings int     `json:"numberOfRatings" firestore:"numberOfRatings"`
}

// IsNode is a "label" that marks this struct (and those that embed it) as
// implementations of the "Base" interface defined in our GraphQL schema.
func (bp *BranchProfile) IsNode() {}

// GetID returns the struct's's ID value
func (bp *BranchProfile) GetID() ID {
	return IDValue(bp.ID)
}

// SetID set's the branch profile's ID
func (bp *BranchProfile) SetID(id string) {
	bp.ID = id
}

// IsEntity ...
func (bp BranchProfile) IsEntity() {}

// BranchProfileEdge is used to serialize GraphQL Relay edges for organization
type BranchProfileEdge struct {
	Cursor        *string                `json:"cursor"`
	Node          *BranchProfile         `json:"node"`
	RecentRatings []*BranchProfileRating `json:"recentRatings,omitempty"`
}

// BranchProfileConnection is used to serialize GraphQL Relay connections for organizations
type BranchProfileConnection struct {
	Edges    []*BranchProfileEdge `json:"edges"`
	PageInfo *PageInfo            `json:"pageInfo"`
}

// BranchProfileNode is a Relay type node that combines a branch profile with other relevant info
type BranchProfileNode struct {
	BranchProfile BranchProfile          `json:"branchProfile,omitempty"`
	RecentRatings []*BranchProfileRating `json:"recentRatings,omitempty"`
}

// BranchProfileInput is used to create or update branch profiles.
type BranchProfileInput struct {
	ID                           *string                 `json:"id"`
	BranchSladeCode              string                  `json:"branchSladeCode"`
	BranchName                   string                  `json:"branchName"`
	OrganizationSladeCode        string                  `json:"organizationSladeCode"`
	Coordinates                  string                  `json:"coordinates"`
	HasTelehealth                bool                    `json:"hasTelehealth"`
	Active                       bool                    `json:"active"`
	Specialties                  []PractitionerSpecialty `json:"specialties,omitempty"`
	SortPriority                 int                     `json:"sortPriority"`
	Profile                      Markdown                `json:"profile"`
	License                      string                  `json:"license"`
	PhotoContentType             ContentType             `json:"photoContentType"`
	PhotoBase64data              string                  `json:"photoBase64data"`
	AverageConsultationPrice     float64                 `json:"averageConsultationPrice"`
	AverageTeleconsultationPrice float64                 `json:"averageTeleconsultationPrice"`
	Phones                       []string                `json:"phones"`
	Emails                       []string                `json:"emails"`
	ProfilePages                 []URL                   `json:"profilePages"`
}

// BranchProfilePayload is used to return the results of creating or updating a
// branch.
type BranchProfilePayload struct {
	BranchProfile *BranchProfile `json:"branchProfile"`
}

// FreeBusy is used to serialize a branch calendar's free/busy.
//
// This information is processed a bit from it's original form in the Google
// Calendar base.
//
// The output from the API looks like this:
//
// {
//   "kind": "calendar#freeBusy",
//   "timeMin": "2020-06-10T11:50:37.000Z",
//   "timeMax": "2020-06-12T11:50:37.000Z",
//   "calendars": {
//     "healthcloud.co.ke_edt6esvi9rckjans7na7euikfo@group.calendar.google.com": {
//     "busy": []
//     }
//   }
// }
type FreeBusy struct {
	TimeMin   time.Time   `json:"timeMin"`
	TimeMax   time.Time   `json:"timeMax"`
	BusySlots []*BusySlot `json:"busySlots"`
}

// PanelProviderPayload is used to return panel providers and associated info e.g time slots.
type PanelProviderPayload struct {
	BranchProfile *BranchProfile `json:"branchProfile"`
	FreeBusy      *FreeBusy      `json:"freeBusy"`
}

// BusySlot is used to serialize a single calendar's busy time slot.
type BusySlot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

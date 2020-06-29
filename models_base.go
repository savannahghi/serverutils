package base

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Sep is a separator, used to create "opaque" IDs
const Sep = "|"

// ID is fulfilled by all stringifiable types.
// A valid Relay ID must fulfill this interface.
type ID interface {
	fmt.Stringer
}

// Node is a Relay (GraphQL Relay) node.
// Any valid type in this server should be a node.
type Node interface {
	IsNode()
	GetID() ID
	SetID(string)
}

// IDValue represents GraphQL object identifiers
type IDValue string

func (val IDValue) String() string { return string(val) }

// Typeof returns the type name for the supplied value
func Typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// MarshalID get's a re-fetchable GraphQL Relay ID that combines an objects's ID with it's type
// and encodes it into an "opaque" Base64 string.
func MarshalID(id string, n Node) ID {
	nodeType := Typeof(n)
	combinedID := fmt.Sprintf("%s%s%s", id, Sep, nodeType)
	return IDValue(base64.StdEncoding.EncodeToString([]byte(combinedID)))
}

// PageInfo is used to add pagination information to Relay edges.
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}

// NewString returns a pointer to the supplied string.
func NewString(s string) *string {
	return &s
}

// AuditLog records changes made to models
type AuditLog struct {
	ID        uuid.UUID
	RecordID  uuid.UUID        // ID of the audited record
	TypeName  string           // type of the audited record
	Operation string           // e.g pre_save, post_save
	When      time.Time        // timestamp of the operation
	UID       string           // UID of the involved user
	JSON      *json.RawMessage // serialized JSON snapshot
}

// Model defines common behavior for our models.
// It is also an ideal place to place hooks that are common to all models
// e.g audit, streaming analytics etc.
// CAUTION: Model should be evolved with cautions, because of migrations.
type Model struct {
	ID string `json:"id" firestore:"id"`

	// All models have a non nullable name field
	// If a derived model does not need this, it should use a placeholder e.g "-"
	Name string `json:"name" firestore:"name,omitempty"`

	// All records have an optional description
	Description string `json:"description" firestore:"description,omitempty"`

	// bug alert! If you add "omitempty" to the firestore struct tag, `false`
	// values will not be saved
	Deleted bool `json:"deleted,omitempty" firestore:"deleted"`

	// This is used for audit tracking but is not saved or serialized
	CreatedByUID string `json:"createdByUID" firestore:"createdByUID,omitempty"`
	UpdatedByUID string `json:"updatedByUID" firestore:"updatedByUID,omitempty"`
}

// IsNode is a "label" that marks this struct (and those that embed it) as
// implementations of the "Base" interface defined in our GraphQL schema.
func (c *Model) IsNode() {}

// GetID returns the struct's ID value
func (c *Model) GetID() ID {
	return IDValue(c.ID)
}

// SetID sets the struct's ID value
func (c *Model) SetID(id string) {
	c.ID = id
}

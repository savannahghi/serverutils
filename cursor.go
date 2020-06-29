package base

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/vmihailenco/msgpack"
)

// Cursor represents an opaque "position" for a record, for use in pagination
type Cursor struct {
	Offset int `json:"offset"`
}

// NewCursor creates a cursor from an offset and ID
func NewCursor(offset int) *Cursor {
	return &Cursor{Offset: offset}
}

// EncodeCursor converts a cursor to a string
func EncodeCursor(cursor *Cursor) string {
	b, err := msgpack.Marshal(cursor)
	if err != nil {
		msg := fmt.Sprintf("unable to encode cursor: %s", err)
		log.Println(msg)
		return msg
	}
	return base64.StdEncoding.EncodeToString(b)
}

// CreateAndEncodeCursor creates a cursor and immediately encodes it.
// It panics if it cannot encode the cursor.
// These cursors use ZERO BASED indexing.
func CreateAndEncodeCursor(offset int) *string {
	c := NewCursor(offset)
	enc := EncodeCursor(c)
	return &enc
}

// DecodeCursor decodes a cursor string to a pointer to a cursor struct
func DecodeCursor(cursor string) (*Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("unable to decode cursor to base64: %w", err)
	}

	var out Cursor
	err = msgpack.Unmarshal(b, &out)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal cursor via msgpack: %w", err)
	}
	return &out, nil
}

// MustDecodeCursor decodes the cursor or panics
func MustDecodeCursor(cursor string) *Cursor {
	decoded, err := DecodeCursor(cursor)
	if err != nil {
		msg := fmt.Sprintf("unable to encode cursor: %s", err)
		log.Panicf(msg)
	}
	return decoded
}

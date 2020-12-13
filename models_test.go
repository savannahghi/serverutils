package base

import (
	"testing"

	"github.com/segmentio/ksuid"
)

func TestUpload_IsEntity(t *testing.T) {
	u := Upload{}
	u.IsEntity()
}

func TestUpload_IsNode(t *testing.T) {
	u := Upload{}
	u.IsNode()
}

func TestUpload_GetID(t *testing.T) {
	randomID := ksuid.New().String()
	u := Upload{
		ID: randomID,
	}
	if u.GetID().String() != randomID {
		t.Errorf(
			"Upload GetID() gave back unexpected ID: %s instead of %s",
			u.GetID().String(),
			randomID,
		)
	}
}

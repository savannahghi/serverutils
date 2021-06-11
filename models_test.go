package go_utils_test

import (
	"testing"

	base "github.com/savannahghi/go_utils"
	"github.com/segmentio/ksuid"
)

func TestUpload_IsEntity(t *testing.T) {
	u := base.Upload{}
	u.IsEntity()
}

func TestUpload_IsNode(t *testing.T) {
	u := base.Upload{}
	u.IsNode()
}

func TestUpload_GetID(t *testing.T) {
	randomID := ksuid.New().String()
	u := base.Upload{
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

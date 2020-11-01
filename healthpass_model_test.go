package base_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestPatient(t *testing.T) {
	name := "testname"
	m1 := base.Patient{
		Name: nil,
	}
	assert.Empty(t, m1.Names())

	m2 := base.Patient{
		Name: []*base.HumanName{
			{
				Text: &name,
			},
		},
	}
	assert.NotEmpty(t, m2.Names())

	m3 := base.Patient{
		Name: []*base.HumanName{
			nil,
		},
	}

	assert.Empty(t, m3.Names())

	m4 := base.Patient{
		Name: []*base.HumanName{
			{
				Text: nil,
			},
		},
	}

	assert.Empty(t, m4.Names())
}

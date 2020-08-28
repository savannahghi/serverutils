package base_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestGenerateRandomWithNDigits(t *testing.T) {
	result, err := base.GenerateRandomWithNDigits(5)
	assert.NotNil(t, result)
	assert.Nil(t, err)
}

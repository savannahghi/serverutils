package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomWithNDigits(t *testing.T) {
	result, err := GenerateRandomWithNDigits(5)
	assert.NotNil(t, result)
	assert.Nil(t, err)
}

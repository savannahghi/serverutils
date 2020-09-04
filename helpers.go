package base

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
)

// GenerateRandomWithNDigits - given a digit generate random numbers
func GenerateRandomWithNDigits(numberOfDigits int) (string, error) {
	rangeEnd := int64(math.Pow10(numberOfDigits) - 1)
	value, _ := rand.Int(rand.Reader, big.NewInt(rangeEnd))
	return strconv.FormatInt(value.Int64(), 10), nil
}

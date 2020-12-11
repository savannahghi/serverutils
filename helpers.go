package base

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"
)

// GenerateRandomWithNDigits - given a digit generate random numbers
func GenerateRandomWithNDigits(numberOfDigits int) (string, error) {
	rangeEnd := int64(math.Pow10(numberOfDigits) - 1)
	value, _ := rand.Int(rand.Reader, big.NewInt(rangeEnd))
	return strconv.FormatInt(value.Int64(), 10), nil
}

// GenerateRandomEmail allows us to get "unique" emails while still keeping
// one main be.well@bewell.co.ke email account
func GenerateRandomEmail() string {
	return fmt.Sprintf("be.well+%v@bewell.co.ke", time.Now().Unix())
}

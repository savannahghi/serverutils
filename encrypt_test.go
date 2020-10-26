package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {

	message1 := "test message 1"

	enc, err := EncryptMessage(message1)
	if err != nil {
		t.Fatalf("unable to encrypt message: %v", err)
	}

	assert.NotEqual(t, message1, *enc)

	dec, err := DecrypMessage(*enc)
	if err != nil {
		t.Fatalf("unable to decrypt message: %v", err)
	}

	assert.NotEqual(t, *enc, *dec)
	assert.Equal(t, message1, *dec)

}

func TestHashTo32Bytes(t *testing.T) {
	key := "this-is-a-test-key"
	hash := hashTo32Bytes(key)
	assert.Equal(t, 32, len(hash))
}

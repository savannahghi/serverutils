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

func TestCreateCoverHash(t *testing.T) {
	cv1 := Cover{
		PayerName:      "payer1",
		PayerSladeCode: 1,
		MemberNumber:   "mem1",
		MemberName:     "name1",
	}

	cv2 := Cover{
		PayerName:      "payer2",
		PayerSladeCode: 2,
		MemberNumber:   "mem2",
		MemberName:     "name2",
	}

	// exactly similar to cv1
	cv3 := Cover{
		PayerName:      "payer1",
		PayerSladeCode: 1,
		MemberNumber:   "mem1",
		MemberName:     "name1",
	}

	// exactly similar to cv2
	cv4 := Cover{
		PayerName:      "payer2",
		PayerSladeCode: 2,
		MemberNumber:   "mem2",
		MemberName:     "name2",
	}

	// similar to cv1 but has different member number
	cv5 := Cover{
		PayerName:      "payer1",
		PayerSladeCode: 1,
		MemberNumber:   "mem11",
		MemberName:     "name1",
	}

	// similar to cv1 but has different member number and member name
	cv6 := Cover{
		PayerName:      "payer1",
		PayerSladeCode: 1,
		MemberNumber:   "mem11",
		MemberName:     "name1",
	}

	cv1Hash := CreateCoverHash(cv1)
	cv2Hash := CreateCoverHash(cv2)
	cv3Hash := CreateCoverHash(cv3)
	cv4Hash := CreateCoverHash(cv4)
	cv5Hash := CreateCoverHash(cv5)
	cv6Hash := CreateCoverHash(cv6)

	assert.NotNil(t, cv1Hash)
	assert.NotNil(t, cv2Hash)
	assert.NotNil(t, cv3Hash)
	assert.NotNil(t, cv4Hash)
	assert.NotNil(t, cv5Hash)
	assert.NotNil(t, cv6Hash)

	assert.NotEqual(t, cv1Hash, cv2Hash)
	assert.Equal(t, cv1Hash, cv3Hash)
	assert.Equal(t, cv2Hash, cv4Hash)

	assert.NotEqual(t, cv1Hash, cv5Hash)
	assert.NotEqual(t, cv3Hash, cv5Hash)
	assert.NotEqual(t, cv1Hash, cv6Hash)
	assert.NotEqual(t, cv3Hash, cv6Hash)
}

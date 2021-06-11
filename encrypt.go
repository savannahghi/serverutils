package go_utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

const messageKeyEnvName = "MESSAGE_KEY"

// GetMessageKey returns the key used to encrypt/decrypt data. The env is not a must
// so will safely use os package to get it. It's the job of services using base
// to enforce if this env is a must
func GetMessageKey() string {
	return os.Getenv(messageKeyEnvName)
}

// HashTo32Bytes will compute a cryptographically useful hash of the input string.
func hashTo32Bytes(input string) []byte {
	data := sha256.Sum256([]byte(input))
	return data[0:]
}

// EncryptMessage takes the message as a string and produces a pointer to a cipher text and an err
func EncryptMessage(plainText string) (*string, error) {

	key := hashTo32Bytes(GetMessageKey())
	encrypted, err := encryptAES(key, []byte(plainText))
	if err != nil {
		return nil, err
	}
	c := base64.URLEncoding.EncodeToString(encrypted)
	return &c, nil
}

func encryptAES(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// create two 'windows' in to the output slice.
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	// populate the IV slice with random data.
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	en := cipher.NewCFBEncrypter(block, iv)

	// note that encrypted is still a window in to the output slice
	en.XORKeyStream(encrypted, data)
	return output, nil
}

// DecrypMessage cipher is the text to be decrypted.
// The function will output the resulting pointer to the plain string and an error
func DecrypMessage(cipher string) (*string, error) {
	encrypted, err := base64.URLEncoding.DecodeString(cipher)
	if err != nil {
		return nil, err
	}
	if len(encrypted) < aes.BlockSize {
		return nil, fmt.Errorf("cipherText is too short. minimum length is 16; got %v", len(encrypted))
	}

	decrypted, err := decryptAES(hashTo32Bytes(GetMessageKey()), encrypted)
	if err != nil {
		return nil, err
	}
	d := string(decrypted)
	return &d, nil
}

func decryptAES(key, data []byte) ([]byte, error) {
	// split the input up in to the IV seed and then the actual encrypted data.
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}

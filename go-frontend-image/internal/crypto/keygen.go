package crypto

import (
	"crypto/rand"
)

func CreateRandomBytes(length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := rand.Read(buf)
	return buf, err
}

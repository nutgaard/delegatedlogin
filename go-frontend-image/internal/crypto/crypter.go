package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"time"
)

type Crypter interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}
type crypter struct {
	key []byte
}

func NewCrypter(key []byte) Crypter {
	return &crypter{key}
}

func (c *crypter) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := IV()
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	return append(nonce, ciphertext...), nil
}

func (c *crypter) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) <= 12 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, ciphertext[:12], ciphertext[12:], nil)
}

func IV() ([]byte, error) {
	storage := make([]byte, 0)
	buf := bytes.NewBuffer(storage)

	err := binary.Write(buf, binary.BigEndian, time.Now().UnixNano())
	if err != nil {
		return nil, err
	}

	random, err := CreateRandomBytes(4)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, random)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

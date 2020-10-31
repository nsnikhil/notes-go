package user

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
	"io"
)

type passwordHasher struct {
	//86, 4096, 32
	saltLength, iterations, keyLength int
}

func newPasswordHasher(saltLength, iterations, keyLength int) *passwordHasher {
	return &passwordHasher{
		saltLength: saltLength,
		iterations: iterations,
		keyLength:  keyLength,
	}
}

func (h *passwordHasher) toSalt(password string) ([]byte, []byte) {
	salt := make([]byte, h.saltLength)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, nil
	}

	return salt, pbkdf2.Key([]byte(password), salt, h.iterations, h.keyLength, sha3.New512)
}

func (h *passwordHasher) toHash(password string, salt []byte) string {
	return h.encodeSalt(pbkdf2.Key([]byte(password), salt, h.iterations, h.keyLength, sha3.New512))
}

func (h *passwordHasher) encodeSalt(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

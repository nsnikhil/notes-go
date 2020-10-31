package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPasswordHasherToSalt(t *testing.T) {
	ph := newPasswordHasher(saltLength, iterations, keyLength)
	_, key := ph.toSalt(password)
	assert.NotNil(t, key)
}

func TestPasswordHasherEncodeSalt(t *testing.T) {
	ph := newPasswordHasher(saltLength, iterations, keyLength)
	_, key := ph.toSalt(password)
	assert.NotNil(t, key)
	assert.NotEmpty(t, ph.encodeSalt(key))
}

func TestPasswordHasherToHash(t *testing.T) {
	ph := newPasswordHasher(saltLength, iterations, keyLength)

	salt, key := ph.toSalt(password)
	assert.NotNil(t, salt)

	h1 := ph.encodeSalt(key)

	h2 := ph.toHash(password, salt)

	assert.Equal(t, h1, h2)
}

package user

import (
	"crypto/ecdsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
	"testing"
)

func TestAuthTokenGenerate(t *testing.T) {
	at, err := newAuthToken(pemString)
	require.NoError(t, err)

	ph := newPasswordHasher(saltLength, iterations, keyLength)

	us, err := newBuilder(ph).name(name).email(email).password(password).rootDirID(dirID).build()
	require.NoError(t, err)

	ac, rf, err := at.generate(us)
	require.NoError(t, err)

	privateKey, err := ssh.ParseRawPrivateKey([]byte(pemString))
	require.NoError(t, err)

	rac, err := jwt.ParseWithClaims(ac, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey.(*ecdsa.PrivateKey).Public(), nil
	})
	require.NoError(t, err)
	assert.Equal(t, name, rac.Claims.(*accessTokenClaims).Name)

	rrf, err := jwt.ParseWithClaims(rf, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey.(*ecdsa.PrivateKey).Public(), nil
	})
	require.NoError(t, err)
	assert.Equal(t, "notes", rrf.Claims.(*refreshTokenClaims).Issuer)
}

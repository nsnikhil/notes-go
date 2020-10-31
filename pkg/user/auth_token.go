package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"notes/pkg/liberr"
	"time"
)

type authToken struct {
	signingKey interface{}
}

func newAuthToken(pemString string) (*authToken, error) {
	privateKey, err := ssh.ParseRawPrivateKey([]byte(pemString))
	if err != nil {
		return nil, liberr.WithArgs(liberr.Operation("authToken.newAuthToken"), err)
	}

	return &authToken{signingKey: privateKey}, nil
}

func (tg *authToken) generate(user *user) (string, string, error) {
	return generateTokens(user, tg.signingKey)
}

func generateTokens(user *user, signingKey interface{}) (string, string, error) {
	now := time.Now()
	accessToken, refreshToken := generateAccessToken(now, user), generateRefreshToken(now, user)

	as, err := accessToken.SignedString(signingKey)
	if err != nil {
		return "", "", liberr.WithArgs(liberr.Operation("authToken.generateTokens"), err)
	}

	rs, err := refreshToken.SignedString(signingKey)
	if err != nil {
		return "", "", liberr.WithArgs(liberr.Operation("authToken.generateTokens"), err)
	}

	return as, rs, nil
}

func generateAccessToken(now time.Time, user *user) *jwt.Token {
	claims := getAccessTokenClaims(now, user.id, user.name, user.email, user.rootDirID)
	return jwt.NewWithClaims(jwt.SigningMethodES384, claims)
}

func generateRefreshToken(now time.Time, user *user) *jwt.Token {
	claims := getRefreshTokenClaims(now, user.id)
	return jwt.NewWithClaims(jwt.SigningMethodES384, claims)
}

type accessTokenClaims struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	RootDirID string `json:"root_dir_id"`
	jwt.StandardClaims
}

func getAccessTokenClaims(now time.Time, id, name, email, rootDirID string) jwt.Claims {
	return accessTokenClaims{
		Id:             id,
		Name:           name,
		Email:          email,
		RootDirID:      rootDirID,
		StandardClaims: getStandardClaims(now, now.Add(time.Hour).Unix(), id),
	}
}

type refreshTokenClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}

func getRefreshTokenClaims(now time.Time, id string) jwt.Claims {
	return refreshTokenClaims{
		Id:             id,
		StandardClaims: getStandardClaims(now, now.AddDate(0, 0, 10).Unix(), id),
	}
}

func getStandardClaims(now time.Time, expiresAt int64, id string) jwt.StandardClaims {
	return jwt.StandardClaims{
		Audience:  "user",
		ExpiresAt: expiresAt,
		Id:        uuid.New().String(),
		IssuedAt:  now.Unix(),
		Issuer:    "notes",
		NotBefore: now.Unix(),
		Subject:   id,
	}
}

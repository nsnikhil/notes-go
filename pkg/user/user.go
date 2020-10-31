package user

import (
	"errors"
	"fmt"
	"notes/pkg/liberr"
	"regexp"
	"time"
	"unicode"
)

type user struct {
	id string

	name  string
	email string

	passwordHash string
	passwordSalt []byte

	rootDirID string

	createdAt time.Time
	updatedAt time.Time
}

// TODO: IS THIS REALLY NEEDED ?
type builder struct {
	user           *user
	passwordHasher *passwordHasher
	err            error
	op             string
}

func newBuilder(passwordHasher *passwordHasher) *builder {
	return &builder{
		user:           &user{},
		passwordHasher: passwordHasher,
	}
}

func (b *builder) id(id string) *builder {
	if b.err != nil {
		return b
	}

	if len(id) == 0 {
		b.err = errors.New("id cannot be empty")
		b.op = "ID"
		return b
	}

	if !isValidUUID(id) {
		b.err = errors.New("invalid user id")
		b.op = "ID"
		return b
	}

	b.user.id = id
	return b
}

func (b *builder) name(name string) *builder {
	if b.err != nil {
		return b
	}

	if len(name) == 0 {
		b.err = errors.New("user name cannot be empty")
		b.op = "Name"
		return b
	}

	b.user.name = name
	return b
}

func (b *builder) email(email string) *builder {
	if b.err != nil {
		return b
	}

	if len(email) == 0 {
		b.err = errors.New("email cannot be empty")
		b.op = "Email"
		return b
	}

	b.user.email = email
	return b
}

func (b *builder) password(password string) *builder {
	if b.err != nil {
		return b
	}

	if len(password) == 0 {
		b.err = errors.New("password cannot be empty")
		b.op = "Password"
		return b
	}

	if !isValidPassword(password) {
		b.err = errors.New("invalid password")
		b.op = "Password"
		return b
	}

	salt, key := b.passwordHasher.toSalt(password)
	hash := b.passwordHasher.encodeSalt(key)

	b.user.passwordSalt = salt
	b.user.passwordHash = hash
	return b
}

func (b *builder) rootDirID(rootDirID string) *builder {
	if b.err != nil {
		return b
	}

	if len(rootDirID) == 0 {
		b.err = errors.New("root dir id cannot be empty")
		b.op = "Password"
		return b
	}

	if !isValidUUID(rootDirID) {
		b.err = errors.New("invalid root dir id")
		b.op = "Password"
		return b
	}

	b.user.rootDirID = rootDirID
	return b
}

func (b *builder) createdAt(createdAt time.Time) *builder {
	if b.err != nil {
		return b
	}

	b.user.createdAt = createdAt
	return b
}

func (b *builder) updatedAt(updatedAt time.Time) *builder {
	if b.err != nil {
		return b
	}

	b.user.updatedAt = updatedAt
	return b
}

func (b *builder) build() (*user, error) {
	if b.err != nil {
		return nil, liberr.WithArgs(liberr.Operation(fmt.Sprintf("builder.%s", b.op)), liberr.ValidationError, b.err)
	}

	return b.user, nil
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	u, l, n, s := 0, 0, 0, 0

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			n++
		case unicode.IsLower(c):
			l++
		case unicode.IsUpper(c):
			u++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			s++
		default:
			return false
		}
	}
	return u > 0 && l > 0 && n > 0 && s > 0
}

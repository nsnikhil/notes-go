package user

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/sha3"
	"io"
	"notes/pkg/liberr"
	"regexp"
	"time"
	"unicode"
)

type User struct {
	id string

	name     string
	email    string
	password string

	rootDirID string

	createdAt time.Time
	updatedAt time.Time
}

type Builder struct {
	*User
	err error
	op  string
}

func NewBuilder() *Builder {
	return &Builder{User: &User{}}
}

func (b *Builder) ID(id string) *Builder {
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

	b.id = id
	return b
}

func (b *Builder) Name(name string) *Builder {
	if b.err != nil {
		return b
	}

	if len(name) == 0 {
		b.err = errors.New("user name cannot be empty")
		b.op = "Name"
		return b
	}

	b.name = name
	return b
}

func (b *Builder) Email(email string) *Builder {
	if b.err != nil {
		return b
	}

	if len(email) == 0 {
		b.err = errors.New("email cannot be empty")
		b.op = "Email"
		return b
	}

	b.email = email
	return b
}

func (b *Builder) Password(password string) *Builder {
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

	b.password = password
	return b
}

func (b *Builder) RootDirID(rootDirID string) *Builder {
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

	b.rootDirID = rootDirID
	return b
}

func (b *Builder) CreatedAt(createdAt time.Time) *Builder {
	if b.err != nil {
		return b
	}

	b.createdAt = createdAt
	return b
}

func (b *Builder) UpdatedAt(updatedAt time.Time) *Builder {
	if b.err != nil {
		return b
	}

	b.updatedAt = updatedAt
	return b
}

func (b *Builder) Build() (*User, error) {
	if b.err != nil {
		return nil, liberr.WithArgs(liberr.Operation(fmt.Sprintf("Builder.%s", b.op)), liberr.ValidationError, b.err)
	}

	return b.User, nil
}

//TODO: USE THIS
func encode(password string) string {
	salt := make([]byte, 84)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return password
	}

	dk := pbkdf2.Key([]byte(password), salt, 4096, 32, sha3.New512)
	return base64.StdEncoding.EncodeToString(dk)
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

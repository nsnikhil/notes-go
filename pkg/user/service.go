package user

import (
	"database/sql"
	"errors"
	"notes/pkg/liberr"
	"notes/pkg/notes/directory"
)

const (
	rootDirectoryName = "/"
)

type Service interface {
	Create(name, email, password string) (string, error)
	Login(email, password string) (string, string, error)
}

// TODO: RENAME
type userService struct {
	ds directory.Service
	st store

	hasher *passwordHasher
	token  *authToken
}

func (us *userService) Create(name, email, password string) (string, error) {
	wrap := func(err error) error { return liberr.WithOp("Service.createUser", err) }

	//TODO: SPIKE IF YOU COULD ACHIEVE THIS VIA EVENTS ?
	dirID, err := us.ds.CreateDirectory(rootDirectoryName)
	if err != nil {
		return "", wrap(err)
	}

	user, err := newBuilder(us.hasher).name(name).email(email).password(password).rootDirID(dirID).build()
	if err != nil {
		return "", wrap(err)
	}

	userID, err := us.st.createUser(user)
	if err != nil {
		return "", wrap(err)
	}

	return userID, nil
}

func (us *userService) Login(email, password string) (string, string, error) {
	wrap := func(err error) error { return liberr.WithOp("Service.Login", err) }

	user, err := us.st.getByEmail(email)
	if err != nil {
		return "", "", wrap(err)
	}

	if user.passwordHash != us.hasher.toHash(password, user.passwordSalt) {
		return "", "", wrap(errors.New("invalid credentials"))
	}

	ac, rf, err := us.token.generate(user)
	if err != nil {
		return "", "", wrap(err)
	}

	return ac, rf, nil
}

func newService(st store, ds directory.Service, hasher *passwordHasher, token *authToken) Service {
	return &userService{
		st:     st,
		ds:     ds,
		hasher: hasher,
		token:  token,
	}
}

func NewService(ds directory.Service, db *sql.DB, saltLength, iterations, keyLength int, pemString string) (Service, error) {
	ust := newUserStore(db)
	ph := newPasswordHasher(saltLength, iterations, keyLength)

	tk, err := newAuthToken(pemString)
	if err != nil {
		return nil, liberr.WithOp("Service.NewService", err)
	}

	return newService(ust, ds, ph, tk), nil
}

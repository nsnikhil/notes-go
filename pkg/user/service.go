package user

import (
	"notes/pkg/liberr"
	"notes/pkg/notes/directory"
)

const (
	rootDirectoryName = "/"
)

type Service interface {
	CreateUser(name, email, password string) (string, error)
}

// TODO: RENAME
type userService struct {
	ds directory.Service
	st Store
}

func (us *userService) CreateUser(name, email, password string) (string, error) {
	wrap := func(err error) error { return liberr.WithOp("Service.CreateUser", err) }

	dirID, err := us.ds.CreateDirectory(rootDirectoryName)
	if err != nil {
		return "", wrap(err)
	}

	user, err := NewBuilder().Name(name).Email(email).Password(password).RootDirID(dirID).Build()
	if err != nil {
		return "", wrap(err)
	}

	userID, err := us.st.CreateUser(user)
	if err != nil {
		return "", wrap(err)
	}

	return userID, nil
}

func NewService(st Store, ds directory.Service) Service {
	return &userService{
		st: st,
		ds: ds,
	}
}

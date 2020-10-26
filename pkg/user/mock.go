package user

import (
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (mock *MockUserService) CreateUser(name, email, password string) (string, error) {
	args := mock.Called(name, email, password)
	return args.String(0), args.Error(1)
}

type MockUserStore struct {
	mock.Mock
}

func (mock *MockUserStore) CreateUser(user *User) (string, error) {
	args := mock.Called(user)
	return args.String(0), args.Error(1)
}

package user

import (
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (mock *MockUserService) Create(name, email, password string) (string, error) {
	args := mock.Called(name, email, password)
	return args.String(0), args.Error(1)
}

func (mock *MockUserService) Login(email, password string) (string, string, error) {
	args := mock.Called(email, password)
	return args.String(0), args.String(1), args.Error(2)
}

type mockUserStore struct {
	mock.Mock
}

func (mock *mockUserStore) createUser(user *user) (string, error) {
	args := mock.Called(user)
	return args.String(0), args.Error(1)
}

func (mock *mockUserStore) getByEmail(email string) (*user, error) {
	args := mock.Called(email)
	return args.Get(0).(*user), args.Error(1)
}

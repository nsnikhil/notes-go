package directory

import "github.com/stretchr/testify/mock"

type MockDirectoryService struct {
	mock.Mock
}

func (mock *MockDirectoryService) CreateDirectory(name string) (string, error) {
	args := mock.Called(name)
	return args.String(0), args.Error(1)
}

type mockDirectoryStore struct {
	mock.Mock
}

func (mock *mockDirectoryStore) createDirectory(directory *directory) (string, error) {
	args := mock.Called(directory)
	return args.String(0), args.Error(1)
}

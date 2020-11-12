package directory

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"notes/pkg/liberr"
	"testing"
)

const (
	dirID = "d87f0cb7-46c0-4501-83fa-e1c5ed5338e6"
)

func TestCreateNewDirectorySuccess(t *testing.T) {
	st := &mockDirectoryStore{}
	st.On("createDirectory", mock.AnythingOfType("*directory.directory")).Return(dirID, nil)

	testCreateNewDirectory(t, nil, "root", newDirectoryService(st))
}

func TestCreateNewDirectoryFailure(t *testing.T) {
	testCases := map[string]struct {
		svc           func() Service
		name          string
		expectedError error
	}{
		"test failure when db fails to save": {
			svc: func() Service {
				st := &mockDirectoryStore{}
				st.On("createDirectory", mock.AnythingOfType("*directory.directory")).Return("", liberr.WithArgs(errors.New("failed to create new directory")))

				return newDirectoryService(st)
			},
			name:          "root",
			expectedError: liberr.WithArgs(liberr.Operation("DirectoryService.createDirectory"), liberr.WithArgs(errors.New("failed to create new directory"))),
		},
		"test failure when name is invalid": {
			svc: func() Service {
				st := &mockDirectoryStore{}
				st.On("createDirectory", mock.AnythingOfType("*directory.directory")).Return(dirID, nil)

				return newDirectoryService(st)
			},
			name: "",
			expectedError: liberr.WithArgs(
				liberr.Operation("DirectoryService.createDirectory"),
				liberr.WithArgs(liberr.Operation("directory.validate"), liberr.ValidationError, liberr.SeverityError,
					errors.New("directory name cannot be empty")),
			),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testCreateNewDirectory(t, testCase.expectedError, testCase.name, testCase.svc())
		})
	}
}

func testCreateNewDirectory(t *testing.T, expectedError error, name string, svc Service) {
	_, err := svc.CreateDirectory(name)
	assert.Equal(t, expectedError, err)
}

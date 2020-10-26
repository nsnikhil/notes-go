package directory_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"notes/pkg/liberr"
	"notes/pkg/notes/directory"
	"testing"
)

const (
	dirID = "d87f0cb7-46c0-4501-83fa-e1c5ed5338e6"
)

func TestCreateNewDirectorySuccess(t *testing.T) {
	st := &directory.MockDirectoryStore{}
	st.On("CreateDirectory", mock.AnythingOfType("*directory.Directory")).Return(dirID, nil)

	testCreateNewDirectory(t, nil, "root", directory.NewDirectoryService(st))
}

func TestCreateNewDirectoryFailure(t *testing.T) {
	testCases := map[string]struct {
		svc           func() directory.Service
		name          string
		expectedError error
	}{
		"test failure when db fails to save": {
			svc: func() directory.Service {
				st := &directory.MockDirectoryStore{}
				st.On("CreateDirectory", mock.AnythingOfType("*directory.Directory")).Return("", liberr.WithArgs(errors.New("failed to create new directory")))

				return directory.NewDirectoryService(st)
			},
			name:          "root",
			expectedError: liberr.WithArgs(liberr.Operation("DirectoryService.CreateDirectory"), liberr.WithArgs(errors.New("failed to create new directory"))),
		},
		"test failure when name is invalid": {
			svc: func() directory.Service {
				st := &directory.MockDirectoryStore{}
				st.On("CreateDirectory", mock.AnythingOfType("*directory.Directory")).Return(dirID, nil)

				return directory.NewDirectoryService(st)
			},
			name: "",
			expectedError: liberr.WithArgs(
				liberr.Operation("DirectoryService.CreateDirectory"),
				liberr.WithArgs(liberr.Operation("Directory.validate"), liberr.ValidationError, liberr.SeverityError,
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

func testCreateNewDirectory(t *testing.T, expectedError error, name string, svc directory.Service) {
	_, err := svc.CreateDirectory(name)
	assert.Equal(t, expectedError, err)
}

package user_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"notes/pkg/liberr"
	"notes/pkg/notes/directory"
	"notes/pkg/user"
	"testing"
)

const (
	dirID  = "d87f0cb7-46c0-4501-83fa-e1c5ed5338e6"
	userID = "86d690dd-92a0-40ac-ad48-110c951e3cb8"
)

func TestCreatNewUserSuccess(t *testing.T) {
	ds := &directory.MockDirectoryService{}
	ds.On("CreateDirectory", "/").Return(dirID, nil)

	st := &user.MockUserStore{}
	st.On("CreateUser", mock.AnythingOfType("*user.User")).Return(userID, nil)

	testServiceCreatNewUser(t, nil, ds, st)
}

func TestCreatNewUserFailure(t *testing.T) {
	testCases := map[string]struct {
		svc           func() directory.Service
		st            func() user.Store
		expectedError error
	}{
		"test failure when directory creation fails": {
			svc: func() directory.Service {
				ds := &directory.MockDirectoryService{}
				ds.On("CreateDirectory", "/").Return("root", liberr.WithArgs(errors.New("failed to create new directory")))
				return ds
			},
			st: func() user.Store {
				return &user.MockUserStore{}
			},
			expectedError: liberr.WithArgs(liberr.Operation("Service.CreateUser"), liberr.WithArgs(errors.New("failed to create new directory"))),
		},
		"test failure when create new user fails": {
			svc: func() directory.Service {
				ds := &directory.MockDirectoryService{}
				ds.On("CreateDirectory", "/").Return(dirID, nil)
				return ds
			},
			st: func() user.Store {
				st := &user.MockUserStore{}
				st.On("CreateUser", mock.AnythingOfType("*user.User")).Return(userID, liberr.WithArgs(errors.New("failed to save new user")))
				return st
			},
			expectedError: liberr.WithArgs(liberr.Operation("Service.CreateUser"), liberr.WithArgs(errors.New("failed to save new user"))),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testServiceCreatNewUser(t, testCase.expectedError, testCase.svc(), testCase.st())
		})
	}
}

func testServiceCreatNewUser(t *testing.T, expectedError error, ds directory.Service, st user.Store) {
	us := user.NewService(st, ds)

	_, err := us.CreateUser(name, email, password)
	assert.Equal(t, expectedError, err)
}

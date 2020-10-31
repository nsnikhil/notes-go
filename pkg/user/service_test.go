package user

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"notes/pkg/liberr"
	"notes/pkg/notes/directory"
	"testing"
)

func TestCreatNewUserSuccess(t *testing.T) {
	ds := &directory.MockDirectoryService{}
	ds.On("CreateDirectory", "/").Return(dirID, nil)

	st := &mockUserStore{}
	st.On("createUser", mock.AnythingOfType("*user.user")).Return(userID, nil)

	testServiceCreatNewUser(t, nil, ds, st)
}

func TestCreatNewUserFailure(t *testing.T) {
	testCases := map[string]struct {
		svc           func() directory.Service
		st            func() store
		expectedError error
	}{
		"test failure when directory creation fails": {
			svc: func() directory.Service {
				ds := &directory.MockDirectoryService{}
				ds.On("CreateDirectory", "/").Return("root", liberr.WithArgs(errors.New("failed to create new directory")))
				return ds
			},
			st: func() store {
				return &mockUserStore{}
			},
			expectedError: liberr.WithArgs(liberr.Operation("Service.createUser"), liberr.WithArgs(errors.New("failed to create new directory"))),
		},
		"test failure when create new user fails": {
			svc: func() directory.Service {
				ds := &directory.MockDirectoryService{}
				ds.On("CreateDirectory", "/").Return(dirID, nil)
				return ds
			},
			st: func() store {
				st := &mockUserStore{}
				st.On("createUser", mock.AnythingOfType("*user.user")).Return(userID, liberr.WithArgs(errors.New("failed to save new user")))
				return st
			},
			expectedError: liberr.WithArgs(liberr.Operation("Service.createUser"), liberr.WithArgs(errors.New("failed to save new user"))),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testServiceCreatNewUser(t, testCase.expectedError, testCase.svc(), testCase.st())
		})
	}
}

func testServiceCreatNewUser(t *testing.T, expectedError error, ds directory.Service, st store) {
	at, err := newAuthToken(pemString)
	require.NoError(t, err)

	ph := newPasswordHasher(86, 4096, 32)

	us := newService(st, ds, ph, at)

	_, err = us.Create(name, email, password)
	assert.Equal(t, expectedError, err)
}

func TestLoginUserSuccess(t *testing.T) {
	hash := "7pbNhPbzAz46dY/7qV2KlzewmnvzAYMhioJku7BgGbE="
	salt := []byte(`\x5fe5782aeef2bc2299ddceda5d29e679c0f2bb5828ecb2a0023909fdeab74879`)

	st := &mockUserStore{}
	cu := &user{
		name:         name,
		email:        email,
		rootDirID:    dirID,
		passwordHash: hash,
		passwordSalt: salt,
	}

	st.On("getByEmail", email).Return(cu, nil)
	testLoginUser(t, nil, st)
}

func TestLoginUserFailure(t *testing.T) {
	testCases := map[string]struct {
		store         func() store
		expectedError error
	}{
		"test failure when store returns error": {
			store: func() store {
				st := &mockUserStore{}
				st.On("getByEmail", email).Return(&user{}, liberr.WithArgs(errors.New("failed to get data")))

				return st
			},
			expectedError: liberr.WithArgs(liberr.Operation("Service.Login"), liberr.WithArgs(errors.New("failed to get data"))),
		},
		"test failure when credentials are invalid": {
			store: func() store {
				st := &mockUserStore{}

				cu := &user{
					name:         name,
					email:        email,
					rootDirID:    dirID,
					passwordHash: "invalid/7qV2KlzewmnvzAYMhioJku7BgGbE=",
					passwordSalt: []byte(`\x5fe5782aeef2bc2299ddceda5d29e679c0f2bb5828ecb2a0023909fdeab74879`),
				}

				st.On("getByEmail", email).Return(cu, nil)

				return st
			},
			expectedError: liberr.WithArgs(liberr.Operation("Service.Login"), errors.New("invalid credentials")),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testLoginUser(t, testCase.expectedError, testCase.store())
		})
	}
}

func testLoginUser(t *testing.T, expectedError error, st store) {
	at, err := newAuthToken(pemString)
	require.NoError(t, err)

	ph := newPasswordHasher(86, 4096, 32)

	us := newService(st, nil, ph, at)

	_, _, err = us.Login(email, password)
	assert.Equal(t, expectedError, err)
}

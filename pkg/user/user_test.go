package user_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"notes/pkg/liberr"
	"notes/pkg/user"
	"testing"
)

const (
	name     = "Test Name"
	email    = "test@test.com"
	password = "Password@1234"

	emptyString = ""

	invalidPasswordOne   = "password@1234"
	invalidPasswordTwo   = "PASSWORD@1234"
	invalidPasswordThree = "Password@"
	invalidPasswordFour  = "Password1"
	invalidPasswordFive  = "Pa@1"
)

func TestCreateNewUserSuccess(t *testing.T) {
	testCreateNewUser(t, nil, name, email, password)
}

func TestCreateNewUserValidationFailure(t *testing.T) {
	testCases := map[string]struct {
		input         func() (string, string, string)
		expectedError error
	}{
		"test failure when name is empty": {
			input: func() (string, string, string) {
				return emptyString, email, password
			},
			expectedError: liberr.WithArgs(liberr.Operation("Builder.Name"), liberr.ValidationError, liberr.SeverityError, errors.New("user name cannot be empty")),
		},

		"test failure when email is empty": {
			input: func() (string, string, string) {
				return name, emptyString, password
			},
			expectedError: liberr.WithArgs(liberr.Operation("Builder.Email"), liberr.ValidationError, liberr.SeverityError, errors.New("email cannot be empty")),
		},

		"test failure when password is empty": {
			input: func() (string, string, string) {
				return name, email, emptyString
			},
			expectedError: liberr.WithArgs(liberr.Operation("Builder.Password"), liberr.ValidationError, liberr.SeverityError, errors.New("password cannot be empty")),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			name, email, password := testCase.input()
			testCreateNewUser(t, testCase.expectedError, name, email, password)
		})
	}
}

func TestCreateNewUserPasswordFailure(t *testing.T) {
	testCases := map[string]struct {
		input         func() (string, string, string)
		expectedError error
	}{
		"test failure for invalid password one": {
			input: func() (string, string, string) {
				return name, email, invalidPasswordOne
			},
			expectedError: liberr.WithArgs(liberr.Operation("Builder.Password"), liberr.ValidationError, liberr.SeverityError, errors.New("invalid password")),
		},

		"test failure for invalid password two": {
			input: func() (string, string, string) {
				return name, email, invalidPasswordTwo
			},
			expectedError: liberr.WithArgs(liberr.Operation("Builder.Password"), liberr.ValidationError, liberr.SeverityError, errors.New("invalid password")),
		},

		"test failure for invalid password three": {
			input: func() (string, string, string) {
				return name, email, invalidPasswordThree
			},
			expectedError: liberr.WithArgs(liberr.Operation("Builder.Password"), liberr.ValidationError, liberr.SeverityError, errors.New("invalid password")),
		},

		"test failure for invalid password four": {
			input: func() (string, string, string) {
				return name, email, invalidPasswordFour
			},
			expectedError: liberr.WithArgs(liberr.Operation("Builder.Password"), liberr.ValidationError, liberr.SeverityError, errors.New("invalid password")),
		},

		"test failure for invalid password five": {
			input: func() (string, string, string) {
				return name, email, invalidPasswordFive
			},
			expectedError: liberr.WithArgs(liberr.Operation("Builder.Password"), liberr.ValidationError, liberr.SeverityError, errors.New("invalid password")),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			name, email, password := testCase.input()
			testCreateNewUser(t, testCase.expectedError, name, email, password)
		})
	}
}

func testCreateNewUser(t *testing.T, expectedError error, name, email, password string) {
	_, err := user.NewBuilder().Name(name).Email(email).Password(password).Build()
	assert.Equal(t, expectedError, err)
}

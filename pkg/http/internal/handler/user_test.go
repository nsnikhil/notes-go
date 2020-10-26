package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"notes/pkg/http/contract"
	"notes/pkg/http/internal/handler"
	mdl "notes/pkg/http/internal/middleware"
	"notes/pkg/liberr"
	reporters "notes/pkg/reporting"
	"notes/pkg/user"
	"testing"
)

const (
	name     = "Nikhil Soni"
	email    = "n.nikhil.ns65@gmail.com"
	password = "Password@1234"

	userID = "86d690dd-92a0-40ac-ad48-110c951e3cb8"
)

func TestCreateUserSuccess(t *testing.T) {
	svc := &user.MockUserService{}
	svc.On("CreateUser", name, email, password).Return(userID, nil)

	req := contract.CreateUserRequest{Name: name, Email: email, Password: password}

	b, err := json.Marshal(req)
	require.NoError(t, err)

	expectedBody := `{"data":{"message":"user created successfully"},"success":true}`

	testCreateUser(t, http.StatusCreated, expectedBody, bytes.NewBuffer(b), svc)
}

func TestCreateUserFailure(t *testing.T) {
	testCases := map[string]struct {
		svc          func() user.Service
		body         func() io.Reader
		expectedCode int
		expectedBody string
	}{
		"test failure when body parsing fails": {
			svc: func() user.Service {
				return &user.MockUserService{}
			},
			body: func() io.Reader {
				return nil
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":{"message":"unexpected end of JSON input"},"success":false}`,
		},
		"test failure when svc call fails fails": {
			svc: func() user.Service {
				svc := &user.MockUserService{}
				svc.On("CreateUser", name, email, password).Return("", liberr.WithArgs(errors.New("failed to create new user")))

				return svc
			},
			body: func() io.Reader {
				req := contract.CreateUserRequest{Name: name, Email: email, Password: password}

				b, err := json.Marshal(req)
				require.NoError(t, err)

				return bytes.NewBuffer(b)
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":{"message":"internal server error"},"success":false}`,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testCreateUser(t, testCase.expectedCode, testCase.expectedBody, testCase.body(), testCase.svc())
		})
	}
}

func testCreateUser(t *testing.T, expectedCode int, expectedBody string, body io.Reader, svc user.Service) {
	lgr := reporters.NewLogger("dev", "debug")

	uh := handler.NewUserHandler(svc)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodPost, "/user/create", body)

	mdl.WithError(lgr, uh.CreateUser)(w, r)

	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
}

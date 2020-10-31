package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	svc.On("Create", name, email, password).Return(userID, nil)

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
				svc.On("Create", name, email, password).Return("", liberr.WithArgs(errors.New("failed to create new user")))

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

func TestLoginSuccess(t *testing.T) {
	at := "eyJhbGciOiJFUzM4NCIsInR5cCI6IkpXVCJ9.eyJpZCI6IiIsIm5hbWUiOiJUZXN0IE5hbWUiLCJlbWFpbCI6InRlc3RAdGVzdC5jb20iLCJyb290X2Rpcl9pZCI6ImQ4N2YwY2I3LTQ2YzAtNDUwMS04M2ZhLWUxYzVlZDUzMzhlNiIsImF1ZCI6InVzZXIiLCJleHAiOjE2MDQxMzA1OTIsImp0aSI6ImU4MmI5MDM0LTBlOWMtNDU5My05ZGQ1LWVkNjc2NmE0YjA4MiIsImlhdCI6MTYwNDEyNjk5MiwiaXNzIjoibm90ZXMiLCJuYmYiOjE2MDQxMjY5OTJ9.6fOocSi6b1ircmPMISlMdIcN7s1dd-SxgtzFluNKPok7yDO6EJZcHF0CmBqAaJHNEvFpqL65OT80Ne1GeV5baphpJRqcY1GEEb7-05ub4AgaiBEOb1E8MJzGCofhnd7c"
	rf := "eyJhbGciOiJFUzM4NCIsInR5cCI6IkpXVCJ9.eyJpZCI6IiIsImF1ZCI6InVzZXIiLCJleHAiOjE2MDQ5OTA5OTIsImp0aSI6IjhhZjcwYTA0LTcyZDktNDVhNi04YjY2LWU3MmE1N2RmMDEwNSIsImlhdCI6MTYwNDEyNjk5MiwiaXNzIjoibm90ZXMiLCJuYmYiOjE2MDQxMjY5OTJ9.PnpayfQxh1R6yOws42nevgmM6_m1i3wScOXaN8SerVdRZqcHO_YFeEgqQB3uC6MTfBomYO7ihP_9_RHLJ9IC547MmlCV1_LsAHCkgYCb8z478lWOu8NIxtKuQHkL50Qa"

	svc := &user.MockUserService{}
	svc.On("Login", email, password).Return(at, rf, nil)

	req := contract.LoginUserRequest{Email: email, Password: password}

	b, err := json.Marshal(req)
	require.NoError(t, err)

	expectedBody := fmt.Sprintf(`{"data":{"access_token":"%s","refresh_token":"%s"},"success":true}`, at, rf)

	testLogin(t, http.StatusOK, expectedBody, bytes.NewBuffer(b), svc)
}

func TestLoginFailure(t *testing.T) {
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
				svc.On("Login", email, password).Return("", "", liberr.WithArgs(errors.New("failed to login")))

				return svc
			},
			body: func() io.Reader {
				req := contract.LoginUserRequest{Email: email, Password: password}

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
			testLogin(t, testCase.expectedCode, testCase.expectedBody, testCase.body(), testCase.svc())
		})
	}
}

func testLogin(t *testing.T, expectedCode int, expectedBody string, body io.Reader, svc user.Service) {
	lgr := reporters.NewLogger("dev", "debug")

	uh := handler.NewUserHandler(svc)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(http.MethodPost, "/user/login", body)

	mdl.WithError(lgr, uh.LoginUser)(w, r)

	assert.Equal(t, expectedCode, w.Code)
	assert.Equal(t, expectedBody, w.Body.String())
}

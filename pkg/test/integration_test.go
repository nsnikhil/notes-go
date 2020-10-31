package test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"notes/pkg/app"
	"notes/pkg/config"
	"notes/pkg/database"
	"notes/pkg/http/contract"
	"testing"
	"time"
)

const (
	name     = "Test Name"
	email    = "test@test.com"
	password = "Password@1234"

	address = "http://127.0.0.1:8080"
)

func TestAPI(t *testing.T) {
	startApp()
	defer truncateTables(t)

	cl := &http.Client{Timeout: time.Minute}

	testPing(t, cl)
	testUserAPI(t, cl)
}

func testPing(t *testing.T, cl *http.Client) {
	req := newRequest(t, http.MethodGet, "ping", nil)
	resp := execRequest(t, cl, req)

	responseData := getData(t, http.StatusOK, resp)
	expectedData := contract.APIResponse{Success: true, Data: "pong"}

	verifyData(t, expectedData, responseData)
}

func testUserAPI(t *testing.T, cl *http.Client) {
	testCreateUser(t, cl)
	testLoginUser(t, cl)
}

func testCreateUser(t *testing.T, cl *http.Client) {
	crq := contract.CreateUserRequest{Name: name, Email: email, Password: password}
	b, err := json.Marshal(&crq)
	require.NoError(t, err)

	req := newRequest(t, http.MethodPost, "user/create", bytes.NewBuffer(b))
	resp := execRequest(t, cl, req)

	responseData := getData(t, http.StatusCreated, resp)
	expectedData := contract.APIResponse{Success: true, Data: map[string]interface{}{"message": "user created successfully"}}

	verifyData(t, expectedData, responseData)
}

func testLoginUser(t *testing.T, cl *http.Client) {
	lrq := contract.LoginUserRequest{Email: email, Password: password}
	b, err := json.Marshal(&lrq)
	require.NoError(t, err)

	req := newRequest(t, http.MethodPost, "user/login", bytes.NewBuffer(b))
	resp := execRequest(t, cl, req)

	responseData := getData(t, http.StatusOK, resp)

	require.True(t, responseData.Success)
	require.True(t, len(responseData.Data.(map[string]interface{})["access_token"].(string)) != 0)
	require.True(t, len(responseData.Data.(map[string]interface{})["refresh_token"].(string)) != 0)
}

func execRequest(t *testing.T, cl *http.Client, req *http.Request) *http.Response {
	resp, err := cl.Do(req)
	require.NoError(t, err)

	return resp
}

func getData(t *testing.T, expectedCode int, resp *http.Response) contract.APIResponse {
	require.Equal(t, expectedCode, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var res contract.APIResponse
	err = json.Unmarshal(b, &res)
	require.NoError(t, err)

	return res
}

func verifyData(t *testing.T, expectedResponse contract.APIResponse, response contract.APIResponse) {
	require.Equal(t, expectedResponse, response)
}

func newRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", address, path), body)
	require.NoError(t, err)

	return req
}

func startApp() {
	configFile := "test.env"
	database.RunMigrations(configFile)
	go app.StartHTTPServer(configFile)
	time.Sleep(time.Second)
}

func truncateTables(t *testing.T) {
	configFile := "test.env"
	cfg := config.NewConfig(configFile)

	db, err := database.NewHandler(cfg.DatabaseConfig()).GetDB()
	require.NoError(t, err)

	_, err = db.Exec(`TRUNCATE documents`)
	require.NoError(t, err)

	_, err = db.Exec(`TRUNCATE directories CASCADE`)
	require.NoError(t, err)

	_, err = db.Exec(`TRUNCATE users`)
	require.NoError(t, err)
}

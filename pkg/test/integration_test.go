package test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
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

const address = "http://127.0.0.1:8080"

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
	verifyResponse(t, http.StatusOK, "pong", resp)
}

func testUserAPI(t *testing.T, cl *http.Client) {
	testCreateUser(t, cl)
	testLoginUser(t, cl)
}

func testCreateUser(t *testing.T, cl *http.Client) {
	crq := contract.CreateUserRequest{Name: "test name", Email: "test@test.com", Password: "Password@1234"}
	b, err := json.Marshal(&crq)
	require.NoError(t, err)

	req := newRequest(t, http.MethodPost, "user/create", bytes.NewBuffer(b))
	resp := execRequest(t, cl, req)
	verifyResponse(t, http.StatusCreated, map[string]interface{}{"message": "user created successfully"}, resp)
}

func testLoginUser(t *testing.T, cl *http.Client) {
	//lrq := contract.LoginUserRequest{Email: "test@test.com", Password: "Password@1234"}
	//b, err := json.Marshal(&lrq)
	//require.NoError(t, err)

	//req := newRequest(t, http.MethodPost, "user/login", bytes.NewBuffer(b))
	//resp := execRequest(t, cl, req)
	//verifyResponse(t, http.StatusOK, nil, resp)
}

func execRequest(t *testing.T, cl *http.Client, req *http.Request) *http.Response {
	resp, err := cl.Do(req)
	require.NoError(t, err)

	return resp
}

func verifyResponse(t *testing.T, expectedCode int, expectedData interface{}, resp *http.Response) {
	assert.Equal(t, expectedCode, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	var res contract.APIResponse
	err = json.Unmarshal(b, &res)
	require.NoError(t, err)

	assert.True(t, res.Success)
	assert.Equal(t, expectedData, res.Data)
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

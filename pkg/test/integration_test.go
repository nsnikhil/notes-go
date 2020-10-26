package test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	testCreateUser(t, cl)
}

func testPing(t *testing.T, cl *http.Client) {
	res, err := cl.Get(fmt.Sprintf("%s/%s", address, "ping"))
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var data contract.APIResponse
	err = json.Unmarshal(b, &data)
	require.NoError(t, err)

	assert.True(t, data.Success)
	assert.Equal(t, "pong", data.Data)
}

func testCreateUser(t *testing.T, cl *http.Client) {
	rd := contract.CreateUserRequest{Name: "test name", Email: "test@test.com", Password: "Password@1234"}

	rqB, err := json.Marshal(&rd)
	require.NoError(t, err)

	res, err := cl.Post(fmt.Sprintf("%s/%s", address, "user/create"), "application/json", bytes.NewBuffer(rqB))
	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	rpB, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var data contract.APIResponse
	err = json.Unmarshal(rpB, &data)
	require.NoError(t, err)

	assert.True(t, data.Success)
	assert.Equal(t, map[string]interface{}{"message": "user created successfully"}, data.Data)
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

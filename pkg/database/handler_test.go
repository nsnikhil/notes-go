package database_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"notes/pkg/config"
	"notes/pkg/database"
	"notes/pkg/liberr"
	"testing"
)

func TestGetDB(t *testing.T) {
	testCases := map[string]struct {
		dbConfig      config.DatabaseConfig
		expectedError error
	}{
		"test get db success": {
			dbConfig:      config.NewConfig("../../local.env").DatabaseConfig(),
			expectedError: nil,
		},
		"test get db failure invalid config": {
			dbConfig:      config.DatabaseConfig{},
			expectedError: liberr.WithArgs(liberr.Operation("Handler.GetDB.sql.Open"), liberr.SeverityError, errors.New("sql: unknown driver \"\" (forgotten import?)")),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testDBHandler(t, testCase.expectedError, testCase.dbConfig)
		})
	}
}

func testDBHandler(t *testing.T, expectedError error, cfg config.DatabaseConfig) {
	handler := database.NewHandler(cfg)
	_, err := handler.GetDB()

	if expectedError != nil {
		assert.Equal(t, expectedError, err)
	} else {
		assert.Nil(t, err)
	}
}

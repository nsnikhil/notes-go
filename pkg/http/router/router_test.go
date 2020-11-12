package router_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"notes/pkg/http/router"
	reporters "notes/pkg/reporting"
	"testing"
)

func TestRouter(t *testing.T) {
	r := router.NewRouter(
		zap.NewNop(),
		&reporters.MockPrometheus{},
	)

	rf := func(method, path string) *http.Request {
		req, err := http.NewRequest(method, path, nil)
		require.NoError(t, err)
		return req
	}

	testCases := map[string]struct {
		name    string
		request *http.Request
	}{
		"test ping route": {
			request: rf(http.MethodGet, "/ping"),
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r.ServeHTTP(w, testCase.request)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}

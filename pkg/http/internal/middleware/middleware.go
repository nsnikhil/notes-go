package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"notes/pkg/http/internal/resperr"
	"notes/pkg/http/internal/util"
	"notes/pkg/liberr"
	reporters "notes/pkg/reporting"
	"time"
)

func WithError(lgr *zap.Logger, handler func(resp http.ResponseWriter, req *http.Request) error) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		err := handler(resp, req)
		if err == nil {
			return
		}

		t, ok := err.(*liberr.Error)
		if ok {
			lgr.Error(t.EncodedStack())
		} else {
			lgr.Error(err.Error())
		}

		util.WriteFailureResponse(resperr.MapError(err), resp)
	}
}

//TODO: ADD MASKING BEFORE LOGGING REQ AND RESP
func WithReqRespLog(lgr *zap.Logger, handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		cr := util.NewCopyWriter(resp)

		handler(cr, req)

		//b, _ := cr.Body()

		//lgr.Sugar().Debug(req)
		//lgr.Sugar().Debug(string(b))
	}
}

func WithResponseHeaders(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		headers := map[string]string{
			"Content-Type":              "application/json",
			"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
			"X-Frame-Options":           "deny",
			"X-Content-Type-Options":    "nosniff",
		}

		for key, value := range headers {
			resp.Header().Set(key, value)
		}

		handler(resp, req)
	}
}

func WithRequestContext(handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// TODO: CHANGE TEMP VALUE
		ctx := context.WithValue(req.Context(), "key", "val")
		handler(resp, req.WithContext(ctx))
	}
}

func WithPrometheus(prometheus reporters.Prometheus, api string, handler http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// TODO CHANGE THIS
		hasError := func(code int) bool {
			return code >= 400 && code <= 600
		}

		start := time.Now()
		prometheus.ReportAttempt(api)

		cr := util.NewCopyWriter(resp)

		handler(cr, req)
		if hasError(cr.Code()) {
			duration := time.Since(start)
			prometheus.Observe(api, duration.Seconds())
			prometheus.ReportFailure(api)
			return
		}

		duration := time.Since(start)
		prometheus.Observe(api, duration.Seconds())

		prometheus.ReportSuccess(api)
	}
}

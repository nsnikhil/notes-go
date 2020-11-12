package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"notes/pkg/http/internal/handler"
	mdl "notes/pkg/http/internal/middleware"
	reporters "notes/pkg/reporting"
)

const ()

func NewRouter(lgr *zap.Logger, prometheus reporters.Prometheus) http.Handler {
	return getChiRouter(lgr, prometheus)
}

func getChiRouter(lgr *zap.Logger, pr reporters.Prometheus) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/ping", handler.PingHandler())
	r.Handle("/metrics", promhttp.Handler())

	return r
}

func withMiddlewares(lgr *zap.Logger, prometheus reporters.Prometheus, api string, handler func(resp http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return mdl.WithReqRespLog(lgr,
		mdl.WithResponseHeaders(
			mdl.WithPrometheus(prometheus, api, handler),
		),
	)
}

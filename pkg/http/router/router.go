package router

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"notes/pkg/http/internal/handler"
	mdl "notes/pkg/http/internal/middleware"
	reporters "notes/pkg/reporting"
	"notes/pkg/user"
)

const (
	pingAPI  = "ping"
	pingPath = "/ping"

	metricPath = "/metrics"

	userAPI  = "user"
	userPath = "/user"

	createAPI  = "create"
	createPath = "/create"
)

func NewRouter(lgr *zap.Logger, prometheus reporters.Prometheus, userService user.Service) http.Handler {
	return getChiRouter(lgr, prometheus, userService)
}

func getChiRouter(lgr *zap.Logger, pr reporters.Prometheus, userService user.Service) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get(pingPath, withMiddlewares(lgr, pr, pingAPI, handler.PingHandler()))
	r.Handle(metricPath, promhttp.Handler())

	registerUserRoutes(r, lgr, pr, userService)

	return r
}

func registerUserRoutes(r chi.Router, lgr *zap.Logger, prometheus reporters.Prometheus, userService user.Service) {
	apiFunc := func(api string) string {
		return fmt.Sprintf("%s_%s", userAPI, api)
	}

	uh := handler.NewUserHandler(userService)

	r.Route(userPath, func(r chi.Router) {
		r.Post(createPath, withMiddlewares(lgr, prometheus, apiFunc(createAPI), mdl.WithError(lgr, uh.CreateUser)))
	})
}

func withMiddlewares(lgr *zap.Logger, prometheus reporters.Prometheus, api string, handler func(resp http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return mdl.WithReqRespLog(lgr,
		mdl.WithResponseHeaders(
			mdl.WithPrometheus(prometheus, api, handler),
		),
	)
}

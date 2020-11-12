package app

import (
	"go.uber.org/zap"
	"io"
	"net/http"
	"notes/pkg/config"
	"notes/pkg/http/router"
	"notes/pkg/http/server"
	reporters "notes/pkg/reporting"
	"os"
)

func initHTTPServer(configFile string) server.Server {
	cfg, lgr, pr := initCommons(configFile)
	rt := initRouter(lgr, pr)
	return server.NewServer(cfg, lgr, rt)
}

func initCommons(configFile string) (config.Config, *zap.Logger, reporters.Prometheus) {
	cfg := config.NewConfig(configFile)

	lgr := initLogger(cfg)
	pr := reporters.NewPrometheus()

	return cfg, lgr, pr
}

func initRouter(lgr *zap.Logger, prometheus reporters.Prometheus) http.Handler {
	return router.NewRouter(lgr, prometheus)
}

func initLogger(cfg config.Config) *zap.Logger {
	return reporters.NewLogger(
		cfg.Env(),
		cfg.LogConfig().Level(),
		getWriters(cfg)...,
	)
}

func getWriters(cfg config.Config) []io.Writer {
	//TODO: MOVE TO CONST
	logSinkMap := map[string]io.Writer{
		"stdout": os.Stdout,
		"file":   reporters.NewExternalLogFile(cfg.LogFileConfig()),
	}

	var writers []io.Writer
	for _, sink := range cfg.LogConfig().Sinks() {
		w, ok := logSinkMap[sink]
		if ok {
			writers = append(writers, w)
		}
	}

	return writers
}

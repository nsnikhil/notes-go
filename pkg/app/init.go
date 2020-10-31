package app

import (
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"notes/pkg/config"
	"notes/pkg/database"
	"notes/pkg/http/router"
	"notes/pkg/http/server"
	"notes/pkg/notes/directory"
	reporters "notes/pkg/reporting"
	"notes/pkg/user"
	"os"
)

func initHTTPServer(configFile string) server.Server {
	cfg, lgr, pr, svc := initCommons(configFile)
	rt := initRouter(lgr, pr, svc)
	return server.NewServer(cfg, lgr, rt)
}

func initCommons(configFile string) (config.Config, *zap.Logger, reporters.Prometheus, user.Service) {
	cfg := config.NewConfig(configFile)

	lgr := initLogger(cfg)
	pr := reporters.NewPrometheus()

	svc := initServices(cfg)

	return cfg, lgr, pr, svc
}

func initRouter(lgr *zap.Logger, prometheus reporters.Prometheus, svc user.Service) http.Handler {
	return router.NewRouter(lgr, prometheus, svc)
}

func initServices(cfg config.Config) user.Service {
	db, err := database.NewHandler(cfg.DatabaseConfig()).GetDB()
	if err != nil {
		log.Fatal(err)
	}

	uc := cfg.UserConfig()

	dsv := directory.NewDirectoryService(db)
	usv, err := user.NewService(dsv, db, uc.SaltLength(), uc.Iterations(), uc.KeyLength(), uc.PemString())
	if err != nil {
		log.Fatal(err)
	}

	return usv
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

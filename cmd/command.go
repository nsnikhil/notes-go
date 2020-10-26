package main

import (
	"log"
	"notes/pkg/app"
)

const (
	grpcServeCommand = "grpc-serve"
	httpServeCommand = "http-serve"
	migrateCommand   = "migrate"
	rollbackCommand  = "rollback"
)

func commands() map[string]func(configFile string) {
	return map[string]func(configFile string){
		//grpcServeCommand: app.StartGRPCServer,
		httpServeCommand: app.StartHTTPServer,
		//migrateCommand:   database.RunMigrations,
		//rollbackCommand:  database.RollBackMigrations,
	}
}

func execute(cmd string, configFile string) {
	run, ok := commands()[cmd]
	if !ok {
		log.Fatal("invalid command")
	}

	run(configFile)
}

package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"notes/pkg/config"
	"notes/pkg/liberr"
	"time"
)

type Handler interface {
	GetDB() (*sql.DB, error)
}

type sqlDBHandler struct {
	cfg config.DatabaseConfig
}

func (dbh *sqlDBHandler) GetDB() (*sql.DB, error) {
	db, err := sql.Open(dbh.cfg.DriverName(), dbh.cfg.Source())
	if err != nil {
		fmt.Println(err)
		return nil, liberr.WithArgs(liberr.Operation("Handler.GetDB.sql.Open"), liberr.SeverityError, err)
	}

	db.SetMaxOpenConns(dbh.cfg.MaxOpenConnections())
	db.SetMaxIdleConns(dbh.cfg.IdleConnections())
	db.SetConnMaxLifetime(time.Minute * time.Duration(dbh.cfg.ConnectionMaxLifetime()))

	if err := db.Ping(); err != nil {
		return nil, liberr.WithArgs(liberr.Operation("Handler.GetDB.db.Ping"), liberr.SeverityError, err)
	}

	return db, nil
}

func NewHandler(cfg config.DatabaseConfig) Handler {
	return &sqlDBHandler{
		cfg: cfg,
	}
}

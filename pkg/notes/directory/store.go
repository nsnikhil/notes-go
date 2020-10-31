package directory

import (
	"database/sql"
	"notes/pkg/liberr"
)

const (
	insertDirectory = `INSERT INTO directories (name) VALUES ($1) RETURNING id`
)

type store interface {
	createDirectory(directory *directory) (string, error)
}

type directoryStore struct {
	db *sql.DB
}

func (ds *directoryStore) createDirectory(directory *directory) (string, error) {
	var id string

	err := ds.db.QueryRow(insertDirectory, directory.name).Scan(&id)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("DirectoryStore.createDirectory"), liberr.InternalError, liberr.SeverityError, err)
	}

	return id, nil
}

func newDirectoryStore(db *sql.DB) store {
	return &directoryStore{
		db: db,
	}
}

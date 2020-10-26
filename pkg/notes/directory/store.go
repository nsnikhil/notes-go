package directory

import (
	"database/sql"
	"notes/pkg/liberr"
)

const (
	insertDirectory = `INSERT INTO directories (name) VALUES ($1) RETURNING id`
)

type Store interface {
	CreateDirectory(directory *Directory) (string, error)
}

type directoryStore struct {
	db *sql.DB
}

func (ds *directoryStore) CreateDirectory(directory *Directory) (string, error) {
	var id string

	err := ds.db.QueryRow(insertDirectory, directory.Name()).Scan(&id)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("DirectoryStore.CreateDirectory"), liberr.InternalError, liberr.SeverityError, err)
	}

	return id, nil
}

func NewDirectoryStore(db *sql.DB) Store {
	return &directoryStore{
		db: db,
	}
}

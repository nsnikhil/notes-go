package directory

import (
	"database/sql"
	"notes/pkg/liberr"
)

type Service interface {
	CreateDirectory(name string) (string, error)
}

type directoryService struct {
	st store
}

func (ds *directoryService) CreateDirectory(name string) (string, error) {
	dir, err := newDirectory(name)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("DirectoryService.createDirectory"), err)
	}

	id, err := ds.st.createDirectory(dir)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("DirectoryService.createDirectory"), err)
	}

	return id, nil
}

func newDirectoryService(st store) Service {
	return &directoryService{
		st: st,
	}
}

func NewDirectoryService(db *sql.DB) Service {
	return newDirectoryService(newDirectoryStore(db))
}

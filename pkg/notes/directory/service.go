package directory

import (
	"notes/pkg/liberr"
)

type Service interface {
	CreateDirectory(name string) (string, error)
}

type directoryService struct {
	st Store
}

func (ds *directoryService) CreateDirectory(name string) (string, error) {
	dir, err := NewDirectory(name)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("DirectoryService.CreateDirectory"), err)
	}

	id, err := ds.st.CreateDirectory(dir)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("DirectoryService.CreateDirectory"), err)
	}

	return id, nil
}

func NewDirectoryService(st Store) Service {
	return &directoryService{
		st: st,
	}
}

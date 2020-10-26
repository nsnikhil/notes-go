package directory

import (
	"errors"
	"notes/pkg/liberr"
	"notes/pkg/notes/document"
	"time"
)

type Directory struct {
	id   string
	name string

	parent      *Directory
	directories []*Directory
	documents   []*document.Document

	createdAt time.Time
	updatedAt time.Time
}

func (d *Directory) Name() string {
	return d.name
}

func NewDirectory(name string) (*Directory, error) {
	if err := validate(name); err != nil {
		return nil, err
	}

	return &Directory{name: name}, nil
}

func validate(name string) error {
	if len(name) == 0 {
		return liberr.WithArgs(liberr.Operation("Directory.validate"), liberr.ValidationError, liberr.SeverityError, errors.New("directory name cannot be empty"))
	}

	return nil
}

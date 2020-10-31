package directory

import (
	"errors"
	"notes/pkg/liberr"
	"notes/pkg/notes/document"
	"time"
)

type directory struct {
	id   string
	name string

	parent      *directory
	directories []*directory
	documents   []*document.Document

	createdAt time.Time
	updatedAt time.Time
}

func newDirectory(name string) (*directory, error) {
	if err := validate(name); err != nil {
		return nil, err
	}

	return &directory{name: name}, nil
}

func validate(name string) error {
	if len(name) == 0 {
		return liberr.WithArgs(liberr.Operation("directory.validate"), liberr.ValidationError, liberr.SeverityError, errors.New("directory name cannot be empty"))
	}

	return nil
}

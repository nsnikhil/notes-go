package document

import (
	"database/sql"
	"notes/pkg/liberr"
)

const (
	insertDocument = `INSERT INTO documents (name, content) VALUES ($1, $2) RETURNING id`
)

type Store interface {
	CreateDocument(document *Document) (string, error)
}

//
type documentStore struct {
	db *sql.DB
}

func (ds *documentStore) CreateDocument(document *Document) (string, error) {
	var id string

	err := ds.db.QueryRow(insertDocument, document.name, document.content).Scan(&id)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("DocumentStore.CreateDocument"), liberr.InternalError, liberr.SeverityError, err)
	}

	return id, nil
}

func NewDocumentStore(db *sql.DB) Store {
	return &documentStore{
		db: db,
	}
}

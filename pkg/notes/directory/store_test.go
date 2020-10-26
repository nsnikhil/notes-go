package directory_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"notes/pkg/notes/directory"
	"regexp"
	"testing"
)

func TestDirectoryStoreCreateDirectorySuccess(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	query := `INSERT INTO directories (name) VALUES ($1) RETURNING id`
	dirName := "root"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(dirName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(dirID))

	ds := directory.NewDirectoryStore(db)

	_, err := ds.CreateDirectory(getDirectory(t, dirName))
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDirectoryStoreCreateDirectoryFailure(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	query := `INSERT INTO directories (name) VALUES ($1) RETURNING id`
	dirName := "root"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(dirName).
		WillReturnError(errors.New("failed to create new directory"))

	ds := directory.NewDirectoryStore(db)

	_, err := ds.CreateDirectory(getDirectory(t, dirName))
	require.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func getDirectory(t *testing.T, name string) *directory.Directory {
	dir, err := directory.NewDirectory(name)
	require.NoError(t, err)

	return dir
}

func getMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock
}

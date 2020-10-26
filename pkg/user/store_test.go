package user_test

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"notes/pkg/user"
	"regexp"
	"testing"
)

func TestCreateNewUserSuccessful(t *testing.T) {
	db, mock := getMockDB(t)

	query := `INSERT INTO USERS (name, email, password, rootdirid) VALUES ($1, $2, $3, $4) RETURNING id`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(name, email, password, dirID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	us := user.NewUserStore(db)

	_, err := us.CreateUser(getUser(t, name, email, password, dirID))
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNewUserFailure(t *testing.T) {
	db, mock := getMockDB(t)

	query := `INSERT INTO USERS (name, email, password, rootdirid) VALUES ($1, $2, $3, $4) RETURNING id`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(name, email, password, dirID).
		WillReturnError(errors.New("failed to create new user"))

	us := user.NewUserStore(db)

	_, err := us.CreateUser(getUser(t, name, email, password, dirID))
	require.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func getUser(t *testing.T, name, email, password, dirID string) *user.User {
	us, err := user.NewBuilder().Name(name).Email(email).Password(password).RootDirID(dirID).Build()
	require.NoError(t, err)

	return us
}

func getMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock
}

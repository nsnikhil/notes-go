package user

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestCreateNewUserSuccessful(t *testing.T) {
	db, mock := getMockDB(t)

	query := `INSERT INTO USERS (name, email, passwordhash, passwordsalt, rootdirid) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(name, email, passwordHash, []byte{}, dirID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	store := newUserStore(db)

	currUser := &user{name: name,
		email:        email,
		rootDirID:    dirID,
		passwordHash: passwordHash,
		passwordSalt: []byte{},
	}

	_, err := store.createUser(currUser)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNewUserFailure(t *testing.T) {
	db, mock := getMockDB(t)

	query := `INSERT INTO USERS (name, email, passwordhash, passwordsalt, rootdirid) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(name, email, passwordHash, []byte{}, dirID).
		WillReturnError(errors.New("failed to create new user"))

	store := newUserStore(db)

	currUser := &user{name: name,
		email:        email,
		rootDirID:    dirID,
		passwordHash: passwordHash,
		passwordSalt: []byte{},
	}

	_, err := store.createUser(currUser)
	require.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserForEmailSuccess(t *testing.T) {
	db, mock := getMockDB(t)

	query := `SELECT id, name, email, passwordhash, passwordsalt, rootdirid FROM users WHERE email = $1`

	rows := sqlmock.NewRows(
		[]string{
			"id", "name", "email", "passwordhash", "passwordsalt", "rootdirid",
		},
	)

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(email).
		WillReturnRows(rows.AddRow("", "", "", "", "", ""))

	us := newUserStore(db)

	_, err := us.getByEmail(email)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserForEmailFailure(t *testing.T) {
	db, mock := getMockDB(t)

	query := `SELECT id, name, email, passwordhash, passwordsalt, rootdirid FROM users WHERE email = $1`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(email).
		WillReturnError(errors.New("failed to get data"))

	us := newUserStore(db)

	_, err := us.getByEmail(email)
	require.Error(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}

func getUser(t *testing.T, name, email, password, dirID string) *user {
	ph := newPasswordHasher(saltLength, iterations, keyLength)
	us, err := newBuilder(ph).name(name).email(email).password(password).rootDirID(dirID).build()
	require.NoError(t, err)

	return us
}

func getMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return db, mock
}

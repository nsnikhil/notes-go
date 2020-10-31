package user

import (
	"database/sql"
	"fmt"
	"notes/pkg/liberr"
)

const (
	insertUser     = `INSERT INTO USERS (name, email, passwordhash, passwordsalt, rootdirid) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	getUserByEmail = `SELECT id, name, email, passwordhash, passwordsalt, rootdirid FROM users WHERE email = $1`
)

type store interface {
	createUser(user *user) (string, error)
	getByEmail(email string) (*user, error)
}

// TODO: RENAME
type userStore struct {
	db *sql.DB
}

func (us *userStore) createUser(user *user) (string, error) {
	var id string

	err := us.db.QueryRow(insertUser, user.name, user.email, user.passwordHash, user.passwordSalt, user.rootDirID).Scan(&id)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("UserStore.createUser"), liberr.InternalError, liberr.SeverityError, err)
	}

	return id, nil
}

func (us *userStore) getByEmail(email string) (*user, error) {

	row := us.db.QueryRow(getUserByEmail, email)
	if row.Err() != nil {
		return nil, liberr.WithArgs(liberr.Operation("UserStore.getByEmail"), liberr.InternalError, liberr.SeverityError, row.Err())
	}

	var user user
	err := row.Scan(&user.id, &user.name, &user.email, &user.passwordHash, &user.passwordSalt, &user.rootDirID)
	if err != nil {
		fmt.Print(err)
		return nil, liberr.WithArgs(liberr.Operation("UserStore.getByEmail"), liberr.InternalError, liberr.SeverityError, row.Err())
	}

	return &user, nil
}

func newUserStore(db *sql.DB) store {
	return &userStore{
		db: db,
	}
}

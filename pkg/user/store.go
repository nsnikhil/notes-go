package user

import (
	"database/sql"
	"notes/pkg/liberr"
)

const (
	insertUser = `INSERT INTO USERS (name, email, password, rootdirid) VALUES ($1, $2, $3, $4) RETURNING id`
)

type Store interface {
	CreateUser(user *User) (string, error)
}

// TODO: RENAME
type userStore struct {
	db *sql.DB
}

func (us *userStore) CreateUser(user *User) (string, error) {
	var id string

	err := us.db.QueryRow(insertUser, user.name, user.email, user.password, user.rootDirID).Scan(&id)
	if err != nil {
		return "", liberr.WithArgs(liberr.Operation("UserStore.CreateUser"), liberr.InternalError, liberr.SeverityError, err)
	}

	return id, nil
}

func NewUserStore(db *sql.DB) Store {
	return &userStore{
		db: db,
	}
}

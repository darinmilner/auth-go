package repository

import (
	"database/sql"
	"errors"
)

var (
	// ErrNoRecord no record found in database error
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials invalid username/password error
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail duplicate email error
	ErrDuplicateEmail = errors.New("models: duplicate email")
	// ErrInactiveAccount inactive account error
	ErrInactiveAccount = errors.New("models: Inactive Account")
)

type DBRepo struct {
	DB *sql.DB
}

type TestDBRepo struct {
	DB *sql.DB
}

//Repo is the wrapper for the db
type Repo struct {
	DB DBRepo
}

//NewRepo returns models with db pool
func NewRepo(db *sql.DB) Repo {
	return Repo{
		DB: DBRepo{DB: db},
	}
}

//NewTestRepo returns test data with db pool
func NewTestRepo(db *sql.DB) Repo {
	return Repo{
		DB: DBRepo{DB: db},
	}
}

package postgres

import (
	"auth-rbac/internal/user"
	"database/sql"

	"github.com/pkg/errors"
)

var _ user.Users = &UserStorage{}

// UserStorage ...
type UserStorage struct {
	statementStorage

	createStmt      *sql.Stmt
	findByIDStmt    *sql.Stmt
	findByEmailStmt *sql.Stmt
}

// NewUserStorage ...
func NewUserStorage(db *DB) (*UserStorage, error) {
	s := &UserStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: createUserQuery, Dst: &s.createStmt},
		{Query: findUserByIDQuery, Dst: &s.findByIDStmt},
		{Query: findUserByLoginQuery, Dst: &s.findByEmailStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can not init statements")
	}

	return s, nil
}

const userFields = "login, password, role, created_at"

const createUserQuery = "INSERT INTO public.users (" + userFields + ") VALUES ($1, $2, $3, now()) RETURNING id"

func (s *UserStorage) Create(u *user.User) error {
	if _, err := s.createStmt.Exec(&u.Login, &u.Password, &u.Role); err != nil {
		return err
	}

	return nil
}

const findUserByIDQuery = "SELECT id, " + userFields + " FROM public.users WHERE id=$1"

func (s *UserStorage) FindUserByID(id int64) (*user.User, error) {
	var u user.User

	row := s.findByIDStmt.QueryRow(id)
	if err := scanUser(row, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

const findUserByLoginQuery = "SELECT id, " + userFields + " FROM public.users WHERE login=$1"

func (s *UserStorage) FindUserByLogin(login string) (*user.User, error) {
	var u user.User

	row := s.findByEmailStmt.QueryRow(login)
	if err := scanUser(row, &u); err != nil {
		return nil, err
	}

	return &u, nil
}

func scanUser(scanner sqlScanner, u *user.User) error {
	return scanner.Scan(&u.ID, &u.Login, &u.Password, &u.Role, &u.CreatedAt)
}

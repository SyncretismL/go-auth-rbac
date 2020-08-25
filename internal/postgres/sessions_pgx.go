package postgres

import (
	"auth-rbac/internal/session"
	"database/sql"

	"github.com/pkg/errors"
)

var _ session.Sessions = &SessionStorage{}

type SessionStorage struct {
	statementStorage

	upsertStmt      *sql.Stmt
	findByIDStmt    *sql.Stmt
	findByTokenStmt *sql.Stmt
}

func NewSessionStorage(db *DB) (*SessionStorage, error) {
	s := &SessionStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: upsertSessionQuery, Dst: &s.upsertStmt},
		{Query: findSessionByIDQuery, Dst: &s.findByIDStmt},
		{Query: findSessionByTokenQuery, Dst: &s.findByTokenStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can not init statements")
	}

	return s, nil
}

const sessionFields = "user_id, token, created_at, valid_until"

const upsertSessionQuery = "INSERT INTO public.sessions(" + sessionFields + ") VALUES ($1, $2, $3, $4)" +
	" ON CONFLICT (user_id) DO UPDATE SET token = excluded.token, created_at = excluded.created_at, valid_until = excluded.valid_until"

func (s *SessionStorage) Upsert(ses *session.Session) error {
	_, err := s.upsertStmt.Exec(&ses.UserID, &ses.Token, &ses.CreatedAt, &ses.ValidUntil)
	if err != nil {
		return err
	}

	return nil
}

const findSessionByIDQuery = "SELECT " + sessionFields + " FROM public.sessions WHERE user_id=$1"

func (s *SessionStorage) FindByUserID(id int64) (*session.Session, error) {
	var session session.Session

	row := s.findByIDStmt.QueryRow(id)
	if err := scanSession(row, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

const findSessionByTokenQuery = "SELECT " + sessionFields + " FROM public.sessions WHERE token=$1"

func (s *SessionStorage) FindByToken(token string) (*session.Session, error) {
	var session session.Session

	row := s.findByTokenStmt.QueryRow(token)
	if err := scanSession(row, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func scanSession(scanner sqlScanner, s *session.Session) error {
	return scanner.Scan(&s.UserID, &s.Token, &s.CreatedAt, &s.ValidUntil)
}

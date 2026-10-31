package store

import (
	"database/sql"
	"net/http"
)

const usersTable = "users"

type Store struct{ db *sql.DB }

func (s *Store) FindByName(r *http.Request) (*sql.Rows, error) {
	name := r.URL.Query().Get("name")
	return s.db.Query("SELECT id, email FROM users WHERE name = '" + name + "'")
}

func (s *Store) FindByEmail(r *http.Request) (*sql.Rows, error) {
	email := r.URL.Query().Get("email")
	return s.db.Query("SELECT id, name FROM "+usersTable+" WHERE email = $1", email)
}

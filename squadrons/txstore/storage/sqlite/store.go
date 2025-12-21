package sqlite

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Store struct {
	Path string
	DB   *sql.DB
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &Store{Path: path, DB: db}, nil
}

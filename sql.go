package sqlbench

import (
	"database/sql"
	"sync"

	// Adding postgres driver
	_ "github.com/lib/pq"
)

type sqlRunner struct {
	dsn  string
	db   *sql.DB
	once sync.Once
}

func (s *sqlRunner) init() error {
	var err error
	s.db, err = sql.Open("postgres", s.dsn)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqlRunner) run(q string) error {
	rows, err := s.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
	}
	return nil
}

func (s *sqlRunner) tag(q string) (string, error) {
	var value string
	err := s.db.QueryRow(q).Scan(&value)

	if err == sql.ErrNoRows || err != nil {
		return "", err
	}
	return value, nil
}

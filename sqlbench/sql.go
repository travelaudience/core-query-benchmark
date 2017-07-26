package sqlbench

import (
	"database/sql"
	"sync"

	"log"

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
		log.Println(err)
		return err
	}
	return nil
}

func (s *sqlRunner) run(q string) error {
	_, err := s.db.Query(q)
	return err
}

func (s *sqlRunner) tag(q string) (string, error) {
	var value string
	err := s.db.QueryRow(q).Scan(&value)

	if err == sql.ErrNoRows || err != nil {
		log.Println("No results or err", err, s.dsn)
		return "", err
	}
	return value, nil
}

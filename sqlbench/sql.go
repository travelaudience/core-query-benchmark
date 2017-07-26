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

func (s *sqlRunner) run(q string) error {
	var err error
	s.once.Do(func() {
		s.db, err = sql.Open("postgres", s.dsn)
	})
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = s.db.Query(q)
	return err
}

func (s *sqlRunner) tag(q string) (string, error) {
	var err error
	s.once.Do(func() {
		s.db, err = sql.Open("postgres", s.dsn)
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	id := 123
	var value string
	err = s.db.QueryRow(q, id).Scan(&value)

	if err == sql.ErrNoRows || err != nil {
		log.Printf("No results or err", err)
		return "", err
	}
	return value, nil
}

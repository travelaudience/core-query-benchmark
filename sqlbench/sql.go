package sqlbench

import (
	"database/sql"
	"sync"

	"log"

	// Adding postgres driver
	_ "github.com/lib/pq"
)

type sqlRunner struct {
	db   *sql.DB
	once sync.Once
}

func (s *sqlRunner) run(dsn string, q string) error {
	var err error
	s.once.Do(func() {
		println(dsn)
		s.db, err = sql.Open("postgres", dsn)
	})
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = s.db.Query(q)
	return err
}

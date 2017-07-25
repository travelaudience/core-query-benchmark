package sqlbench

import "database/sql"
import "sync"

type sqlRunner struct {
	db   *sql.DB
	once sync.Once
}

func (s *sqlRunner) run(dsn string, q string) error {
	var err error
	s.once.Do(func() {
		s.db, err = sql.Open("postgres", dsn)
	})
	if err != nil {
		return err
	}
	_, err = s.db.Query(q)
	return err
}

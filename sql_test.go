package sqlbench

import (
	"database/sql"
	"sync"
	"time"
)

type sqlRunnerTest struct {
	dsn  string
	db   *sql.DB
	once sync.Once
}

func (s *sqlRunnerTest) init() error {
	return nil
}

func (s *sqlRunnerTest) run(q string) error {
	time.Sleep(time.Millisecond * 10)
	return nil
}

func (s *sqlRunnerTest) tag(q string) (string, error) {
	return "mytag", nil
}

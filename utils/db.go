package utils

import
(
	"database/sql"
 	_ "github.com/lib/pq"
	"time"
	"errors"
)

func DbConnect(connString string) (*sql.DB, error) {
	count := 0
	ticker := time.NewTicker(1*time.Second)
	for range ticker.C {
		db, err := sql.Open("postgres", connString)
		if err == nil {
			return db, nil
		}
		count++
		if count > 10 {
			return nil, err
		}
	}
	return nil, errors.New("impossible error")
}


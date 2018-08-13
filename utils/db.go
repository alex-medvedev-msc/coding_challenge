package utils

import
(
	"database/sql"
	// db drivers in go are always imported like that
 	_ "github.com/lib/pq"
	"time"
	"errors"
)

// DbConnect waits for db to be ready up to 10 seconds and returns valid connection or error
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
			break
		}
	}
	return nil, errors.New("cannot connect to db, 10 attempts failed")
}


package utils

import
(
	"database/sql"
 	_ "github.com/lib/pq"
)

func DbConnect(connString string) (*sql.DB, error) {
	return sql.Open("postgres", connString)
}

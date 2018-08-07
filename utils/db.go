package utils

import "database/sql"

func DbConnect(connString string) (*sql.DB, error) {
	return sql.Open("postgres", connString)
}

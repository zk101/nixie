package storage

import "database/sql"

// SQL provides an interface for SQL Server
type SQL interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// EOF

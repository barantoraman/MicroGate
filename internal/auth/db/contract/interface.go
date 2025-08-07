package contract

import "database/sql"

type DBConnection interface {
	DB() *sql.DB
	Close()
}

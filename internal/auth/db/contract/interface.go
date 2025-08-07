package contract

import "database/sql"

type DBConnection interface {
	Close()
	DB() *sql.DB
}

package mock

import (
	"database/sql"

	dbContract "github.com/barantoraman/microgate/internal/auth/db/contract"
)

type mockConnection struct{}

func NewMockConnection() (dbContract.DBConnection, error) {
	return &mockConnection{}, nil
}

func (m *mockConnection) DB() *sql.DB { return nil }

func (m *mockConnection) Close() {}

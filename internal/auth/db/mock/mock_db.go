package mock

import (
	"database/sql"

	dbContract "github.com/barantoraman/microgate/internal/auth/db/contract"
	"github.com/barantoraman/microgate/pkg/config"
)

type mockConnection struct{}

func NewMockConnection(cfg config.AuthServiceConfigurations) (dbContract.DBConnection, error) {
	return &mockConnection{}, nil
}

func (m *mockConnection) DB() *sql.DB { return nil }

func (m *mockConnection) Close() {}

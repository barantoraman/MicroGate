package db

import (
	"errors"

	dbContract "github.com/barantoraman/microgate/internal/auth/db/contract"
	"github.com/barantoraman/microgate/internal/auth/db/mock"
	"github.com/barantoraman/microgate/internal/auth/db/postgres"
	"github.com/barantoraman/microgate/pkg/config"
)

const (
	PqType       = "postgres"
	InMemoryType = "inmemory"
)

func GetDatabase(cfg config.AuthServiceConfigurations) (dbContract.DBConnection, error) {
	switch cfg.DBType {
	case PqType:
		return postgres.NewPostgresConnection(cfg)
	case InMemoryType:
		return mock.NewMockConnection()
	default:
		return nil, errors.New("invalid database type")
	}
}

package db

import (
	"errors"

	dbContract "github.com/barantoraman/microgate/internal/task/db/contract"
	"github.com/barantoraman/microgate/internal/task/db/mock"
	"github.com/barantoraman/microgate/internal/task/db/mongo"
	"github.com/barantoraman/microgate/pkg/config"
)

const (
	MongoType    = "mongo"
	InMemoryType = "inmemory"
)

func GetDatabase(cfg config.TaskServiceConfigurations) (dbContract.DBConnection, error) {
	switch cfg.DBType {
	case MongoType:
		return mongo.NewMongoConnection(cfg)
	case InMemoryType:
		return mock.NewMockConnection()
	default:
		return nil, errors.New("invalid database type")
	}
}

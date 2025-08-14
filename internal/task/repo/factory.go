package repo

import (
	"github.com/barantoraman/microgate/internal/task/db"
	dbContract "github.com/barantoraman/microgate/internal/task/db/contract"
	repoContract "github.com/barantoraman/microgate/internal/task/repo/contract"
	"github.com/barantoraman/microgate/internal/task/repo/mock"
	"github.com/barantoraman/microgate/internal/task/repo/task"
	"github.com/barantoraman/microgate/pkg/config"
)

func GetRepository(cfg config.TaskServiceConfigurations, conn dbContract.DBConnection) repoContract.TaskRepository {
	switch cfg.DBType {
	case db.MongoType:
		return task.NewTaskRepository(conn)
	case db.InMemoryType:
		return mock.NewMockTaskRepository()
	default:
		return nil
	}
}

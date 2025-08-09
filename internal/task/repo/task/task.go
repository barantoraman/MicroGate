package task

import (
	"context"
	"errors"

	dbContract "github.com/barantoraman/microgate/internal/task/db/contract"
	repoContract "github.com/barantoraman/microgate/internal/task/repo/contract"
	"github.com/barantoraman/microgate/internal/task/repo/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrRecordNotFound = errors.New("record not found")

type taskRepository struct {
	collection *mongo.Collection
	db         *mongo.Client
}

func NewTaskRepository(conn dbContract.DBConnection) repoContract.TaskRepository {
	coll := conn.DB().Database("microgate").Collection("task")
	return &taskRepository{
		collection: coll,
		db:         conn.DB(),
	}
}

func (t *taskRepository) CreateTask(ctx context.Context, task *entity.Task) error {
	panic("unimplemented")
}

func (t *taskRepository) DeleteTask(ctx context.Context, taskID string, userID int64) error {
	panic("unimplemented")
}

func (t *taskRepository) ListTask(ctx context.Context, userID int64) ([]entity.Task, error) {
	panic("unimplemented")
}

func (t *taskRepository) ServiceStatus(ctx context.Context) error {
	panic("unimplemented")
}

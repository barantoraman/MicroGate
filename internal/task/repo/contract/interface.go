package contract

import (
	"context"

	"github.com/barantoraman/microgate/internal/task/repo/entity"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *entity.Task) error
	ListTask(ctx context.Context, userID int64) ([]entity.Task, error)
	DeleteTask(ctx context.Context, taskID string, userID int64) error
	ServiceStatus(ctx context.Context) error
}

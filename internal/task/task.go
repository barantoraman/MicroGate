package task

import (
	"context"

	repoContract "github.com/barantoraman/microgate/internal/task/repo/contract"
	"github.com/barantoraman/microgate/internal/task/repo/entity"
	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
)

type taskService struct {
	taskRepository repoContract.TaskRepository
}

func NewService(taskRepo repoContract.TaskRepository, logger loggerContract.Logger) Service {
	return &taskService{
		taskRepository: taskRepo,
	}
}

func (t *taskService) CreateTask(ctx context.Context, task entity.Task) (string, error) {
}

func (t *taskService) ListTask(ctx context.Context, userID int64) ([]entity.Task, error) {
	panic("unimplemented")
}

func (t *taskService) DeleteTask(ctx context.Context, taskID string, userID int64) error {
	panic("unimplemented")
}

func (t *taskService) ServiceStatus(ctx context.Context) (int, error) {
	panic("unimplemented")
}

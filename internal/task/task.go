package task

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	repoContract "github.com/barantoraman/microgate/internal/task/repo/contract"
	"github.com/barantoraman/microgate/internal/task/repo/entity"
	taskPkg "github.com/barantoraman/microgate/internal/task/repo/task"
	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
	"github.com/barantoraman/microgate/pkg/validator"
)

type taskService struct {
	taskRepository repoContract.TaskRepository
	logger         loggerContract.Logger
}

func NewService(taskRepo repoContract.TaskRepository, logger loggerContract.Logger) Service {
	return &taskService{
		taskRepository: taskRepo,
		logger:         logger,
	}
}

func (t *taskService) CreateTask(ctx context.Context, task entity.Task) (string, error) {
	v := validator.New()
	if taskPkg.ValidateTask(v, task); !v.Valid() {
		t.logger.Error("validation error")
		return "", fmt.Errorf("validation error")
	}
	if err := t.taskRepository.CreateTask(ctx, &task); err != nil {
		t.logger.Error("failed to create task")
		return "", errors.New("failed to create task")
	}
	return task.Id.Hex(), nil
}

func (t *taskService) ListTask(ctx context.Context, userID int64) ([]entity.Task, error) {
	tasks, err := t.taskRepository.ListTask(ctx, userID)
	if err != nil {
		t.logger.Error("failed to get tasks")
		return nil, errors.New("failed to get tasks")
	}
	return tasks, nil
}

func (t *taskService) DeleteTask(ctx context.Context, taskID string, userID int64) error {
	if err := t.taskRepository.DeleteTask(ctx, taskID, userID); err != nil {
		t.logger.Error("failed to delete task")
		return errors.New("failed to delete task")
	}
	return nil
}

func (t *taskService) ServiceStatus(ctx context.Context) (int, error) {
	if err := t.taskRepository.ServiceStatus(ctx); err != nil {
		t.logger.Error("task service status error")
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

package mock

import (
	"context"
	"errors"
	"sync"

	repoContract "github.com/barantoraman/microgate/internal/task/repo/contract"
	"github.com/barantoraman/microgate/internal/task/repo/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockTaskRepository struct {
	mu    sync.RWMutex
	tasks []entity.Task
}

func NewTaskRepository() repoContract.TaskRepository {
	return &mockTaskRepository{
		mu:    sync.RWMutex{},
		tasks: make([]entity.Task, 0),
	}
}

func (m *mockTaskRepository) CreateTask(ctx context.Context, task *entity.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	task.Id = primitive.NewObjectID()
	m.tasks = append(m.tasks, *task)
	return nil
}

func (m *mockTaskRepository) DeleteTask(ctx context.Context, taskID string, userID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	id, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return errors.New("invalid ObjectID")
	}

	for i, t := range m.tasks {
		if t.Id == id && t.UserID == userID {
			m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("record not found")
}

func (m *mockTaskRepository) ListTask(ctx context.Context, userID int64) ([]entity.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var filtered []entity.Task
	for _, t := range m.tasks {
		if t.UserID == userID {
			filtered = append(filtered, t)
		}
	}
	return filtered, nil
}

func (m *mockTaskRepository) ServiceStatus(ctx context.Context) error {
	return nil
}

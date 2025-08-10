package mock

import (
	"context"
	"testing"

	"github.com/barantoraman/microgate/internal/task/repo/entity"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMockTaskRepository(t *testing.T) {
	repo := NewMockTaskRepository()
	ctx := context.Background()

	t.Run("CreateTask - Success", func(t *testing.T) {
		task := entity.Task{
			UserID:      123,
			Title:       "Title",
			Description: "Description",
			Status:      "Status",
		}
		err := repo.CreateTask(ctx, &task)
		require.NoError(t, err)
		require.NotZero(t, task.Id, "Task ID should be set after creation")
	})

	t.Run("ListTask - Returns correct tasks", func(t *testing.T) {
		tasks, err := repo.ListTask(ctx, 123)
		require.NoError(t, err)
		require.Len(t, tasks, 1)
		require.Equal(t, "Title", tasks[0].Title)
		require.Equal(t, int64(123), tasks[0].UserID)
	})

	t.Run("DeleteTask - Success", func(t *testing.T) {
		tasks, _ := repo.ListTask(ctx, 123)
		err := repo.DeleteTask(ctx, tasks[0].Id.Hex(), 123)
		require.NoError(t, err)

		// Deleted task should not be in the list
		tasksAfter, _ := repo.ListTask(ctx, 123)
		require.Len(t, tasksAfter, 0)
	})

	t.Run("DeleteTask - Invalid ObjectID", func(t *testing.T) {
		err := repo.DeleteTask(ctx, "invalid", 123)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid ObjectID")
	})

	t.Run("DeleteTask - Record not found", func(t *testing.T) {
		err := repo.DeleteTask(ctx, primitive.NewObjectID().Hex(), 999) // wrong userID
		require.Error(t, err)
		require.Contains(t, err.Error(), "record not found")
	})

	t.Run("ServiceStatus - Always OK", func(t *testing.T) {
		err := repo.ServiceStatus(ctx)
		require.NoError(t, err)
	})
}

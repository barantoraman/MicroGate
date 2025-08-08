package mock

import (
	"context"
	"testing"

	"github.com/barantoraman/microgate/internal/auth/repo/entity"
	"github.com/stretchr/testify/require"
)

func TestMockRepository_CreateUser(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	t.Run("creates a new user successfully", func(t *testing.T) {
		user := &entity.User{
			Email:        "peter@parker.com",
			Password:     "spiderman",
			PasswordHash: []byte{10, 20, 30},
		}

		err := repo.CreateUser(ctx, user)
		require.NoError(t, err)
		require.Equal(t, int64(1), user.UserID)
	})

	t.Run("fails to create user with duplicate email", func(t *testing.T) {
		user := &entity.User{
			Email:        "peter@parker.com",
			Password:     "spiderman",
			PasswordHash: []byte{1},
		}

		err := repo.CreateUser(ctx, user)
		require.ErrorIs(t, err, entity.ErrDuplicateEmail)
	})
}

func TestMockRepository_GetUser(t *testing.T) {
	repo := NewMockRepository()
	ctx := context.Background()

	user := &entity.User{
		Email:        "clark@kent.com",
		Password:     "superman",
		PasswordHash: []byte{9, 9, 9},
	}
	_ = repo.CreateUser(ctx, user)

	t.Run("gets an existing user", func(t *testing.T) {
		result, err := repo.GetUser(ctx, "clark@kent.com")
		require.NoError(t, err)
		require.Equal(t, user.Email, result.Email)
	})

	t.Run("returns error for non-existing user", func(t *testing.T) {
		result, err := repo.GetUser(ctx, "anonymous@user.com")
		require.Nil(t, result)
		require.ErrorIs(t, err, entity.ErrRecordNotFound)
	})
}

func TestMockRepository_ServiceStatus(t *testing.T) {
	repo := NewMockRepository()
	err := repo.ServiceStatus(context.Background())
	require.NoError(t, err)
}

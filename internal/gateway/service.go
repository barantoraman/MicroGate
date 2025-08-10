package gateway

import (
	"context"

	authEntity "github.com/barantoraman/microgate/internal/auth/repo/entity"
	taskEntity "github.com/barantoraman/microgate/internal/task/repo/entity"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
)

type Service interface {
	AddTask(ctx context.Context, task taskEntity.Task) (string, error)
	ListTask(ctx context.Context, userID int64) ([]taskEntity.Task, error)
	DeleteTask(ctx context.Context, taskID string, userID int64) error

	SignUp(ctx context.Context, user authEntity.User) (int64, tokenPkg.Token, error)
	Login(ctx context.Context, user authEntity.User) (int64, tokenPkg.Token, error)
	Logout(ctx context.Context, token tokenPkg.Token) error
}

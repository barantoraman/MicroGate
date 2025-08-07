package contract

import (
	"context"

	"github.com/barantoraman/microgate/internal/auth/repo/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUser(ctx context.Context, email string) (*entity.User, error)
	ServiceStatus(ctx context.Context) error
}

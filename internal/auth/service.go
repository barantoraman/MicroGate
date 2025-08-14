package auth

import (
	"context"

	"github.com/barantoraman/microgate/internal/auth/repo/entity"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
)

type Service interface {
	IsAuth(ctx context.Context, token tokenPkg.Token) (tokenPkg.Token, error)
	SignUp(ctx context.Context, user entity.User) (int64, tokenPkg.Token, error)
	Login(ctx context.Context, user entity.User) (int64, tokenPkg.Token, error)
	Logout(ctx context.Context, token tokenPkg.Token) error
	ServiceStatus(ctx context.Context) (int, error)
}

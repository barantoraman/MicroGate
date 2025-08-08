package auth

import (
	"context"

	cacheContract "github.com/barantoraman/microgate/internal/auth/cache/contract"
	repoContract "github.com/barantoraman/microgate/internal/auth/repo/contract"
	"github.com/barantoraman/microgate/internal/auth/repo/entity"
	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
)

type authService struct {
	userRepository repoContract.UserRepository
	store          cacheContract.Store
	logger         loggerContract.Logger
}

func NewService(userRepo repoContract.UserRepository, store cacheContract.Store, logger loggerContract.Logger) Service {
	return &authService{
		userRepository: userRepo,
		store:          store,
		logger:         logger,
	}
}

func (a *authService) IsAuth(ctx context.Context, sessionToken tokenPkg.Token) (tokenPkg.Token, error) {
}
func (a *authService) SignUp(ctx context.Context, user entity.User) (int64, tokenPkg.Token, error) {}
func (a *authService) Login(ctx context.Context, user entity.User) (int64, tokenPkg.Token, error)  {}
func (a *authService) Logout(ctx context.Context, sessionToken tokenPkg.Token) error               {}
func (a *authService) ServiceStatus(ctx context.Context) (int, error)                              {}

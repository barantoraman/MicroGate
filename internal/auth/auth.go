package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	cacheContract "github.com/barantoraman/microgate/internal/auth/cache/contract"
	repoContract "github.com/barantoraman/microgate/internal/auth/repo/contract"
	"github.com/barantoraman/microgate/internal/auth/repo/entity"
	userPkg "github.com/barantoraman/microgate/internal/auth/repo/user"
	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
	"github.com/barantoraman/microgate/pkg/validator"
	"go.uber.org/zap"
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
	v := validator.New()
	tokenPkg.ValidateTokenPlaintext(v, sessionToken.PlainText)
	if !v.Valid() {
		a.logger.Error("failed to validate token",
			zap.String("service", "auth"),
			zap.Any("validationErrors", v.Errors),
		)
		return tokenPkg.Token{}, fmt.Errorf("failed to validate token")
	}
	tkn, err := a.store.Get(ctx, sessionToken.PlainText)
	if err != nil {
		a.logger.Error("session is not available",
			zap.String("service", "auth"),
			zap.Error(err),
		)
		return tokenPkg.Token{}, errors.New("session is not available")
	}
	return tkn, nil
}

func (a *authService) SignUp(ctx context.Context, user entity.User) (int64, tokenPkg.Token, error) {
	if err := user.Set(user.Password); err != nil {
		a.logger.Error("failed to hash password",
			zap.String("service", "auth"),
			zap.Int64("userId", user.UserID),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, errors.New("failed to hash password")
	}

	v := validator.New()
	userPkg.ValidateUser(v, &user)
	if !v.Valid() {
		a.logger.Warn("failed to validate user",
			zap.String("service", "auth"),
			zap.Int64("userId", user.UserID),
			zap.Any("validationErrors", v.Errors),
		)
		return 0, tokenPkg.Token{}, errors.New("failed to validate user")
	}

	if err := a.userRepository.CreateUser(ctx, &user); err != nil {
		a.logger.Error("failed to create a new user",
			zap.String("service", "auth"),
			zap.Int64("userId", user.UserID),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, errors.New("failed to create a new user")
	}

	sessionToken, err := tokenPkg.GenerateToken(user.UserID, time.Minute*60, tokenPkg.ScopeAuthentication)
	if err != nil {
		a.logger.Error("failed to generate a new token",
			zap.String("service", "auth"),
			zap.Int64("userId", user.UserID),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, errors.New("user created but, failed to create a new token")
	}

	if err := a.store.Set(ctx, sessionToken); err != nil {
		a.logger.Error("failed to set session token to redis",
			zap.String("service", "auth"),
			zap.Int64("userId", user.UserID),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, errors.New("failed to set session token to redis")
	}
	return user.UserID, *sessionToken, nil
}

func (a *authService) Login(ctx context.Context, user entity.User) (int64, tokenPkg.Token, error) {
	if err := user.Set(user.Password); err != nil {
		a.logger.Error("failed to hash password",
			zap.String("service", "auth"),
			zap.String("email", user.Email),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, errors.New("failed to hash password")
	}

	v := validator.New()
	userPkg.ValidateUser(v, &user)
	if !v.Valid() {
		a.logger.Warn("failed to validate user",
			zap.String("service", "auth"),
			zap.String("email", user.Email),
			zap.Any("validationErrors", v.Errors),
		)
		return 0, tokenPkg.Token{}, errors.New("failed to validate user")
	}

	usr, err := a.userRepository.GetUser(ctx, user.Email)
	if err != nil {
		if errors.Is(err, entity.ErrRecordNotFound) {
			a.logger.Warn("user not found",
				zap.String("service", "auth"),
				zap.String("email", user.Email),
			)
			return 0, tokenPkg.Token{}, errors.New("user not found")
		}
		a.logger.Error("error retrieving user",
			zap.String("service", "auth"),
			zap.String("email", user.Email),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, err
	}

	match, err := usr.Matches(user.Password)
	if err != nil {
		a.logger.Error("error matching password",
			zap.String("service", "auth"),
			zap.Int64("userId", usr.UserID),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, errors.New("user not found")
	}
	if !match {
		a.logger.Warn("wrong password attempt",
			zap.String("service", "auth"),
			zap.Int64("userId", usr.UserID),
		)
		return 0, tokenPkg.Token{}, errors.New("wrong password")
	}

	sessionToken, err := tokenPkg.GenerateToken(usr.UserID, time.Minute*60, tokenPkg.ScopeAuthentication)
	if err != nil {
		a.logger.Error("failed to generate token",
			zap.String("service", "auth"),
			zap.Int64("userId", usr.UserID),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, errors.New("failed to generate token")
	}

	if err := a.store.Set(ctx, sessionToken); err != nil {
		a.logger.Error("failed to set session token",
			zap.String("service", "auth"),
			zap.Int64("userId", usr.UserID),
			zap.Error(err),
		)
		return 0, tokenPkg.Token{}, errors.New("failed to set session token")
	}

	a.logger.Info("user login successful",
		zap.String("service", "auth"),
		zap.Int64("userId", usr.UserID),
	)

	return usr.UserID, *sessionToken, nil
}

func (a *authService) Logout(ctx context.Context, sessionToken tokenPkg.Token) error {
	v := validator.New()
	tokenPkg.ValidateTokenPlaintext(v, sessionToken.PlainText)
	if !v.Valid() {
		a.logger.Error("failed to validate token",
			zap.Any("errors", v.Errors),
			zap.String("service", "auth"),
			zap.String("method", "Logout"),
		)
		return errors.New("failed to validate token")
	}

	hash := sha256.Sum256([]byte(sessionToken.PlainText))
	sessionToken.Hash = hash[:]
	if err := a.store.Delete(ctx, string(sessionToken.Hash)); err != nil {
		a.logger.Error("failed to delete session",
			zap.Error(err),
			zap.String("service", "auth"),
			zap.String("method", "Logout"),
		)
		return errors.New("failed to delete session")
	}

	a.logger.Info("logout successful",
		zap.String("service", "auth"),
		zap.String("method", "Logout"),
	)
	return nil
}

func (a *authService) ServiceStatus(ctx context.Context) (int, error) {
	if err := a.userRepository.ServiceStatus(ctx); err != nil {
		a.logger.Error("user repository service status check failed",
			zap.Error(err),
			zap.String("service", "auth"),
			zap.String("method", "ServiceStatus"),
		)
		return http.StatusInternalServerError, errors.New("internal server error")
	}

	a.logger.Info("service status OK",
		zap.String("service", "auth"),
		zap.String("method", "ServiceStatus"),
	)
	return http.StatusOK, nil
}

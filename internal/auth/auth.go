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
		// TODO
		a.logger.Error("failed to validate token")
		return tokenPkg.Token{}, fmt.Errorf("failed to validate token")
	}
	tkn, err := a.store.Get(ctx, sessionToken.PlainText)
	if err != nil {
		// TODO
		a.logger.Error("session is not available")
		return tokenPkg.Token{}, errors.New("session is not available")
	}
	return tkn, nil
}

func (a *authService) SignUp(ctx context.Context, user entity.User) (int64, tokenPkg.Token, error) {
	if err := user.Set(user.Password); err != nil {
		// TODO
		//a.logger.Error("failed to hash password %v", user.UserID)
		return 0, tokenPkg.Token{}, errors.New("failed to hash password")
	}

	v := validator.New()
	userPkg.ValidateUser(v, &user)
	if !v.Valid() {
		// TODO
		//a.logger.Error("failed to user validation")
		return 0, tokenPkg.Token{}, errors.New("failed to validate user")
	}

	if err := a.userRepository.CreateUser(ctx, &user); err != nil {
		// TODO
		// a.logger.Error("failed to create a new user")
		return 0, tokenPkg.Token{}, errors.New("failed to create a new user")
	}

	sessionToken, err := tokenPkg.GenerateToken(user.UserID, time.Minute*60, tokenPkg.ScopeAuthentication)
	if err != nil {
		// TODO
		// a.logger.Error("failed to generate a new token")
		return 0, tokenPkg.Token{}, errors.New("user created but, failed to create a new token")
	}

	if err := a.store.Set(ctx, sessionToken); err != nil {
		// TODO
		// a.logger.Error("failed to generate a new token")
		return 0, tokenPkg.Token{}, errors.New("failed to set session token to redis")
	}
	return user.UserID, *sessionToken, nil
}

func (a *authService) Login(ctx context.Context, user entity.User) (int64, tokenPkg.Token, error) {
	if err := user.Set(user.Password); err != nil {
		// TODO: log
		return 0, tokenPkg.Token{}, errors.New("failed to hash password")
	}

	v := validator.New()
	userPkg.ValidateUser(v, &user)
	if !v.Valid() {
		// TODO log
		return 0, tokenPkg.Token{}, errors.New("failed to validate user")
	}

	usr, err := a.userRepository.GetUser(ctx, user.Email)
	if err != nil {
		if errors.Is(err, entity.ErrRecordNotFound) {
			//TODO log here..
			return 0, tokenPkg.Token{}, errors.New("user not found")
		}
		return 0, tokenPkg.Token{}, err
	}

	match, err := usr.Matches(user.Password)
	if err != nil {
		//TODO log here..
		return 0, tokenPkg.Token{}, errors.New("user not found")
	}
	if !match {
		//TODO log here..
		return 0, tokenPkg.Token{}, errors.New("wrong password")
	}

	sessionToken, err := tokenPkg.GenerateToken(usr.UserID, time.Minute*60, tokenPkg.ScopeAuthentication)
	if err != nil {
		//TODO log here..
		return 0, tokenPkg.Token{}, errors.New("failed to generate token")
	}

	if err := a.store.Set(ctx, sessionToken); err != nil {
		// TODO: log here
		return 0, tokenPkg.Token{}, errors.New("failed to set session token")
	}
	return usr.UserID, *sessionToken, nil
}

func (a *authService) Logout(ctx context.Context, sessionToken tokenPkg.Token) error {
	v := validator.New()
	tokenPkg.ValidateTokenPlaintext(v, sessionToken.PlainText)
	if !v.Valid() {
		// TODO: log here
		return errors.New("failed to validate token")
	}

	hash := sha256.Sum256([]byte(sessionToken.PlainText))
	sessionToken.Hash = hash[:]
	if err := a.store.Delete(ctx, string(sessionToken.Hash)); err != nil {
		// TODO: log here
		return errors.New("failed to delete session")
	}
	return nil
}

func (a *authService) ServiceStatus(ctx context.Context) (int, error) {
	// TODO: log here
	if err := a.userRepository.ServiceStatus(ctx); err != nil {
		return http.StatusInternalServerError, errors.New("")
	}
	return http.StatusOK, nil
}

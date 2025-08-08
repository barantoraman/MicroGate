package auth

import (
	dbContract "github.com/barantoraman/microgate/internal/auth/db/contract"
	repoContract "github.com/barantoraman/microgate/internal/auth/repo/contract"
	"github.com/barantoraman/microgate/internal/auth/repo/mock"
	"github.com/barantoraman/microgate/internal/auth/repo/user"
	"github.com/barantoraman/microgate/pkg/config"
)

const (
	pqType   = "postgres"
	mockType = "inmemory"
)

func GetUserRepository(cfg config.AuthServiceConfigurations, conn dbContract.DBConnection) repoContract.UserRepository {
	switch cfg.DBType {
	case pqType:
		return user.NewUserRepository(conn)
	case mockType:
		return mock.NewMockRepository()
	default:
		return nil
	}
}

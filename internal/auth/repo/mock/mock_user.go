package mock

import (
	"context"
	"sync"

	repoContract "github.com/barantoraman/microgate/internal/auth/repo/contract"
	"github.com/barantoraman/microgate/internal/auth/repo/entity"
)

type mockRepository struct {
	mu    sync.RWMutex
	users []entity.User
}

func NewMockRepository() repoContract.UserRepository {
	return &mockRepository{
		mu:    sync.RWMutex{},
		users: make([]entity.User, 0),
	}
}

func (m *mockRepository) CreateUser(ctx context.Context, user *entity.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, u := range m.users {
		if u.Email == user.Email {
			return entity.ErrDuplicateEmail
		}
	}

	user.UserID = int64(len(m.users) + 1)
	m.users = append(m.users, *user)
	return nil
}

func (m *mockRepository) GetUser(ctx context.Context, email string) (*entity.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, u := range m.users {
		if u.Email == email {
			ucp := u
			return &ucp, nil
		}
	}
	return nil, entity.ErrRecordNotFound
}

func (m *mockRepository) ServiceStatus(ctx context.Context) error {
	return nil
}

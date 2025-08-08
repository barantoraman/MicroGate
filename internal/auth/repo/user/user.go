package user

import (
	"context"
	"database/sql"
	"time"

	dbContract "github.com/barantoraman/microgate/internal/auth/db/contract"
	repoContract "github.com/barantoraman/microgate/internal/auth/repo/contract"
	"github.com/barantoraman/microgate/internal/auth/repo/entity"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(conn dbContract.DBConnection) repoContract.UserRepository {
	return &userRepository{
		db: conn.DB(),
	}
}

func (u *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id`

	args := []any{user.Email, user.PasswordHash}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := u.db.QueryRowContext(ctx, query, args...).Scan(&user.UserID); err != nil {
		switch {
		// Buraya dikkat
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return entity.ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (u *userRepository) GetUser(ctx context.Context, email string) (*entity.User, error) {
	panic("unimplemented")
}

func (u *userRepository) ServiceStatus(ctx context.Context) error {
	panic("unimplemented")
}

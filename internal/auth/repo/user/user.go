package user

import (
	"context"
	"database/sql"
	"errors"
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
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return entity.ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (u *userRepository) GetUser(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user entity.User
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, entity.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepository) ServiceStatus(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := u.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

package contract

import (
	"context"

	tokenPkg "github.com/barantoraman/microgate/pkg/token"
)

type Store interface {
	Get(ctx context.Context, sessionTokenHash string) (tokenPkg.Token, error)
	Set(ctx context.Context, sessionToken *tokenPkg.Token) error
	Delete(ctx context.Context, token string) error
}

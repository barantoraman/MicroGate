package ctx

import (
	"context"
	"net/http"

	"github.com/barantoraman/microgate/internal/auth/repo/entity"
)

// TODO: pkg/ is not the right place..
type contextKey string

const userContextKey = contextKey("user")

func SetUser(r *http.Request, u *entity.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, u)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *entity.User {
	user, ok := r.Context().Value(userContextKey).(*entity.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

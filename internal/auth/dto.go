package auth

import (
	"github.com/barantoraman/microgate/internal/auth/repo/entity"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
)

// IsAuth - Req & Resp
type (
	IsAuthRequest struct {
		Token tokenPkg.Token `json:"token"`
	}
	IsAuthResponse struct {
		Token tokenPkg.Token `json:"token,omitempty"`
		Err   string         `json:"err,omitempty"`
	}
)

// SignUp - Req & Resp
type (
	SignUpRequest struct {
		User entity.User `json:"user"`
	}
	SignUpResponse struct {
		UserId int64          `json:"userId,omitempty"`
		Token  tokenPkg.Token `json:"token,omitempty"`
		Err    string         `json:"err,omitempty"`
	}
)

// LogIn - Req & Resp
type (
	LoginRequest struct {
		User entity.User `json:"user"`
	}
	LoginResponse struct {
		UserId int64          `json:"userId,omitempty"`
		Token  tokenPkg.Token `json:"token,omitempty"`
		Err    string         `json:"err,omitempty"`
	}
)

// Logout - Req & Resp
type (
	LogoutRequest struct {
		Token tokenPkg.Token `json:"token"`
	}
	LogoutResponse struct {
		Err string `json:"err,omitempty"`
	}
)

// Service Status - Req & Resp
type (
	ServiceStatusRequest  struct{}
	ServiceStatusResponse struct {
		Code int    `json:"code"`
		Err  string `json:"err,omitempty"`
	}
)

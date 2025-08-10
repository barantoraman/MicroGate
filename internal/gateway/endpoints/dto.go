package endpoints

import (
	authEntity "github.com/barantoraman/microgate/internal/auth/repo/entity"
	taskEntity "github.com/barantoraman/microgate/internal/task/repo/entity"
	"github.com/barantoraman/microgate/pkg/token"
)

type (
	AddTaskRequest struct {
		Task taskEntity.Task `json:"task"`
	}
	AddTaskResponse struct {
		TaskID string `json:"task_id,omitempty"`
		Err    string `json:"err,omitempty"`
	}
)

type (
	ListTaskRequest struct {
		UserID int64 `json:"user_id"`
	}
	ListTaskResponse struct {
		Tasks []taskEntity.Task `json:"tasks,omitempty"`
		Err   string            `json:"err,omitempty"`
	}
)

type (
	DeleteTaskRequest struct {
		TaskID string `json:"task_id,omitempty"`
		UserID int64  `json:"user_id"`
	}
	DeleteTaskResponse struct {
		Err string `json:"err,omitempty"`
	}
)

type (
	SignUpRequest struct {
		User authEntity.User `json:"user"`
	}
	SignUpResponse struct {
		UserID int64       `json:"user_id,omitempty"`
		Token  token.Token `json:"token,omitempty"`
		Err    string      `json:"err,omitempty"`
	}
)

type (
	LoginRequest struct {
		User authEntity.User `json:"user"`
	}
	LoginResponse struct {
		UserID int64       `json:"user_id,omitempty"`
		Token  token.Token `json:"token"`
		Err    string      `json:"err,omitempty"`
	}
)

type (
	LogoutRequest struct {
		Token token.Token `json:"token"`
	}
	LogoutResponse struct {
		Err string `json:"err,omitempty"`
	}
)

package task

import "github.com/barantoraman/microgate/internal/task/repo/entity"

type (
	CreateTaskRequest struct {
		Task entity.Task `json:"task"`
	}
	CreateTaskResponse struct {
		TaskID string `json:"task_id,omitempty"`
		Err    string `json:"err,omitempty"`
	}
)

type (
	ListTaskRequest struct {
		UserID int64 `json:"user_id"`
	}
	ListTaskResponse struct {
		Tasks []entity.Task `json:"tasks,omitempty"`
		Err   string        `json:"err,omitempty"`
	}
)

type (
	DeleteTaskRequest struct {
		TaskID string `json:"task_id"`
		UserID int64  `json:"user_id"`
	}
	DeleteTaskResponse struct {
		Err string `json:"err,omitempty"`
	}
)

type (
	ServiceStatusRequest struct {
	}
	ServiceStatusResponse struct {
		Code int    `json:"code,omitempty"`
		Err  string `json:"err,omitempty"`
	}
)

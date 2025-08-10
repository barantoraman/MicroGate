package gateway

import (
	"context"
	"errors"
	"fmt"

	authPb "github.com/barantoraman/microgate/internal/auth/pb"
	authEntity "github.com/barantoraman/microgate/internal/auth/repo/entity"
	taskPb "github.com/barantoraman/microgate/internal/task/pb"
	taskEntity "github.com/barantoraman/microgate/internal/task/repo/entity"
	taskPkg "github.com/barantoraman/microgate/internal/task/repo/task"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
	"github.com/barantoraman/microgate/pkg/validator"
)

type apiGatewayService struct {
	taskClient taskPb.TaskServiceClient
	authClient authPb.AuthClient
}

func (a *apiGatewayService) AddTask(ctx context.Context, task taskEntity.Task) (string, error) {
	v := validator.New()
	taskPkg.ValidateTask(v, task)
	if !v.Valid() {
		return "", errors.New("failed to validate event")
	}

	pbTask := &taskPb.Task{
		Id:          task.Id.Hex(),
		UserId:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
	}

	resp, err := a.taskClient.CreateTask(ctx, &taskPb.CreateTaskRequest{
		Task: pbTask,
	})
	if err != nil {
		return resp.TaskId, fmt.Errorf("failed to validate event: %v", v.Errors)
	}
	return resp.TaskId, nil

}

// DeleteTask implements Service.
func (a *apiGatewayService) DeleteTask(ctx context.Context, taskID string, userID int64) error {
	panic("unimplemented")
}

// ListTask implements Service.
func (a *apiGatewayService) ListTask(ctx context.Context, userID int64) ([]taskEntity.Task, error) {
	panic("unimplemented")
}

// Login implements Service.
func (a *apiGatewayService) Login(ctx context.Context, user authEntity.User) (int64, tokenPkg.Token, error) {
	panic("unimplemented")
}

// Logout implements Service.
func (a *apiGatewayService) Logout(ctx context.Context, token tokenPkg.Token) error {
	panic("unimplemented")
}

// SignUp implements Service.
func (a *apiGatewayService) SignUp(ctx context.Context, user authEntity.User) (int64, tokenPkg.Token, error) {
	panic("unimplemented")
}

func NewApiGatewayService(taskClient taskPb.TaskServiceClient, authClient authPb.AuthClient) Service {
	return &apiGatewayService{
		taskClient: taskClient,
		authClient: authClient,
	}
}

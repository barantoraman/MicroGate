package gateway

import (
	"context"

	authPb "github.com/barantoraman/microgate/internal/auth/pb"
	authEntity "github.com/barantoraman/microgate/internal/auth/repo/entity"
	taskPb "github.com/barantoraman/microgate/internal/task/pb"
	taskEntity "github.com/barantoraman/microgate/internal/task/repo/entity"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
)

type apiGatewayService struct {
	taskClient taskPb.TaskServiceClient
	authClient authPb.AuthClient
}

// AddTask implements Service.
func (a *apiGatewayService) AddTask(ctx context.Context, task taskEntity.Task) (string, error) {
	panic("unimplemented")
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

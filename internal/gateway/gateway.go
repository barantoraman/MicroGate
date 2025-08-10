package gateway

import (
	"context"
	"errors"
	"fmt"

	authPb "github.com/barantoraman/microgate/internal/auth/pb"
	authEntity "github.com/barantoraman/microgate/internal/auth/repo/entity"
	userPkg "github.com/barantoraman/microgate/internal/auth/repo/user"
	taskPb "github.com/barantoraman/microgate/internal/task/pb"
	taskEntity "github.com/barantoraman/microgate/internal/task/repo/entity"
	taskPkg "github.com/barantoraman/microgate/internal/task/repo/task"
	tokenPkg "github.com/barantoraman/microgate/pkg/token"
	"github.com/barantoraman/microgate/pkg/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type apiGatewayService struct {
	taskClient taskPb.TaskServiceClient
	authClient authPb.AuthClient
}

func (a *apiGatewayService) AddTask(ctx context.Context, task taskEntity.Task) (string, error) {
	v := validator.New()
	taskPkg.ValidateTask(v, task)
	if !v.Valid() {
		return "", errors.New("failed to validate task")
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
		return resp.TaskId, fmt.Errorf("failed to validate task: %v", v.Errors)
	}
	return resp.TaskId, nil

}

// ListTask implements Service.
func (a *apiGatewayService) ListTask(ctx context.Context, userID int64) ([]taskEntity.Task, error) {
	resp, err := a.taskClient.ListTask(ctx, &taskPb.ListTaskRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, errors.New("failed to list task")
	}

	var tasks []taskEntity.Task
	for _, e := range resp.Tasks {
		objId, err := primitive.ObjectIDFromHex(e.Id)
		if err != nil {
			return nil, errors.New("failed to convert object id")
		}
		task := taskEntity.Task{
			Id:          objId,
			UserID:      e.UserId,
			Title:       e.Title,
			Description: e.Description,
			Status:      e.Status,
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// DeleteTask implements Service.
func (a *apiGatewayService) DeleteTask(ctx context.Context, taskID string, userID int64) error {
	_, err := a.taskClient.DeleteTask(ctx, &taskPb.DeleteTaskRequest{
		TaskId: taskID,
		UserId: userID,
	})
	if err != nil {
		return errors.New("failed to delete task")
	}
	return nil
}

// SignUp implements Service.
func (a *apiGatewayService) SignUp(ctx context.Context, user authEntity.User) (int64, tokenPkg.Token, error) {
	err := user.Set(user.Password)
	if err != nil {
		return 0, tokenPkg.Token{}, errors.New("failed to hash password")
	}

	v := validator.New()
	userPkg.ValidateUser(v, &user)
	if !v.Valid() {
		return 0, tokenPkg.Token{}, fmt.Errorf("failed to user: %v", v.Errors)
	}

	usr := authPb.User{
		Email:    user.Email,
		Password: user.Password,
	}

	resp, err := a.authClient.SignUp(ctx, &authPb.SignUpRequest{User: &usr})
	if err != nil {
		return 0, tokenPkg.Token{}, errors.New("failed to signup")
	}

	tkn := tokenPkg.Token{
		PlainText: resp.Token.PlaintText,
		Hash:      resp.Token.Hash,
		UserID:    resp.Token.UserId,
		Expiry:    resp.Token.Expiry.AsTime(),
		Scope:     resp.Token.Scope,
	}
	return resp.UserId, tkn, nil
}

// Login implements Service.
func (a *apiGatewayService) Login(ctx context.Context, user authEntity.User) (int64, tokenPkg.Token, error) {
	err := user.Set(user.Password)
	if err != nil {
		return 0, tokenPkg.Token{}, fmt.Errorf("failed to hash password %w", err)
	}

	v := validator.New()
	userPkg.ValidateUser(v, &user)
	if !v.Valid() {
		return 0, tokenPkg.Token{}, fmt.Errorf("failed to user: %v", v.Errors)
	}

	usr := authPb.User{
		Email:    user.Email,
		Password: user.Password,
	}

	resp, err := a.authClient.Login(ctx, &authPb.LoginRequest{
		User: &usr,
	})
	if err != nil {
		return resp.UserId, tokenPkg.Token{}, errors.New(resp.Err)
	}

	tkn := tokenPkg.Token{
		PlainText: resp.Token.PlaintText,
		Hash:      resp.Token.Hash,
		UserID:    resp.Token.UserId,
		Expiry:    resp.Token.Expiry.AsTime(),
		Scope:     resp.Token.Scope,
	}

	return resp.UserId, tkn, nil
}

// Logout implements Service.
func (a *apiGatewayService) Logout(ctx context.Context, token tokenPkg.Token) error {
	panic("unimplemented")
}

func NewApiGatewayService(taskClient taskPb.TaskServiceClient, authClient authPb.AuthClient) Service {
	return &apiGatewayService{
		taskClient: taskClient,
		authClient: authClient,
	}
}

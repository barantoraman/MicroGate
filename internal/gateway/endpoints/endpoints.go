package endpoints

import (
	"context"

	"github.com/barantoraman/microgate/internal/gateway"
	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	AddTaskEndpoint    endpoint.Endpoint
	ListTaskEndpoint   endpoint.Endpoint
	DeleteTaskEndpoint endpoint.Endpoint

	SignUpEndpoint endpoint.Endpoint
	LoginEndpoint  endpoint.Endpoint
	LogoutEndpoint endpoint.Endpoint
}

func New(s gateway.Service) Set {
	return Set{
		AddTaskEndpoint:    MakeAddTaskEndpoint(s),
		ListTaskEndpoint:   MakeListTaskEndpoint(s),
		DeleteTaskEndpoint: MakeDeleteTaskEndpoint(s),

		SignUpEndpoint: MakeSignUpEndpoint(s),
		LoginEndpoint:  MakeLoginEndpoint(s),
		LogoutEndpoint: MakeLogoutEndpoint(s),
	}
}

func MakeAddTaskEndpoint(s gateway.Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(AddTaskRequest)
		taskID, err := s.AddTask(ctx, req.Task)
		if err != nil {
			return AddTaskResponse{TaskID: taskID, Err: err.Error()}, err
		}
		return AddTaskResponse{TaskID: taskID, Err: ""}, nil
	}
}

func MakeListTaskEndpoint(s gateway.Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(ListTaskRequest)

		tasks, err := s.ListTask(ctx, req.UserID)
		if err != nil {
			return ListTaskResponse{Tasks: tasks, Err: err.Error()}, err
		}
		return ListTaskResponse{Tasks: tasks, Err: ""}, nil
	}
}

func MakeDeleteTaskEndpoint(s gateway.Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DeleteTaskRequest)
		if err := s.DeleteTask(ctx, req.TaskID, req.UserID); err != nil {
			return DeleteTaskResponse{Err: err.Error()}, err
		}
		return DeleteTaskResponse{Err: ""}, nil
	}
}

func MakeSignUpEndpoint(s gateway.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SignUpRequest)

		userID, sessionToken, err := s.SignUp(ctx, req.User)
		if err != nil {
			return SignUpResponse{UserID: userID, Token: sessionToken, Err: err.Error()}, err
		}
		return SignUpResponse{UserID: userID, Token: sessionToken, Err: ""}, nil
	}
}

func MakeLoginEndpoint(s gateway.Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(LoginRequest)

		userId, sessionToken, err := s.Login(ctx, req.User)
		if err != nil {
			return LoginResponse{UserID: userId, Token: sessionToken, Err: err.Error()}, err
		}
		return LoginResponse{UserID: userId, Token: sessionToken, Err: ""}, nil
	}
}

func MakeLogoutEndpoint(s gateway.Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(LogoutRequest)

		err := s.Logout(ctx, req.Token)
		if err != nil {
			return LogoutResponse{Err: err.Error()}, err
		}
		return LogoutResponse{Err: ""}, nil
	}
}

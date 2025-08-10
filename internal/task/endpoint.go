package task

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	CreateTaskEndpoint    endpoint.Endpoint
	ListTaskEndpoint      endpoint.Endpoint
	DeleteTaskEndpoint    endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func New(s Service) Set {
	return Set{
		CreateTaskEndpoint:    MakeCreateTaskEndpoint(s),
		ListTaskEndpoint:      MakeListTaskEndpoint(s),
		DeleteTaskEndpoint:    MakeDeleteTaskEndpoint(s),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(s),
	}
}

func MakeCreateTaskEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(CreateTaskRequest)

		taskID, err := s.CreateTask(ctx, req.Task)
		if err != nil {
			return CreateTaskResponse{TaskID: taskID, Err: err.Error()}, err
		}
		return CreateTaskResponse{TaskID: taskID, Err: ""}, err
	}
}

func MakeListTaskEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(ListTaskRequest)

		tasks, err := s.ListTask(ctx, req.UserID)
		if err != nil {
			return ListTaskResponse{Tasks: tasks, Err: err.Error()}, err
		}
		return ListTaskResponse{Tasks: tasks, Err: ""}, nil
	}
}

func MakeDeleteTaskEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DeleteTaskRequest)

		if err := s.DeleteTask(ctx, req.TaskID, req.UserID); err != nil {
			return DeleteTaskResponse{Err: err.Error()}, err
		}
		return DeleteTaskResponse{Err: ""}, nil
	}
}

func MakeServiceStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		_ = request.(ServiceStatusRequest)
		code, err := s.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, err
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

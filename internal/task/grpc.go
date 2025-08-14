package task

import (
	"context"

	"github.com/barantoraman/microgate/internal/task/pb"
	"github.com/barantoraman/microgate/internal/task/repo/entity"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	pb.UnimplementedTaskServiceServer
	createTask    grpcTransport.Handler
	listTask      grpcTransport.Handler
	deleteTask    grpcTransport.Handler
	serviceStatus grpcTransport.Handler
}

func NewGRPCServer(ep Set) pb.TaskServiceServer {
	return &gRPCServer{
		createTask: grpcTransport.NewServer(
			ep.CreateTaskEndpoint,
			decodeCreateTaskRequest,
			encodeCreateTaskResponse),
		listTask: grpcTransport.NewServer(
			ep.ListTaskEndpoint,
			decodeListTaskRequest,
			encodeListTaskResponse),
		deleteTask: grpcTransport.NewServer(
			ep.DeleteTaskEndpoint,
			decodeDeleteTaskRequest,
			encodeDeleteTaskResponse),
		serviceStatus: grpcTransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeServiceStatusRequest,
			encodeServiceStatusResponse),
	}
}

func (g *gRPCServer) CreateTask(ctx context.Context, r *pb.CreateTaskRequest) (*pb.CreateTaskReply, error) {
	_, resp, err := g.createTask.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateTaskReply), nil
}

func (g *gRPCServer) ListTask(ctx context.Context, r *pb.ListTaskRequest) (*pb.ListTaskReply, error) {
	_, resp, err := g.listTask.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ListTaskReply), nil
}

func (g *gRPCServer) DeleteTask(ctx context.Context, r *pb.DeleteTaskRequest) (*pb.DeleteTaskReply, error) {
	_, resp, err := g.deleteTask.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.DeleteTaskReply), nil
}

func (g *gRPCServer) ServiceStatus(ctx context.Context, r *pb.ServiceStatusRequest) (*pb.ServiceStatusReply, error) {
	_, resp, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ServiceStatusReply), nil
}

func decodeCreateTaskRequest(_ context.Context, req any) (any, error) {
	request := req.(*pb.CreateTaskRequest)
	task := entity.Task{
		UserID:      request.Task.UserId,
		Title:       request.Task.Title,
		Description: request.Task.Description,
		Status:      request.Task.Status,
	}
	return CreateTaskRequest{Task: task}, nil
}

func encodeCreateTaskResponse(_ context.Context, res any) (any, error) {
	reply := res.(CreateTaskResponse)
	return &pb.CreateTaskReply{TaskId: reply.TaskID, Err: reply.Err}, nil
}

func decodeListTaskRequest(_ context.Context, req any) (any, error) {
	request := req.(*pb.ListTaskRequest)
	return ListTaskRequest{UserID: request.UserId}, nil
}

func encodeListTaskResponse(_ context.Context, res any) (any, error) {
	reply := res.(ListTaskResponse)

	var tasks []*pb.Task
	for _, e := range reply.Tasks {
		task := &pb.Task{
			Id:          e.Id.Hex(),
			UserId:      e.UserID,
			Title:       e.Title,
			Description: e.Description,
			Status:      e.Status,
		}
		tasks = append(tasks, task)
	}
	return &pb.ListTaskReply{Tasks: tasks, Err: ""}, nil
}

func decodeDeleteTaskRequest(_ context.Context, req any) (any, error) {
	request := req.(*pb.DeleteTaskRequest)
	return DeleteTaskRequest{TaskID: request.TaskId, UserID: request.UserId}, nil
}

func encodeDeleteTaskResponse(_ context.Context, res any) (any, error) {
	reply := res.(DeleteTaskResponse)
	return &pb.DeleteTaskReply{Err: reply.Err}, nil
}

func decodeServiceStatusRequest(_ context.Context, req any) (any, error) {
	_ = req.(*pb.ServiceStatusRequest)
	return ServiceStatusRequest{}, nil
}

func encodeServiceStatusResponse(_ context.Context, res any) (any, error) {
	reply := res.(ServiceStatusResponse)
	return &pb.ServiceStatusReply{Code: int32(reply.Code), Err: reply.Err}, nil
}

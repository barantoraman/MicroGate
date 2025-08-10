package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	taskPb "github.com/barantoraman/microgate/internal/task/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTaskClient(grpcAddr string) (taskPb.TaskServiceClient, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("fail to dial %w", err)
	}
	client := taskPb.NewTaskServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, _ := client.ServiceStatus(ctx, &taskPb.ServiceStatusRequest{})

	if resp.Code != http.StatusOK {
		return nil, errors.New(resp.Err)
	}

	return client, nil
}

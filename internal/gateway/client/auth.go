package client

import (
	"context"
	"fmt"
	"time"

	authPb "github.com/barantoraman/microgate/internal/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthClient(grpcAddr string) (authPb.AuthClient, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("fail to dial %s", err.Error())
	}
	client := authPb.NewAuthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.ServiceStatus(ctx, &authPb.ServiceStatusRequest{})
	if err != nil {
		return nil, fmt.Errorf("auth service error %s", err.Error())
	}
	return client, nil
}

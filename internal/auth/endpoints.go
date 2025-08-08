package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	IsAuthEndpoint        endpoint.Endpoint
	SignUpEndpoint        endpoint.Endpoint
	LoginEndpoint         endpoint.Endpoint
	LogoutEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
}

func New(s Service) Set {
	return Set{
		IsAuthEndpoint:        MakeIsAuthEndpoint(s),
		SignUpEndpoint:        MakeSignUpEndpoint(s),
		LoginEndpoint:         MakeLoginEndpoint(s),
		LogoutEndpoint:        MakeLogoutEndpoint(s),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(s),
	}
}

func MakeIsAuthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {}
}

func MakeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {}
}

func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {}
}

func MakeLogoutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {}
}

func MakeServiceStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {}
}

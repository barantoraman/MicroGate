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
	return func(ctx context.Context, request any) (any, error) {
		req := request.(IsAuthRequest)

		tkn, err := s.IsAuth(ctx, req.Token)
		if err != nil {
			return IsAuthResponse{Token: tkn, Err: err.Error()}, err
		}
		return IsAuthResponse{Token: tkn, Err: ""}, err
	}
}

func MakeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(SignUpRequest)

		userID, sessionToken, err := s.SignUp(ctx, req.User)
		if err != nil {
			return SignUpResponse{UserId: userID, Token: sessionToken, Err: err.Error()}, err
		}
		return SignUpResponse{UserId: userID, Token: sessionToken, Err: ""}, nil
	}
}

func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(LoginRequest)

		userId, sessionToken, err := s.Login(ctx, req.User)
		if err != nil {
			return LoginResponse{UserId: 0, Token: sessionToken, Err: err.Error()}, err
		}
		return LoginResponse{UserId: userId, Token: sessionToken, Err: ""}, nil
	}
}

func MakeLogoutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(LogoutRequest)

		err := s.Logout(ctx, req.Token)
		if err != nil {
			return LogoutResponse{Err: err.Error()}, err
		}
		return LogoutResponse{Err: ""}, nil
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

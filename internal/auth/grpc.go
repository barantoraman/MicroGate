package auth

import (
	"context"

	"github.com/barantoraman/microgate/internal/auth/pb"
	grpcTransport "github.com/go-kit/kit/transport/grpc"
)

type gRPCServer struct {
	pb.UnimplementedAuthServer
	isAuth        grpcTransport.Handler
	signUp        grpcTransport.Handler
	login         grpcTransport.Handler
	logout        grpcTransport.Handler
	serviceStatus grpcTransport.Handler
}

func NewGRPCServer(ep Set) pb.AuthServer {
	return &gRPCServer{
		isAuth: grpcTransport.NewServer(
			ep.IsAuthEndpoint,
			decodeIsAuthRequest,
			encodeIsAuthResponse),
		signUp: grpcTransport.NewServer(
			ep.SignUpEndpoint,
			decodeSignUpRequest,
			encodeSignUpResponse),
		login: grpcTransport.NewServer(
			ep.LoginEndpoint,
			decodeLoginRequest,
			encodeLoginResponse),
		logout: grpcTransport.NewServer(
			ep.LogoutEndpoint,
			decodeLogoutRequest,
			encodeLogoutResponse),
		serviceStatus: grpcTransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeServiceStatusRequest,
			encodeServiceStatusResponse),
	}
}

func (g *gRPCServer) IsAuth(ctx context.Context, r *pb.IsAuthRequest) (*pb.IsAuthReply, error) {
	_, resp, err := g.isAuth.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.IsAuthReply), nil
}

func (g *gRPCServer) SignUp(ctx context.Context, r *pb.SignUpRequest) (*pb.SignUpReply, error) {
	_, resp, err := g.signUp.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SignUpReply), nil
}

func (g *gRPCServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginReply, error) {
	_, resp, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LoginReply), nil
}

func (g *gRPCServer) Logout(ctx context.Context, r *pb.LogoutRequest) (*pb.LogoutReply, error) {
	_, resp, err := g.logout.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LogoutReply), nil
}

func (g *gRPCServer) ServiceStatus(ctx context.Context, r *pb.ServiceStatusRequest) (*pb.ServiceStatusReply, error) {
	_, resp, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ServiceStatusReply), nil
}

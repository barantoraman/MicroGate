package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/barantoraman/microgate/internal/auth"
	"github.com/barantoraman/microgate/internal/auth/cache"
	"github.com/barantoraman/microgate/internal/auth/db"
	"github.com/barantoraman/microgate/internal/auth/pb"
	"github.com/barantoraman/microgate/internal/auth/repo"
	"github.com/barantoraman/microgate/pkg/config"
	"github.com/barantoraman/microgate/pkg/logger"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	logger := logger.GetLogger("debug")

	var cfg config.AuthServiceConfigurations
	loader := config.GetLoader()
	if err := loader.GetConfigByKey("auth_service", &cfg); err != nil {
		logger.Fatal("failed to get config", zap.Error(err))
	}

	conn, err := db.GetDatabase(cfg)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := repo.GetUserRepository(cfg, conn)
	redis := cache.GetStore(ctx, cfg, logger)
	service := auth.NewService(repo, redis, logger)
	eps := auth.New(service)
	gRPCServer := auth.NewGRPCServer(eps)
	healthServer := health.NewServer()

	grpcAddr := net.JoinHostPort(cfg.GRPCHost, cfg.GRPCPort)

	var g group.Group
	{
		grpcListener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Fatal("error during grpc listen", zap.Error(err))
		}

		g.Add(func() error {
			logger.Debug("starting gRPC server", zap.String("addr", grpcAddr))
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterAuthServer(baseServer, gRPCServer)
			healthpb.RegisterHealthServer(baseServer, healthServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			logger.Info("shutting down gRPC server")
			grpcListener.Close()
		})
	}

	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				logger.Info("received signal", zap.String("signal", sig.String()))
				cancel()
				return fmt.Errorf("signal received: %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	if err := g.Run(); err != nil {
		logger.Error("server stopped", zap.Error(err))
	}
}

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
	cacheContract "github.com/barantoraman/microgate/internal/auth/cache/contract"
	"github.com/barantoraman/microgate/internal/auth/db"
	dbContract "github.com/barantoraman/microgate/internal/auth/db/contract"
	"github.com/barantoraman/microgate/internal/auth/pb"
	"github.com/barantoraman/microgate/internal/auth/repo"
	"github.com/barantoraman/microgate/pkg/config"
	"github.com/barantoraman/microgate/pkg/logger"
	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	logger := logger.GetLogger("debug")

	cfg := loadConfig(logger)
	conn := initDatabase(cfg, logger)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redis := initCache(ctx, cfg, logger)
	gRPCServer, healthServer := initGRPCServer(cfg, conn, redis, logger)

	var g group.Group
	startGRPCServer(&g, cfg, gRPCServer, healthServer, logger)
	initSignalHandler(&g, cancel, logger)

	if err := g.Run(); err != nil {
		logger.Error("server stopped", zap.Error(err))
	}
}

func loadConfig(logger loggerContract.Logger) config.AuthServiceConfigurations {
	var cfg config.AuthServiceConfigurations
	if err := config.GetLoader().GetConfigByKey("auth_service", &cfg); err != nil {
		logger.Fatal("failed to get config", zap.Error(err))
	}
	return cfg
}

func initDatabase(cfg config.AuthServiceConfigurations, logger loggerContract.Logger) dbContract.DBConnection {
	conn, err := db.GetDatabase(cfg)
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}
	return conn
}

func initCache(ctx context.Context, cfg config.AuthServiceConfigurations, logger loggerContract.Logger) cacheContract.Store {
	return cache.GetStore(ctx, cfg, logger)
}

func initGRPCServer(cfg config.AuthServiceConfigurations, conn dbContract.DBConnection, store cacheContract.Store, logger loggerContract.Logger) (pb.AuthServer, *health.Server) {
	repo := repo.GetUserRepository(cfg, conn)
	service := auth.NewService(repo, store, logger)
	eps := auth.New(service)
	return auth.NewGRPCServer(eps), health.NewServer()
}

func startGRPCServer(g *group.Group, cfg config.AuthServiceConfigurations, authServer pb.AuthServer, healthServer *health.Server, logger loggerContract.Logger) {
	grpcAddr := net.JoinHostPort(cfg.GRPCHost, cfg.GRPCPort)

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Fatal("error during grpc listen", zap.Error(err))
	}

	g.Add(func() error {
		logger.Debug("starting gRPC server", zap.String("addr", grpcAddr))
		baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
		pb.RegisterAuthServer(baseServer, authServer)
		healthpb.RegisterHealthServer(baseServer, healthServer)
		return baseServer.Serve(grpcListener)
	}, func(error) {
		logger.Info("shutting down gRPC server")
		grpcListener.Close()
	})
}

func initSignalHandler(g *group.Group, cancel context.CancelFunc, logger loggerContract.Logger) {
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

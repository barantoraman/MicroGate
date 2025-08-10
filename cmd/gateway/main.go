package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	authPb "github.com/barantoraman/microgate/internal/auth/pb"
	"github.com/barantoraman/microgate/internal/gateway"
	"github.com/barantoraman/microgate/internal/gateway/client"
	"github.com/barantoraman/microgate/internal/gateway/endpoints"
	"github.com/barantoraman/microgate/internal/gateway/transport"
	taskPb "github.com/barantoraman/microgate/internal/task/pb"
	"github.com/barantoraman/microgate/pkg/config"
	"github.com/barantoraman/microgate/pkg/logger"
	loggerContract "github.com/barantoraman/microgate/pkg/logger/contract"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
)

func main() {
	logger := logger.GetLogger("debug")

	cfg := loadConfig(logger)
	taskClient, authClient := initClients(cfg, logger)

	service := gateway.NewApiGatewayService(taskClient, authClient)
	eps := endpoints.New(service)
	httpHandler := transport.NewHTTPHandler(eps)

	var g group.Group
	initHTTPServer(&g, cfg, httpHandler, logger)
	initSignalHandler(&g)

	if err := g.Run(); err != nil {
		logger.Error("server stopped", zap.Error(err))
	}
}

func loadConfig(logger loggerContract.Logger) config.ApiGatewayServiceConfigurations {
	var cfg config.ApiGatewayServiceConfigurations
	if err := config.GetLoader().GetConfigByKey("api_gateway_service", &cfg); err != nil {
		logger.Fatal("failed to get config")
	}
	return cfg
}

func initClients(cfg config.ApiGatewayServiceConfigurations, logger loggerContract.Logger) (taskPb.TaskServiceClient, authPb.AuthClient) {
	grpcAddrAuth := net.JoinHostPort(cfg.AuthServiceHost, cfg.AuthServicePort)
	grpcAddrTask := net.JoinHostPort(cfg.TaskServiceHost, cfg.TaskServicePort)

	taskClient, err := client.NewTaskClient(grpcAddrTask)
	if err != nil {
		logger.Fatal("failed to get task client")
	}

	authClient, err := client.NewAuthClient(grpcAddrAuth)
	if err != nil {
		logger.Fatal("failed to get auth client")
	}

	return taskClient, authClient
}

func initHTTPServer(g *group.Group, cfg config.ApiGatewayServiceConfigurations, handler http.Handler, logger loggerContract.Logger) {
	httpAddr := net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort)

	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		logger.Fatal("error during http listen")
	}

	g.Add(func() error {
		return http.Serve(httpListener, handler)
	}, func(error) {
		httpListener.Close()
	})
}

func initSignalHandler(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

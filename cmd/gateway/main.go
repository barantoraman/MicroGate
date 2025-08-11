package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/barantoraman/microgate/internal/gateway"
	"github.com/barantoraman/microgate/internal/gateway/client"
	"github.com/barantoraman/microgate/internal/gateway/endpoints"
	"github.com/barantoraman/microgate/internal/gateway/transport"
	"github.com/barantoraman/microgate/pkg/config"
	"github.com/barantoraman/microgate/pkg/logger"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
)

func main() {
	logger := logger.GetLogger("debug")

	var cfg config.ApiGatewayServiceConfigurations
	err := config.GetLoader().GetConfigByKey("api_gateway_service", &cfg)
	if err != nil {
		logger.Fatal("failed to get config",
			zap.String("service", "api-gateway"),
			zap.Error(err),
		)
	}

	var grpcAddrAuth = net.JoinHostPort(cfg.AuthServiceHost, cfg.AuthServicePort)
	var grpcAddrTask = net.JoinHostPort(cfg.TaskServiceHost, cfg.TaskServicePort)
	var httpAddr = net.JoinHostPort(cfg.HTTPHost, cfg.HTTPPort)

	taskClient, err := client.NewTaskClient(grpcAddrTask)
	if err != nil {
		logger.Fatal("failed to get task client",
			zap.String("service", "api-gateway"),
			zap.String("target", "task-service"),
			zap.String("address", grpcAddrTask),
			zap.Error(err),
		)
	}

	authClient, err := client.NewAuthClient(grpcAddrAuth)
	if err != nil {
		logger.Fatal("failed to get auth client",
			zap.String("service", "api-gateway"),
			zap.String("target", "auth-service"),
			zap.String("address", grpcAddrAuth),
			zap.Error(err),
		)
	}

	service := gateway.NewApiGatewayService(taskClient, authClient)
	eps := endpoints.New(service)
	httpHandler := transport.NewHTTPHandler(eps)

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		if err != nil {
			logger.Fatal("error during http listen",
				zap.String("service", "api-gateway"),
				zap.String("address", httpAddr),
				zap.Error(err),
			)
		}

		g.Add(func() error {
			logger.Info("HTTP server started",
				zap.String("service", "api-gateway"),
				zap.String("address", httpAddr),
			)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			logger.Info("HTTP server stopped",
				zap.String("service", "api-gateway"),
			)
			httpListener.Close()
		})
	}
	{
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
	if err := g.Run(); err != nil {
		logger.Error("server stopped unexpectedly",
			zap.String("service", "api-gateway"),
			zap.Error(err),
		)
	}
}

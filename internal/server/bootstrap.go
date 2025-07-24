// internal/server/bootstrap.go
package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go-code-runner-microservice/api-gateway/internal/config"
	"go-code-runner-microservice/api-gateway/internal/logger"
	"go-code-runner-microservice/api-gateway/internal/middleware"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/coding_tests"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/company_auth"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/executor"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/problems"
)

func Run() {
	// Load environment variables
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Use basic logging before logger is initialized
		panic("failed to load configuration: " + err.Error())
	}

	// Initialize logger
	logConfig := logger.Config{
		Level:       cfg.Logging.Level,
		Environment: cfg.Logging.Environment,
		ServiceName: "api-gateway",
	}
	if err := logger.Initialize(logConfig); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	log := logger.Get()
	log.Info("starting api-gateway service",
		zap.String("version", "1.0.0"), // Add version from build info
		zap.String("environment", cfg.Logging.Environment),
		zap.String("log_level", cfg.Logging.Level),
	)

	// Create gRPC dial options with logging interceptor
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.UnaryClientLoggingInterceptor()),
	}

	// Initialize gRPC clients with logging
	executorClient, err := executor.NewClientWithOptions(cfg.ExecutorServiceAddress, grpcOpts...)
	if err != nil {
		log.Fatal("failed to connect to executor service",
			zap.String("address", cfg.ExecutorServiceAddress),
			zap.Error(err),
		)
	}
	defer executorClient.Close()
	log.Info("connected to executor service", zap.String("address", cfg.ExecutorServiceAddress))

	problemsClient, err := problems.NewClientWithOptions(cfg.ExecutorServiceAddress, grpcOpts...)
	if err != nil {
		log.Fatal("failed to connect to problems service",
			zap.String("address", cfg.ExecutorServiceAddress),
			zap.Error(err),
		)
	}
	defer problemsClient.Close()
	log.Info("connected to problems service", zap.String("address", cfg.ExecutorServiceAddress))

	codingTestsClient, err := coding_tests.NewClientWithOptions(cfg.ExecutorServiceAddress, grpcOpts...)
	if err != nil {
		log.Fatal("failed to connect to coding tests service",
			zap.String("address", cfg.ExecutorServiceAddress),
			zap.Error(err),
		)
	}
	defer codingTestsClient.Close()
	log.Info("connected to coding tests service", zap.String("address", cfg.ExecutorServiceAddress))

	companyAuthClient, err := company_auth.NewClientWithOptions(cfg.CompanyAuthAddress, grpcOpts...)
	if err != nil {
		log.Fatal("failed to connect to company auth service",
			zap.String("address", cfg.CompanyAuthAddress),
			zap.Error(err),
		)
	}
	defer companyAuthClient.Close()
	log.Info("connected to company auth service", zap.String("address", cfg.CompanyAuthAddress))

	// Create router
	r := NewRouter(executorClient, problemsClient, codingTestsClient, companyAuthClient)

	// Create HTTP server
	addr := ":" + cfg.ServerPort
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info("starting HTTP server", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutdown signal received, initiating graceful shutdown")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("shutting down HTTP server")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited successfully")
}

package server

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"go-code-runner-microservice/api-gateway/internal/config"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/coding_tests"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/company_auth"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/executor"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/problems"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	logger := log.New(os.Stdout, "API-GATEWAY: ", log.LstdFlags|log.Lmicroseconds)
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to laod configuration: %v1", err)
	}

	executorClient, err := executor.NewClient(cfg.ExecutorServiceAddress)
	logger.Println("executor service address: ", cfg.ExecutorServiceAddress)
	if err != nil {
		logger.Fatalf("failed to connect to executor service: %v", err)
	}
	defer executorClient.Close()

	problemsClient, err := problems.NewClient(cfg.ExecutorServiceAddress)
	logger.Println("problems service address: ", cfg.ExecutorServiceAddress)
	if err != nil {
		logger.Fatalf("failed to connect to problems service: %v", err)
	}
	defer problemsClient.Close()

	codingTestsClient, err := coding_tests.NewClient(cfg.ExecutorServiceAddress)
	logger.Println("coding tests service address: ", cfg.ExecutorServiceAddress)
	if err != nil {
		logger.Fatalf("failed to connect to coding tests service: %v", err)
	}
	defer codingTestsClient.Close()

	companyAuthClient, err := company_auth.NewClient(cfg.CompanyAuthAddress)
	logger.Println("company auth service address: ", cfg.CompanyAuthAddress)
	if err != nil {
		logger.Fatalf("failed to connect to company auth service: %v", err)
	}
	defer companyAuthClient.Close()

	r := NewRouter(executorClient, problemsClient, codingTestsClient, companyAuthClient)

	addr := ":" + cfg.ServerPort
	logger.Printf("starting HTTP server on %s", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutdown signal received, initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Println("Shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exiting")
}

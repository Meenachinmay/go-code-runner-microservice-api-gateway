package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-code-runner-microservice/api-gateway/internal/handler"
	"go-code-runner-microservice/api-gateway/internal/middleware"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/coding_tests"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/company_auth"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/executor"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/problems"
)

func NewRouter(
	executorClient *executor.Client,
	problemsClient *problems.Client,
	codingTestsClient *coding_tests.Client,
	companyAuthClient *company_auth.Client) *gin.Engine {

	r := gin.New()

	r.Use(middleware.ErrorHandlingMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Correlation-ID"}
	config.ExposeHeaders = []string{"X-Request-ID", "X-Correlation-ID"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	r.GET("/health", handler.MakeHealthHandler())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/execute", handler.MakeExecuteHandler(executorClient))
		v1.GET("/execute/job/:job_id", handler.MakeJobStatusHandler(executorClient))

		v1.GET("/problems", handler.MakeListProblemsHandler(problemsClient))
		v1.GET("/problems/:id", handler.MakeGetProblemHandler(problemsClient))
		v1.GET("/problems/:id/test-cases", handler.MakeGetTestCasesByProblemIDHandler(problemsClient))

		codingTests := v1.Group("/tests")
		{
			codingTests.GET("/:test_id/verify", handler.MakeVerifyTestHandler(codingTestsClient))
			codingTests.POST("/:test_id/start", handler.MakeStartTestHandler(codingTestsClient))
			codingTests.POST("/:test_id/submit", handler.MakeSubmitTestHandler(codingTestsClient))
			codingTests.POST("/generate", handler.MakeGenerateTestHandler(codingTestsClient))
			codingTests.GET("/company/:company_id", handler.MakeGetCompanyTestsHandler(codingTestsClient))
		}

		// Company authentication routes
		companyHandler := handler.NewCompanyHandler(companyAuthClient)
		companies := v1.Group("/companies")
		{
			companies.POST("/register", companyHandler.Register)
			companies.POST("/login", companyHandler.Login)
			companies.POST("/api-key", companyHandler.GenerateAPIKey)
			companies.POST("/client-id", companyHandler.GenerateClientID)
		}
	}

	return r

}

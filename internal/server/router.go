package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-code-runner-microservice/api-gateway/internal/handler"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/executor"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/problems"
)

func NewRouter(executorClient *executor.Client, problemsClient *problems.Client) *gin.Engine {

	r := gin.New()
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
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
	}

	return r

}

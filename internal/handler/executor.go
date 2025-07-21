package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-code-runner-microservice/api-gateway/internal/model"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/executor"
)

func MakeExecuteHandler(executorClient *executor.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.ExecuteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.ExecuteResponse{
				Success: false,
				Error:   "Invalid request payload: " + err.Error(),
			})
			return
		}

		if req.Language != "go" {
			c.JSON(http.StatusBadRequest, model.ExecuteResponse{
				Success: false,
				Error:   "Unsupported language. Only 'go' is supported.",
			})
			return
		}

		resp, err := executorClient.Execute(c.Request.Context(), req.Language, req.Code, req.ProblemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ExecuteResponse{
				Success: false,
				Error:   "Failed to executor: " + err.Error(),
			})
			return
		}

		if !resp.Success {
			c.JSON(http.StatusInternalServerError, model.ExecuteResponse{
				Success: false,
				Error:   resp.Error,
			})
			return
		}

		c.JSON(http.StatusAccepted, model.ExecuteResponse{
			Success: true,
			JobID:   resp.JobId,
			Message: resp.Message,
		})
	}
}

func MakeJobStatusHandler(executorClient *executor.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("job_id")

		resp, err := executorClient.GetJobStatus(c.Request.Context(), jobID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get job status: " + err.Error(),
			})
			return
		}

		if !resp.Success {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   resp.Error,
			})
			return
		}

		response := model.JobStatusResponse{
			Success: true,
			JobID:   resp.JobId,
			Status:  resp.Status,
			Output:  resp.Output,
			Error:   resp.Error,
		}

		if len(resp.TestResults) > 0 {
			response.TestResults = make([]model.TestResult, len(resp.TestResults))
			for i, tr := range resp.TestResults {
				response.TestResults[i] = model.TestResult{
					TestCaseID:     int(tr.TestCaseId),
					Input:          tr.Input,
					ExpectedOutput: tr.ExpectedOutput,
					ActualOutput:   tr.ActualOutput,
					Error:          tr.Error,
					Passed:         tr.Passed,
				}
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

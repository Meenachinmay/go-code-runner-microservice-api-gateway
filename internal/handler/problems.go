package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-code-runner-microservice/api-gateway/internal/model"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/problems"
)

func MakeListProblemsHandler(problemsClient *problems.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := problemsClient.ListProblems(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ListProblemsResponse{
				Success: false,
				Error:   "Failed to list problemResponses: " + err.Error(),
			})
			return
		}

		problemResponses := make([]model.ProblemResponse, len(resp.Problems))
		for i, p := range resp.Problems {
			problemResponses[i] = model.ProblemResponse{
				ID:          int(p.Id),
				Title:       p.Title,
				Description: p.Description,
				Difficulty:  p.Difficulty,
			}
			if p.CreatedAt != nil {
				problemResponses[i].CreatedAt = p.CreatedAt.AsTime().String()
			}
			if p.UpdatedAt != nil {
				problemResponses[i].UpdatedAt = p.UpdatedAt.AsTime().String()
			}
		}

		c.JSON(http.StatusOK, model.ListProblemsResponse{
			Success:  true,
			Problems: problemResponses,
		})
	}
}

func MakeGetProblemHandler(problemsClient *problems.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.GetProblemResponse{
				Success: false,
				Error:   "Invalid problem ID: " + err.Error(),
			})
			return
		}

		resp, err := problemsClient.GetProblem(c.Request.Context(), int32(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.GetProblemResponse{
				Success: false,
				Error:   "Failed to get problem: " + err.Error(),
			})
			return
		}

		problem := model.ProblemResponse{
			ID:          int(resp.Problem.Id),
			Title:       resp.Problem.Title,
			Description: resp.Problem.Description,
			Difficulty:  resp.Problem.Difficulty,
		}
		if resp.Problem.CreatedAt != nil {
			problem.CreatedAt = resp.Problem.CreatedAt.AsTime().String()
		}
		if resp.Problem.UpdatedAt != nil {
			problem.UpdatedAt = resp.Problem.UpdatedAt.AsTime().String()
		}

		c.JSON(http.StatusOK, model.GetProblemResponse{
			Success: true,
			Problem: problem,
		})
	}
}

func MakeGetTestCasesByProblemIDHandler(problemsClient *problems.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.GetTestCasesByProblemIDResponse{
				Success: false,
				Error:   "Invalid problem ID: " + err.Error(),
			})
			return
		}

		resp, err := problemsClient.GetTestCasesByProblemID(c.Request.Context(), int32(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.GetTestCasesByProblemIDResponse{
				Success: false,
				Error:   "Failed to get test cases: " + err.Error(),
			})
			return
		}

		testCases := make([]model.TestCaseResponse, len(resp.TestCases))
		for i, tc := range resp.TestCases {
			testCases[i] = model.TestCaseResponse{
				ID:             int(tc.Id),
				ProblemID:      int(tc.ProblemId),
				Input:          tc.Input,
				ExpectedOutput: tc.ExpectedOutput,
				IsHidden:       tc.IsHidden,
			}
			if tc.CreatedAt != nil {
				testCases[i].CreatedAt = tc.CreatedAt.AsTime().String()
			}
			if tc.UpdatedAt != nil {
				testCases[i].UpdatedAt = tc.UpdatedAt.AsTime().String()
			}
		}

		c.JSON(http.StatusOK, model.GetTestCasesByProblemIDResponse{
			Success:   true,
			TestCases: testCases,
		})
	}
}

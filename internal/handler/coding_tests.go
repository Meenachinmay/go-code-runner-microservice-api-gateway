package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-code-runner-microservice/api-gateway/internal/model"
	"go-code-runner-microservice/api-gateway/internal/service/grpc/coding_tests"
)

// MakeVerifyTestHandler creates a handler for verifying a coding test
func MakeVerifyTestHandler(codingTestsClient *coding_tests.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		testID := c.Param("test_id")
		if testID == "" {
			c.JSON(http.StatusBadRequest, model.VerifyTestResponse{
				Success: false,
				Error:   "Test ID is required",
			})
			return
		}

		resp, err := codingTestsClient.VerifyTest(c.Request.Context(), testID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.VerifyTestResponse{
				Success: false,
				Error:   "Failed to verify test: " + err.Error(),
			})
			return
		}

		// Convert proto CodingTest to model.CodingTest
		test := model.CodingTest{
			ID:                  resp.Test.Id,
			CompanyID:           int(resp.Test.CompanyId),
			ProblemID:           int(resp.Test.ProblemId),
			Status:              resp.Test.Status,
			TestDurationMinutes: int(resp.Test.TestDurationMinutes),
		}

		if resp.Test.CandidateName != "" {
			candidateName := resp.Test.CandidateName
			test.CandidateName = &candidateName
		}

		if resp.Test.CandidateEmail != "" {
			candidateEmail := resp.Test.CandidateEmail
			test.CandidateEmail = &candidateEmail
		}

		if resp.Test.SubmissionCode != "" {
			submissionCode := resp.Test.SubmissionCode
			test.SubmissionCode = &submissionCode
		}

		if resp.Test.PassedPercentage != 0 {
			passedPercentage := int(resp.Test.PassedPercentage)
			test.PassedPercentage = &passedPercentage
		}

		// Handle timestamps if they exist
		if resp.Test.StartedAt != nil {
			startedAt := resp.Test.StartedAt.AsTime()
			test.StartedAt = &startedAt
		}

		if resp.Test.CompletedAt != nil {
			completedAt := resp.Test.CompletedAt.AsTime()
			test.CompletedAt = &completedAt
		}

		if resp.Test.ExpiresAt != nil {
			test.ExpiresAt = resp.Test.ExpiresAt.AsTime()
		}

		if resp.Test.CreatedAt != nil {
			test.CreatedAt = resp.Test.CreatedAt.AsTime()
		}

		if resp.Test.UpdatedAt != nil {
			test.UpdatedAt = resp.Test.UpdatedAt.AsTime()
		}

		c.JSON(http.StatusOK, model.VerifyTestResponse{
			Success: true,
			Test:    test,
		})
	}
}

// MakeStartTestHandler creates a handler for starting a coding test
func MakeStartTestHandler(codingTestsClient *coding_tests.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		testID := c.Param("test_id")
		if testID == "" {
			c.JSON(http.StatusBadRequest, model.StartTestResponse{
				Success: false,
				Error:   "Test ID is required",
			})
			return
		}

		var req model.StartTestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.StartTestResponse{
				Success: false,
				Error:   "Invalid request payload: " + err.Error(),
			})
			return
		}

		resp, err := codingTestsClient.StartTest(c.Request.Context(), testID, req.CandidateName, req.CandidateEmail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.StartTestResponse{
				Success: false,
				Error:   "Failed to start test: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, model.StartTestResponse{
			Success: true,
			Message: resp.Message,
		})
	}
}

// MakeSubmitTestHandler creates a handler for submitting a coding test
func MakeSubmitTestHandler(codingTestsClient *coding_tests.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		testID := c.Param("test_id")
		if testID == "" {
			c.JSON(http.StatusBadRequest, model.SubmitTestResponse{
				Success: false,
				Error:   "Test ID is required",
			})
			return
		}

		var req model.SubmitTestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.SubmitTestResponse{
				Success: false,
				Error:   "Invalid request payload: " + err.Error(),
			})
			return
		}

		resp, err := codingTestsClient.SubmitTest(c.Request.Context(), testID, req.Code, int32(req.PassedPercentage))
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.SubmitTestResponse{
				Success: false,
				Error:   "Failed to submit test: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, model.SubmitTestResponse{
			Success: true,
			Message: resp.Message,
		})
	}
}

// MakeGenerateTestHandler creates a handler for generating a coding test
func MakeGenerateTestHandler(codingTestsClient *coding_tests.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.GenerateTestRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, model.GenerateTestResponse{
				Success: false,
				Error:   "Invalid request payload: " + err.Error(),
			})
			return
		}

		resp, err := codingTestsClient.GenerateTest(
			c.Request.Context(), 
			int32(req.CompanyID), 
			int32(req.ProblemID), 
			int32(req.ExpiresInHours),
			*req.ClientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.GenerateTestResponse{
				Success: false,
				Error:   "Failed to generate test: " + err.Error(),
			})
			return
		}

		// Convert proto CodingTest to model.CodingTest
		test := model.CodingTest{
			ID:                  resp.Test.Id,
			CompanyID:           int(resp.Test.CompanyId),
			ProblemID:           int(resp.Test.ProblemId),
			Status:              resp.Test.Status,
			TestDurationMinutes: int(resp.Test.TestDurationMinutes),
		}

		// Handle timestamps if they exist
		if resp.Test.ExpiresAt != nil {
			test.ExpiresAt = resp.Test.ExpiresAt.AsTime()
		}

		if resp.Test.CreatedAt != nil {
			test.CreatedAt = resp.Test.CreatedAt.AsTime()
		}

		if resp.Test.UpdatedAt != nil {
			test.UpdatedAt = resp.Test.UpdatedAt.AsTime()
		}

		c.JSON(http.StatusOK, model.GenerateTestResponse{
			Success: true,
			Test:    test,
			Link:    resp.Link,
		})
	}
}

// MakeGetCompanyTestsHandler creates a handler for getting company tests
func MakeGetCompanyTestsHandler(codingTestsClient *coding_tests.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyIDStr := c.Param("company_id")
		companyID, err := strconv.Atoi(companyIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.GetCompanyTestsResponse{
				Success: false,
				Error:   "Invalid company ID: " + err.Error(),
			})
			return
		}

		resp, err := codingTestsClient.GetCompanyTests(c.Request.Context(), int32(companyID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.GetCompanyTestsResponse{
				Success: false,
				Error:   "Failed to get company tests: " + err.Error(),
			})
			return
		}

		tests := make([]model.CodingTest, len(resp.Tests))
		for i, t := range resp.Tests {
			test := model.CodingTest{
				ID:                  t.Id,
				CompanyID:           int(t.CompanyId),
				ProblemID:           int(t.ProblemId),
				Status:              t.Status,
				TestDurationMinutes: int(t.TestDurationMinutes),
			}

			if t.CandidateName != "" {
				candidateName := t.CandidateName
				test.CandidateName = &candidateName
			}

			if t.CandidateEmail != "" {
				candidateEmail := t.CandidateEmail
				test.CandidateEmail = &candidateEmail
			}

			if t.SubmissionCode != "" {
				submissionCode := t.SubmissionCode
				test.SubmissionCode = &submissionCode
			}

			if t.PassedPercentage != 0 {
				passedPercentage := int(t.PassedPercentage)
				test.PassedPercentage = &passedPercentage
			}

			// Handle timestamps if they exist
			if t.StartedAt != nil {
				startedAt := t.StartedAt.AsTime()
				test.StartedAt = &startedAt
			}

			if t.CompletedAt != nil {
				completedAt := t.CompletedAt.AsTime()
				test.CompletedAt = &completedAt
			}

			if t.ExpiresAt != nil {
				test.ExpiresAt = t.ExpiresAt.AsTime()
			}

			if t.CreatedAt != nil {
				test.CreatedAt = t.CreatedAt.AsTime()
			}

			if t.UpdatedAt != nil {
				test.UpdatedAt = t.UpdatedAt.AsTime()
			}

			tests[i] = test
		}

		c.JSON(http.StatusOK, model.GetCompanyTestsResponse{
			Success: true,
			Tests:   tests,
		})
	}
}

package model

import (
	"time"
)

// Problem represents a coding problem
type Problem struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Difficulty  string    `json:"difficulty" db:"difficulty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProblemResponse is used for API responses with string timestamps
type ProblemResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// ListProblemsResponse is the response for listing problems
type ListProblemsResponse struct {
	Success  bool             `json:"success"`
	Problems []ProblemResponse `json:"problems"`
	Error    string           `json:"error,omitempty"`
}

// GetProblemResponse is the response for getting a single problem
type GetProblemResponse struct {
	Success bool           `json:"success"`
	Problem ProblemResponse `json:"problem"`
	Error   string         `json:"error,omitempty"`
}

type TestCase struct {
	ID             int       `json:"id" db:"id"`
	ProblemID      int       `json:"problem_id" db:"problem_id"`
	Input          string    `json:"input" db:"input"`
	ExpectedOutput string    `json:"expected_output" db:"expected_output"`
	IsHidden       bool      `json:"is_hidden" db:"is_hidden"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type TestResult struct {
	TestCaseID     int    `json:"test_case_id"`
	Input          string `json:"input,omitempty"`
	ExpectedOutput string `json:"expected_output,omitempty"`
	ActualOutput   string `json:"actual_output"`
	Error          string `json:"error,omitempty"`
	Passed         bool   `json:"passed"`
}

type ExecutionResults struct {
	Success     bool         `json:"success"`
	TestResults []TestResult `json:"test_results"`
}

// ExecuteRequest is the request for executing code
type ExecuteRequest struct {
	Language  string `json:"language" binding:"required"`
	Code      string `json:"code" binding:"required"`
	ProblemID int    `json:"problem_id,omitempty"`
}

// ExecuteResponse is the response for executing code
type ExecuteResponse struct {
	Success bool   `json:"success"`
	JobID   string `json:"job_id,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// JobStatusResponse is the response for checking job status
type JobStatusResponse struct {
	Success     bool         `json:"success"`
	JobID       string       `json:"job_id"`
	Status      string       `json:"status"`
	Output      string       `json:"output,omitempty"`
	Error       string       `json:"error,omitempty"`
	TestResults []TestResult `json:"test_results,omitempty"`
}

type Company struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	APIKey       *string   `json:"api_key,omitempty"`
	ClientID     *string   `json:"client_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CodingTest struct {
	ID                  string     `json:"id" db:"id"`
	CompanyID           int        `json:"company_id" db:"company_id"`
	ProblemID           int        `json:"problem_id" db:"problem_id"`
	CandidateName       *string    `json:"candidate_name" db:"candidate_name"`
	CandidateEmail      *string    `json:"candidate_email" db:"candidate_email"`
	Status              string     `json:"status" db:"status"`
	StartedAt           *time.Time `json:"started_at" db:"started_at"`
	CompletedAt         *time.Time `json:"completed_at" db:"completed_at"`
	ExpiresAt           time.Time  `json:"expires_at" db:"expires_at"`
	TestDurationMinutes int        `json:"test_duration_minutes" db:"test_duration_minutes"`
	SubmissionCode      *string    `json:"submission_code" db:"submission_code"`
	PassedPercentage    *int       `json:"passed_percentage" db:"passed_percentage"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// TestCaseResponse is used for API responses with string timestamps
type TestCaseResponse struct {
	ID             int    `json:"id"`
	ProblemID      int    `json:"problem_id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	IsHidden       bool   `json:"is_hidden"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// GetTestCasesByProblemIDResponse is the response for getting test cases for a problem
type GetTestCasesByProblemIDResponse struct {
	Success   bool              `json:"success"`
	TestCases []TestCaseResponse `json:"test_cases"`
	Error     string            `json:"error,omitempty"`
}

const (
	TestStatusPending   = "pending"
	TestStatusStarted   = "started"
	TestStatusCompleted = "completed"
	TestStatusExpired   = "expired"
)

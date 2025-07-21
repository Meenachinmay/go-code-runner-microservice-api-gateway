package problems

import (
	"context"
	"fmt"
	problemspb "go-code-runner-microservice/api-gateway/go-code-runner-microservice/proto/problems/v1"
	baseClient "go-code-runner-microservice/api-gateway/internal/service/grpc"
)

type Client struct {
	client problemspb.ProblemServiceClient
	base   *baseClient.Client
}

func NewClient(address string) (*Client, error) {
	base, err := baseClient.NewClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create problems client: %w", err)
	}

	return &Client{
		client: problemspb.NewProblemServiceClient(base.Connection()),
		base:   base,
	}, nil
}

func (c *Client) Close() error {
	return c.base.Close()
}

func (c *Client) GetProblem(ctx context.Context, id int32) (*problemspb.GetProblemResponse, error) {
	req := &problemspb.GetProblemRequest{
		Id: id,
	}

	return c.client.GetProblem(ctx, req)
}

func (c *Client) ListProblems(ctx context.Context) (*problemspb.ListProblemsResponse, error) {
	req := &problemspb.ListProblemsRequest{}

	return c.client.ListProblems(ctx, req)
}

func (c *Client) GetTestCasesByProblemID(ctx context.Context, problemID int32) (*problemspb.GetTestCasesByProblemIDResponse, error) {
	req := &problemspb.GetTestCasesByProblemIDRequest{
		ProblemId: problemID,
	}

	return c.client.GetTestCasesByProblemID(ctx, req)
}
